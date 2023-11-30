package v1_1

import (
	"encoding/json"
	"fmt"
	"regexp"
)

type InteractionTemplate struct {
	FType    string `json:"f_type"`
	FVersion string `json:"f_version"`
	ID       string `json:"id"`
	Data     Data   `json:"data"`
}

type Data struct {
	Type         string       `json:"type"`
	Interface    string       `json:"interface"`
	Messages     []Message    `json:"messages"`
	Cadence      Cadence      `json:"cadence"`
	Dependencies []Dependency `json:"dependencies"`
	Parameters   []Parameter  `json:"parameters"`
}

type Message struct {
	Key  string `json:"key"`
	I18n []I18n `json:"i18n"`
}

type I18n struct {
	Tag         string `json:"tag"`
	Translation string `json:"translation"`
}

type Cadence struct {
	Body        string       `json:"body"`
	NetworkPins []NetworkPin `json:"network_pins"`
}

type NetworkPin struct {
	Network string `json:"network"`
	PinSelf string `json:"pin_self"`
}

type Dependency struct {
	Contracts []Contract `json:"contracts"`
}

type Contract struct {
	Contract string    `json:"contract"`
	Networks []Network `json:"networks"`
}

type Network struct {
	Network                  string    `json:"network"`
	Address                  string    `json:"address"`
	DependencyPinBlockHeight int64     `json:"dependency_pin_block_height"`
	DependencyPin            PinDetail `json:"dependency_pin"`
}

type PinDetail struct {
	Pin                string   `json:"pin"`
	PinSelf            string   `json:"pin_self"`
	PinContractName    string   `json:"pin_contract_name"`
	PinContractAddress string   `json:"pin_contract_address"`
	Imports            []Import `json:"imports"`
}

type Import struct {
	Pin                string   `json:"pin"`
	PinSelf            string   `json:"pin_self"`
	PinContractName    string   `json:"pin_contract_name"`
	PinContractAddress string   `json:"pin_contract_address"`
	Imports            []Import `json:"imports"` // Recursive imports, if any
}

type Parameter struct {
	Label    string    `json:"label"`
	Index    int       `json:"index"`
	Type     string    `json:"type"`
	Messages []Message `json:"messages"`
	Balance  string    `json:"balance"`
}

func (t *InteractionTemplate) IsScript() bool {
	return t.Data.Type == "script"
}

func (t *InteractionTemplate) IsTransaction() bool {
	return t.Data.Type == "transaction"
}

func (t *InteractionTemplate) GetAndReplaceCadenceImports(networkName string) (string, error) {
	var cadence string

	for _, Dependence := range t.Data.Dependencies {
		for _, contract := range Dependence.Contracts {
			contractName := contract.Contract
			var dependencyAddress string
			for _, network := range contract.Networks {
				if network.Network == networkName {
					dependencyAddress = network.Address
					break
				}
			}
			if dependencyAddress == "" {
				return "", fmt.Errorf("network %s not found for contract %s in dependencies", networkName, contractName)
			}

			re, err := regexp.Compile(`import "(.+?)"`)
			if err != nil {
				return "", fmt.Errorf("invalid regex pattern: %v", err)
			}

			replacement := fmt.Sprintf("import %s from %s", contractName, dependencyAddress)
			cadence = re.ReplaceAllString(t.Data.Cadence.Body, replacement)
		}
	}

	return cadence, nil

}

func ParseFlix(template string) (*InteractionTemplate, error) {
	var flowTemplate InteractionTemplate

	err := json.Unmarshal([]byte(template), &flowTemplate)
	if err != nil {
		return nil, err
	}

	return &flowTemplate, nil
}
