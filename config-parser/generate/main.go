/*
Copyright 2019 HAProxy Technologies

Licensed under the Apache License, Version 2.0 (the "License");
you may not use this file except in compliance with the License.
You may obtain a copy of the License at

    http://www.apache.org/licenses/LICENSE-2.0

Unless required by applicable law or agreed to in writing, software
distributed under the License is distributed on an "AS IS" BASIS,
WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
See the License for the specific language governing permissions and
limitations under the License.
*/

package main

import (
	"log"
	"os"
	"path"
	"path/filepath"
	"strings"

	"github.com/google/renameio/maybe"
)

//nolint:gochecknoglobals
var configFile = ConfigFile{}

func main() {
	dir, err := filepath.Abs(filepath.Dir(os.Args[0]))
	if err != nil {
		log.Fatal(err)
	}
	if strings.HasSuffix(dir, "generate") {
		dir, err = filepath.Abs(filepath.Dir(os.Args[0]) + "/..")
		if err != nil {
			log.Fatal(err)
		}
	}
	if len(os.Args) > 1 {
		dir = os.Args[1]
	}
	log.Println(dir)

	configFile.Section = map[string][]string{}

	generateTypes(dir, "")
	generateTypesGeneric(dir)
	generateTypesOther(dir)
	// spoe
	generateTypes(dir, "spoe/")

	filePath := path.Join(dir, "tests", "configs", "haproxy_generated.cfg.go")
	res, err := GoFmt([]byte(configFile.String()))
	CheckErr(err)
	err = maybe.WriteFile(filePath, res, 0o644)
	CheckErr(err)

	configFile.StringFiles(path.Join(dir, "tests", "integration"))
}
