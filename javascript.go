package flixkit

import (
	"bytes"
	"embed"
	"fmt"
	"strings"
	"text/template"
	"unicode"
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

func GenerateJavaScript(flix *FlowInteractionTemplate, templateLocation string, isLocal bool) (string, error) {
    tmpl, err := template.ParseFS(templateFiles, "templates/*.tmpl")
    if err != nil {
        fmt.Println("Error executing template:", err)
        return "", err
    }

    methodName := TitleToMethodName(flix.Data.Messages.Title.I18N["en-US"])
    description := flix.Data.Messages.Description.I18N["en-US"]
    data := TemplateData{
        Version: flix.FVersion,
        Parameters: TransformArguments(flix.Data.Arguments),
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


func TitleToMethodName(s string) string {
	s = strings.TrimSpace(s)
	var result string
	upperNext := false

	for _, r := range s {
		if r == ' ' || r == '_' || r == '-' {
			upperNext = true
		} else {
			if upperNext {
				result += string(unicode.ToUpper(r))
			} else {
				result += string(unicode.ToLower(r))
			}
			upperNext = false
		}
	}

	return result
}

func TransformArguments(args Arguments) []SimpleParameter {
	simpleArgs := []SimpleParameter{}
	for name, arg := range args {
        isArray, cType, jsType := IsArrayParameter(arg)
        desciption := arg.Messages.Title.I18N["en-US"]
        if isArray {
            simpleArgs = append(simpleArgs, SimpleParameter{Name: name, CadType: cType, JsType: jsType, FclType: "Array(t." + cType + ")", Description: desciption})
        } else {
            jsType := ConvertCadenceTypeToJS(arg.Type)
            simpleArgs = append(simpleArgs, SimpleParameter{Name: name, CadType: arg.Type, JsType: jsType, FclType: arg.Type, Description: desciption})
        }
	}
	return simpleArgs
}


func IsArrayParameter(arg Argument) (bool, string, string) {
    isArray := arg.Type[0] == '[' && arg.Type[len(arg.Type)-1] == ']'
    if (!isArray) {
        return isArray, "", ""
    }
    cadenceType := arg.Type[1 : len(arg.Type)-1]
    jsType := "Array<" + ConvertCadenceTypeToJS(cadenceType) + ">"
    return isArray, cadenceType, jsType
}

// ConvertCadenceTypeToJS takes a Cadence type as a string and returns its JavaScript equivalent
func ConvertCadenceTypeToJS(cadenceType string) string {
    switch cadenceType {
    case "Int":
        return "Number"
    case "Int8":
        return "Number"
    case "Int16":
        return "Number"
    case "Int32":
        return "Number"
    case "Int64":
        return "BigInt"
    case "Int128":
        return "BigInt"
    case "Int256":
        return "BigInt"
    case "UInt":
        return "string"
    case "UInt8":
        return "string"
    case "UInt16":
        return "string"
    case "UInt32":
        return "string"
    case "UInt64":
        return "string"
    case "UInt128":
        return "string"
    case "UInt256":
        return "string"
    case "UFix64":
        return "string"
    case "Fix64":
        return "BigInt"
    case "String":
        return "string"
    case "Character":
        return "string"
    case "Bool":
        return "boolean"
    case "Address":
        return "string"
    case "Void":
        return "void"
    default:
        // For composite and resource types, you can customize further.
        // For now, let's just return 'any' for unknown or complex types
        return "any"
    }
}
