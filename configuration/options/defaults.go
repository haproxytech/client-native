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

const (
	// DefaultUseValidation sane default using validation in client native
	DefaultUseValidation = true
	// DefaultPersistentTransactions sane default using persistent transactions in client native
	DefaultPersistentTransactions = true
	// DefaultValidateConfigurationFile is used to validate HAProxy configuration file
	DefaultValidateConfigurationFile = true

	// DefaultConfigurationFile sane default for path to haproxy configuration file
	DefaultConfigurationFile = "/etc/haproxy/haproxy.cfg"

	// DefaultHaproxy sane default for path to haproxy executable
	DefaultHaproxy = "/usr/sbin/haproxy"

	// DefaultTransactionDir sane default for path for transactions
	DefaultTransactionDir = "/etc/haproxy/transactions"
)
