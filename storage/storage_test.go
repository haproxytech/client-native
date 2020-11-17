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
	"io/ioutil"
	"path/filepath"
	"reflect"
	"testing"
)

func createDir(config string, createFile bool) (dirname, file string, err error) {
	dirname, err = ioutil.TempDir("/tmp", "storage")
	if err != nil {
		return "", "", err
	}
	if !createFile {
		return dirname, "", nil
	}
	f, err := ioutil.TempFile(dirname, "")
	if err != nil {
		return "", "", err
	}
	if config != "" {
		_, err = f.WriteString(config)
		if err != nil {
			return "", "", err
		}
	}
	return dirname, filepath.Base(f.Name()), nil
}

func readFile(name string) (string, error) {
	b, err := ioutil.ReadFile(name)
	if err != nil {
		return "", err
	}
	return string(b), nil
}

func TestNew(t *testing.T) {
	s, _ := New("/tmp", MapsType)
	sNoDir, _ := New("", MapsType)
	tests := []struct {
		name     string
		dirname  string
		fileType StorageFileType
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

func TestStorage_GetAll(t *testing.T) {
	conf1 := `key1 val1
key2 val2`
	dirWithFile, file, err := createDir(conf1, true)
	if err != nil {
		t.Error(err.Error())
	}
	defer func() {
		_ = remove(file)
		_ = remove(dirWithFile)
	}()

	dirWithoutFile, _, err := createDir("", false)
	if err != nil {
		t.Error(err.Error())
	}
	defer func() {
		_ = remove(dirWithoutFile)
	}()

	tests := []struct {
		name    string
		dirname string
		want    []string
		wantErr bool
	}{
		{
			name:    "Should return created file names",
			dirname: dirWithFile,
			want:    []string{filepath.Join(dirWithFile, file)},
			wantErr: false,
		},
		{
			name:    "Should return an error if no files in directory",
			dirname: dirWithoutFile,
			want:    nil,
			wantErr: true,
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
				dirname: tt.dirname,
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
	dirWithFile, file, err := createDir(conf1, true)
	if err != nil {
		t.Error(err.Error())
	}
	defer func() {
		_ = remove(file)
		_ = remove(dirWithFile)
	}()

	dirWithoutFile, noFile, err := createDir("", false)
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
			name:     "Should return error if file doesn't exists",
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
			got, err := s.Get(tt.fileName)
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
	dirWithFile, file, err := createDir("", true)
	if err != nil {
		t.Error(err.Error())
	}
	defer func() {
		_ = remove(file)
		_ = remove(dirWithFile)
	}()

	dirWithoutFile, noFile, err := createDir("", false)
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
	dirWithFile, file, err := createDir(conf1, true)
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
