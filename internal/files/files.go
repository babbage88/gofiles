package files

import (
	"io/fs"
	"os"
	"path"

	"github.com/babbage88/gofiles/internal/pretty"
)

func ListAllFiles(dir string) []string {
	root := os.DirFS(dir)

	allFiles, err := fs.Glob(root, "*")

	if err != nil {
		pretty.PrintError(err.Error())
	}

	var files []string
	for _, v := range allFiles {
		files = append(files, path.Join(dir, v))
	}
	return files
}

func PrintAllFiles(dir string) {
	files := ListAllFiles(dir)
	for _, file := range files {
		pretty.Print(file)
	}
}
