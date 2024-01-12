package flixkitv2

import (
	"encoding/json"
	"fmt"
	"net/url"
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