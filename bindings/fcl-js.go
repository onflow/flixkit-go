package bindings

import (
	"bytes"
	"os"
	"path/filepath"
	"runtime"
	"sort"
	"text/template"

	"github.com/onflow/flixkit-go"
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
	Version         string
	Parameters      []simpleParameter
	Title           string
	Description     string
	Location        string
	IsScript        bool
	IsLocalTemplate bool
}

type FclJSGenerator struct {
	TemplateDir string
}

func NewFclJSGenerator() *FclJSGenerator {
	_, currentFilePath, _, _ := runtime.Caller(0)
	baseDir := filepath.Dir(currentFilePath)
	templateDir := filepath.Join(baseDir, "templates")

	return &FclJSGenerator{
		TemplateDir: templateDir,
	}
}

func (g FclJSGenerator) Generate(flix *flixkit.FlowInteractionTemplate, templateLocation string, isLocal bool) (string, error) {
	templateFiles, err := getAllFiles(g.TemplateDir)
	if err != nil {
		return "", err
	}
	tmpl, err := template.ParseFiles(templateFiles...)
	if err != nil {
		return "", err
	}

	methodName := strcase.LowerCamelCase(flix.Data.Messages.GetTitleValue("Request"))
	description := flix.GetDescription()
	data := templateData{
		Version:         flix.FVersion,
		Parameters:      transformArguments(flix.Data.Arguments),
		Title:           methodName,
		Description:     description,
		Location:        templateLocation,
		IsScript:        flix.IsScript(),
		IsLocalTemplate: isLocal,
	}

	var buf bytes.Buffer
	err = tmpl.Execute(&buf, data)
	return buf.String(), err
}

func transformArguments(args flixkit.Arguments) []simpleParameter {
	simpleArgs := []simpleParameter{}
	var keys []string
	// get keys for sorting
	for k := range args {
		keys = append(keys, k)
	}

	sort.SliceStable(keys, func(i, j int) bool {
		return args[keys[i]].Index < args[keys[j]].Index
	})
	for _, key := range keys {
		arg := args[key]
		isArray, cType, jsType := isArrayParameter(arg)
		desciption := arg.Messages.GetTitleValue("")
		if isArray {
			simpleArgs = append(simpleArgs, simpleParameter{Name: key, CadType: cType, JsType: jsType, FclType: "Array(t." + cType + ")", Description: desciption})
		} else {
			jsType := convertCadenceTypeToJS(arg.Type)
			simpleArgs = append(simpleArgs, simpleParameter{Name: key, CadType: arg.Type, JsType: jsType, FclType: arg.Type, Description: desciption})
		}
	}
	return simpleArgs
}

func isArrayParameter(arg flixkit.Argument) (isArray bool, cType string, jsType string) {
	if arg.Type == "" || arg.Type[0] != '[' {
		return false, "", ""
	}
	cadenceType := arg.Type[1 : len(arg.Type)-1]
	javascriptType := "Array<" + convertCadenceTypeToJS(cadenceType) + ">"
	return true, cadenceType, javascriptType
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

func getAllFiles(dir string) ([]string, error) {
	var files []string
	err := filepath.Walk(dir, func(path string, info os.FileInfo, err error) error {
		if err != nil {
			return err
		}
		// If it's a directory, skip it
		if info.IsDir() {
			return nil
		}
		files = append(files, path)
		return nil
	})
	if err != nil {
		return nil, err
	}
	return files, nil
}
