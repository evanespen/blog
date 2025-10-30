package main

import (
	"io"
	"os"
	"path/filepath"
)

// copyFile copies a file from src to dst.
// It reads the source file and writes the content to the destination file.
// Parameters:
//   - src: The path to the source file.
//   - dst: The path to the destination file.
//
// Returns:
//   - Any error encountered during the process.
func copyFile(src, dst string) error {
	sourceFile, err := os.Open(src)
	if err != nil {
		return err
	}
	defer sourceFile.Close()

	destFile, err := os.Create(dst)
	if err != nil {
		return err
	}
	defer destFile.Close()

	_, err = io.Copy(destFile, sourceFile)
	if err != nil {
		return err
	}

	return nil
}

// copyDir copies a directory and its contents from src to dst.
// It walks through the source directory and copies each file and subdirectory to the destination.
// Parameters:
//   - src: The path to the source directory.
//   - dst: The path to the destination directory.
//
// Returns:
//   - Any error encountered during the process.
func copyDir(src, dst string) error {
	return filepath.Walk(src, func(path string, info os.FileInfo, err error) error {
		if err != nil {
			return err
		}

		relPath, err := filepath.Rel(src, path)
		if err != nil {
			return err
		}

		destPath := filepath.Join(dst, relPath)

		if info.IsDir() {
			return os.MkdirAll(destPath, info.Mode())
		}

		return copyFile(path, destPath)
	})
}
