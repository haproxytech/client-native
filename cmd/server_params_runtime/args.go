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
	Licence     string
	LicencePath string
	Directory   string
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
				return fmt.Errorf("missing licence file after -l")
			}
			a.LicencePath = os.Args[i+1]
			i++
		case strings.HasPrefix(val, "--licence="):
			a.LicencePath = strings.TrimPrefix(val, "--licence=")
		default:
			selector = val
		}
	}

	if a.LicencePath != "" {
		var licence []byte
		licence, err = os.ReadFile(a.LicencePath)
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
		a.Licence = s.String()
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

			if strings.HasSuffix(path, serverParamFileName) {
				fmt.Println(path) //nolint:forbidigo
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
