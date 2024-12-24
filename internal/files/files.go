package files

import (
	"fmt"
	"io/fs"
	"os"
	"path"
	"path/filepath"

	"github.com/babbage88/gofiles/internal/pretty"
)

type FileInfo struct {
	FullName string
	Size     int64
	IsDir    bool
}

// ListFiles recursively lists all files and directories in the specified root path using WalkDir.
func ListFiles(rootPath string) ([]FileInfo, error) {
	var files []FileInfo

	// WalkDir to traverse the directory tree
	err := filepath.WalkDir(rootPath, func(path string, d os.DirEntry, err error) error {
		if err != nil {
			return err
		}

		// Retrieve the file information
		info, err := d.Info()
		if err != nil {
			return err
		}

		files = append(files, FileInfo{
			FullName: path,
			Size:     info.Size(),
			IsDir:    d.IsDir(),
		})
		return nil
	})

	if err != nil {
		return nil, err
	}
	return files, nil
}

func GlobAllFiles(dir string, recurse bool) []string {
	root := os.DirFS(dir)
	pattern := "*"

	if recurse {
		pattern = "*/**"
	}
	allFiles, err := fs.Glob(root, pattern)

	if err != nil {
		pretty.PrintError(err.Error())
	}

	var files []string
	for _, v := range allFiles {
		files = append(files, path.Join(dir, v))
	}
	return files
}

func PrintAllFiles(dir string, recurse bool) {
	files, err := ListFiles(dir)
	if err != nil {
		pretty.PrintError(err.Error())
	}
	for _, file := range files {
		pretty.Print(fmt.Sprint("Name: ", file.FullName, "Size: ", file.Size))
	}
}
