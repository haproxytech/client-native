package misc

import (
	"fmt"
	"io/ioutil"
	"os"
	"path/filepath"
)

func CheckOrCreateWritableDirectory(dirname string) (string, error) {

	if dirname == "" {
		return dirname, fmt.Errorf("no storage dir or empty dirname specified")
	}

	dirname = filepath.Clean(dirname)

	info, err := os.Stat(dirname)
	if err != nil {
		if os.IsNotExist(err) {
			// try to create the directory
			err = os.MkdirAll(dirname, 0750)
			if err != nil {
				return "", fmt.Errorf("missing directory, unable to create: %w", err)
			}
		} else {
			return "", fmt.Errorf("accessing directory: %w", err)
		}
	} else { // file or directory exists, make sure it's a directory we have access to
		if !info.IsDir() {
			return "", fmt.Errorf("specified storage directory %s already exists and is a file, refusing to overwrite", dirname)
		} else {
			f, err := ioutil.TempFile(dirname, ".dataplaneapi")
			if err != nil {
				return "", fmt.Errorf("error using storage directory: %w", err)
			}
			os.Remove(f.Name())
		}
	}
	return dirname, nil
}
