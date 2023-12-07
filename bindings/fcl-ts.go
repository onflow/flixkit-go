package bindings

import (
	"bytes"
	"fmt"
	"path/filepath"
	"text/template"

	"github.com/onflow/flixkit-go"
	bindings "github.com/onflow/flixkit-go/bindings/templates"
	v1_1 "github.com/onflow/flixkit-go/flixkitv1_1"
	"github.com/stoewer/go-strcase"
)

func NewFclTSGenerator() *FclGenerator {
	templates := []string{
		bindings.GetTsFclMainTemplate(),
		bindings.GetTsFclScriptTemplate(),
		bindings.GetTsFclTxTemplate(),
		bindings.GetTsFclParamsTemplate(),
	}

	return &FclGenerator{
		Templates: templates,
	}
}

func (g FclGenerator) GenerateTS(flixString string, templateLocation string) (string, error) {
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

func GetTemplateDataV1_1(flix *v1_1.InteractionTemplate, templateLocation string, isLocal bool) templateData {
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

// GetRelativePath computes the relative path from generated file to flix json file.
// This path is used in the binding file to reference the flix json file.
func GetRelativePath(configFile, bindingFile string) (string, error) {
	relPath, err := filepath.Rel(filepath.Dir(bindingFile), configFile)
	if err != nil {
		return "", err
	}

	// If the file is in the same directory and doesn't start with "./", prepend it.
	if !filepath.IsAbs(relPath) && relPath[0] != '.' {
		relPath = "./" + relPath
	}

	// Currently binding files are js, we need to convert the path to unix style
	return filepath.ToSlash(relPath), nil
}
