package bindings

func GetJsFclParamsTemplate() string {
	const template = `{{define "params"}}
 {{- if len .Parameters -}}
  {
    {{- range $index, $ele := .Parameters -}}
      {{if $index}}, {{end}}{{.Name}}
    {{- end -}}
  }
  {{- end -}}
{{end}}
`

	return template

}
