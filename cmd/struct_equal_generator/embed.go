package main

import (
	_ "embed"
)

//go:embed header.tmpl
var tmplHeader string

//go:embed generate.tmpl
var tmplEqualAndDiff string

//go:embed test.tmpl
var tmplCompareTest string
