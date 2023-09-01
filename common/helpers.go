package common

import (
	"fmt"
	"regexp"
	"strings"
	"unicode"
)

type FlowInteractionTemplate struct {
	FType    string `json:"f_type"`
	FVersion string `json:"f_version"`
	ID       string `json:"id"`
	Data     Data   `json:"data"`
}

func (t *FlowInteractionTemplate) IsScript() bool {
	return t.Data.Type == "script"
}

func (t *FlowInteractionTemplate) IsTransaction() bool {
	return t.Data.Type == "transaction"
}

func (t *FlowInteractionTemplate) GetAndReplaceCadenceImports(networkName string) (string, error) {
	cadence := t.Data.Cadence

	for dependencyAddress, contracts := range t.Data.Dependencies {
		for contractName, networks := range contracts {
			network, ok := networks[networkName]
			if !ok {
				return "", fmt.Errorf("network %s not found for contract %s", networkName, contractName)
			}

			pattern := fmt.Sprintf(`import\s*%s\s*from\s*%s`, contractName, dependencyAddress)
			re, err := regexp.Compile(pattern)
			if err != nil {
				return "", fmt.Errorf("invalid regex pattern: %v", err)
			}

			replacement := fmt.Sprintf("import %s from %s", contractName, network.Address)
			cadence = re.ReplaceAllString(cadence, replacement)
		}
	}

	return cadence, nil
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

func IsLocalTemplate(templateLocation string) bool {
	return strings.HasPrefix(templateLocation, "/") || 
		strings.HasPrefix(templateLocation, "./") || 
		strings.HasPrefix(templateLocation, "../")
}

