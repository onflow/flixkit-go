package contracts

import (
	"github.com/onflow/flowkit/v2/config"
)

func GetCoreContracts() map[string]map[string]string {
	// TODO: this should be from core contracts module
	return map[string]map[string]string{
		"FungibleToken": {
			config.MainnetNetwork.Name:  "0xf233dcee88fe0abe",
			config.TestnetNetwork.Name:  "0x9a0766d93b6608b7",
			config.EmulatorNetwork.Name: "0xee82856bf20e2aa6",
		},
		"NonFungibleToken": {
			config.MainnetNetwork.Name:  "0x1d7e57aa55817448",
			config.TestnetNetwork.Name:  "0x631e88ae7f1d7c20",
			config.EmulatorNetwork.Name: "0xf8d6e0586b0a20c7",
		},
		"MetadataViews": {
			config.MainnetNetwork.Name:  "0x1d7e57aa55817448",
			config.TestnetNetwork.Name:  "0x631e88ae7f1d7c20",
			config.EmulatorNetwork.Name: "0xf8d6e0586b0a20c7",
		},
		"FlowToken": {
			config.MainnetNetwork.Name:  "0x1654653399040a61",
			config.TestnetNetwork.Name:  "0x7e60df042a9c0868",
			config.EmulatorNetwork.Name: "0x0ae53cb6e3f42a79",
		},
		"HybridCustody": {
			config.MainnetNetwork.Name:  "0xd8a7e05a7ac670c0",
			config.TestnetNetwork.Name:  "0x294e44e1ec6993c6",
			config.EmulatorNetwork.Name: "0xf8d6e0586b0a20c7",
		},
		"CapabilityDelegator": {
			config.MainnetNetwork.Name:  "0xd8a7e05a7ac670c0",
			config.TestnetNetwork.Name:  "0x294e44e1ec6993c6",
			config.EmulatorNetwork.Name: "0xf8d6e0586b0a20c7",
		},
		"CapabilityFactory": {
			config.MainnetNetwork.Name:  "0xd8a7e05a7ac670c0",
			config.TestnetNetwork.Name:  "0x294e44e1ec6993c6",
			config.EmulatorNetwork.Name: "0xf8d6e0586b0a20c7",
		},
		"CapabilityFilter": {
			config.MainnetNetwork.Name:  "0xd8a7e05a7ac670c0",
			config.TestnetNetwork.Name:  "0x294e44e1ec6993c6",
			config.EmulatorNetwork.Name: "0xf8d6e0586b0a20c7",
		},
	}
}

func GetCoreContractForNetwork(contractName string, network string) string {
	// TODO: this should be from core contracts module
	c := GetCoreContracts()
	return c[contractName][network]
}
