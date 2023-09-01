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
		simpleArgs = append(simpleArgs, SimpleParameter{Name: name, Type: arg.Type})
	}
	return simpleArgs
}

func GetTemplateReference(templateLocation string) (string, bool, error) {
    var err error
    templatePath := templateLocation
    IsLocal := common.IsLocalTemplate(templateLocation)
    if (IsLocal) {
        parsedURL, err := url.Parse(templateLocation)
        if err != nil {
            return templatePath, IsLocal, err
        }
        templatePath = "./" + path.Base(parsedURL.Path)   
    }
    return templatePath, IsLocal, err
}