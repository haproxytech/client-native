/*
Copyright 2024 HAProxy Technologies

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

import "fmt"

type preferredTimeSuffix struct {
	suffix string
}

func (u preferredTimeSuffix) Set(p *ConfigurationOptions) error {
	switch u.suffix {
	case "nearest":
		p.PreferredTimeSuffix = "d"
	case "none", "ms", "s", "m", "h", "d":
		p.PreferredTimeSuffix = u.suffix
	default:
		return fmt.Errorf("invalid PreferredTimeSuffix value '%s'", u.suffix)
	}
	return nil
}

// PreferredTimeSuffix allows specifying which time unit will be favored when
// serializing Time values in the configuration.
// Allowed values: nearest (default), none, ms, s, m, h, d.
func PreferredTimeSuffix(suffix string) ConfigurationOption {
	return preferredTimeSuffix{
		suffix: suffix,
	}
}
