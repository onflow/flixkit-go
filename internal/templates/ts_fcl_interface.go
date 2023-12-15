package templates

func GetTsFclInterfaceTemplate() string {
	const template = `{{ define "interface" }}
{{- if len .Parameters -}}
interface {{ .ParametersPrefixName }}Params {
{{- range .Parameters }}
  {{ .Name }}: {{ .JsType }}; {{- if .Description }} // {{ .Description }} {{- end }}
{{- end }}
}
{{ end }}
{{ end }}
`

	return template
}
