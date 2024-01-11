package flixkit

import (
	"bytes"
	"fmt"
	"sort"
	"text/template"

	v1 "github.com/onflow/flixkit-go/flixkit/v1"
	v1_1 "github.com/onflow/flixkit-go/flixkit/v1_1"
	"github.com/onflow/flixkit-go/internal/templates"
	"github.com/stoewer/go-strcase"
)

func NewFclTSGenerator() Binder {
	t := []string{
		templates.GetTsFclMainTemplate(),
		templates.GetTsFclScriptTemplate(),
		templates.GetTsFclTxTemplate(),
		templates.GetTsFclParamsTemplate(),
		templates.GetTsFclInterfaceTemplate(),
	}

	return &FclGenerator{
		Templates: t,
	}
}

func NewFclJSGenerator() Binder {
	t := []string{
		templates.GetJsFclMainTemplate(),
		templates.GetJsFclScriptTemplate(),
		templates.GetJsFclTxTemplate(),
		templates.GetJsFclParamsTemplate(),
	}

	return &FclGenerator{
		Templates: t,
	}
}

func (g *FclGenerator) Generate(flixString string, templateLocation string) (string, error) {
	tmpl, err := parseTemplates(g.Templates)
	if err != nil {
		return "", err
	}
	if flixString == "" {
		return "", fmt.Errorf("no flix template provided")
	}

	ver, err := GetTemplateVersion(flixString)
	if err != nil {
		return "", fmt.Errorf("invalid flix template version, %w", err)
	}

	isLocal := !isUrl(templateLocation)
	var data templateData
	switch ver {
	case "1.0.0":

		flix, err := v1.ParseFlix(flixString)
		if err != nil {
			return "", err
		}
		data = getTemplateDataV1_0(flix, templateLocation, isLocal)
	case "1.1.0":
		flix, err := v1_1.ParseFlix(flixString)
		if err != nil {
			return "", err
		}
		data = getTemplateDataV1_1(flix, templateLocation, isLocal)
	default:
		return "", fmt.Errorf("invalid flix template version: %s", ver)
	}

	data.FclVersion = GetFlixFclCompatibility(ver)
	var buf bytes.Buffer
	err = tmpl.Execute(&buf, data)
	return buf.String(), err
}

type simpleParameter struct {
	Name        string
	JsType      string
	Description string
	FclType     string
	CadType     string
}

type templateData struct {
	FclVersion           string
	Version              string
	Parameters           []simpleParameter
	ParametersPrefixName string
	Output               simpleParameter
	Title                string
	Description          string
	Location             string
	IsScript             bool
	IsLocalTemplate      bool
}

type FclGenerator struct {
	Templates []string
}

var _ Binder = (*FclGenerator)(nil)

type FlixParameter struct {
	Name string
	Type string
}

func GetFlixFclCompatibility(flixVersion string) string {
	compatibleVersions := map[string]string{
		"1.0.0": "1.3.0",
		"1.1.0": "1.9.0",
		// add more versions here
	}
	v, ok := compatibleVersions[flixVersion]
	if !ok {
		// default to latest if flix version not configured
		return "1.9.0"
	}
	return v
}

func getTemplateDataV1_1(flix *v1_1.InteractionTemplate, templateLocation string, isLocal bool) templateData {
	var msgs v1_1.InteractionTemplateMessages = flix.Data.Messages
	title := msgs.GetTitle("Request")
	methodName := strcase.LowerCamelCase(title)
	description := msgs.GetDescription("")
	var sp simpleParameter

	if flix.Data.Type == "script" {
		oTemp := flix.Data.Output
		if flix.Data.Output == nil {
			// if output is not defined, add default output
			oTemp = &v1_1.Parameter{
				Label:    "result",
				Type:     "String",
				Messages: v1_1.InteractionTemplateMessages{},
			}
		}
		o := transformParameters([]v1_1.Parameter{*oTemp})
		if len(o) > 0 {
			sp = o[0]
		}
	}
	data := templateData{
		Version:              flix.FVersion,
		Parameters:           transformParameters(flix.Data.Parameters),
		ParametersPrefixName: strcase.UpperCamelCase(title),
		Output:               sp,
		Title:                methodName,
		Description:          description,
		Location:             templateLocation,
		IsScript:             flix.IsScript(),
		IsLocalTemplate:      isLocal,
	}
	return data
}

