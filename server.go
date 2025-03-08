package main

import (
	"log/slog"
	"net/http"
	"os"

	"github.com/babbage88/smbplusplus/internal/swaggerui"
	"github.com/minio/minio-go/v7"
	"github.com/prometheus/client_golang/prometheus/promhttp"
	"gopkg.in/yaml.v3"
)

type SmbPlusSquaredServerOption func(p *SmbPlusSquaredServer)

type ISmbPlusSquaredServer interface {
	New(opts ...SmbPlusSquaredServerOption) *SmbPlusSquaredServer
	NewFromEnv(e string) *SmbPlusSquaredServer
	Start()
}

func WithEnvFile(s string) SmbPlusSquaredServerOption {
	return func(g *SmbPlusSquaredServer) {
		g.ConfigFile = s
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
	LocalPath string           `json:"localPath" yaml:"localPath"`
	S3Bucket  minio.BucketInfo `json:"s3bucket" yaml:"s3Bucket"`
}

type SmbPlusSquaredServer struct {
	SquaredShares []SquaredShare `json:"localShares" yaml:"localShares"`
	ListenAddr    string         `json:"listenAddr" yaml:"listenAddress"`
	ConfigFile    string         `json:"configFile" yaml:"configFile"`
}

func New(opts ...SmbPlusSquaredServerOption) *SmbPlusSquaredServer {
	const (
		envFile  = ".env"
		listAddr = ":4200"
	)
	srv := &SmbPlusSquaredServer{
		ListenAddr: listAddr,
	}

	for _, opt := range opts {
		opt(srv)
	}

	return srv
}

func NewSmbPlusServerFromConfig(config string) (*SmbPlusSquaredServer, error) {
	var server *SmbPlusSquaredServer
	configFile, err := os.ReadFile(config)
	if err != nil {
		slog.Error("error when attemping to read/open config file.", slog.String("error", err.Error()))
		return &SmbPlusSquaredServer{}, err

	}
	err = yaml.Unmarshal(configFile, server)
	if err != nil {
		slog.Error("error when attemping to parse config file.", slog.String("error", err.Error()))
		return &SmbPlusSquaredServer{}, err
	}
	return server, err
}

func (g *SmbPlusSquaredServer) Start() error {
	mux := http.NewServeMux()
	mux.Handle("/metrics", promhttp.Handler())
	// Add Swagger UI handler
	mux.Handle("/swaggerui/", http.StripPrefix("/swaggerui", swaggerui.ServeSwaggerUI(swaggerSpec)))
	slog.Info("Starting api http server.")
	err := http.ListenAndServe(g.ListenAddr, mux)
	if err != nil {
		slog.Error("Failed to start server", slog.String("Error", err.Error()))
	}
	return err
}
