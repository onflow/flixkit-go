package bindings

import (
	"bytes"
	"fmt"
	"sort"

	"github.com/onflow/flixkit-go"
	bindings "github.com/onflow/flixkit-go/bindings/templates"
	v1_1 "github.com/onflow/flixkit-go/flixkitv1_1"
	"github.com/stoewer/go-strcase"
)

func NewFclJSGenerator() *FclGenerator {
	templates := []string{
		bindings.GetJsFclMainTemplate(),
		bindings.GetJsFclScriptTemplate(),
		bindings.GetJsFclTxTemplate(),
		bindings.GetJsFclParamsTemplate(),
	}

	return &FclGenerator{
		Templates: templates,
	}
}

func (g FclGenerator) GenerateJS(flixString string, templateLocation string) (string, error) {
	tmpl, err := parseTemplates(g.Templates)
	if err != nil {
		return "", err
	}
	if flixString == "" {
		return "", fmt.Errorf("no flix template provided")
	}
	isLocal := !isUrl(templateLocation)

	ver, err := flixkit.GetTemplateVersion(flixString)
	if err != nil {
		return "", fmt.Errorf("invalid flix template version, %s", err)
	}
	var data templateData
	data.FclVersion = GetFlixFclCompatibility(ver)
	if ver == "1.0.0" {
		flix, err := flixkit.ParseFlix(flixString)
		if err != nil {
			return "", err
		}
		data = getTemplateDataV1_0(flix, templateLocation, isLocal)

	} else if ver == "1.1.0" {
		flix, err := v1_1.ParseFlix(flixString)
		if err != nil {
			return "", err
		}
		data = getTemplateDataV1_1(flix, templateLocation, isLocal)

	} else {
		return "", fmt.Errorf("invalid flix template version, support v1.0.0 and v1.1.0")
	}
	data.FclVersion = GetFlixFclCompatibility(ver)

	var buf bytes.Buffer
	err = tmpl.Execute(&buf, data)
	return buf.String(), err
}

func getTemplateDataV1_1(flix *v1_1.InteractionTemplate, templateLocation string, isLocal bool) templateData {
	var msgs v1_1.InteractionTemplateMessages = flix.Data.Messages
	methodName := strcase.LowerCamelCase(msgs.GetTitle("Request"))
	description := msgs.GetDescription("")
	data := templateData{
		Version:         flix.FVersion,
		Parameters:      transformParameters(flix.Data.Parameters),
		Title:           methodName,
		Description:     description,
		Location:        templateLocation,
		IsScript:        flix.IsScript(),
		IsLocalTemplate: isLocal,
	}
	return data
}

func getTemplateDataV1_0(flix *flixkit.FlowInteractionTemplate, templateLocation string, isLocal bool) templateData {
	methodName := strcase.LowerCamelCase(flix.Data.Messages.GetTitleValue("Request"))
	description := flix.GetDescription()
	data := templateData{
		Version:         flix.FVersion,
		Parameters:      transformArguments(flix.Data.Arguments),
		Title:           methodName,
		Description:     description,
		Location:        templateLocation,
		IsScript:        flix.IsScript(),
		IsLocalTemplate: isLocal,
	}
	return data
}

func transformParameters(args []v1_1.Parameter) []simpleParameter {
	simpleArgs := []simpleParameter{}
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

func transformArguments(args flixkit.Arguments) []simpleParameter {
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
