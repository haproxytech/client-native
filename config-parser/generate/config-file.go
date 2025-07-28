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
	"path"
	"sort"
	"strings"

	"github.com/google/renameio/maybe"
	"golang.org/x/text/cases"
	"golang.org/x/text/language"
)

type AliasTestData struct {
	Alias string
	Test  string
}

type DefaultTestData struct {
	Test    string
	Default string
}

type TableTestData struct {
	Test  string
	Table string
}

type Data struct { //nolint:maligned
	ParserMultiple     bool
	ParserSections     []string
	ParserName         string
	ParserSecondName   string
	StructName         string
	ParserType         string
	ParserTypeOverride string
	NoInit             bool
	NoName             bool
	NoParse            bool
	NoGet              bool
	NoSections         bool
	IsInterface        bool
	Dir                string
	ModeOther          bool
	TestOK             []string
	TestOKDefaults     []string
	TestOKFrontend     []string
	TestOKBackend      []string
	TestOKEscaped      []string
	TestFail           []string
	TestFailEscaped    []string
	TestAliasOK        []AliasTestData
	TestAliasFail      []AliasTestData
	TestDefault        []DefaultTestData
	TestTableOK        []TableTestData
	TestSkip           bool
	DataDir            string
	Deprecated         bool
	HasAlias           bool
	HasDefault         bool
	HasTable           bool
	Model              string // model name for swagger
	Doc                string // url to docs
}

type ConfigFile struct {
	Section    map[string][]string
	SectionAll map[string][]string
	Tests      strings.Builder
}

func (c *ConfigFile) AddParserData(parser Data) { //nolint:gocognit
	sections := parser.ParserSections
	testOK := parser.TestOK
	testOKDefaults := parser.TestOKDefaults
	testOKFrontend := parser.TestOKFrontend
	testOKBackend := parser.TestOKBackend
	TestOKEscaped := parser.TestOKEscaped
	if len(sections) == 0 && !parser.NoSections {
		log.Fatalf("parser %s does not have any section defined", parser.ParserName)
	}
	var lines []string
	var linesDefaults []string
	var linesFrontend []string
	var linesBackend []string
	var lines2 []string
	for _, section := range sections {
		_, ok := c.Section[section]
		if !ok {
			c.Section[section] = []string{}
			if c.SectionAll == nil {
				c.SectionAll = make(map[string][]string)
			}
			c.SectionAll[section] = []string{}
		}
		// line = testOK[0]
		if parser.ParserMultiple {
			lines = testOK
			//nolint:gosimple
			for _, line := range testOK {
				c.Section[section] = append(c.Section[section], line)
			}
			// c.Section[s] = append(c.Section[s], testOK...)
		} else {
			lines = []string{testOK[0]}
			c.Section[section] = append(c.Section[section], testOK[0])
		}
		if len(parser.TestTableOK) > 0 {
			for _, line := range parser.TestTableOK {
				lines = append(lines, line.Table)
				c.Section[section] = append(c.Section[section], line.Table)
			}
		}
		c.SectionAll[section] = append(c.SectionAll[section], testOK...)
		if section == "defaults" {
			if parser.ParserMultiple {
				linesDefaults = testOKDefaults
				//nolint:gosimple
				for _, line := range testOKDefaults {
					c.Section[section] = append(c.Section[section], line)
				}
			} else if len(testOKDefaults) > 0 {
				linesDefaults = []string{testOKDefaults[0]}
				c.Section[section] = append(c.Section[section], testOKDefaults[0])
			}
			c.SectionAll[section] = append(c.SectionAll[section], testOKDefaults...)
		}
		if section == "frontend" {
			if parser.ParserMultiple {
				linesFrontend = testOKFrontend
				//nolint:gosimple
				for _, line := range testOKFrontend {
					c.Section[section] = append(c.Section[section], line)
				}
			} else if len(testOKFrontend) > 0 {
				linesFrontend = []string{testOKFrontend[0]}
				c.Section[section] = append(c.Section[section], testOKFrontend[0])
			}
			c.SectionAll[section] = append(c.SectionAll[section], testOKFrontend...)
		}
		if section == "backend" {
			if parser.ParserMultiple {
				linesBackend = testOKBackend
				//nolint:gosimple
				for _, line := range testOKBackend {
					c.Section[section] = append(c.Section[section], line)
				}
			} else if len(testOKBackend) > 0 {
				linesBackend = []string{testOKBackend[0]}
				c.Section[section] = append(c.Section[section], testOKBackend[0])
			}
			c.SectionAll[section] = append(c.SectionAll[section], testOKBackend...)
		}
		if parser.ParserMultiple {
			lines2 = TestOKEscaped
			//nolint:gosimple
			for _, line := range TestOKEscaped {
				c.Section[section] = append(c.Section[section], line)
			}
			// c.Section[s] = append(c.Section[s], TestOKEscaped...)
		} else if len(TestOKEscaped) > 0 {
			lines2 = []string{TestOKEscaped[0]}
			c.Section[section] = append(c.Section[section], TestOKEscaped[0])
		}
		c.SectionAll[section] = append(c.SectionAll[section], TestOKEscaped...)
	}
	if len(lines) == 0 && len(lines2) == 0 {
		if parser.NoSections {
			return
		} else if !parser.Deprecated {
			log.Fatalf("parser %s does not have any tests defined", parser.ParserName)
		}
	}
	if !parser.NoSections {
		for _, line := range lines {
			c.Tests.WriteString(fmt.Sprintf("  {`  %s\n`, %d},\n", line, len(sections)))
		}
		for _, line := range linesDefaults {
			c.Tests.WriteString(fmt.Sprintf("  {`  %s\n`, 1},\n", line))
		}
		for _, line := range linesFrontend {
			c.Tests.WriteString(fmt.Sprintf("  {`  %s\n`, 1},\n", line))
		}
		for _, line := range linesBackend {
			c.Tests.WriteString(fmt.Sprintf("  {`  %s\n`, 1},\n", line))
		}
	}
}

