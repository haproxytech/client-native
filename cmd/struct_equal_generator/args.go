package main

import (
	_ "embed"
	"errors"
	"fmt"
	"io/fs"
	"log"
	"os"
	"path"
	"path/filepath"
	"strings"
)

type Args struct {
	License     string
	Directory   string
	LicensePath string
	Selector    string
	Files       []string
}

func (a *Args) Parse() error { //nolint:gocognit
	selector, err := os.Getwd()
	if err != nil {
		log.Panic(err)
	}
	for i := 1; i < len(os.Args); i++ {
		val := os.Args[i]
		switch {
		case val == "-l":
			if i+1 >= len(os.Args) {
				return errors.New("missing licence file after -l")
			}
			a.LicensePath = os.Args[i+1]
			i++
		case strings.HasPrefix(val, "--licence="):
			a.LicensePath = strings.TrimPrefix(val, "--licence=")
		default:
			selector = val
		}
	}

	if a.LicensePath != "" {
		var licence []byte
		licence, err = os.ReadFile(a.LicensePath)
		if err != nil {
			return err
		}
		lines := strings.Split(string(licence), "\n")
		var s strings.Builder
		for _, line := range lines {
			s.WriteString("// ")
			s.WriteString(line)
			s.WriteString("\n")
		}
		a.License = s.String()
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
