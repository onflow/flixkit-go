package flixkit

import (
	"bytes"
	"encoding/binary"
	"encoding/hex"
	"encoding/json"
	"fmt"
	"regexp"
	"sort"

	"github.com/ethereum/go-ethereum/rlp"
	"github.com/onflow/flixkit-go/core_contracts"
	"golang.org/x/crypto/sha3"
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

func (t *FlowInteractionTemplate) GetAndReplaceCadenceImports(networkName string) (string, error) {
	var cadence string

	for dependencyAddress, contracts := range t.Data.Dependencies {
		for contractName, networks := range contracts {
			var networkAddress string
			network, ok := networks[networkName]
			networkAddress = network.Address
			if !ok {
				coreContractAddress := core_contracts.GetCoreContractForNetwork(contractName, networkName)
				if coreContractAddress == "" {
					return "", fmt.Errorf("network %s not found for contract %s in dependencies", networkName, contractName)
				}
				networkAddress = coreContractAddress
			}

			pattern := fmt.Sprintf(`import\s*%s\s*from\s*%s`, contractName, dependencyAddress)
			re, err := regexp.Compile(pattern)
			if err != nil {
				return "", fmt.Errorf("invalid regex pattern: %v", err)
			}

			replacement := fmt.Sprintf("import %s from %s", contractName, networkAddress)
			cadence = re.ReplaceAllString(t.Data.Cadence, replacement)
		}
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

func messagesToRlp(messages Messages) []interface{} {
	values := make([]interface{}, 0)

	if messages.Title != nil {
		var titleValue []interface{}
		titleValue = append(titleValue, ShaHex("title", "title"))
		if messages.Title.I18N != nil {
			var langTitle []interface{}
			for k, v := range messages.Title.I18N {
				var anotherNesting []interface{}
				anotherNesting = append(anotherNesting, ShaHex(k, "I18N"))
				anotherNesting = append(anotherNesting, ShaHex(v, "I18N"))
				langTitle = append(langTitle, anotherNesting)
			}
			titleValue = append(titleValue, langTitle)
		}
		values = append(values, titleValue)
	}

	if messages.Description != nil {
		var descValue []interface{}
		descValue = append(descValue, ShaHex("description", "description"))
		if messages.Description.I18N != nil {
			var langDesc []interface{}
			for k, v := range messages.Description.I18N {
				var anotherNesting []interface{}
				anotherNesting = append(anotherNesting, ShaHex(k, "I18N"))
				anotherNesting = append(anotherNesting, ShaHex(v, "I18N"))
				langDesc = append(langDesc, anotherNesting)
			}
			descValue = append(descValue, langDesc)
		}
		values = append(values, descValue)
	}

	return values
}

func argumentsToRlp(arguments Arguments) []interface{} {
	values := make([]interface{}, 0)
	sortedArguments := arguments.SortArguments()
	for _, argument := range sortedArguments {

		var args []interface{}
		args = append(args, ShaHex(argument.Key, "key"))

		var arg []interface{}
		arg = append(arg, ShaHex(fmt.Sprint(argument.Index), "index"))
		arg = append(arg, ShaHex(argument.Type, "type"))
		arg = append(arg, ShaHex(argument.Balance, "balance"))
		arg = append(arg, messagesToRlp(argument.Messages))
		args = append(args, arg)
		values = append(values, args)
	}
	return values
}

func dependenciesToRlp(Dependencies Dependencies) []interface{} {
	values := make([]interface{}, 0)
	keys := SortMapKeys(Dependencies)
	for _, key := range keys {
		var deps []interface{}
		value := Dependencies[key]
		deps = append(deps, ShaHex(key, "key"))
		deps = append(deps, contractsToRlp(value))
		values = append(values, deps)
	}
	return values
}

func contractsToRlp(Contracts Contracts) []interface{} {
	values := make([]interface{}, 0)
	keys := SortMapKeys(Contracts)
	for _, key := range keys {
		value := Contracts[key]
		var contracts []interface{}
		contracts = append(contracts, ShaHex(key, "key"))
		contracts = append(contracts, networksToRlp(value))
		values = append(values, contracts)
	}
	return values
}

func networksToRlp(Networks Networks) []interface{} {
	values := make([]interface{}, 0)
	keys := SortMapKeys(Networks)
	for _, key := range keys {
		value := Networks[key]
		var networks []interface{}
		networks = append(networks, ShaHex(key, "key"))
		networks = append(networks, networkToRlp(value))
		values = append(values, networks)
	}
	return values
}

func networkToRlp(network Network) []interface{} {
	values := make([]interface{}, 0)
	values = append(values, ShaHex(network.Address, "address"))
	values = append(values, ShaHex(network.Contract, "contract"))
	values = append(values, ShaHex(network.FqAddress, "fq_address"))
	values = append(values, ShaHex(network.Pin, "pin"))
	values = append(values, ShaHex(fmt.Sprint(network.PinBlockHeight), "pin_block_height"))
	return values
}

func (flix FlowInteractionTemplate) EncodeRLP() (result string, err error) {
	var buffer bytes.Buffer // Create a new buffer

	input := []interface{}{
		ShaHex(flix.FType, "f-type"),
		ShaHex(flix.FVersion, "f-version"),
		ShaHex(flix.Data.Type, "type"),
		ShaHex(flix.Data.Interface, "interface"),
		messagesToRlp(flix.Data.Messages),
		ShaHex(flix.Data.Cadence, "cadence"),
		dependenciesToRlp(flix.Data.Dependencies),
		argumentsToRlp(flix.Data.Arguments),
	}

	//	msg := dependenciesToRlp(flix.Data.Dependencies)
	//prettyJSON, _ := json.MarshalIndent(input, "", "    ")
	//fmt.Println(string(prettyJSON))

	err = rlp.Encode(&buffer, input)
	if err != nil {
		return "", err
	}
	hexString := hex.EncodeToString(buffer.Bytes())
	fullyHashed := ShaHex(hexString, "input")

	return fullyHashed, nil

}

func GenerateFlixID(flix *FlowInteractionTemplate) (string, error) {
	rlpOutput, err := flix.EncodeRLP()
	if err != nil {
		return "", err
	}
	return string(rlpOutput), nil
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

type ArgumentKey struct {
	Key string
	Argument
}

type ArgumentKeys []ArgumentKey

func (args Arguments) SortArguments() ArgumentKeys {
	keys := make(ArgumentKeys, 0, len(args))
	for key, argument := range args {
		keys = append(keys, ArgumentKey{Key: key, Argument: argument})
	}

	// Use sort.Slice instead of sort.Sort
	sort.Slice(keys, func(i, j int) bool {
		return keys[i].Index < keys[j].Index
	})

	return keys
}

type MapKeySorter[T any] func(map[string]T) []string

func SortMapKeys[T any](m map[string]T) []string {
	keys := make([]string, 0, len(m))
	for k := range m {
		keys = append(keys, k)
	}
	sort.Strings(keys)
	return keys
}
