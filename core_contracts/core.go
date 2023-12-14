package core_contracts

import (
	"github.com/onflow/flow-cli/flowkit/config"
)

func GetCoreContracts() map[string]map[string]string {
	// TODO: this should be from core contracts module
	c := make(map[string]map[string]string)
	c["FungibleToken"] = make(map[string]string)
	c["FungibleToken"][config.MainnetNetwork.Name] = "0xf233dcee88fe0abe"
	c["FungibleToken"][config.TestnetNetwork.Name] = "0x9a0766d93b6608b7"
	c["FungibleToken"][config.EmulatorNetwork.Name] = "0xee82856bf20e2aa6"
	c["NonFungibleToken"] = make(map[string]string)
	c["NonFungibleToken"][config.MainnetNetwork.Name] = "0x1d7e57aa55817448"
	c["NonFungibleToken"][config.TestnetNetwork.Name] = "0x631e88ae7f1d7c20"
	c["NonFungibleToken"][config.EmulatorNetwork.Name] = "0xf8d6e0586b0a20c7"
	c["MetadataViews"] = make(map[string]string)
	c["MetadataViews"][config.MainnetNetwork.Name] = "0x1d7e57aa55817448"
	c["MetadataViews"][config.TestnetNetwork.Name] = "0x631e88ae7f1d7c20"
	c["MetadataViews"][config.EmulatorNetwork.Name] = "0xf8d6e0586b0a20c7"
	c["FlowToken"] = make(map[string]string)
	c["FlowToken"][config.MainnetNetwork.Name] = "0x1654653399040a61"
	c["FlowToken"][config.TestnetNetwork.Name] = "0x7e60df042a9c0868"
	c["FlowToken"][config.EmulatorNetwork.Name] = "0x0ae53cb6e3f42a79"
	c["HybridCustody"] = make(map[string]string)
	c["HybridCustody"][config.MainnetNetwork.Name] = "0xd8a7e05a7ac670c0"
	c["HybridCustody"][config.TestnetNetwork.Name] = "0x294e44e1ec6993c6"
	c["HybridCustody"][config.EmulatorNetwork.Name] = "0xf8d6e0586b0a20c7"
	c["CapabilityDelegator"] = make(map[string]string)
	c["CapabilityDelegator"][config.MainnetNetwork.Name] = "0xd8a7e05a7ac670c0"
	c["CapabilityDelegator"][config.TestnetNetwork.Name] = "0x294e44e1ec6993c6"
	c["CapabilityDelegator"][config.EmulatorNetwork.Name] = "0xf8d6e0586b0a20c7"
	c["CapabilityFactory"] = make(map[string]string)
	c["CapabilityFactory"][config.MainnetNetwork.Name] = "0xd8a7e05a7ac670c0"
	c["CapabilityFactory"][config.TestnetNetwork.Name] = "0x294e44e1ec6993c6"
	c["CapabilityFactory"][config.EmulatorNetwork.Name] = "0xf8d6e0586b0a20c7"
	c["CapabilityFilter"] = make(map[string]string)
	c["CapabilityFilter"][config.MainnetNetwork.Name] = "0xd8a7e05a7ac670c0"
	c["CapabilityFilter"][config.TestnetNetwork.Name] = "0x294e44e1ec6993c6"
	c["CapabilityFilter"][config.EmulatorNetwork.Name] = "0xf8d6e0586b0a20c7"

	return c
}

func GetCoreContractForNetwork(contractName string, network string) string {
	// TODO: this should be from core contracts module
	c := GetCoreContracts()
	return c[contractName][network]
}
