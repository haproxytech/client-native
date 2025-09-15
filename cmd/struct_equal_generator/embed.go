package main

import (
	_ "embed"
)

//go:embed header.tmpl
var tmplHeader string

//go:embed test.tmpl
var tmplCompareTest string
