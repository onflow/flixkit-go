package js

import (
	"bytes"
	"embed"
	_ "embed"
	"fmt"
	"net/url"
	"path"
	"text/template"

	"github.com/onflow/flixkit-go/common"
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

func GenerateJavaScript(flix *common.FlowInteractionTemplate, templateLocation string) (string, error) {
    tmpl, err := template.ParseFS(templateFiles, "templates/*.tmpl")
    if err != nil {
        fmt.Println("Error executing template:", err)
        return "", err
    }

    templatePath, IsLocal, _ := GetTemplateReference(templateLocation)
    methodName := common.TitleToMethodName(flix.Data.Messages.Title.I18N["en-US"])
    data := TemplateData{
        Version: flix.FVersion,
        Parameters: TransformArguments(flix.Data.Arguments),
        Title: methodName,
        Location: templatePath,
        IsScript: flix.IsScript(),
        IsLocalTemplate: IsLocal,
    }

    var buf bytes.Buffer
    err = tmpl.Execute(&buf, data)
    return buf.String(), err    
}


func TransformArguments(args common.Arguments) []SimpleParameter {
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


func IsArrayParameter(arg common.Argument) (bool, string) {
    isArray := arg.Type[0] == '[' && arg.Type[len(arg.Type)-1] == ']'
    if (!isArray) {
        return isArray, ""
    }
    return isArray, arg.Type[1 : len(arg.Type)-1]
}

func GetTemplateReference(templateLocation string) (string, bool, error) {
    var err error
    templatePath := templateLocation
    isLocal := common.IsLocalTemplate(templateLocation)
    if isLocal {
        parsedURL, err := url.Parse(templateLocation)
        if err != nil {
            return templatePath, isLocal, err
        }
        templatePath = "./" + path.Base(parsedURL.Path)   
    }
    return templatePath, isLocal, err
}
