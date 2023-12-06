package v1_1

import (
	"encoding/binary"
	"encoding/hex"
	"encoding/json"
	"fmt"
	"regexp"
	"strings"

	"github.com/onflow/cadence/runtime/ast"
	"golang.org/x/crypto/sha3"
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
	DependencyPinBlockHeight uint64    `json:"dependency_pin_block_height"`
	DependencyPin            PinDetail `json:"dependency_pin"`
}

type PinDetail struct {
	Pin                string      `json:"pin"`
	PinSelf            string      `json:"pin_self"`
	PinContractName    string      `json:"pin_contract_name"`
	PinContractAddress string      `json:"pin_contract_address"`
	Imports            []PinDetail `json:"imports"`
}

func (p *PinDetail) CalculatePin(blockHeight uint64) {
	var a []string
	a = append(a, ShaHex(p.PinContractAddress, "address"))
	a = append(a, ShaHex(p.PinContractName, "address"))
	a = append(a, ShaHex(p.PinSelf, "pin_self"))
	a = append(a, ShaHex(fmt.Sprint(blockHeight), "pin_block_height"))
	hash := ShaHex(strings.Join(a, ""), "calculate_pin")
	p.Pin = hash
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
}

func (t *InteractionTemplate) Init() {
	t.FType = "InteractionTemplate"
	if t.FVersion == "" {
		t.FVersion = "1.1.0"
	}
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

func (template *InteractionTemplate) ParsePragma(program *ast.Program) error {
	pragmas := program.PragmaDeclarations()
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

func (template *InteractionTemplate) ProcessParameters(program *ast.Program) error {
	if program == nil {
		return fmt.Errorf("no cadence program provided")
	}
	var parameterList []*ast.Parameter
	functionDeclaration := program.FunctionDeclarations()
	// only interested in main function of script
	for _, d := range functionDeclaration {
		if d.Identifier.String() == "main" {
			parameterList = d.ParameterList.Parameters
		}
	}

	if program.SoleTransactionDeclaration() != nil && program.SoleTransactionDeclaration().ParameterList != nil {
		parameterList = program.SoleTransactionDeclaration().ParameterList.Parameters
	}

	if parameterList != nil && len(template.Data.Parameters) == 0 {
		template.Data.Parameters = make([]Parameter, 0)
	}

	// use existing parameter of template or create new one
	for i, param := range parameterList {
		var tempParam *Parameter
		if hasValueAtIndex(template.Data.Parameters, i) {
			tempParam = &template.Data.Parameters[i]
			// verify that the parameter name matches,
			// could happen if dev inputted param data incorrectly
			if tempParam.Label != param.Identifier.String() {
				return fmt.Errorf("parameter name mismatch, expected %s, got %s", tempParam.Label, param.Identifier.String())
			}
		} else {
			template.Data.Parameters = append(template.Data.Parameters, *tempParam)
		}
		tempParam.Label = param.Identifier.String()
		tempParam.Type = param.TypeAnnotation.Type.String()
		tempParam.Index = i
	}

	return nil
}

func hasValueAtIndex(arr []Parameter, index int) bool {
	if len(arr) == 0 {
		return false
	}
	if index >= 0 && index < len(arr) {
		return true
	}
	return false
}

func (template *InteractionTemplate) DetermineCadenceType(program *ast.Program) error {
	funcs := program.FunctionDeclarations()
	trans := program.TransactionDeclarations()
	var t string
	if len(funcs) > 0 {
		t = "script"
	} else if len(trans) > 0 {
		t = "transaction"
	} else {
		return fmt.Errorf("no function or transaction declarations found")
	}
	template.Data.Type = t
	return nil
}

func (template *InteractionTemplate) ProcessImports(cadenceCode string) {
	// Define a regular expression to match the "import ContractName from 0xContractName" pattern
	pattern := regexp.MustCompile(`import\s+(\w+)\s+from\s+0x\w+`)
	// Replace the matched pattern with "import \"ContractName\""
	replaced := pattern.ReplaceAllString(cadenceCode, `import "$1"`)
	template.Data.Cadence.Body = replaced
}

func ShaHex(value interface{}, debugKey string) string {

	// Convert the value to a byte array
	data, err := convertToBytes(value)
	if err != nil {
		if debugKey != "" {
			fmt.Printf("%30s value=%v hex=%x\n", debugKey, value, err.Error())
		}
		return ""
	}

	// Calculate the SHA-3 hash
	hash := sha3.Sum256(data)

	// Convert the hash to a hexadecimal string
	hashHex := hex.EncodeToString(hash[:])

	return hashHex
}

func convertToBytes(value interface{}) ([]byte, error) {
	switch v := value.(type) {
	case []byte:
		return v, nil
	case string:
		return []byte(v), nil
	case int:
		buf := make([]byte, 4)
		binary.BigEndian.PutUint32(buf, uint32(v))
		return buf, nil
	case uint64:
		buf := make([]byte, 8)
		binary.BigEndian.PutUint64(buf, v)
		return buf, nil
	default:
		return nil, fmt.Errorf("unsupported type %T", v)
	}
}
