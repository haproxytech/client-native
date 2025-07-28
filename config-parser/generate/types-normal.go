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
	"fmt"
	"log"
	"os"
	"path"
	"strings"

	"github.com/haproxytech/client-native/v6/config-parser/common"
)

func generateTypes(dir string, dataDir string) { //nolint:gocognit
	dat, err := os.ReadFile(dataDir + "types/types.go")
	if err != nil {
		log.Println(err)
	}
	lines := common.StringSplitIgnoreEmpty(string(dat), '\n')

	parserData := Data{}

	for _, line := range lines {
		parserData.DataDir = dataDir
		if strings.HasPrefix(line, "// model:") {
			parserData.Model = strings.TrimSpace(strings.TrimPrefix(line, "// model:"))
		}
		if strings.HasPrefix(line, "//doc:") {
			parserData.Doc = strings.TrimSpace(strings.TrimPrefix(line, "//doc:"))
		}
		if strings.HasPrefix(line, "//deprecated:") {
			parserData.Deprecated = true
		}
		if strings.HasPrefix(line, "//sections:") {
			s := strings.Split(line, ":")
			parserData.ParserSections = strings.Split(s[1], ",")
		}
		if strings.HasPrefix(line, "//no:sections") {
			parserData.NoSections = true
		}
		if strings.HasPrefix(line, "//name:") {
			data := common.StringSplitIgnoreEmpty(line, ':')
			items := common.StringSplitIgnoreEmpty(data[1], ' ')
			parserData.ParserName = data[1]
			if len(items) > 1 {
				parserData.ParserName = items[0]
				parserData.ParserSecondName = items[1]
			}
		}
		if strings.HasPrefix(line, "// name:") {
			data := common.StringSplitIgnoreEmpty(line, ':')
			items := common.StringSplitIgnoreEmpty(data[1], ' ')
			parserData.ParserName = data[1]
			if len(items) > 1 {
				parserData.ParserName = items[0]
				parserData.ParserSecondName = items[1]
			}
		}
		if strings.HasPrefix(line, "//is:multiple") {
			parserData.ParserMultiple = true
		}
		if strings.HasPrefix(line, "//no:init") {
			parserData.NoInit = true
		}
		if strings.HasPrefix(line, "//no:name") {
			parserData.NoName = true
		}
		if strings.HasPrefix(line, "//no:parse") {
			parserData.NoParse = true
		}
		if strings.HasPrefix(line, `//test:quote_ok`) && !parserData.Deprecated {
			data := strings.SplitN(line, ":", 3)
			parserData.TestOKEscaped = append(parserData.TestOKEscaped, data[2])
		}
		if strings.HasPrefix(line, `//test:"defaults-ok"`) && !parserData.Deprecated {
			data := strings.SplitN(line, ":", 3)
			parserData.TestOKDefaults = append(parserData.TestOKDefaults, data[2])
		}
		if strings.HasPrefix(line, `//test:"frontend-ok"`) && !parserData.Deprecated {
			data := strings.SplitN(line, ":", 3)
			parserData.TestOKFrontend = append(parserData.TestOKFrontend, data[2])
		}
		if strings.HasPrefix(line, `//test:"backend-ok"`) && !parserData.Deprecated {
			data := strings.SplitN(line, ":", 3)
			parserData.TestOKBackend = append(parserData.TestOKBackend, data[2])
		}
		if strings.HasPrefix(line, "//test:ok") && !parserData.Deprecated {
			data := strings.SplitN(line, ":", 3)
			parserData.TestOK = append(parserData.TestOK, data[2])
		}
		if strings.HasPrefix(line, `//test:quote_fail`) && !parserData.Deprecated {
			data := strings.SplitN(line, ":", 3)
			parserData.TestFailEscaped = append(parserData.TestFailEscaped, data[2])
		}
		if strings.HasPrefix(line, "//test:fail") && !parserData.Deprecated {
			data := strings.SplitN(line, ":", 3)
			parserData.TestFail = append(parserData.TestFail, data[2])
		}
		if strings.HasPrefix(line, "//test:expected-ok") && !parserData.Deprecated {
			data := strings.SplitN(line, ":", 4)
			tableTestData := TableTestData{
				Test:  data[2],
				Table: data[3],
			}
			parserData.TestTableOK = append(parserData.TestTableOK, tableTestData)
			parserData.HasTable = true
		}
		if strings.HasPrefix(line, "//test:alias") {
			data := strings.SplitN(line, ":", 5)
			aliasTestData := AliasTestData{
				Alias: data[2],
				Test:  data[4],
			}
			switch data[3] {
			case "ok":
				parserData.TestAliasOK = append(parserData.TestAliasOK, aliasTestData)
			case "fail":
				parserData.TestAliasFail = append(parserData.TestAliasFail, aliasTestData)
			default:
				log.Fatalf("not able to process line %s", line)
			}
		}
		if strings.HasPrefix(line, "//has-alias:true") {
			parserData.HasAlias = true
		}

		if !strings.HasPrefix(line, "type ") {
			continue
		}

		if parserData.ParserName == "" {
			parserData = Data{}
			continue
		}
		data := common.StringSplitIgnoreEmpty(line, ' ')
		parserData.StructName = data[1]
		parserData.ParserType = data[1]

		filename := parserData.ParserName
		if parserData.ParserSecondName != "" {
			filename = fmt.Sprintf("%s %s", filename, parserData.ParserSecondName)
		}

		filePath := path.Join(dir, dataDir, "parsers", cleanFileName(filename)+"_generated.go")
		executeTemplate(TemplateTypeNormal, &parserData, filePath)

		// parserData.TestFail = append(parserData.TestFail, "") parsers should not get empty line!
		parserData.TestFail = append(parserData.TestFail, "---")
		parserData.TestFail = append(parserData.TestFail, "--- ---")

		filePath = path.Join(dir, dataDir, "tests", cleanFileName(filename)+"_generated_test.go")
		executeTemplate(TemplateTypeTest, &parserData, filePath)

		configFile.AddParserData(parserData)
		parserData = Data{}
	}
}
