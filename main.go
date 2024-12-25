package main

import (
	"embed"
	"fmt"
	"html/template"
	"net/http"
	"os"

	"github.com/babbage88/gofiles/internal/files"
	"github.com/babbage88/gofiles/internal/pretty"
	"github.com/joho/godotenv"
)

//go:embed all:templates
var viewtmpl embed.FS

//go:embed static/*
var staticfs embed.FS

type TemplateData struct {
	Files []files.FileInfo
}

func serveFilesTemplate(w http.ResponseWriter, r *http.Request, scannedFiles []files.FileInfo) {
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

func main() {
	err := godotenv.Load(".env")
	if err != nil {
		msg := fmt.Sprint("Error loading .env file: ", err.Error())
		pretty.PrintError(msg)
	}
	fdldir := os.Getenv("FILES_DIR")
	srvport := fmt.Sprint(":", os.Getenv("LISTEN_PORT"))

	scanned, err := files.ListOnlyFiles(fdldir)
	if err != nil {
		pretty.PrintError(fmt.Sprintf("Error scanning files: %v", err))
		return
	}
	pretty.Print(fdldir)
	pretty.Print(srvport)

	msg := fmt.Sprintf("Total count: %d", len(scanned))
	pretty.Print(msg)

	fs := http.FileServer(http.Dir(fdldir))
	// Serve static files
	http.Handle("/static/", http.FileServer(http.FS(staticfs)))
	http.Handle("/files/", http.StripPrefix("/files/", fs))
	http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		serveFilesTemplate(w, r, scanned)
	})

	pretty.Print("Listening on " + srvport + "...")
	err = http.ListenAndServe(srvport, nil)
	if err != nil {
		pretty.PrintError(err.Error())
	}
}
