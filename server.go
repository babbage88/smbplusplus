package main

import (
	"embed"
	"fmt"
	"net/http"
	"os"
	"strings"

	"github.com/babbage88/smbplusplus/internal/pretty"
	"github.com/joho/godotenv"
)

type SmbPlusSquaredServerOption func(p *SmbPlusSquaredServer)

type ISmbPlusSquaredServer interface {
	New(opts ...SmbPlusSquaredServerOption) *SmbPlusSquaredServer
	NewFromEnv(e string) *SmbPlusSquaredServer
	Start()
}

func WithEnvFile(s string) SmbPlusSquaredServerOption {
	return func(g *SmbPlusSquaredServer) {
		g.EnvFile = s
	}
}

func WithFilesDir(s string) SmbPlusSquaredServerOption {
	return func(g *SmbPlusSquaredServer) {
		g.FilesDir = s
	}
}

func WithListenAddr(s string) SmbPlusSquaredServerOption {
	return func(g *SmbPlusSquaredServer) {
		g.ListenAddr = s
	}
}

func WithStaticFiles(e *embed.FS) SmbPlusSquaredServerOption {
	return func(g *SmbPlusSquaredServer) {
		g.StaticFiles = *e
	}
}

func WithTemplateFiles(e *embed.FS) SmbPlusSquaredServerOption {
	return func(g *SmbPlusSquaredServer) {
		g.TemplateFiles = *e
	}
}

type SmbPlusSquaredServer struct {
	FilesDir      string   `json:"filesDir"`
	EnvFile       string   `json:"envFile"`
	ListenAddr    string   `json:"listenAddr"`
	StaticFiles   embed.FS `json:"staticFs"`
	TemplateFiles embed.FS `json:"templateFs"`
}

func New(opts ...SmbPlusSquaredServerOption) *SmbPlusSquaredServer {
	const (
		envFile  = ".env"
		filesDir = "/mnt/files/htfiles"
		listAddr = ":4100"
	)
	srv := &SmbPlusSquaredServer{
		EnvFile:       envFile,
		FilesDir:      filesDir,
		ListenAddr:    listAddr,
		StaticFiles:   staticfs,
		TemplateFiles: viewtmpl,
	}

	for _, opt := range opts {
		opt(srv)
	}

	return srv
}

func NewFromEnv(e string) *SmbPlusSquaredServer {
	g := &SmbPlusSquaredServer{
		EnvFile:       e,
		TemplateFiles: viewtmpl,
		StaticFiles:   staticfs,
	}

	err := godotenv.Load(e)
	if err != nil {
		msg := fmt.Sprint("Error loading .env file: ", err.Error())
		pretty.PrintError(msg)
	}
	g.FilesDir = os.Getenv("FILES_DIR")
	port := os.Getenv("LISTEN_PORT")
	if strings.HasPrefix(port, ":") {
		g.ListenAddr = port
	} else {
		g.ListenAddr = fmt.Sprint(":", os.Getenv("LISTEN_PORT"))
	}

	pretty.Print(g.FilesDir)
	pretty.Print(g.ListenAddr)
	return g
}

func (g *SmbPlusSquaredServer) Start() {
	fs := http.FileServer(http.Dir(g.FilesDir))
	// Serve static files
	http.Handle("/static/", http.FileServer(http.FS(staticfs)))
	http.Handle("/files/", http.StripPrefix("/files/", fs))
	http.Handle("/test/", http.FileServer(http.FS(testfile)))
	http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		g.ServeTemplatesAndScanFiles(w, r)
	})

	pretty.Print("Listening on " + g.ListenAddr + "...")
	err := http.ListenAndServeTLS(g.ListenAddr, "cert.pem", "privkey.pem", nil)
	if err != nil {
		pretty.PrintError(err.Error())
	}
}
