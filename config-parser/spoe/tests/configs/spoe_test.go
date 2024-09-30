package configs

import (
	"strings"
	"testing"

	"github.com/haproxytech/client-native/v5/config-parser/spoe"
)

func TestWholeConfigs(t *testing.T) {
	tests := []struct {
		Name, Config string
	}{
		{"configBasic1", configBasic1},
	}
	for _, config := range tests {
		t.Run(config.Name, func(t *testing.T) {
			p := spoe.Parser{}
			//var buffer bytes.Buffer
			//buffer.WriteString(config.Config)
			//_ = p.Process(&buffer)
			err := p.ParseData(configBasic1)
			if err != nil {
				t.Fatal(err)
			}
			result := p.String()
			if result != config.Config {
				compare(t, config.Config, result)
				t.Fatalf("configurations does not match")
			}
		})
	}
}

func compare(t *testing.T, configOriginal, configResult string) {
	original := strings.Split(configOriginal, "\n")
	result := strings.Split(configResult, "\n")
	if len(original) != len(result) {
		t.Logf("not the same size: original: %d, result: %d", len(original), len(result))
		return
	}
	for index, line := range original {
		if line != result[index] {
			t.Logf("line %d: '%s' != '%s'", index+3, line, result[index])
		}
	}
}
