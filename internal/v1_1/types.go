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
	ID              string `json:"id" cadence:"id"`
	Data            Data   `json:"data" cadence:"data"`
	CadenceBodyHash string `json:"cadence_body_hash" cadence:"cadenceBodyHash"`
	Status          string `json:"status" cadence:"status"`
}

type Data struct {
	Type         string       `json:"type" cadence:"type"`
	Interface    string       `json:"interface" cadence:"interface"`
	Messages     []Message    `json:"messages" cadence:"messages"`
	Cadence      Cadence      `json:"cadence" cadence:"cadence"`
	Dependencies []Dependency `json:"dependencies" cadence:"dependencies"`
	Parameters   []Parameter  `json:"parameters" cadence:"parameters"`
}

type Message struct {
	Key  string `json:"key" cadence:"key"`
	I18n []I18n `json:"i18n" cadence:"il8n"`
}

type I18n struct {
	Tag         string `json:"tag" cadence:"tag"`
	Translation string `json:"translation" cadence:"translation"`
}

type Cadence struct {
	Body        string       `json:"body" cadence:"body"`
	NetworkPins []NetworkPin `json:"network_pins" cadence:"networkPins"`
}

type NetworkPin struct {
	Network string `json:"network" cadence:"network"`
	PinSelf string `json:"pin_self" cadence:"pinSelf"`
}

type Dependency struct {
	Contracts []Contract `json:"contracts" cadence:"contracts"`
}

type Contract struct {
	Contract string    `json:"contract" cadence:"contract"`
	Networks []Network `json:"networks" cadence:"networks"`
}

type Network struct {
	Network                  string        `json:"network" cadence:"network"`
	Address                  string        `json:"address" cadence:"address"`
	DependencyPinBlockHeight uint64        `json:"dependency_pin_block_height" cadence:"dependencyPinBlockHeight"`
	DependencyPin            DependencyPin `json:"dependency_pin" cadence:"dependencyPin"`
}

// Updated from PinDetail
type DependencyPin struct {
	Pin                string   `json:"pin" cadence:"pin"`
	PinSelf            string   `json:"pin_self" cadence:"pinSelf"`
	PinContractName    string   `json:"pin_contract_name" cadence:"pinContractName"`
	PinContractAddress string   `json:"pin_contract_address" cadence:"pinContractAddress"`
	Imports            []Import `json:"imports" cadence:"imports"`
}

type Import struct {
	Pin                string   `json:"pin" cadence:"pin"`
	PinSelf            string   `json:"pin_self" cadence:"pinSelf"`
	PinContractName    string   `json:"pin_contract_name" cadence:"pinContractName"`
	PinContractAddress string   `json:"pin_contract_address" cadence:"pinContractAddress"`
	Imports            []Import `json:"imports" cadence:"imports"`
}

type Parameter struct {
	Label    string    `json:"label" cadence:"label"`
	Index    int       `json:"index" cadence:"index"`
	Type     string    `json:"type" cadence:"type"`
	Messages []Message `json:"messages" cadence:"messages"`
}

func ParseJSON(flixJSON []byte) (InteractionTemplate, error) {
	var flix InteractionTemplate
	err := json.Unmarshal(flixJSON, &flix)
	if err != nil {
		return InteractionTemplate{}, nil
	}

	return flix, nil
}

func ParseCadence(flixCadence cadence.Value) (InteractionTemplate, error) {
	return InteractionTemplate{}, nil
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
