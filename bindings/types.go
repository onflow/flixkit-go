package bindings

import (
	"fmt"
	"net/url"
	"text/template"

	"github.com/onflow/flixkit-go"
	v1_1 "github.com/onflow/flixkit-go/flixkitv1_1"
	"github.com/stoewer/go-strcase"
)

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
	Title                string
	Description          string
	Location             string
	IsScript             bool
	IsLocalTemplate      bool
}

type FclGenerator struct {
	Templates []string
}

type FlixParameter struct {
	Name string
	Type string
}

func isUrl(str string) bool {
	u, err := url.Parse(str)
	return err == nil && u.Scheme != "" && u.Host != ""
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
	fmt.Println("hello v1_1")
	var msgs v1_1.InteractionTemplateMessages = flix.Data.Messages
	title := msgs.GetTitle("Request")
	methodName := strcase.LowerCamelCase(title)
	description := msgs.GetDescription("")
	fmt.Println("found xxx", description, "scirption", title)
	data := templateData{
		Version:              flix.FVersion,
		Parameters:           transformParameters(flix.Data.Parameters),
		ParametersPrefixName: strcase.UpperCamelCase(title),
		Title:                methodName,
		Description:          description,
		Location:             templateLocation,
		IsScript:             flix.IsScript(),
		IsLocalTemplate:      isLocal,
	}
	return data
}

func getTemplateDataV1_0(flix *flixkit.FlowInteractionTemplate, templateLocation string, isLocal bool) templateData {
	title := flix.Data.Messages.GetTitleValue("Request")
	methodName := strcase.LowerCamelCase(title)
	description := flix.GetDescription()

	data := templateData{
		Version:              flix.FVersion,
		Parameters:           transformArguments(flix.Data.Arguments),
		ParametersPrefixName: strcase.UpperCamelCase(title),
		Title:                methodName,
		Description:          description,
		Location:             templateLocation,
		IsScript:             flix.IsScript(),
		IsLocalTemplate:      isLocal,
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
