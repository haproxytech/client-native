/*
Copyright 2022 HAProxy Technologies

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

type backups struct {
	backupsNumber int
}

func (u backups) Set(p *ConfigurationOptions) error {
	p.BackupsNumber = u.backupsNumber
	return nil
}

// Backups sets number of backups for configuration file
func Backups(number int) ConfigurationOption {
	return backups{
		backupsNumber: number,
	}
}

type backupsDir struct {
	backupsDir string
}

func (u backupsDir) Set(p *ConfigurationOptions) error {
	p.BackupsDir = u.backupsDir
	return nil
}

func BackupsDir(bckDir string) ConfigurationOption {
	return backupsDir{
		backupsDir: bckDir,
	}
}
