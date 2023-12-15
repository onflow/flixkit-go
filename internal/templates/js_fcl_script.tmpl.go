package templates

func GetJsFclScriptTemplate() string {
	const template = `{{define "script"}}export async function {{.Title}}( 
{{- template "params" .}}) {
  const info = await fcl.query({
    template: flixTemplate,
    {{ if len .Parameters -}}
    args: (arg, t) => [
      {{- range $index, $ele := .Parameters -}}
        {{if $index}}, {{end}}arg({{.Name}}, t.{{.FclType}})
      {{- end -}}
      ]
    {{- end }}
  });

  return info
}{{end}}
`

	return template

}
