package configs

import (
	"fmt"
	"log"
	"testing"

	parser "github.com/haproxytech/client-native/v5/config-parser"
	"github.com/haproxytech/client-native/v5/config-parser/spoe"
	"github.com/haproxytech/client-native/v5/config-parser/types"
)

func TestVersion(t *testing.T) {
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
			version, err := p.Get("", parser.Comments, parser.CommentsSectionName, "# _version", true)
			if err != nil {
				t.Fatal(err)
			}
			log.Println(version)
			v, ok := version.(*types.ConfigVersion)
			if !ok {
				t.Fatal("version fetch failed")
			}
			if v.Value != 3 {
				t.Fatalf("version mismatch: has %d want %d", v.Value, 3)
			}
		})
	}
}

func TestEvent2(t *testing.T) {
	p := spoe.Parser{}
	err := p.LoadData("spoe.cfg")
	if err != nil {
		fmt.Println(err)
		return
	}
	aNames, err := p.SectionsGet("[ip-reputation]", parser.SPOEAgent)
	if err != nil {
		fmt.Println(err)
		return
	}
	for _, name := range aNames {
		fmt.Println(name)
	}
	log.Println(p.Get("", parser.Comments, parser.CommentsSectionName, "# _version", true))
	log.Println(p.String())

	s, err := p.SectionsGet("[ip-reputation]", "spoe-agent")
	log.Println(s, err)
	/*
		params := spoe.Params{
			SpoeDir:        "/tmp/spoe",
			TransactionDir: "/tmp/haproxy-spoe",
		}
		spoeConfig, _ := spoe.NewSpoe(params)
		client, err := spoeConfig.GetSingleSpoe("spoe.cfg")
		if err != nil {
			log.Fatal(err)
		}
		_, agents, err := client.GetAgents("[ip-reputation]", "")
		if err != nil {
			fmt.Println(err)
		}
		for _, a := range agents {
			fmt.Println(a)
		}
		v, agent, err := client.GetAgent("[ip-reputation]", "iprep-agent", "")
		if err != nil {
			fmt.Println(err)
		}*/
}
