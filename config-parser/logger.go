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

package parser

import (
	"errors"
)

func (p *configParser) SetLoggerState(active bool) error {
	if p.Options.Logger == nil {
		return errors.New("logger is not set")
	}
	p.Options.Logger.Debugf("%slogger set to state: %v", p.Options.LogPrefix, active)
	p.Options.Log = active
	return nil
}
