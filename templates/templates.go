package templates

import (
	"bytes"
	"text/template"
)

// WriteTemplate takes in the template and the definition of its variables
// and returns a filled-out template.
func WriteTemplate(templateStr string, templateStruct interface{}) (string, error) {

	tmpl, err := template.
		New("test").
		Funcs(template.FuncMap{"index1": func(i int) int { return i + 1 }}).
		Parse(templateStr)
	if err != nil {
		return "", err
	}

	var tmpBuf bytes.Buffer
	err = tmpl.Execute(&tmpBuf, templateStruct)
	if err != nil {
		return "", err
	}
	return tmpBuf.String(), nil
}
