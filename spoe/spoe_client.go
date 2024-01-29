// Copyright 2019 HAProxy Technologies
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//     http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.
//

package spoe

import (
	"bytes"
	"fmt"
	"io"
	"os"
	"path"
	"path/filepath"
	"sync"

	"github.com/google/renameio"

	conf "github.com/haproxytech/client-native/v6/configuration"
	"github.com/haproxytech/client-native/v6/misc"
)

type Spoe interface {
	GetSingleSpoe(name string) (*SingleSpoe, error)
	GetAll() ([]string, error)
	Delete(name string) error
	Create(name string, readCloser io.ReadCloser) (string, error)
	Get(name string) (string, error)
}

type spoeclient struct {
	clients    map[string]*SingleSpoe
	initParams Params
	mu         sync.Mutex
}

func NewSpoe(params Params) (Spoe, error) {
	var err error
	params.SpoeDir, err = misc.CheckOrCreateWritableDirectory(params.SpoeDir)
	if err != nil {
		return nil, err
	}
	params.TransactionDir, err = misc.CheckOrCreateWritableDirectory(params.TransactionDir)
	if err != nil {
		return nil, err
	}
	c := spoeclient{}

	files, _ := c.getSpoeFiles(params.SpoeDir)

	prm := Params{
		TransactionDir:         params.TransactionDir,
		BackupsNumber:          params.BackupsNumber,
		PersistentTransactions: params.PersistentTransactions,
		UseValidation:          params.UseValidation,
		SpoeDir:                params.SpoeDir,
		SkipFailedTransactions: params.PersistentTransactions,
	}
	c.clients = make(map[string]*SingleSpoe)
	for _, f := range files {
		err := c.addClient(f, prm)
		if err != nil {
			return nil, err
		}
	}
	c.initParams = params
	return &c, nil
}

// GetSingleSpoe returns single SPOE client by its file name if found, otherwise returns an error
func (c *spoeclient) GetSingleSpoe(name string) (*SingleSpoe, error) {
	c.mu.Lock()
	client, ok := c.clients[name]
	c.mu.Unlock()
	if ok {
		return client, nil
	}
	return nil, fmt.Errorf("client %s not configured", name)
}

// GetAll returns array of configured spoe files or nil if any not found.
func (c *spoeclient) GetAll() ([]string, error) {
	files, err := c.getSpoeFiles(c.initParams.SpoeDir)
	if err != nil {
		return nil, err
	}
	return files, nil
}

// Delete deletes one SPOE file by its name
func (c *spoeclient) Delete(name string) error {
	name = c.setFilePath(name)
	if err := os.Remove(name); err != nil {
		if os.IsNotExist(err) {
			return nil
		}
		return err
	}
	c.removeClient(name)
	return nil
}

// Create creates a new SPOE file with its entries, returns an error if file already exists
func (c *spoeclient) Create(name string, readCloser io.ReadCloser) (string, error) {
	name = c.setFilePath(name)
	if _, err := os.Stat(name); err == nil {
		return "", fmt.Errorf("file %s already exists. You should delete an existing file first", name)
	}
	b, err := io.ReadAll(readCloser)
	if err != nil {
		return "", err
	}
	version := "# _version="
	if !bytes.Contains(b, []byte(version)) {
		version = fmt.Sprintf("%s%s", version, "1\n")
		b = append([]byte(version), b...)
	}
	err = renameio.WriteFile(name, b, 0o644)
	if err != nil {
		return "", err
	}
	err = c.addClient(name, c.initParams)
	if err != nil {
		return "", err
	}
	return name, nil
}

// Get returns single file name or error if not found.
func (c *spoeclient) Get(name string) (string, error) {
	files, err := c.GetAll()
	if err != nil {
		return "", err
	}

	for _, f := range files {
		fn := c.getFileName(f)
		if fn == name {
			client, err := c.GetSingleSpoe(name)
			if err != nil {
				return "", conf.NewConfError(conf.ErrObjectDoesNotExist, fmt.Sprintf("spoe file %s does not exist", name))
			}
			return client.Parser.String(), nil
		}
	}
	return "", conf.NewConfError(conf.ErrObjectDoesNotExist, fmt.Sprintf("spoe file %s does not exist", name))
}

// getSpoeFiles returns files list in dir or error
func (c *spoeclient) getSpoeFiles(dir string) ([]string, error) {
	fis, err := os.ReadDir(dir)
	if err != nil {
		return nil, err
	}
	files := []string{}
	for _, f := range fis {
		path := path.Join(dir, "/", f.Name())
		fi, err := os.Stat(path)
		if err != nil {
			return nil, err
		}
		if fi.Mode().IsRegular() {
			files = append(files, path)
		}
	}
	if len(files) == 0 {
		return files, fmt.Errorf("no SPOE files in dir: %s", dir)
	}
	return files, nil
}

func (c *spoeclient) setFilePath(name string) string {
	return path.Join(c.initParams.SpoeDir, "/", name)
}

// getFileName returns file name
func (c *spoeclient) getFileName(fullPath string) string {
	return filepath.Base(fullPath)
}

func (c *spoeclient) addClient(file string, params Params) error {
	params.ConfigurationFile = file
	client, err := newSingleSpoe(params)
	if err != nil {
		return err
	}
	name := c.getFileName(file)
	c.mu.Lock()
	c.clients[name] = client
	c.mu.Unlock()
	return nil
}

func (c *spoeclient) removeClient(name string) {
	name = c.getFileName(name)
	delete(c.clients, name)
}
