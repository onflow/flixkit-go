package bindings

func GetJsFclScriptTemplate() string {
	const template = `{{define "script"}}export async function {{.Title}}( 
 {{- if len .Parameters -}}
  {
    {{- range $index, $ele := .Parameters -}}
      {{if $index}}, {{end}}{{.Name}}
    {{- end -}}
  }
  {{- end -}}
) {
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
