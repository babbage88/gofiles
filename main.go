package main

import (
	"fmt"
	"html/template"
	"net/http"
	"os"
	"path/filepath"

	"github.com/babbage88/gofiles/internal/files"
	"github.com/babbage88/gofiles/internal/pretty"
	"github.com/joho/godotenv"
)

type TemplateData struct {
	Files []files.FileInfo
}

func serveFilesTemplate(w http.ResponseWriter, r *http.Request, scannedFiles []files.FileInfo) {
	lp := filepath.Join("templates", "layout.html")
	fp := filepath.Join("templates", "example.html")

	tmpl, err := template.ParseFiles(lp, fp)
	if err != nil {
		http.Error(w, "Error loading templates", http.StatusInternalServerError)
		return
	}

	data := TemplateData{Files: scannedFiles}
	err = tmpl.ExecuteTemplate(w, "layout", data)
	if err != nil {
		http.Error(w, "Error rendering template", http.StatusInternalServerError)
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
