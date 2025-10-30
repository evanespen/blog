package main

import (
	"fmt"
	"log"
	"os"
	"path/filepath"

	"github.com/bep/godartsass/v2"
)

type importResolver struct {
	baseDir string
}

func (t importResolver) CanonicalizeURL(url string) (string, error) {
	fullPath := filepath.Join("/", url)
	return fullPath, nil
}

func (t importResolver) Load(url string) (godartsass.Import, error) {
	fullPath := filepath.Join(t.baseDir, url)
	data, err := os.ReadFile(fullPath)
	if err != nil {
		return godartsass.Import{}, fmt.Errorf("cannot read %s: %v", fullPath, err)
	}

	return godartsass.Import{
		Content:      string(data),
		SourceSyntax: godartsass.SourceSyntaxSCSS,
	}, nil
}

func compileSCSS() (string, error) {
	data, _ := os.ReadFile("styles/main.scss")

	args := godartsass.Args{
		Source:       string(data),
		URL:          "styles/main.scss",
		IncludePaths: []string{"styles/"},
		ImportResolver: importResolver{
			baseDir: "styles/",
		},
		OutputStyle:             godartsass.OutputStyleExpanded,
		EnableSourceMap:         false,
		SourceMapIncludeSources: false,
	}

	transpiler, err := godartsass.Start(godartsass.Options{})
	if err != nil {
		log.Fatal(err)
	}

	css, err := transpiler.Execute(args)

	if err != nil {
		log.Fatal(err)
	}

	log.Println("CSS compiled")

	return css.CSS, nil
}
