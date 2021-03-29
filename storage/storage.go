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

package storage

import (
	"crypto"
	"crypto/ecdsa"
	"crypto/rsa"
	"crypto/x509"
	"encoding/pem"
	"fmt"
	"io"
	"io/ioutil"
	"os"
	"path/filepath"
	"strings"

	"github.com/google/renameio"

	conf "github.com/haproxytech/client-native/v2/configuration"
	"github.com/haproxytech/client-native/v2/misc"
)

type FileType string

const (
	MapsType             FileType = "maps"
	SSLType              FileType = "certs"
	SpoeType             FileType = "spoe"
	SpoeTransactionsType FileType = "spoe-transactions"
	BackupsType          FileType = "backups"
	TransactionsType     FileType = "transactions"
)

type Storage interface {
	GetAll() ([]string, error)
	Get(name string) (string, error)
	Delete(name string) error
	Replace(name string, config string) (string, error)
	Create(name string, contents io.ReadCloser) (string, error)
}

type storage struct {
	dirname  string
	fileType FileType
}

func New(dirname string, fileType FileType) (Storage, error) {
	dirname, err := misc.CheckOrCreateWritableDirectory(dirname)
	if err != nil {
		return nil, err
	}
	switch fileType { //nolint:exhaustive
	case MapsType, SSLType:
		return &storage{
			dirname:  dirname,
			fileType: fileType,
		}, nil
	default:
		return nil, fmt.Errorf("fileType is not valid %s", fileType)
	}
}

func (s *storage) GetAll() ([]string, error) {
	fis, err := readDir(s.dirname)
	if err != nil {
		return nil, err
	}
	if len(fis) == 0 {
		return nil, conf.NewConfError(conf.ErrObjectDoesNotExist, fmt.Sprintf("no files in dir: %s", s.dirname))
	}
	files := []string{}
	for _, fi := range fis {
		file := filepath.Join(s.dirname, fi.Name())
		switch s.fileType { //nolint:exhaustive
		case SSLType:
			raw, err := readFile(file)
			if err != nil {
				return nil, err
			}
			err = s.validatePEM([]byte(raw))
			if err != nil {
				return nil, err
			}
			files = append(files, file)
		case MapsType:
			files = append(files, file)
		}
	}
	return files, nil
}

func (s *storage) Get(name string) (string, error) {
	f, err := getFile(s.dirname, name)
	if err != nil {
		return "", err
	}
	if f == "" {
		return "", conf.NewConfError(conf.ErrObjectDoesNotExist, fmt.Sprintf("file %s doesn't exists in dir: %s", name, s.dirname))
	}
	return f, nil
}

func (s *storage) Delete(name string) error {
	f, err := s.Get(name)
	if err != nil {
		return err
	}
	return remove(f)
}

func (s storage) Replace(name string, config string) (string, error) {
	f, err := getFile(s.dirname, name)
	if err != nil {
		return "", err
	}
	switch s.fileType { //nolint:exhaustive
	case SSLType:
		err = s.validatePEM([]byte(config))
		if err != nil {
			return "", err
		}
	case MapsType:
	}

	err = renameio.WriteFile(f, []byte(config), 0644)
	if err != nil {
		return "", err
	}
	return f, nil
}

func (s *storage) Create(name string, readCloser io.ReadCloser) (string, error) {
	name = misc.SanitizeFilename(name)
	if s.fileType == MapsType {
		if !strings.HasSuffix(name, ".map") {
			name = fmt.Sprintf("%s.map", name)
		}
	}
	f := filepath.Join(s.dirname, name)
	if _, err := os.Stat(f); err == nil {
		return "", conf.NewConfError(conf.ErrObjectAlreadyExists, fmt.Sprintf("file %s already exists", f))
	}
	b, err := ioutil.ReadAll(readCloser)
	if err != nil {
		return "", err
	}

	switch s.fileType { //nolint:exhaustive
	case SSLType:
		err = s.validatePEM(b)
		if err != nil {
			return "", err
		}
	case MapsType:
	}
	err = renameio.WriteFile(f, b, 0644)
	if err != nil {
		return "", err
	}
	return f, nil
}

func readFile(name string) (string, error) {
	b, err := ioutil.ReadFile(name)
	if err != nil {
		return "", err
	}
	return string(b), nil
}

func getFile(dirname, name string) (string, error) {
	name = misc.SanitizeFilename(name)
	if name == "" {
		return "", fmt.Errorf("no file name")
	}
	f := filepath.Join(dirname, name)
	if _, err := os.Stat(f); os.IsNotExist(err) {
		return "", conf.NewConfError(conf.ErrObjectDoesNotExist, fmt.Sprintf("file %s doesn't exists in dir: %s", name, dirname))
	}
	return f, nil
}

func remove(name string) error {
	if name == "" {
		return conf.NewConfError(conf.ErrObjectDoesNotExist, fmt.Sprintf("file %s doesn't exist", name))
	}
	err := os.Remove(name)
	if err != nil {
		if !os.IsNotExist(err) {
			return err
		}
	}
	return nil
}

func readDir(dirname string) ([]os.FileInfo, error) {
	files, err := ioutil.ReadDir(dirname)
	if err != nil {
		return nil, err
	}
	return files, nil
}

func (s storage) validatePEM(raw []byte) error {
	crtPool := x509.NewCertPool()
	ok := crtPool.AppendCertsFromPEM(raw)
	if !ok {
		return fmt.Errorf("failed to parse certificate")
	}
	// HAProxy requires private and public key in same pem file
	hasPrivateKey := false
	for {
		block, rest := pem.Decode(raw)
		if block == nil {
			break
		}
		raw = rest

		if block.Type != "CERTIFICATE" {
			_, err := parsePrivateKey(block.Bytes)
			if err != nil {
				return fmt.Errorf("reading private key error: %w", err)
			}
			hasPrivateKey = true
		}
	}
	if !hasPrivateKey {
		return fmt.Errorf("missing private key")
	}
	return nil
}

func parsePrivateKey(der []byte) (crypto.PrivateKey, error) {
	if key, err := x509.ParsePKCS1PrivateKey(der); err == nil {
		return key, nil
	}
	if key, err := x509.ParsePKCS8PrivateKey(der); err == nil {
		switch key := key.(type) {
		case *rsa.PrivateKey, *ecdsa.PrivateKey:
			return key, nil
		default:
			return nil, fmt.Errorf("found unknown private key type in PKCS#8 wrapping")
		}
	}
	if key, err := x509.ParseECPrivateKey(der); err == nil {
		return key, nil
	}
	return nil, fmt.Errorf("failed to parse private key")
}
