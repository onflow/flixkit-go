package v1_1

import (
	"encoding/json"

	"github.com/bjartek/underflow"
	"github.com/onflow/cadence"
)

// TODO: Add Proper Values
const (
	DefaultMainnetRegistryId  = "A.f8d6e0586b0a20c7.FLIXSchema_v1_1_0"
	DefaultTestnetRegistryId  = "A.f8d6e0586b0a20c7.FLIXSchema_v1_1_0"
	DefaultEmulatorRegistryId = "A.f8d6e0586b0a20c7.FLIXSchema_v1_1_0"
)

type InteractionTemplate struct {
	FType    string `json:"f_type"`
	FVersion string `json:"f_version"`
	ID       string `json:"id"`
	Data     Data   `json:"data"`
}

type FLIX struct {
	ID              string `json:"id"`
	Data            Data   `json:"data"`
	CadenceBodyHash string `json:"cadence_body_hash"`
	Status          string `json:"status"`
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
	Network                  string        `json:"network"`
	Address                  string        `json:"address"`
	DependencyPinBlockHeight uint64        `json:"dependency_pin_block_height"`
	DependencyPin            DependencyPin `json:"dependency_pin"`
}

// Updated from PinDetail
type DependencyPin struct {
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
}

func ParseJSON(flixJSON []byte) (InteractionTemplate, error) {
	var flix InteractionTemplate
	err := json.Unmarshal(flixJSON, &flix)
	if err != nil {
		return InteractionTemplate{}, nil
	}

	return flix, nil
}

// TODO: Fix Cadence Body Hash
// TODO: Cleaner custom network passing
func (f InteractionTemplate) AsCadance(status string, network string) (cadence.Value, error) {
	var registryId string
	switch network {
	case "mainnet":
		registryId = DefaultMainnetRegistryId
	case "testnet":
		registryId = DefaultTestnetRegistryId
	case "emulator":
		registryId = DefaultEmulatorRegistryId
	default:
		registryId = network
	}

	resolver := func(s string) (string, error) {
		return registryId + "." + s, nil
	}

	flix := FLIX{
		ID:              f.ID,
		Data:            f.Data,
		CadenceBodyHash: f.Data.Cadence.Body,
		Status:          status,
	}

	return underflow.InputToCadence(flix, resolver)
}

func (f InteractionTemplate) AsJSON() ([]byte, error) {
	return json.Marshal(f)
}

func (f InteractionTemplate) ReplaceImports() {

}

func (f InteractionTemplate) CreateBindings() (string, error) {
	return "", nil
}
