package templates

func GetJsFclMainTemplate() string {
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
{{"\n"}}/**
* {{.Description}}{{"\n"}}
   {{- if gt (len .Parameters) 0 -}}
* @param {Object} Parameters - parameters for the cadence
    {{- range $index, $ele := .Parameters -}}
{{"\n"}}* @param {{"{"}}{{$ele.JsType}}{{"}"}} Parameters.{{$ele.Name}} - {{$ele.Description}}: {{$ele.CadType}}
    {{- end -}}
    {{ else -}}
* No parameters needed.
 {{- end -}}
 {{- if not .IsScript -}}
{{"\n"}}* @returns {Promise<string>} - returns a promise which resolves to the transaction id
 {{- end -}}
{{- "\n"}}*/
{{if .IsScript}}
{{- template "script" .}}
{{else}}
{{- template "tx" .}}
{{- end}}




`

	return template
}
