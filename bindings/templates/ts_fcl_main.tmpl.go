package bindings

func GetTsFclMainTemplate() string {
	const template = `/**
    This binding file was auto generated based on FLIX template v{{.Version}}. 
    Changes to this file might get overwritten.
    Note fcl version {{.FclVersion}} or higher is required to use templates. 
**/

import * as fcl from "@onflow/fcl"
{{- if .IsLocalTemplate }}
import flixTemplate from "{{.Location}}"
{{- else}}
const flixTemplate = "{{.Location}}"
{{- end}}
{{"\n"}}
{{ template "interface" . }}
/**
* {{.Title}}: {{.Description}}
{{- range $param := .Parameters }}
* @param {{$param.JsType}} {{$param.Name}} - {{$param.Description}}
{{- end }}
{{- if not .IsScript }}
* @returns {Promise<string>} - Returns a promise that resolves to the transaction ID
{{- end }}
*/

{{if .IsScript}}
{{- template "script" .}}
{{else}}
{{- template "tx" .}}
{{- end}}




`

	return template
}
