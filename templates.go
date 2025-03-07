package main

import (
	"fmt"
	"html/template"
	"net/http"

	"github.com/babbage88/smbplusplus/internal/files"
	"github.com/babbage88/smbplusplus/internal/pretty"
)

type TemplateData struct {
	Files []files.FileInfo
}

func (g *SmbPlusSquaredServer) ServeTemplatesAndScanFiles(w http.ResponseWriter, r *http.Request) {
	scannedFiles, err := files.ListOnlyFiles(g.FilesDir)
	if err != nil {
		pretty.PrintError(err.Error())
	}
	msg := fmt.Sprintf("Total count: %d", len(scannedFiles))
	pretty.Print(msg)

	tmpl, err := template.New("layout").ParseFS(viewtmpl, "templates/*.html")
	if err != nil {
		http.Error(w, fmt.Sprintf("Error parsing templates: %v", err), http.StatusInternalServerError)
		return
	}

	fmt.Println("Scanned Files:")
	for _, file := range scannedFiles {
		msg := fmt.Sprintf("FullName: %s, RelativeName: %s, Size: %d, IsDir: %v\n", file.FullName, file.RelativeName, file.Size, file.IsDir)
		pretty.Print(msg)
	}

	// Prepare data
	data := TemplateData{Files: scannedFiles}
	fmt.Printf("TemplateData has %d files\n", len(data.Files))

	// Render template
	err = tmpl.ExecuteTemplate(w, "layout", data)
	if err != nil {
		http.Error(w, fmt.Sprintf("Error rendering template: %v", err), http.StatusInternalServerError)
		return
	}
}