func getTemplateDataV1_0(flix *v1.FlowInteractionTemplate, templateLocation string, isLocal bool) templateData {
	title := flix.Data.Messages.GetTitleValue("Request")
	methodName := strcase.LowerCamelCase(title)
	description := flix.GetDescription()
	var sp simpleParameter
	// version 1.0 does not support output parameters, add default output
	if flix.Data.Type == "script" {
		sp = simpleParameter{
			Name:   "result",
			JsType: "string",
		}
	}
	data := templateData{
		Version:              flix.FVersion,
		Parameters:           transformArguments(flix.Data.Arguments),
		ParametersPrefixName: strcase.UpperCamelCase(title),
		Title:                methodName,
		Description:          description,
		Location:             templateLocation,
		IsScript:             flix.IsScript(),
		IsLocalTemplate:      isLocal,
		Output:               sp,
	}
	return data
}

func convertCadenceTypeToJS(cadenceType string) string {
	// need to determine js type based on fcl supported types
	// looking at fcl types and how arguments work as parameters
	// https://github.com/onflow/fcl-js/blob/master/packages/types/src/types.js
	switch cadenceType {
	case "Bool":
		return "boolean"
	case "Void":
		return "void"
	case "Dictionary":
		return "object"
	case "Struct":
		return "object"
	case "Enum":
		return "object"
	default:
		return "string"
	}
}

func parseTemplates(templates []string) (*template.Template, error) {
	baseTemplate := template.New("base")

	for _, tmplStr := range templates {
		_, err := baseTemplate.Parse(tmplStr)
		if err != nil {
			return nil, err
		}
	}

	return baseTemplate, nil
}

func transformParameters(args []v1_1.Parameter) []simpleParameter {
	simpleArgs := []simpleParameter{}
	if len(args) == 0 {
		return simpleArgs
	}
	sort.Slice(args, func(i, j int) bool {
		return args[i].Index < args[j].Index
	})

	for _, arg := range args {
		isArray, cType, jsType := isArrayParameter(FlixParameter{Name: arg.Label, Type: arg.Type})
		var msgs v1_1.InteractionTemplateMessages = arg.Messages
		desciption := msgs.GetDescription("")
		if isArray {
			simpleArgs = append(simpleArgs, simpleParameter{Name: arg.Label, CadType: cType, JsType: jsType, FclType: "Array(t." + cType + ")", Description: desciption})
		} else {
			jsType := convertCadenceTypeToJS(arg.Type)
			simpleArgs = append(simpleArgs, simpleParameter{Name: arg.Label, CadType: arg.Type, JsType: jsType, FclType: arg.Type, Description: desciption})
		}
	}
	return simpleArgs
}

func transformArguments(args v1.Arguments) []simpleParameter {
	simpleArgs := []simpleParameter{}
	var keys []string
	// get keys for sorting
	for k := range args {
		keys = append(keys, k)
	}

	sort.SliceStable(keys, func(i, j int) bool {
		return args[keys[i]].Index < args[keys[j]].Index
	})
	for _, key := range keys {
		arg := args[key]
		isArray, cType, jsType := isArrayParameter(FlixParameter{Name: key, Type: arg.Type})
		desciption := arg.Messages.GetTitleValue("")
		if isArray {
			simpleArgs = append(simpleArgs, simpleParameter{Name: key, CadType: cType, JsType: jsType, FclType: "Array(t." + cType + ")", Description: desciption})
		} else {
			jsType := convertCadenceTypeToJS(arg.Type)
			simpleArgs = append(simpleArgs, simpleParameter{Name: key, CadType: arg.Type, JsType: jsType, FclType: arg.Type, Description: desciption})
		}
	}
	return simpleArgs
}

func isArrayParameter(arg FlixParameter) (isArray bool, cType string, jsType string) {
	if arg.Type == "" || arg.Type[0] != '[' {
		return false, "", ""
	}
	cadenceType := arg.Type[1 : len(arg.Type)-1]
	javascriptType := "Array<" + convertCadenceTypeToJS(cadenceType) + ">"
	return true, cadenceType, javascriptType
}
