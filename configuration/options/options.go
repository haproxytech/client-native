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

type ConfigurationOptions struct {
	ConfigurationFile               string
	Haproxy                         string
	TransactionDir                  string
	BackupsDir                      string
	PersistentTransactions          bool
	SkipFailedTransactions          bool
	BackupsNumber                   int
	UseModelsValidation             bool
	SkipConfigurationFileValidation bool // opposite of previously available ValidateConfigurationFile
	MasterWorker                    bool
	UseMd5Hash                      bool

	// ValidateCmd allows specifying a custom script to validate the transaction file.
	// The injected environment variable DATAPLANEAPI_TRANSACTION_FILE must be used to get the location of the file.
	ValidateCmd               string
	ValidateConfigFilesBefore []string
	ValidateConfigFilesAfter  []string
}

type ConfigurationOption interface {
	Set(p *ConfigurationOptions) error
}
