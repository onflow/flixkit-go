package flixkit

import (
	"bytes"
	"embed"
	"text/template"

	"github.com/stoewer/go-strcase"
)

type SimpleParameter struct {
	Name string
	JsType string
    Description string
    FclType string
    CadType string
}

type TemplateData struct {
    Version     string
    Parameters []SimpleParameter
    Title       string
    Description string
    Location    string
    IsScript    bool
    IsLocalTemplate bool
}

//go:embed templates/*.tmpl
var templateFiles embed.FS

type JavaScriptGenerator struct{}

func (g JavaScriptGenerator) Generate(flix *FlowInteractionTemplate, templateLocation string, isLocal bool) (string, error) {
    tmpl, err := template.ParseFS(templateFiles, "templates/*.tmpl")
    if err != nil {
        return "", err
    }

    methodName := strcase.LowerCamelCase(getMessageValue(flix.Data.Messages, "Request"))
    description := flix.GetDescription()
    data := TemplateData{
        Version: flix.FVersion,
        Parameters: transformArguments(flix.Data.Arguments),
        Title: methodName,
        Description: description,
        Location: templateLocation,
        IsScript: flix.IsScript(),
        IsLocalTemplate: isLocal,
    }

    var buf bytes.Buffer
    err = tmpl.Execute(&buf, data)
    return buf.String(), err    
}

func getMessageValue(messages Messages, placeholder string) string {
    s := placeholder
    if messages.Title != nil && 
        messages.Title.I18N != nil {
        // TODO: relying on en-US for now, future we need to know what language to use
        value, exists := messages.Title.I18N["en-US"]
        if exists {
            s = value
        } 
    }
    return s
}

func transformArguments(args Arguments) []SimpleParameter {
	simpleArgs := []SimpleParameter{}
	for name, arg := range args {
        isArray, cType, jsType := isArrayParameter(arg)
        desciption := getMessageValue(arg.Messages, "")
        if isArray {
            simpleArgs = append(simpleArgs, SimpleParameter{Name: name, CadType: cType, JsType: jsType, FclType: "Array(t." + cType + ")", Description: desciption})
        } else {
            jsType := convertCadenceTypeToJS(arg.Type)
            simpleArgs = append(simpleArgs, SimpleParameter{Name: name, CadType: arg.Type, JsType: jsType, FclType: arg.Type, Description: desciption})
        }
	}
	return simpleArgs
}


func isArrayParameter(arg Argument) (bool, string, string) {
    if arg.Type == "" || arg.Type[0] != '[' {
        return false, "", ""
    }
    cadenceType := arg.Type[1 : len(arg.Type)-1]
    jsType := "Array<" + convertCadenceTypeToJS(cadenceType) + ">"
    return true, cadenceType, jsType
}

func convertCadenceTypeToJS(cadenceType string) string {
    // need to determine js type based on fcl supported types
    // looking at fcl types and how arguments work as parameters
    // https://github.com/onflow/fcl-js/blob/master/packages/types/src/types.js
    switch cadenceType {
    case "Bool":
        return "boolean"
    case "Void":
        return "void" // return type only
    case "Dictionary":
        return "object" // TODO: support Collection type, test to see what fcl 
    case "Struct":
        return "object" // TODO: support Composite type, test to see what fcl 
    case "Enum":
        return "object" // TODO: support Composite type, test to see what fcl 
    default:
        return "string"
    }
}
