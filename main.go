package main

import (
	"log/slog"

	"github.com/babbage88/gofiles/internal/files"
	"github.com/babbage88/gofiles/internal/pretty"
)

func main() {
	pretty.Print("Test Print All Files")
	slog.Info("Test Print All Files")
	files.PrintAllFiles("/home/jtrahan")
}
