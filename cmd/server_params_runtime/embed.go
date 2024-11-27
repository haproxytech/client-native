package main

import (
	_ "embed"
)

//go:embed generate.tmpl
var tmplResetEnumDisabledFields string
