package main

import (
	"bytes"
	"regexp"
	"text/template"
)

func applyTemplate(filenameTmpl string, inv map[string]string) string {
	t := template.New("template")
	template.Must(t.Parse(filenameTmpl))
	var buf bytes.Buffer
	t.Execute(&buf, inv)
	return buf.String()
}

func fetchVideoIDFromFilename(filename string) string {
	re, _ := regexp.Compile("\\[([a-z]{2}?\\d+)\\]")
	result := re.FindAllStringSubmatch(filename, -1)

	return result[0][1]
}
