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
	Type string
}

type TemplateData struct {
    Version     string
    Parameters []SimpleParameter
    Title       string
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
    data := TemplateData{
        Version: flix.FVersion,
        Parameters: TransformArguments(flix.Data.Arguments),
        Title: methodName,
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
        isArray, arrayType := IsArrayParameter(arg)
        if isArray {
            simpleArgs = append(simpleArgs, SimpleParameter{Name: name, Type: "Array(t." + arrayType + ")"})
        } else {
            simpleArgs = append(simpleArgs, SimpleParameter{Name: name, Type: arg.Type})
        }
	}
	return simpleArgs
}


func IsArrayParameter(arg Argument) (bool, string) {
    isArray := arg.Type[0] == '[' && arg.Type[len(arg.Type)-1] == ']'
    if (!isArray) {
        return isArray, ""
    }
    return isArray, arg.Type[1 : len(arg.Type)-1]
}

