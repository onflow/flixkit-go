package internal

import (
	"encoding/hex"
	"encoding/json"
	"fmt"
	"net/url"
	"os"
)

func getTemplateVersion(template string) (string, error) {
	type FlowInteractionTemplateVersion struct {
		FVersion string `json:"f_version"`
	}
	var flowTemplate FlowInteractionTemplateVersion

	err := json.Unmarshal([]byte(template), &flowTemplate)
	if err != nil {
		return "", err
	}

	if flowTemplate.FVersion == "" {
		return "", fmt.Errorf("version not found")
	}

	return flowTemplate.FVersion, nil
}

func isArrayParameter(arg FlixParameter) (isArray bool, cType string, jsType string) {
	if arg.Type == "" || arg.Type[0] != '[' {
		return false, "", ""
	}
	cadenceType := arg.Type[1 : len(arg.Type)-1]
	javascriptType := "Array<" + convertCadenceTypeToJS(cadenceType) + ">"
	return true, cadenceType, javascriptType
}

func isUrl(str string) bool {
	u, err := url.Parse(str)
	return err == nil && u.Scheme != "" && u.Host != ""
}

type flixQueryTypes string

const (
	flixName flixQueryTypes = "name"
	flixPath flixQueryTypes = "path"
	flixId   flixQueryTypes = "id"
	flixUrl  flixQueryTypes = "url"
	flixJson flixQueryTypes = "json"
)

func isHex(str string) bool {
	if len(str) != 64 {
		return false
	}
	_, err := hex.DecodeString(str)
	return err == nil
}

func isPath(path string) bool {
	_, err := os.Stat(path)
	return err == nil
}

func isJson(str string) bool {
	var js json.RawMessage
	return json.Unmarshal([]byte(str), &js) == nil
}

func getType(s string) flixQueryTypes {
	switch {
	case isPath(s):
		return flixPath
	case isHex(s):
		return flixId
	case isUrl(s):
		return flixUrl
	case isJson(s):
		return flixJson
	default:
		return flixName
	}
}
