package main

import (
	"bytes"
	"text/template"
)

func applyTemplate(filenameTmpl string, inv map[string]string) string {
	t := template.New("template")
	template.Must(t.Parse(filenameTmpl))
	var buf bytes.Buffer
	t.Execute(&buf, inv)
	return buf.String()
}
