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
	"io/fs"
	"os"
	"path/filepath"
	"strings"

	"github.com/google/renameio"

	conf "github.com/haproxytech/client-native/v4/configuration"
	"github.com/haproxytech/client-native/v4/misc"
)

type FileType string

const (
	GeneralType          FileType = "general"
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
	GetContents(name string) (string, error)
	GetRawContents(name string) (io.ReadCloser, error)
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
	case MapsType, SSLType, GeneralType, BackupsType:
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
	files := []string{}
	for _, fi := range fis {
		file := filepath.Join(s.dirname, fi.Name())
		switch s.fileType { //nolint:exhaustive
		case SSLType:
			noErrors := true
			raw, err := readFile(file)
			if err != nil {
				noErrors = noErrors && false
			}
			err = s.validatePEM([]byte(raw))
			if err != nil {
				noErrors = noErrors && false
			}
			if noErrors {
				files = append(files, file)
			}
		case MapsType, GeneralType:
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
		return "", conf.NewConfError(conf.ErrObjectDoesNotExist, fmt.Sprintf("file %s doesn't exist in dir: %s", name, s.dirname))
	}
	return f, nil
}

func (s *storage) GetContents(name string) (string, error) {
	f, err := getFile(s.dirname, name)
	if err != nil {
		return "", err
	}
	return readFile(f)
}

func (s *storage) GetRawContents(name string) (io.ReadCloser, error) {
	fname, err := getFile(s.dirname, name)
	if err != nil {
		return nil, err
	}
	return os.Open(fname)
}

func (s *storage) Delete(name string) error {
	f, err := s.Get(name)
	if err != nil {
		return err
	}
	return s.remove(f)
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

	err = renameio.WriteFile(f, []byte(config), 0o644)
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

	switch s.fileType { //nolint:exhaustive
	case SSLType:
		return s.createSSL(f, readCloser)
	case MapsType, GeneralType:
		return s.createFile(f, readCloser)
	}
	return f, nil
}

func (s *storage) createSSL(name string, readCloser io.ReadCloser) (string, error) {
	b, err := io.ReadAll(readCloser)
	if err != nil {
		return "", err
	}
	err = s.validatePEM(b)
	if err != nil {
		return "", err
	}
	err = renameio.WriteFile(name, b, 0o644)
	if err != nil {
		return "", err
	}
	return name, nil
}

func (s *storage) createFile(name string, readCloser io.ReadCloser) (string, error) {
	b, err := io.ReadAll(readCloser)
	if err != nil {
		return "", err
	}
	err = renameio.WriteFile(name, b, 0o644)
	if err != nil {
		return "", err
	}
	return name, nil
}

func (s *storage) remove(name string) error {
	switch s.fileType { //nolint:exhaustive
	case SSLType, MapsType, GeneralType:
		return remove(name)
	}

	return nil
}

func readFile(name string) (string, error) {
	b, err := os.ReadFile(name)
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
		return "", conf.NewConfError(conf.ErrObjectDoesNotExist, fmt.Sprintf("file %s doesn't exist in dir: %s", name, dirname))
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
	entries, err := os.ReadDir(dirname)
	if err != nil {
		return nil, err
	}
	infos := make([]fs.FileInfo, 0, len(entries))
	for _, entry := range entries {
		info, err := entry.Info()
		if err != nil {
			return nil, err
		}
		infos = append(infos, info)
	}
	return infos, nil
}

func (s storage) validatePEM(raw []byte) error {
	crtPool := x509.NewCertPool()
	ok := crtPool.AppendCertsFromPEM(raw)
	if !ok {
		return fmt.Errorf("failed to parse certificate")
	}
	// HAProxy requires private and public key in same pem file
	hasCertificate := false
	hasPrivateKey := false
	for {
		block, rest := pem.Decode(raw)
		if block == nil {
			break
		}
		if block.Type == "CERTIFICATE" {
			hasCertificate = true
		} else { // check all other block types for the key, ignoring non-key blocks
			_, err := parsePrivateKey(block.Bytes)
			if err == nil {
				hasPrivateKey = true
			}
		}
		raw = rest
	}
	if !(hasCertificate && hasPrivateKey) {
		return fmt.Errorf("file should contain both certificate and private key")
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
