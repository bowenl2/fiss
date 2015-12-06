package main

import (
	"archive/zip"
	"fmt"
	"io"
	"os"
	"path/filepath"
	"time"
)

func MakeArchive(basePath string) (string, error) {
	zipPath := filepath.Join(
		os.TempDir(),
		fmt.Sprintf("%v.zip", time.Now().UnixNano()))
	zipFile, err := os.Create(zipPath)
	if err != nil {
		return "", err
	}
	defer zipFile.Close()

	zipWriter := zip.NewWriter(zipFile)
	defer zipWriter.Close()

	err = filepath.Walk(basePath,
		func(p string, info os.FileInfo, err error) error {
			if !info.Mode().IsRegular() || err != nil {
				return nil
			}

			relPath, err := filepath.Rel(basePath, p)
			if err != nil {
				return err
			}
			zipPath := filepath.ToSlash(relPath)

			fileWriter, err := zipWriter.Create(zipPath)
			if err != nil {
				return err
			}

			input, err := os.Open(p)
			if err != nil {
				return err
			}
			defer input.Close()

			_, err = io.Copy(fileWriter, input)
			return err
		})

	return zipPath, err
}
