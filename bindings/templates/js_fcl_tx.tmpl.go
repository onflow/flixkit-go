package bindings

func GetJsFclTxTemplate() string {
	const template = `{{define "tx"}}export async function {{.Title}}({ 
 {{- if len .Parameters -}}
    {{- range $index, $ele := .Parameters -}}
      {{if $index}}, {{end}}{{.Name}}
    {{- end -}}
  {{- end -}}
}) {
  const transactionId = await fcl.mutate({
    template: flixTemplate,
    {{ if len .Parameters -}}
    args: (arg, t) => [
      {{- range $index, $ele := .Parameters -}}
        {{if $index}}, {{end}}arg({{.Name}}, t.{{.FclType}})
      {{- end -}}
      ]
    {{- end }}
  });

  return transactionId
}{{end}}
`

	return template

}
