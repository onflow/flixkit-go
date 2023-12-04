package v1_1

import (
	"encoding/json"
	"fmt"
	"regexp"

	"github.com/onflow/cadence/runtime/ast"
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

	// Compile regular expression to match and capture contract names
	re := regexp.MustCompile(`import\s*"([^"]+)"`)

	// Find all matches and their captured groups
	matches := re.FindAllStringSubmatch(t.Data.Cadence.Body, -1)

	if len(matches) == 0 {
		return t.Data.Cadence.Body, nil
	}
	for _, match := range matches {
		contractName := match[1]
		var dependencyAddress string
		for _, Dependence := range t.Data.Dependencies {
			for _, contract := range Dependence.Contracts {
				if contract.Contract == contractName {
					for _, network := range contract.Networks {
						if network.Network == networkName {
							dependencyAddress = network.Address
							break
						}
					}
					break
				}
			}
		}

		if dependencyAddress == "" {
			return "", fmt.Errorf("network %s not found for contract %s in dependencies", networkName, contractName)
		}

		replacement := fmt.Sprintf("import %s from %s", contractName, dependencyAddress)
		cadence = re.ReplaceAllString(t.Data.Cadence.Body, replacement)
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

type PragmaDeclaration struct {
	Expression InteractionExpression `json:"Expression"`
}

type InteractionExpression struct {
	InvokedExpression IdentifierExpression    `json:"InvokedExpression"`
	Arguments         []Argument              `json:"Arguments"`
	Value             string                  `json:"Value"`  // Used for string expressions
	Type              string                  `json:"Type"`   // Used for string expressions
	Values            []InteractionExpression `json:"Values"` // Used for array expressions
}

type IdentifierExpression struct {
	Identifier Identifier `json:"Identifier"`
}

type Identifier struct {
	Identifier string `json:"Identifier"`
}

type Argument struct {
	Expression InteractionExpression `json:"Expression"`
	Label      string                `json:"Label"`
}

func ParsePragma(pragmas []*ast.PragmaDeclaration, template *InteractionTemplate) error {
	if len(pragmas) == 0 {
		return nil
	}

	for _, prag := range pragmas {
		var pragmaDeclaration PragmaDeclaration
		jsonData, err := prag.MarshalJSON()
		if err != nil {
			return err
		}
		err = json.Unmarshal([]byte(jsonData), &pragmaDeclaration)
		if err != nil {
			return err
		}
		if pragmaDeclaration.Expression.InvokedExpression.Identifier.Identifier == "interaction" {
			pragmaInfo := flatten(pragmaDeclaration)
			if template.FVersion == "" {
				template.FVersion = pragmaInfo.meta["version"]
			}
			if pragmaInfo.meta["title"] != "" {
				template.Data.Messages = append(template.Data.Messages, Message{
					Key: "title",
					I18n: []I18n{
						{
							Tag:         pragmaInfo.meta["language"],
							Translation: pragmaInfo.meta["title"],
						},
					},
				})
			}
			if pragmaInfo.meta["description"] != "" {
				template.Data.Messages = append(template.Data.Messages, Message{
					Key: "description",
					I18n: []I18n{
						{
							Tag:         pragmaInfo.meta["language"],
							Translation: pragmaInfo.meta["description"],
						},
					},
				})
			}
			if pragmaInfo.parameters != nil {
				for i, paramInfo := range pragmaInfo.parameters {
					param := Parameter{
						Label: paramInfo.params["name"],
						Index: i,
					}
					if paramInfo.params["title"] != "" {
						param.Messages = append(param.Messages, Message{
							Key: "title",
							I18n: []I18n{
								{
									Tag:         pragmaInfo.meta["language"],
									Translation: paramInfo.params["title"],
								},
							},
						})
					}
					if paramInfo.params["description"] != "" {
						param.Messages = append(param.Messages, Message{
							Key: "description",
							I18n: []I18n{
								{
									Tag:         pragmaInfo.meta["language"],
									Translation: paramInfo.params["description"],
								},
							},
						})
					}
					template.Data.Parameters = append(template.Data.Parameters, param)
				}
			}

		}

	}

	return nil
}

type parametermetadata struct {
	params map[string]string
}
type metadata struct {
	meta       map[string]string
	parameters []parametermetadata
}

func flatten(pragma PragmaDeclaration) metadata {
	var nameValuePairs map[string]string
	var parameterPairs []parametermetadata
	nameValuePairs = make(map[string]string)
	parameterPairs = make([]parametermetadata, 0)

	for _, arg := range pragma.Expression.Arguments {
		// For regular arguments
		if arg.Expression.Value != "" {
			nameValuePairs[arg.Label] = arg.Expression.Value
		}

		// For arguments that contain arrays (like parameters)
		if len(arg.Expression.Values) > 0 {
			for _, param := range arg.Expression.Values {
				p := parametermetadata{
					params: make(map[string]string),
				}
				for _, paramArg := range param.Arguments {
					p.params[paramArg.Label] = paramArg.Expression.Value
				}
				parameterPairs = append(parameterPairs, p)
			}
		}
	}
	return metadata{nameValuePairs, parameterPairs}
}
