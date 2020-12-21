package misc

import (
	"fmt"
	"io/ioutil"
	"os"
	"path/filepath"
)

func CheckOrCreateWritableDirectory(dirname string) (string, error) {

	dirname = filepath.Clean(dirname)

	if dirname == "" {
		return dirname, fmt.Errorf("no storage dir or empty dirname specified")
	}

	info, err := os.Stat(dirname)
	if err != nil {
		if os.IsNotExist(err) {
			// try to create the directory
			err = os.MkdirAll(dirname, 0750)
			if err != nil {
				return "", fmt.Errorf("missing directory, unable to create")
			}
		} else { // file or directory exists, make sure it's a directory we have access to
			if !info.IsDir() {
				return "", fmt.Errorf("specified storage directory %s already exists and is a file, refusing to overwrite", dirname)
			} else {
				f, err := ioutil.TempFile(dirname, ".dataplaneapi")
				if err != nil {
					return "", fmt.Errorf("error using storage directory %s", dirname)
				}
				os.Remove(f.Name())
			}
		}
	}
	return dirname, nil
}
