package core_contracts

func GetCoreContracts() map[string]map[string]string {
	// TODO: this should be from core contracts module
	c := make(map[string]map[string]string)
	c["FungibleToken"] = make(map[string]string)
	c["FungibleToken"]["mainnet"] = "0xf233dcee88fe0abe"
	c["FungibleToken"]["testnet"] = "0x9a0766d93b6608b7"
	c["FungibleToken"]["emulator"] = "0xee82856bf20e2aa6"
	c["NonFungibleToken"] = make(map[string]string)
	c["NonFungibleToken"]["mainnet"] = "0x1d7e57aa55817448"
	c["NonFungibleToken"]["testnet"] = "0x631e88ae7f1d7c20"
	c["NonFungibleToken"]["emulator"] = "0xf8d6e0586b0a20c7"
	c["MetadataViews"] = make(map[string]string)
	c["MetadataViews"]["mainnet"] = "0x1d7e57aa55817448"
	c["MetadataViews"]["testnet"] = "0x631e88ae7f1d7c20"
	c["FlowToken"] = make(map[string]string)
	c["FlowToken"]["mainnet"] = "0x1654653399040a61"
	c["FlowToken"]["testnet"] = "0x7e60df042a9c0868"
	c["FlowToken"]["emulator"] = "0x0ae53cb6e3f42a79"
	return c
}

func GetCoreContractForNetwork(contractName string, network string) string {
	// TODO: this should be from core contracts module
	c := GetCoreContracts()
	return c[contractName][network]
}
