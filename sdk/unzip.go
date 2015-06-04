package sdk

import (
	"archive/zip"
	"io"
	"os"
	"path/filepath"
	"strings"
)

func unzip(file *os.File, size int64, dest string, renderer ProgressRenderer) (err error) {
	r, err := zip.NewReader(file, size)
	if err != nil {
		return
	}

	if err = os.MkdirAll(dest, 0755); err != nil {
		return
	}

	files := float32(len(r.File))

	for i, f := range r.File {
		if err = extractAndWriteZipFile(f, dest); err != nil {
			return
		}
		renderer.Progress(float32(i) / files)
	}
	return
}

func extractAndWriteZipFile(f *zip.File, dest string) (err error) {
	name, ok := stripPrefix(f.Name)
	if !ok {
		return
	}
	path := filepath.Join(dest, name)
	info := f.FileInfo()

	if info.IsDir() {
		return os.MkdirAll(path, info.Mode())
	}

	srcFile, err := f.Open()
	if err != nil {
		return
	}
	defer srcFile.Close()

	destFile, err := os.OpenFile(path, os.O_WRONLY|os.O_CREATE|os.O_TRUNC, info.Mode())
	if err != nil {
		return
	}
	defer destFile.Close()

	_, err = io.Copy(destFile, srcFile)
	return
}

func stripPrefix(path string) (newPath string, ok bool) {
	newPath = filepath.Clean(path)
	if index := strings.IndexRune(newPath, os.PathSeparator); index != -1 {
		return newPath[index+1 : len(newPath)], true
	}
	return "", false
}
