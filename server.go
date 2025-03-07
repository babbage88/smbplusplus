package main

import (
	"fmt"
	"net/http"
	"os"
	"strings"

	"github.com/babbage88/smbplusplus/internal/pretty"
	"github.com/joho/godotenv"
	"github.com/minio/minio-go/v7"
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

func WithSquaredShares(s []SquaredShare) SmbPlusSquaredServerOption {
	return func(g *SmbPlusSquaredServer) {
		g.SquaredShares = s
	}
}

func WithListenAddr(s string) SmbPlusSquaredServerOption {
	return func(g *SmbPlusSquaredServer) {
		g.ListenAddr = s
	}
}

type SquaredShare struct {
	LocalPath string           `json:"localPath"`
	S3Bucket  minio.BucketInfo `json:"s3Bucket"`
}

type SmbPlusSquaredServer struct {
	SquaredShares []SquaredShare `json:"localShares"`
	EnvFile       string         `json:"envFile"`
	ListenAddr    string         `json:"listenAddr"`
}

func New(opts ...SmbPlusSquaredServerOption) *SmbPlusSquaredServer {
	const (
		envFile  = ".env"
		listAddr = ":4200"
	)
	srv := &SmbPlusSquaredServer{
		EnvFile:    envFile,
		ListenAddr: listAddr,
	}

	for _, opt := range opts {
		opt(srv)
	}

	return srv
}

func NewFromEnv(e string) *SmbPlusSquaredServer {
	g := &SmbPlusSquaredServer{
		EnvFile: e,
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
	http.Handle("/files/", http.StripPrefix("/files/", fs))
	http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		g.ServeTemplatesAndScanFiles(w, r)
	})

	pretty.Print("Listening on " + g.ListenAddr + "...")
	err := http.ListenAndServeTLS(g.ListenAddr, "cert.pem", "privkey.pem", nil)
	if err != nil {
		pretty.PrintError(err.Error())
	}
}
