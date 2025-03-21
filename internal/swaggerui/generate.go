// This program downloads the dist assets for the current swagger-ui version and places them into the embed directory

//go:build ignore
// +build ignore

package main

import (
	"archive/tar"
	"bytes"
	"compress/gzip"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"io/fs"
	"log"
	"net/http"
	"os"
	"path/filepath"
	"regexp"
	"strings"
)

type releaseResp []struct {
	// TagName is a release tag name
	TagName string `json:"tag_name"`
}

func main() {
	log.SetFlags(0)
	releases := releaseResp{}
	// get the releases so we can download the latest one
	req, _ := http.NewRequest("GET", "https://api.github.com/repos/swagger-api/swagger-ui/releases", nil)
	req.Header.Set("Accept", "application/vnd.github.v3+json")
	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		log.Fatalf("error getting release list: %v", err)
	}
	if resp.StatusCode != http.StatusOK {
		log.Fatalf("got status [%s] on release list download", resp.Status)
	}
	if err := json.NewDecoder(resp.Body).Decode(&releases); err != nil {
		log.Fatalf("error decoding response: %v", err)
	}
	resp.Body.Close()
	if len(releases) == 0 {
		log.Fatal("somehow got no releases, nothing to do")
	}
	tag := releases[0].TagName

	current, err := os.ReadFile("current_version.txt")
	if err != nil {
		switch {
		case errors.Is(err, fs.ErrNotExist):
			// no problem, just do it
		default:
			log.Fatalf("unable to check version in current_version.txt: %v", err)
		}
	}
	cv := string(bytes.TrimRight(current, "\n"))

	if cv == tag {
		log.Print("version is current, nothing to do")
		os.Exit(0)
	}

	log.Printf("downloading release %s...", tag)

	resp, err = http.Get(fmt.Sprintf("https://github.com/swagger-api/swagger-ui/archive/%s.tar.gz", tag))
	if err != nil {
		log.Fatalf("error downloading release archive: %v", err)
	}
	defer resp.Body.Close()
	if resp.StatusCode != http.StatusOK {
		log.Fatalf("got status [%s] on release archive download", resp.Status)
	}
	zr, err := gzip.NewReader(resp.Body)
	if err != nil {
		log.Fatalf("error opening file as gzip: %v", err)
	}
	if err := os.RemoveAll("embed"); err != nil {
		log.Fatalf("error removing old embed directory")
	}
	if err := os.Mkdir("embed", 0o700); err != nil {
		log.Fatalf("error recreating embed directory")
	}
	tr := tar.NewReader(zr)
	for {
		header, err := tr.Next()
		if err == io.EOF {
			break
		}
		if err != nil {
			log.Fatalf("tar parsing error: %v", err)
		}
		if header.Typeflag == tar.TypeReg {
			// got a file, remove version directory
			fname := header.Name[strings.Index(header.Name, `/`):]
			if strings.HasPrefix(fname, `/dist`) {
				fname = strings.TrimPrefix(fname, `/dist`)
				out, err := os.Create(filepath.Join("embed", fname))
				if err != nil {
					log.Fatalf("error create output file: %v", err)
				}
				if _, err := io.Copy(out, tr); err != nil {
					log.Fatalf("error writing output file: %v", err)
				}
			}
		}
	}
	// replace the hard-coded JSON file with a generic file and disable the topbar
	initFile, err := os.ReadFile(filepath.Join("embed", "swagger-initializer.js"))
	if err != nil {
		log.Fatalf("error opening swagger-initializer.js for templating :%v", err)
	}
	// replace the url
	//newInit := strings.Replace(
	//string(initFile),
	//`url: "https://petstore.swagger.io/v2/swagger.json"`,
	//`url: "./swagger_spec"`,
	//-1,
	//)
	// Modify the swagger-initializer.js file
	if err != nil {
		log.Fatalf("error opening swagger-initializer.js for templating: %v", err)
	}

	// Replace the `url` value
	modifiedInit := regexp.MustCompile(`url:\s*"[^"]*"`).
		ReplaceAllLiteral(initFile, []byte(`url: "./swagger_spec"`))

	// Ensure required plugins are included
	if !regexp.MustCompile(`SwaggerUIStandalonePreset`).Match(modifiedInit) {
		modifiedInit = regexp.MustCompile(`presets:\s*\[.*?\]`).
			ReplaceAllLiteral(modifiedInit, []byte(`presets: [SwaggerUIBundle.presets.apis, SwaggerUIStandalonePreset]`))
	}

	if !regexp.MustCompile(`SwaggerUIBundle.plugins.DownloadUrl`).Match(modifiedInit) {
		modifiedInit = regexp.MustCompile(`plugins:\s*\[.*?\]`).
			ReplaceAllLiteral(modifiedInit, []byte(`plugins: [SwaggerUIBundle.plugins.DownloadUrl]`))
	}

	// Ensure the layout is set to StandaloneLayout
	modifiedInit = regexp.MustCompile(`layout:\s*"[^"]*"`).
		ReplaceAllLiteral(modifiedInit, []byte(`layout: "StandaloneLayout"`))

	// Write the modified file back
	newInitFile, err := os.Create(filepath.Join("embed", "swagger-initializer.js"))
	if err != nil {
		log.Fatalf("error re-creating swagger-initializer.js file: %v", err)
	}
	defer newInitFile.Close()

	if _, err := newInitFile.Write(modifiedInit); err != nil {
		log.Fatalf("unable to write to swagger-initializer.js: %v", err)
	}

	newcv, err := os.Create("current_version.txt")
	if err != nil {
		log.Fatalf("can't update current_version.txt: %v", err)
	}
	defer newcv.Close()
	newcv.WriteString(tag)
	log.Printf("updated swaggerui from %s => %s, please check templated swagger-initializer.js and retag repo", cv, tag)
}
