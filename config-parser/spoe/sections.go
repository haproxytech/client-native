package spoe

import (
	parser "github.com/haproxytech/client-native/v6/config-parser"
	configparser "github.com/haproxytech/client-native/v6/config-parser/parsers"
	"github.com/haproxytech/client-native/v6/config-parser/parsers/extra"
	"github.com/haproxytech/client-native/v6/config-parser/parsers/simple"
	"github.com/haproxytech/client-native/v6/config-parser/spoe/parsers"
)

func addParser(psrs map[string]parser.ParserInterface, sequence *[]parser.Section, p parser.ParserInterface) {
	p.Init()
	psrs[p.GetParserName()] = p
	*sequence = append(*sequence, parser.Section(p.GetParserName()))
}

func createParsers(ps map[string]parser.ParserInterface, sequence []parser.Section) *parser.Parsers {
	addParser(ps, &sequence, &parsers.SPOESection{Name: "spoe-agent"})
	addParser(ps, &sequence, &parsers.SPOESection{Name: "spoe-group"})
	addParser(ps, &sequence, &parsers.SPOESection{Name: "spoe-message"})
	addParser(ps, &sequence, &extra.UnProcessed{})

	for _, psr := range ps {
		psr.Init()
	}

	return &parser.Parsers{Parsers: ps, ParserSequence: sequence}
}

func getStartParser() *parser.Parsers {
	p := map[string]parser.ParserInterface{}
	sequence := []parser.Section{}
	addParser(p, &sequence, &extra.ConfigVersion{})
	addParser(p, &sequence, &extra.Comments{})
	return createParsers(p, sequence)
}

func getSPOEAgentParser() *parser.Parsers {
	p := map[string]parser.ParserInterface{}
	sequence := []parser.Section{}
	addParser(p, &sequence, &simple.String{Name: "groups"})
	addParser(p, &sequence, &configparser.Log{})
	addParser(p, &sequence, &simple.Number{Name: "maxconnrate"})
	addParser(p, &sequence, &simple.Number{Name: "maxerrrate"})
	addParser(p, &sequence, &simple.Number{Name: "max-frame-size"})
	addParser(p, &sequence, &simple.Number{Name: "max-waiting-frames"})
	addParser(p, &sequence, &simple.String{Name: "messages"})
	addParser(p, &sequence, &simple.Option{Name: "async"})
	addParser(p, &sequence, &simple.Option{Name: "continue-on-error"})
	addParser(p, &sequence, &simple.Option{Name: "dontlog-normal"})
	addParser(p, &sequence, &simple.Option{Name: "force-set-var"})
	addParser(p, &sequence, &simple.Option{Name: "pipelining"})
	addParser(p, &sequence, &simple.Option{Name: "send-frag-payload"})
	addParser(p, &sequence, &simple.TimeTwoWords{Keywords: []string{"option", "set-on-error"}})
	addParser(p, &sequence, &simple.TimeTwoWords{Keywords: []string{"option", "set-process-time"}})
	addParser(p, &sequence, &simple.TimeTwoWords{Keywords: []string{"option", "set-total-time"}})
	addParser(p, &sequence, &simple.TimeTwoWords{Keywords: []string{"option", "var-prefix"}})
	addParser(p, &sequence, &simple.String{Name: "register-var-names"})
	addParser(p, &sequence, &simple.TimeTwoWords{Keywords: []string{"timeout", "hello"}})
	addParser(p, &sequence, &simple.TimeTwoWords{Keywords: []string{"timeout", "idle"}})
	addParser(p, &sequence, &simple.TimeTwoWords{Keywords: []string{"timeout", "processing"}})
	addParser(p, &sequence, &simple.Word{Name: "use-backend"})
	return createParsers(p, sequence)
}

func getSPOEGroupParser() *parser.Parsers {
	p := map[string]parser.ParserInterface{}
	sequence := []parser.Section{}
	addParser(p, &sequence, &simple.String{Name: "messages"})
	return createParsers(p, sequence)
}

func getSPOEMessageParser() *parser.Parsers {
	p := map[string]parser.ParserInterface{}
	sequence := []parser.Section{}
	addParser(p, &sequence, &configparser.ACL{})
	addParser(p, &sequence, &simple.String{Name: "args"})
	addParser(p, &sequence, &parsers.Event{})
	return createParsers(p, sequence)
}
