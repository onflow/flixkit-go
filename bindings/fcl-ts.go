package bindings

import (
	"bytes"
	"fmt"
	"path/filepath"

	"github.com/onflow/flixkit-go"
	bindings "github.com/onflow/flixkit-go/bindings/templates"
	v1_1 "github.com/onflow/flixkit-go/flixkitv1_1"
)

func NewFclTSGenerator() *FclGenerator {
	templates := []string{
		bindings.GetTsFclMainTemplate(),
		bindings.GetTsFclScriptTemplate(),
		bindings.GetTsFclTxTemplate(),
		bindings.GetTsFclParamsTemplate(),
		bindings.GetTsFclInterfaceTemplate(),
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
