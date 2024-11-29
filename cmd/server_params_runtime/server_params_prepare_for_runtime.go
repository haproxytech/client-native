package main

import (
	"errors"
)

const (
	FuncDoNotSendDisabledFields          = "prepareForRuntimeDoNotSendDisabledFields"
	FuncDoNotSendEnabledFields           = "prepareForRuntimeDoNotSendEnabledFields"
	FuncPrepareFieldsForRuntimeAddServer = "PrepareFieldsForRuntimeAddServer"
)

var ServerParamsPrepareForRuntimeMap = map[string]string{ //nolint:gochecknoglobals
	"AgentCheck":       FuncDoNotSendDisabledFields,
	"Backup":           FuncDoNotSendDisabledFields,
	"Check":            FuncDoNotSendDisabledFields,
	"CheckSendProxy":   FuncDoNotSendDisabledFields,
	"CheckSsl":         FuncDoNotSendDisabledFields,
	"CheckViaSocks4":   FuncDoNotSendDisabledFields,
	"ForceSslv3":       FuncDoNotSendDisabledFields,
	"Sslv3":            FuncDoNotSendDisabledFields,
	"ForceTlsv10":      FuncDoNotSendDisabledFields,
	"Tlsv10":           FuncDoNotSendDisabledFields,
	"ForceTlsv11":      FuncDoNotSendDisabledFields,
	"Tlsv11":           FuncDoNotSendDisabledFields,
	"ForceTlsv12":      FuncDoNotSendDisabledFields,
	"Tlsv12":           FuncDoNotSendDisabledFields,
	"ForceTlsv13":      FuncDoNotSendDisabledFields,
	"Tlsv13":           FuncDoNotSendDisabledFields,
	"Maintenance":      FuncDoNotSendDisabledFields,
	"NoSslv3":          FuncDoNotSendEnabledFields,
	"NoTlsv10":         FuncDoNotSendEnabledFields,
	"NoTlsv11":         FuncDoNotSendEnabledFields,
	"NoTlsv12":         FuncDoNotSendEnabledFields,
	"NoTlsv13":         FuncDoNotSendEnabledFields,
	"NoVerifyhost":     FuncDoNotSendEnabledFields,
	"SendProxy":        FuncDoNotSendDisabledFields,
	"SendProxyV2":      FuncDoNotSendDisabledFields,
	"SendProxyV2Ssl":   FuncDoNotSendDisabledFields,
	"SendProxyV2SslCn": FuncDoNotSendDisabledFields,
	"Ssl":              FuncDoNotSendDisabledFields,
	"SslReuse":         FuncDoNotSendDisabledFields,
	"Stick":            FuncDoNotSendDisabledFields,
	"Tfo":              FuncDoNotSendDisabledFields,
	"TLSTickets":       FuncDoNotSendDisabledFields,
}

func checkMissingEnumFields(allFields []string) ([]string, error) {
	missingFields := []string{}
	for _, field := range allFields {
		// check that all enum
		// Enum: [enabled disabled]"
		// fields have an entry in the ServerParamsPrepareForRuntimeMap
		f, ok := ServerParamsPrepareForRuntimeMap[field]
		if !ok || (f != FuncDoNotSendEnabledFields && f != FuncDoNotSendDisabledFields) {
			missingFields = append(missingFields, field)
		}
	}
	if len(missingFields) > 0 {
		return missingFields, errors.New("missing enum fields")
	}
	return missingFields, nil
}

func listEmptyDisabledFields(allFields []string) []string {
	var emptyDisabledFields []string
	for _, field := range allFields {
		if f, ok := ServerParamsPrepareForRuntimeMap[field]; ok {
			if f == FuncDoNotSendDisabledFields {
				emptyDisabledFields = append(emptyDisabledFields, field)
			}
		}
	}
	return emptyDisabledFields
}

func listEmtpyEnabledFields(allFields []string) []string {
	var emptyEnabledFields []string
	for _, field := range allFields {
		if f, ok := ServerParamsPrepareForRuntimeMap[field]; ok {
			if f == FuncDoNotSendEnabledFields {
				emptyEnabledFields = append(emptyEnabledFields, field)
			}
		}
	}
	return emptyEnabledFields
}
