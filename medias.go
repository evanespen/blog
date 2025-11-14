package main

import (
	"errors"
	"io/fs"
	"log"
	"os"
	"path"
	"path/filepath"
	"strings"
)

// copyMedias copies media files from the posts directory to the build/medias directory.
// It creates the build/medias directory if it doesn't exist.
// It walks through the posts directory and copies all .jpg, .jpeg, .png, and .mp4 files.
// Returns any error encountered during the process.
func copyMedias() error {
	if err := os.MkdirAll("build/medias", os.ModePerm); err != nil {
		log.Fatal("Error creating directory:", err)
		return err
	}

	filepath.WalkDir("posts/", func(s string, d fs.DirEntry, err error) error {
		if filepath.Ext(s) == ".jpg" || filepath.Ext(s) == ".jpeg" || filepath.Ext(s) == ".png" || filepath.Ext(s) == ".mp4" {
			newPath := strings.ReplaceAll(s, "posts/", "build/medias/")

			if _, err := os.Stat(newPath); err == nil {
				log.Println("Media", newPath, "already handled")
			} else if errors.Is(err, os.ErrNotExist) {
				err := os.Link(s, newPath)
				if err != nil {
					log.Fatal("Failed to handle media", s)
				}
				log.Println("Copyied media from", s, "to", newPath)
			}
		}

		if d.IsDir() {
			log.Println("Skipping directory:", s)
			copyDir(d.Name(), path.Join("/build/medias", d.Name()))
			return nil
		}
		return nil
	})
	return nil
}
