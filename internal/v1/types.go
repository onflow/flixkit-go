package v1

import (
	"encoding/json"
	"fmt"
	"regexp"
)

type Network struct {
	Address        string `json:"address"`
	FqAddress      string `json:"fq_address"`
	Contract       string `json:"contract"`
	Pin            string `json:"pin"`
	PinBlockHeight uint64 `json:"pin_block_height"`
}

type Argument struct {
	Index    int      `json:"index"`
	Type     string   `json:"type"`
	Messages Messages `json:"messages"`
	Balance  string   `json:"balance"`
}

type Title struct {
	I18N map[string]string `json:"i18n"`
}

type Description struct {
	I18N map[string]string `json:"i18n"`
}

type Messages struct {
	Title       *Title       `json:"title,omitempty"`
	Description *Description `json:"description,omitempty"`
}

type Dependencies map[string]Contracts
type Contracts map[string]Networks
type Networks map[string]Network
type Arguments map[string]Argument

type Data struct {
	Type         string       `json:"type"`
	Interface    string       `json:"interface"`
	Messages     Messages     `json:"messages"`
	Cadence      string       `json:"cadence"`
	Dependencies Dependencies `json:"dependencies"`
	Arguments    Arguments    `json:"arguments"`
}

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

func ParseFlix(template string) (*FlowInteractionTemplate, error) {
	var flowTemplate FlowInteractionTemplate

	err := json.Unmarshal([]byte(template), &flowTemplate)
	if err != nil {
		return nil, err
	}

	return &flowTemplate, nil
}

func (t *FlowInteractionTemplate) ReplaceCadenceImports(networkName string) (string, error) {
	var cadence = t.Data.Cadence

	// Compile regular expression to match imports
	re := regexp.MustCompile(`import\s*(\w+)\s*from\s*(0x\w+)`)
	matches := re.FindAllStringSubmatch(cadence, -1)

	for _, match := range matches {
		if len(match) != 3 {
			continue
		}
		contractName := match[1]
		dependencyAddress := match[2]

		// Check if dependency exists
		contracts, ok := t.Data.Dependencies[dependencyAddress]
		if !ok {
			return "", fmt.Errorf("network %s not found for contract %s in dependencies", networkName, contractName)
		}

		// Check if contract exists in dependency
		networks, ok := contracts[contractName]
		if !ok {
			return "", fmt.Errorf("contract %s not found in dependencies", contractName)
		}

		// Check if network exists for contract
		network, ok := networks[networkName]
		if !ok {
			return "", fmt.Errorf("network %s not found for contract %s in dependencies", networkName, contractName)
		}

		pattern := fmt.Sprintf(`import\s*%s\s*from\s*%s`, contractName, dependencyAddress)
		re, err := regexp.Compile(pattern)
		if err != nil {
			return "", fmt.Errorf("invalid regex pattern: %v", err)
		}

		replacement := fmt.Sprintf("import %s from %s", contractName, network.Address)
		cadence = re.ReplaceAllString(cadence, replacement)
	}

	return cadence, nil
}

func (t *FlowInteractionTemplate) GetDescription() string {
	s := ""
	if t.Data.Messages.Description != nil &&
		t.Data.Messages.Description.I18N != nil {

		// relying on en-US for now, future we need to know what language to use
		value, exists := t.Data.Messages.Description.I18N["en-US"]
		if exists {
			s = value
		}
	}
	return s
}

func (msgs *Messages) GetTitleValue(placeholder string) string {
	s := placeholder
	if msgs.Title != nil &&
		msgs.Title.I18N != nil {
		// relying on en-US for now, future we need to know what language to use
		value, exists := msgs.Title.I18N["en-US"]
		if exists {
			s = value
		}
	}
	return s
}
