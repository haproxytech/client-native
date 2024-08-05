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

import "github.com/haproxytech/go-logger"

type logging struct {
	log    logger.Format
	prefix string
}

func (u logging) Set(p *Parser) error {
	p.Logger = u.log
	p.Log = true
	if u.prefix != "" {
		p.LogPrefix = u.prefix + " "
	}
	return nil
}

// Logger takes acceptable logger that will be used for logging
func Logger(log logger.Format) ParserOption {
	return LoggerWithPrefix(log, "")
}

// Logger takes acceptable logger that will be used for logging, prefix can be defined to distinguish log messages generated in this package
func LoggerWithPrefix(log logger.Format, prefix string) ParserOption {
	return logging{
		log:    log,
		prefix: prefix,
	}
}
