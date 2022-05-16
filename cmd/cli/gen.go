package main

import (
	"bytes"
	"html/template"
)

const tpl = `type {{ .Name }} struct {
{{- range .Fields }}
	{{ .FieldName }} {{ .Type }}
{{- end }}
}`

type templateArgs struct {
	Name string
	Fields []Field
}

type Field struct {
	FieldName string
	Type string
}

func genStruct(structName string, fieldNames []string, types []string) string {
	t, err := template.New("struct").Parse(tpl)
	if err != nil {
		panic(err)
	}

	fields := make([]Field, len(fieldNames))
	for i, v := range fieldNames {
		fields[i] = Field{
			FieldName: v,
			Type:      types[i],
		}
	}

	tplArgs := templateArgs{
		Name: structName,
		Fields: fields,
	}

	s := bytes.NewBufferString("")
	if err := t.Execute(s, tplArgs); err != nil {
		panic(err)
	}

	return s.String()
}
