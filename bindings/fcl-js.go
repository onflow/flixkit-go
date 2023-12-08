package bindings

import (
	"bytes"
	"fmt"

	"github.com/onflow/flixkit-go"
	bindings "github.com/onflow/flixkit-go/bindings/templates"
	v1_1 "github.com/onflow/flixkit-go/flixkitv1_1"
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
