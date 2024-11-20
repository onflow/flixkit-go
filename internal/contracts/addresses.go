package contracts

// CoreContractAddresses maps network names to their core contract addresses
type CoreContractAddresses struct {
	FungibleToken    string
	FlowToken        string
	NonFungibleToken string
	MetadataViews    string
}

var networkAddresses = map[string]CoreContractAddresses{
	"mainnet": {
		FungibleToken:    "0xf233dcee88fe0abe",
		FlowToken:        "0x1654653399040a61",
		NonFungibleToken: "0x1d7e57aa55817448",
		MetadataViews:    "0x1d7e57aa55817448",
	},
	"testnet": {
		FungibleToken:    "0x9a0766d93b6608b7",
		FlowToken:        "0x7e60df042a9c0868",
		NonFungibleToken: "0x631e88ae7f1d7c20",
		MetadataViews:    "0x631e88ae7f1d7c20",
	},
	"emulator": {
		FungibleToken:    "0xee82856bf20e2aa6",
		FlowToken:        "0x0ae53cb6e3f42a79",
		NonFungibleToken: "0xf8d6e0586b0a20c7",
		MetadataViews:    "0xf8d6e0586b0a20c7",
	},
}

// GetCoreContractAddress returns the address for a core contract on the specified network
func GetCoreContractAddress(network, contractName string) string {
	addresses, ok := networkAddresses[network]
	if !ok {
		return ""
	}

	switch contractName {
	case "FungibleToken":
		return addresses.FungibleToken
	case "FlowToken":
		return addresses.FlowToken
	case "NonFungibleToken":
		return addresses.NonFungibleToken
	case "MetadataViews":
		return addresses.MetadataViews
	default:
		return ""
	}
}
