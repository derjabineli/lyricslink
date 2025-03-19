package main

import (
	"html/template"
	"strings"
)

func lyricSheetToHTML(lyricSheet string) template.HTML {
	lyricSheet = strings.ReplaceAll(lyricSheet, "\\n", "\n")
	return template.HTML(strings.ReplaceAll(lyricSheet, "\n", "<br>"))
}