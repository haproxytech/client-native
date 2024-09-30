/*
Copyright 2021 HAProxy Technologies

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

package options

import (
	"io"
	"strings"
)

type filename struct {
	Path string
}

func (f filename) Set(p *Parser) error {
	p.Path = f.Path
	return nil
}

// Reader takes path where configuration is stored
func Path(path string) ParserOption {
	return filename{
		Path: path,
	}
}

type reader struct {
	Reader io.Reader
}

func (f reader) Set(p *Parser) error {
	p.Reader = f.Reader
	return nil
}

// Reader takes io.Reader that will be used to parse data
func Reader(ioReader io.Reader) ParserOption {
	return reader{
		Reader: ioReader,
	}
}

// String takes string that will be used to parse data
func String(configuration string) ParserOption {
	return reader{
		Reader: strings.NewReader(configuration),
	}
}
