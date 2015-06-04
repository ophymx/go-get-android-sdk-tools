package sdk

import (
	"archive/tar"
	"compress/gzip"
	"io"
	"os"
	"path/filepath"
)

func untar(file *os.File, size int64, path string, renderer ProgressRenderer) (err error) {
	// TODO(abic): wrap file reader with a renderer reader to show progress
	gunzip, err := gzip.NewReader(file)
	if err != nil {
		return
	}
	defer gunzip.Close()

	r := tar.NewReader(gunzip)
	var hdr *tar.Header
	for hdr, err = r.Next(); err == nil && hdr != nil; hdr, err = r.Next() {
		if err = extractAndWriteTarFile(hdr, r, path); err != nil {
			return
		}
	}
	if err == io.EOF {
		return nil
	}

	return
}

func extractAndWriteTarFile(header *tar.Header, reader io.Reader, dest string) (err error) {
	name, ok := stripPrefix(header.Name)
	if !ok {
		return
	}
	path := filepath.Join(dest, name)
	info := header.FileInfo()

	if info.IsDir() {
		return os.MkdirAll(path, info.Mode())
	}
	if info.Mode()&os.ModeSymlink != 0 {
		return os.Symlink(header.Linkname, path)
	}

	destFile, err := os.OpenFile(path, os.O_WRONLY|os.O_CREATE|os.O_TRUNC, info.Mode())
	if err != nil {
		return
	}
	defer destFile.Close()

	_, err = io.Copy(destFile, reader)
	return
}
