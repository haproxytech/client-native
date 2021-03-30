package main

import (
	"bufio"
	"flag"
	"fmt"
	"os"
	"path"
	"path/filepath"
	"sort"
	"strings"

	"gopkg.in/yaml.v2"
)

var cache map[string]interface{}

func error(msg string) {
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
			error(err.Error())
		}
		defer fileHandle.Close()

		fileScanner := bufio.NewScanner(fileHandle)
		value := ""
		for fileScanner.Scan() {
			value += fileScanner.Text() + "\n"
		}

		err = yaml.Unmarshal([]byte(value), &m)
		if err != nil {
			fmt.Println(refValue)
			fmt.Println("WARNING: ", err)
			return refValue
		}
		cache[filePath] = m
	}
	retVal := make(map[interface{}]interface{})
	if m[keyPath[1:]] != nil {
		retVal = m[keyPath[1:]].(map[interface{}]interface{})
	} else {
		fmt.Println(refValue)
		fmt.Println(keyPath)
		fmt.Println(m[keyPath[1:]])
	}

	retValByte, err := yaml.Marshal(retVal)
	if err != nil {
		warn("Error encoding YAML")
		return refValue
	}
	retValStr := string(retValByte)
	indentedRetValStr := ""
	indentedLine := ""
	for _, line := range strings.Split(retValStr, "\n") {
		if strings.TrimSpace(line) != "" {
			indentedLine = prefix + "  " + line + "\n"
			indentedRetValStr += indentedLine
		}
	}

	return indentedRetValStr[:len(indentedRetValStr)-1]
}

func main() {
	inputFilePtr := flag.String("file", "", "Source file")

	flag.Parse()

	if *inputFilePtr == "" {
		error("Input file not specified, please specify")
	}
	// sanity checks
	if _, err := os.Stat(strings.TrimSpace(*inputFilePtr)); os.IsNotExist(err) {
		error("File " + *inputFilePtr + " does not exist")
	}

	cache = make(map[string]interface{})

	absPath := filepath.Dir(*inputFilePtr)
	fileHandle, err := os.Open(*inputFilePtr)
	if err != nil {
		error(err.Error())
	}
	defer fileHandle.Close()

	fileScanner := bufio.NewScanner(fileHandle)
	var result strings.Builder

	type tag struct {
		Name        string `yaml:"name,omitempty"`
		Description string `yaml:"description,omitempty"`
	}
	type tags []tag
	var ts tags = tags{}
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
						error(err.Error())
					}
					sort.Slice(ts, func(i, j int) bool {
						return ts[i].Name < ts[j].Name
					})
					result.WriteString("tags:")
					result.WriteString("\n")

					b, _ := yaml.Marshal(&ts)
					result.WriteString(string(b))
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
	fmt.Println(result.String())
}
