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
	FullName     string `json:"fullName"`     // Absolute path
	Size         int64  `json:"size"`         // Size in bytes
	IsDir        bool   `json:"isDir"`        // Is it a directory?
	RelativeName string `json:"relativeName"` // Path relative to fdldir
}

// ListOnlyFiles lists only files (not directories) in the specified rootPath.
func ListOnlyFiles(rootPath string) ([]FileInfo, error) {
	var files []FileInfo

	// WalkDir to traverse the directory tree
	err := filepath.WalkDir(rootPath, func(path string, d os.DirEntry, err error) error {
		if err != nil {
			msg := fmt.Sprint("Error walking dir: ", err.Error())
			pretty.PrintError(msg)
			return err
		}

		// Retrieve the file information
		info, err := d.Info()
		if err != nil {
			msg := fmt.Sprint("Error getting fileInfo: ", err.Error())
			pretty.PrintError(msg)
			return err
		}

		// Only append if it's not a directory
		if !d.IsDir() {
			relativeName, err := filepath.Rel(rootPath, path)
			if err != nil {
				msg := fmt.Sprintf("Error calculating relative name for %s: %v", path, err)
				pretty.PrintError(msg)
				return err
			}

			files = append(files, FileInfo{
				FullName:     path,
				Size:         info.Size(),
				IsDir:        false,
				RelativeName: relativeName,
			})
		}

		return nil
	})
	if err != nil {
		return nil, err
	}

	return files, nil
}

// ListFiles recursively lists all files and directories in the specified root path using WalkDir.
func ListFiles(rootPath string) ([]FileInfo, error) {
	var files []FileInfo

	// WalkDir to traverse the directory tree
	err := filepath.WalkDir(rootPath, func(path string, d os.DirEntry, err error) error {
		if err != nil {
			msg := fmt.Sprint("Error walking dir: ", err.Error())
			pretty.PrintError(msg)
		}

		// Retrieve the file information
		info, err := d.Info()
		if err != nil {
			msg := fmt.Sprint("Error getting fileInfo: ", err.Error())
			pretty.PrintError(msg)

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
