package main

import (
	"io/fs"
	"path/filepath"
	"strings"
)

var ignorePaths = []string{
	".terragrunt-cache",
	".terraform",
}

func ignorePath(p string, ignore []string) bool {
	for _, v := range ignore {
		if strings.Contains(p, v) {
			return true
		}
	}

	return false
}

func FindFilesByExt(dir string, ext string) ([]string, error) {
	files := make([]string, 0)

	err := filepath.Walk(dir, func(path string, info fs.FileInfo, err error) error {
		if err != nil {
			return err
		}

		if filepath.Ext(path) == ext && !ignorePath(path, ignorePaths) {
			files = append(files, path)
		}

		return nil
	})

	if err != nil {
		return []string{}, err
	}

	return files, nil
}
