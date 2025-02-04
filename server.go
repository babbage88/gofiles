package main

import (
	"embed"
	"fmt"
	"net/http"
	"os"
	"strings"

	"github.com/babbage88/gofiles/internal/pretty"
	"github.com/joho/godotenv"
)

type GoFileServerOption func(p *GoFileServer)

type IGoFileServer interface {
	New(opts ...GoFileServerOption) *GoFileServer
	NewFromEnv(e string) *GoFileServer
	Start()
}

func WithEnvFile(s string) GoFileServerOption {
	return func(g *GoFileServer) {
		g.EnvFile = s
	}
}

func WithFilesDir(s string) GoFileServerOption {
	return func(g *GoFileServer) {
		g.FilesDir = s
	}
}

func WithListenAddr(s string) GoFileServerOption {
	return func(g *GoFileServer) {
		g.ListenAddr = s
	}
}

func WithStaticFiles(e *embed.FS) GoFileServerOption {
	return func(g *GoFileServer) {
		g.StaticFiles = *e
	}
}

func WithTemplateFiles(e *embed.FS) GoFileServerOption {
	return func(g *GoFileServer) {
		g.TemplateFiles = *e
	}
}

type GoFileServer struct {
	FilesDir      string   `json:"filesDir"`
	EnvFile       string   `json:"envFile"`
	ListenAddr    string   `json:"listenAddr"`
	StaticFiles   embed.FS `json:"staticFs"`
	TemplateFiles embed.FS `json:"templateFs"`
}

func New(opts ...GoFileServerOption) *GoFileServer {
	const (
		envFile  = ".env"
		filesDir = "/mnt/files/htfiles"
		listAddr = ":4100"
	)
	srv := &GoFileServer{
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

func NewFromEnv(e string) *GoFileServer {
	g := &GoFileServer{
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

func (g *GoFileServer) Start() {
	fs := http.FileServer(http.Dir(g.FilesDir))
	// Serve static files
	http.Handle("/static/", http.FileServer(http.FS(staticfs)))
	http.Handle("/files/", http.StripPrefix("/files/", fs))
	http.Handle("/test/", http.FileServer(http.FS(testfile)))
	http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		g.ServeTemplatesAndScanFiles(w, r)
	})

	pretty.Print("Listening on " + g.ListenAddr + "...")
	err := http.ListenAndServe(g.ListenAddr, nil)
	if err != nil {
		pretty.PrintError(err.Error())
	}
}