func (c *ConfigFile) String() string {
	var result strings.Builder

	result.WriteString(license)
	result.WriteString("package configs\n\n")
	result.WriteString("const generatedConfig = `# _version=1\n# HAProxy Technologies\n# https://www.haproxy.com/\n# sections are in alphabetical order (except global & default) for code generation\n\n")

	first := true
	sectionNames := make([]string, len(c.Section)-1)
	index := 0
	for sectionName := range c.Section {
		if sectionName == "global" {
			continue
		}
		sectionNames[index] = sectionName
		index++
	}
	sort.Strings(sectionNames)

	writeSection := func(sectionName string) {
		if !first {
			result.WriteString("\n")
		} else {
			first = false
		}
		result.WriteString(sectionName)
		result.WriteString(" test\n")
		lines := c.Section[sectionName]
		for _, line := range lines {
			result.WriteString("  ")
			result.WriteString(line)
			result.WriteString("\n")
		}
	}

	writeSection("global")
	for _, sectionName := range sectionNames {
		writeSection(sectionName)
	}
	result.WriteString("`\n\n")

	result.WriteString("var configTests = []configTest{")
	result.WriteString(c.Tests.String())
	result.WriteString("}")

	result.WriteString("\n")
	return result.String()
}

func (c *ConfigFile) StringFiles(baseFolder string) {
	files := map[string][]byte{}

	header := license + "\npackage integration_test\n\n"

	sectionTypes := make([]string, len(c.SectionAll)-1)
	index := 0
	for sectionName := range c.SectionAll {
		if sectionName == "global" {
			continue
		}
		sectionTypes[index] = sectionName
		index++
	}
	sort.Strings(sectionTypes)
	usedNiceNames := map[string]struct{}{}

	writeSection := func(sectionType string) {
		lines := c.SectionAll[sectionType]
		file := files[sectionType]
		if len(file) < 1 {
			file = append(file, []byte(header)...)
		}
		for _, line := range lines {
			niceName := getNiceName(sectionType) + "_" + getNiceName(line)
			exists := true
			for exists {
				_, exists = usedNiceNames[niceName]
				if exists {
					niceName += "_"
				}
			}
			usedNiceNames[niceName] = struct{}{}
			sectionName := " test"
			if sectionType == "global" || sectionType == "traces" {
				sectionName = ""
			}
			oneTest := "const " + niceName + " = `\n" + sectionType + sectionName + "\n" + "  " + line + "\n`" + "\n"
			file = append(file, []byte(oneTest)...)
			files[sectionType] = file
		}
	}

	sectionTypes = append(sectionTypes, "global") //nolint:makezero
	for _, sectionName := range sectionTypes {
		writeSection(sectionName)
		for name, data := range files {
			filePath := path.Join(baseFolder, name+"_data_test.go")
			log.Println(filePath)
			saveFile(filePath, string(data))
		}
		files = map[string][]byte{}

		var testFile strings.Builder
		testName := cases.Title(language.Und, cases.NoLower).String(getNiceName(sectionName))
		testFile.WriteString(strings.Replace(testHeader, "TestWholeConfigsSections", "TestWholeConfigsSections"+testName, 1))
		sortedNames := []string{}
		for name := range usedNiceNames {
			if strings.HasPrefix(name, sectionName) {
				sortedNames = append(sortedNames, name)
			}
		}
		sort.Strings(sortedNames)
		for _, name := range sortedNames {
			testFile.WriteString(fmt.Sprintf("		{\"%s\", %s},\n", name, name))
		}
		testFile.WriteString(testFooter)
		saveFile(path.Join(baseFolder, sectionName+"_test.go"), testFile.String())
	}
}

func saveFile(pathFile, data string) {
	err := maybe.WriteFile(pathFile, []byte(data), 0o644)
	CheckErr(err)
}

//nolint:gochecknoglobals
var testHeader = license + `
package integration_test

import (
	"bytes"
	"testing"

	parser "github.com/haproxytech/client-native/v6/config-parser"
	"github.com/haproxytech/client-native/v6/config-parser/options"
)

func TestWholeConfigsSections(t *testing.T) {
	t.Parallel()
	tests := []struct {
		Name, Config string
	}{
`

//nolint:gochecknoglobals
var testFooter = `	}
	for _, config := range tests {
		t.Run(config.Name, func(t *testing.T) {
			t.Parallel()
			var buffer bytes.Buffer
			buffer.WriteString(config.Config)
			p, err := parser.New(options.Reader(&buffer))
			if err != nil {
				t.Fatal(err.Error())
			}
			result := p.String()
			if result != config.Config {
				compare(t, config.Config, result)
				t.Error("======== ORIGINAL =========")
				t.Error(config.Config)
				t.Error("======== RESULT ===========")
				t.Error(result)
				t.Error("===========================")
				t.Fatal("configurations does not match")
			}
		})
	}
}
`
