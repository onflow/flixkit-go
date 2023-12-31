package templates

func GetTsFclParamsTemplate() string {
	const template = `{{ define "params" }}
{{- if len .Parameters -}}
{
    {{- range $index, $ele := .Parameters -}}
      {{if $index}}, {{end}}{{.Name}}
    {{- end -}}
  }: {{ .ParametersPrefixName }}Params
{{- end -}}
{{ end }}
	
`

	return template

}
