package main

import (
	"bufio"
	"bytes"
	"flag"
	"fmt"
	"html/template"
	"os"
	"path"
	"path/filepath"
	"sort"
	"strings"

	"github.com/go-openapi/swag/mangling"
	"github.com/haproxytech/client-native/v6/configuration/parents"
	"gopkg.in/yaml.v3"
)

var cache map[string]interface{} //nolint:gochecknoglobals

func errorExit(msg string) {
	fmt.Fprintf(os.Stderr, "ERROR: %v\n", msg)
	os.Exit(1)
}

func warn(msg string) {
	fmt.Fprintf(os.Stdout, "WARNING: %v\n", msg)
}

func expandRef(refValue string, absPath string, prefix string) string {
	words := strings.SplitN(refValue, "#", 2)

	if len(words) != 2 {
		warn("Invalid ref: " + refValue)
		return refValue
	}

	filePath := path.Join(absPath, words[0])
	keyPath := words[1]

	m, ok := cache[filePath].(map[string]interface{})
	if !ok {
		fileHandle, err := os.Open(filePath)
		if err != nil {
			errorExit(err.Error())
		}
		defer fileHandle.Close()

		fileScanner := bufio.NewScanner(fileHandle)
		value := ""
		for fileScanner.Scan() {
			value += fileScanner.Text() + "\n"
		}

		err = yaml.Unmarshal([]byte(value), &m)
		if err != nil {
			fmt.Println(refValue)         //nolint:forbidigo
			fmt.Println("WARNING: ", err) //nolint:forbidigo
			return refValue
		}
		cache[filePath] = m
	}
	retVal := make(map[string]interface{})
	if m[keyPath[1:]] != nil {
		retVal = m[keyPath[1:]].(map[string]interface{})
	} else {
		fmt.Println(refValue)       //nolint:forbidigo
		fmt.Println(keyPath)        //nolint:forbidigo
		fmt.Println(m[keyPath[1:]]) //nolint:forbidigo
	}

	buf := bytes.Buffer{}
	enc := yaml.NewEncoder(&buf)
	enc.SetIndent(2)
	// Can set default indent here on the encoder
	err := enc.Encode(&retVal)
	if err != nil {
		warn("Error encoding YAML")
		return refValue
	}
	retValStr := buf.String()

	var indentedRetValStr string
	var indentedLine string
	for _, line := range strings.Split(retValStr, "\n") {
		if strings.TrimSpace(line) != "" {
			indentedLine = prefix + "" + line + "\n"
			indentedRetValStr += indentedLine
		}
	}

	return indentedRetValStr[:len(indentedRetValStr)-1]
}

func main() { //nolint:gocognit
	inputFilePtr := flag.String("file", "", "Source file")

	flag.Parse()

	if *inputFilePtr == "" {
		errorExit("Input file not specified, please specify")
	}
	// sanity checks
	if _, err := os.Stat(strings.TrimSpace(*inputFilePtr)); os.IsNotExist(err) {
		errorExit("File " + *inputFilePtr + " does not exist")
	}

	cache = make(map[string]interface{})

	absPath := filepath.Dir(*inputFilePtr)
	fileHandle, err := os.Open(*inputFilePtr)
	if err != nil {
		errorExit(err.Error())
	}
	defer fileHandle.Close()

	fileScanner := bufio.NewScanner(fileHandle)
	var result strings.Builder

	type tag struct {
		Name        string `yaml:"name,omitempty"`
		Description string `yaml:"description,omitempty"`
	}
	type tags []tag
	ts := tags{}
	var tagResult strings.Builder
	for fileScanner.Scan() {
		line := fileScanner.Text()
		switch {
		case strings.HasPrefix(strings.TrimSpace(line), "$ref:"):
			refValue := strings.TrimSpace(strings.TrimSpace(line)[5:])
			refValue = refValue[1 : len(refValue)-1]
			if strings.HasPrefix(refValue, "#") {
				result.WriteString(line)
				result.WriteString("\n")
			} else {
				prefix := ""
				for _, char := range line {
					if string(char) == " " {
						prefix += " "
					} else {
						break
					}
				}
				result.WriteString(expandRef(refValue, absPath, prefix))
				result.WriteString("\n")
			}
		case strings.HasPrefix(strings.TrimSpace(line), "tags:"):
			for fileScanner.Scan() {
				tagLine := fileScanner.Text()
				if !strings.HasPrefix(strings.TrimSpace(tagLine), "security:") {
					tagResult.WriteString(tagLine)
					tagResult.WriteString("\n")
				} else {
					str := tagResult.String()
					err = yaml.Unmarshal([]byte(str), &ts)
					if err != nil {
						errorExit(err.Error())
					}
					sort.Slice(ts, func(i, j int) bool {
						return ts[i].Name < ts[j].Name
					})
					result.WriteString("tags:")
					result.WriteString("\n")

					b, _ := yaml.Marshal(&ts)

					for _, line := range strings.Split(strings.TrimRight(string(b), "\n"), "\n") {
						result.WriteString("  " + line + "\n")
					}
					result.WriteString("security:")
					result.WriteString("\n")
					break
				}
			}
		default:
			result.WriteString(line)
			result.WriteString("\n")
		}
	}

	tmplRes := expandChildren(result.String())
	res := tmplRes.String()
	res = strings.TrimSuffix(res, "\n")

	fmt.Println(res) //nolint:forbidigo
}

// to expand models for nested children
// update:
// - configuration/parents/constants.go
// - congiuration/parents/parents.go
// and specification/haproxy_spec.yaml (template)

func expandChildren(src string) bytes.Buffer {
	funcMap := template.FuncMap{
		"parents":  parents.Parents,
		"toGoName": mangling.NameMangler.ToGoName,
	}

	tmpl := template.Must(template.New("").Funcs(funcMap).Parse(src))
	var result bytes.Buffer
	err := tmpl.Execute(&result, nil)
	if err != nil {
		errorExit(err.Error())
	}
	return result
}
