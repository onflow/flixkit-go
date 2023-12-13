package bindings

func GetTsFclTxTemplate() string {
	const template = `{{define "tx"}}export async function {{.Title}}(
  {{- template "params" .}}): Promise<string> {
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
