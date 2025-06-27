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
	"bytes"
	"io"
	"io/ioutil"
	"path/filepath"
	"reflect"
	"testing"

	"github.com/haproxytech/client-native/v6/misc"
	"github.com/stretchr/testify/require"
)

func TestNew(t *testing.T) {
	s, _ := New("/tmp", MapsType)
	sNoDir, _ := New("", MapsType)
	tests := []struct {
		name     string
		dirname  string
		fileType FileType
		want     Storage
		wantErr  bool
	}{
		{
			name:     "Should return object when dir specified",
			dirname:  "/tmp",
			fileType: MapsType,
			want:     s,
			wantErr:  false,
		},
		{
			name:     "Should return an error when dirname not specified",
			dirname:  "",
			fileType: MapsType,
			want:     sNoDir,
			wantErr:  true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := New(tt.dirname, tt.fileType)
			if (err != nil) != tt.wantErr {
				t.Errorf("New() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("New() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestNewBackupStorage(t *testing.T) {
	s, _ := New("/tmp", BackupsType)
	sNoDir, _ := New("", BackupsType)
	tests := []struct {
		name     string
		dirname  string
		fileType FileType
		want     Storage
		wantErr  bool
	}{
		{
			name:     "Should return object when dir specified",
			dirname:  "/tmp",
			fileType: BackupsType,
			want:     s,
			wantErr:  false,
		},
		{
			name:     "Should return an error when dirname not specified",
			dirname:  "",
			fileType: BackupsType,
			want:     sNoDir,
			wantErr:  true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := New(tt.dirname, tt.fileType)
			if (err != nil) != tt.wantErr {
				t.Errorf("New() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("New() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestStorage_GetAll(t *testing.T) {
	conf1 := `key1 val1
key2 val2`
	dirWithFile, file, err := misc.CreateTempDir(conf1, true)
	if err != nil {
		t.Error(err.Error())
	}
	defer func() {
		_ = remove(file)
		_ = remove(dirWithFile)
	}()

	dirWithoutFile, _, err := misc.CreateTempDir("", false)
	if err != nil {
		t.Error(err.Error())
	}
	defer func() {
		_ = remove(dirWithoutFile)
	}()

	tests := []struct {
		name     string
		dirname  string
		fileType FileType
		want     []string
		wantErr  bool
	}{
		{
			name:     "Should return created file names",
			dirname:  dirWithFile,
			fileType: MapsType,
			want:     []string{filepath.Join(dirWithFile, file)},
			wantErr:  false,
		},
		{
			name:    "Should return an error if no files in directory",
			dirname: dirWithoutFile,
			want:    []string{},
			wantErr: false,
		},
		{
			name:    "Should return an error when directory not specified",
			dirname: "",
			want:    nil,
			wantErr: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			s := storage{
				dirname:  tt.dirname,
				fileType: tt.fileType,
			}
			got, err := s.GetAll()
			if (err != nil) != tt.wantErr {
				t.Errorf("Storage.GetAll() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("Storage.GetAll() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestStorage_Get(t *testing.T) {
	conf1 := `key1 val1
key2 val2`
	dirWithFile, file, err := misc.CreateTempDir(conf1, true)
	if err != nil {
		t.Error(err.Error())
	}
	defer func() {
		_ = remove(file)
		_ = remove(dirWithFile)
	}()

	dirWithoutFile, noFile, err := misc.CreateTempDir("", false)
	if err != nil {
		t.Error(err.Error())
	}
	defer func() {
		_ = remove(dirWithoutFile)
	}()

	tests := []struct {
		name     string
		dirname  string
		fileName string
		want     string
		wantErr  bool
	}{
		{
			name:     "Should return file name",
			dirname:  dirWithFile,
			fileName: file,
			want:     filepath.Join(dirWithFile, file),
			wantErr:  false,
		},
		{
			name:     "Should return error if file doesn't exist",
			dirname:  dirWithoutFile,
			fileName: noFile,
			want:     noFile,
			wantErr:  true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			s := storage{
				dirname: tt.dirname,
			}
			got, _, err := s.Get(tt.fileName)
			if (err != nil) != tt.wantErr {
				t.Errorf("Storage.Get() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if got != tt.want {
				t.Errorf("Storage.Get() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestStorage_Delete(t *testing.T) {
	dirWithFile, file, err := misc.CreateTempDir("", true)
	if err != nil {
		t.Error(err.Error())
	}
	defer func() {
		_ = remove(file)
		_ = remove(dirWithFile)
	}()

	dirWithoutFile, noFile, err := misc.CreateTempDir("", false)
	if err != nil {
		t.Error(err.Error())
	}
	defer func() {
		_ = remove(dirWithoutFile)
	}()

	tests := []struct {
		name    string
		dirname string
		file    string
		wantErr bool
	}{
		{
			name:    "",
			file:    file,
			wantErr: false,
			dirname: dirWithFile,
		},
		{
			name:    "",
			file:    noFile,
			wantErr: true,
			dirname: dirWithoutFile,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			s := storage{
				dirname: tt.dirname,
			}
			if err := s.Delete(tt.file); (err != nil) != tt.wantErr {
				t.Errorf("Storage.Delete() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

func TestStorage_Replace(t *testing.T) {
	conf1 := `key1 val1
	key2 val2`
	dirWithFile, file, err := misc.CreateTempDir(conf1, true)
	if err != nil {
		t.Error(err.Error())
	}
	defer func() {
		_ = remove(file)
		_ = remove(dirWithFile)
	}()
	tests := []struct {
		name       string
		dirname    string
		file       string
		newcontent string
		want       string
		wantErr    bool
	}{
		{
			name:       "Should update file with new content",
			newcontent: "newcontent",
			dirname:    dirWithFile,
			file:       file,
			want:       filepath.Join(dirWithFile, file),
			wantErr:    false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			s := storage{
				dirname: tt.dirname,
			}
			got, err := s.Replace(tt.file, tt.newcontent)
			if (err != nil) != tt.wantErr {
				t.Errorf("Storage.Replace() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if got != tt.want {
				t.Errorf("Storage.Replace() = %v, want %v", got, tt.want)
			}
			newcontent, err := readFile(tt.want)
			if err != nil {
				t.Errorf("Storage.Replace() error reading file = %v, want %v", err, tt.want)
				return
			}
			if newcontent != tt.newcontent {
				t.Errorf("Storage.Replace() = newcontent %v, want %v", newcontent, tt.newcontent)
				return
			}
		})
	}
}

func Test_storage_Create(t *testing.T) {
	conf1 := `key1 val1
	key2 val2`
	dirWithFile, file, err := misc.CreateTempDir(conf1, true, "*.map")
	if err != nil {
		t.Error(err.Error())
	}
	defer func() {
		_ = remove(file)
		_ = remove(dirWithFile)
	}()

	dirWithoutFile, _, err := misc.CreateTempDir("", false, "*.map")
	if err != nil {
		t.Error(err.Error())
	}
	defer func() {
		_ = remove(dirWithoutFile)
	}()

	tests := []struct {
		name       string
		dirname    string
		fileType   FileType
		file       string
		readCloser io.ReadCloser
		want       string
		wantErr    bool
	}{
		{
			name:       "Should return an error if file exists",
			dirname:    dirWithFile,
			fileType:   MapsType,
			readCloser: nil,
			file:       file,
			want:       "",
			wantErr:    true,
		},
		{
			name:       "Should create file if not exists",
			dirname:    dirWithoutFile,
			fileType:   MapsType,
			readCloser: ioutil.NopCloser(bytes.NewReader([]byte("hello world"))),
			file:       "newfile.map",
			want:       filepath.Join(dirWithoutFile, "newfile.map"),
			wantErr:    false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			s := storage{
				dirname:  tt.dirname,
				fileType: tt.fileType,
			}
			got, _, err := s.Create(tt.file, tt.readCloser)
			if (err != nil) != tt.wantErr {
				t.Errorf("storage.Create() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if got != tt.want {
				t.Errorf("storage.Create() = %v, want %v", got, tt.want)
			}
		})
	}
}

func readPem(path string) ([]byte, error) {
	p, err := filepath.Abs("../storage/test-certs/" + path)
	if err != nil {
		return nil, err
	}
	f, err := readFile(p)
	if err != nil {
		return nil, err
	}
	return []byte(f), nil
}

func Test_storage_validatePEM(t *testing.T) {
	invalidOnlyPublic, err := readPem("invalid/only-public.pem")
	require.NoError(t, err)

	invalidOnlyPrivate, err := readPem("invalid/only-private.pem")
	require.NoError(t, err)

	validChain, err := readPem("valid/OK-key_crt_int1_int2.pem")
	require.NoError(t, err)

	validSelfSignedWithoutRSA, err := readPem("invalid/self-signed-without-rsa.pem")
	require.NoError(t, err)

	tests := []struct {
		name    string
		dirname string
		content []byte
		wantErr bool
	}{
		{
			name:    "Should fail with only public pem",
			content: invalidOnlyPublic,
			dirname: "/tmp",
			wantErr: true,
		},
		{
			name:    "Should fail with only private pem",
			content: invalidOnlyPrivate,
			dirname: "/tmp",
			wantErr: true,
		},
		{
			name:    "Should pass with public, private and intermediate pems with correct order",
			content: validChain,
			dirname: "/tmp",
			wantErr: false,
		},
		{
			name:    "Should pass with self-signed cert without RSA",
			content: validSelfSignedWithoutRSA,
			dirname: "/tmp",
			wantErr: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			s := storage{
				dirname: tt.dirname,
			}

			if err := s.validatePEM(tt.content); (err != nil) != tt.wantErr {
				t.Errorf("storage.validatePEM() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
		})
	}
}

func Test_storage_getCertificatesInfo(t *testing.T) {
	tests := []struct {
		dirname  string
		filename string
	}{
		{
			filename: "only-public.pem",
			dirname:  "invalid",
		},
		{
			filename: "only-private.pem",
			dirname:  "invalid",
		},
		{
			filename: "NOK-crt_key_int1.pem",
			dirname:  "invalid",
		},
		{
			filename: "NOK-int1_int2.pem",
			dirname:  "invalid",
		},
		{
			filename: "self-signed-without-rsa.pem",
			dirname:  "invalid",
		},
		{
			filename: "OK-key_crt_int1_int2.pem",
			dirname:  "valid",
		},
		{
			filename: "OK-crt_key_int1_int2.pem",
			dirname:  "valid",
		},
		{
			filename: "OK-int1_key_crt_int2.pem",
			dirname:  "valid",
		},
		{
			filename: "selfsigned1.pem",
			dirname:  "valid",
		},
	}
	for _, tt := range tests {
		t.Run(tt.filename, func(t *testing.T) {
			s := storage{
				dirname: "../storage/test-certs/" + tt.dirname,
			}

			info, err := s.GetCertificatesInfo(tt.filename)

			if err != nil || info == nil {
				t.Errorf("storage.GetCertificatesInfo() error = %v", err)
				t.Logf("%+v", info)
				return
			}

			if tt.dirname == "valid" {
				require.NotEmpty(t, info.Sha1FingerPrint)
				require.NotEmpty(t, info.Sha256FingerPrint)
				require.NotEmpty(t, info.Subject)
				require.NotEmpty(t, info.Serial)
			}
		})
	}
}
