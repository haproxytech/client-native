package main

import (
	_ "embed"
	"fmt"
	"io/fs"
	"log"
	"os"
	"path"
	"path/filepath"
	"strings"
)

type Args struct {
	Directory string
	Selector  string
	Files     []string
}

func (a *Args) Parse() error { //nolint:gocognit,unparam
	selector, err := os.Getwd()
	if err != nil {
		log.Panic(err)
	}
	for i := 1; i < len(os.Args); i++ {
		selector = os.Args[i]
	}

	isDirectory := false
	file, err := os.Open(selector)
	if err == nil {
		var fileInfo fs.FileInfo
		fileInfo, err = file.Stat()
		if err == nil {
			isDirectory = fileInfo.IsDir()
		}
	}

	if selector == "*" || selector == "." || isDirectory {
		err = filepath.Walk(selector, func(path string, info os.FileInfo, err error) error {
			if err != nil {
				return err
			}
			if info.IsDir() {
				return nil
			}
			if strings.HasSuffix(path, "_compare.go") {
				return nil
			}
			if strings.HasSuffix(path, "_test.go") {
				return nil
			}
			if strings.HasSuffix(path, "_easyjson.go") {
				return nil
			}
			if strings.HasSuffix(path, "_generated.go") {
				return nil
			}
			fmt.Println(path) //nolint:forbidigo
			if strings.HasSuffix(path, ".go") {
				a.Files = append(a.Files, path)
			}
			return nil
		})
		if err != nil {
			log.Panic(err)
		}
	} else {
		a.Files = append(a.Files, selector)
	}
	a.Selector = selector
	if isDirectory {
		a.Directory = selector
	} else {
		a.Directory = path.Dir(selector)
	}

	return nil
}
