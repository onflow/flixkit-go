package javascript

import (
	"bytes"
	"embed"
	_ "embed"
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
}

//go:embed js_script_template.tpl
var scriptTmpl embed.FS
//go:embed js_tx_template.tpl
var txTmpl embed.FS

func GenerateJavaScript(flix *common.FlowInteractionTemplate, templatePath string) (string, error) {
    var buffer bytes.Buffer
    var t *template.Template
    var err error

    switch flix.Data.Type {
        case "script":
            t, err = template.ParseFS(scriptTmpl, "js_script_template.tpl")
        default:
            t, err = template.ParseFS(txTmpl, "js_tx_template.tpl")
    }

    if err != nil {
        return "", err
    }

    methodName := common.TitleToMethodName(flix.Data.Messages.Title.I18N["en-US"])
    data := TemplateData{
        Parameters: TransformArguments(flix.Data.Arguments),
        Title: methodName,
        Location: templatePath,
    }
    t.Execute(&buffer, data)
    return buffer.String(), nil
}


func TransformArguments(args common.Arguments) []SimpleParameter {
	simpleArgs := []SimpleParameter{}
	for name, arg := range args {
		simpleArgs = append(simpleArgs, SimpleParameter{Name: name, Type: arg.Type})
	}
	return simpleArgs
}
