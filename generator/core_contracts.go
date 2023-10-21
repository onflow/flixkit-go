package generator

import (
	"github.com/onflow/flixkit-go"
)

// GetContractInformation returns the contract information for a given contract name
// todo: this information should be generated so it can be updated easily
// todo: this should be moved to a separate package or maybe this inforation already exists somewhere else
func GetContractInformation(contractName string) flixkit.Networks {
	var contracts = map[string]flixkit.Networks{
		"FungibleToken": {
			"mainnet": {
				Address:        "0xf233dcee88fe0abe",
				FqAddress:      "A.0xf233dcee88fe0abe.FungibleToken",
				Contract:       "FungibleToken",
				Pin:            "83c9e3d61d3b5ebf24356a9f17b5b57b12d6d56547abc73e05f820a0ae7d9cf5",
				PinBlockHeight: 34166296,
			},
			"testnet": {
				Address:        "0x9a0766d93b6608b7",
				FqAddress:      "A.0x9a0766d93b6608b7.FungibleToken",
				Contract:       "FungibleToken",
				Pin:            "83c9e3d61d3b5ebf24356a9f17b5b57b12d6d56547abc73e05f820a0ae7d9cf5",
				PinBlockHeight: 74776482,
			},
		},
		"NonFungibleToken": {
			"mainnet": {
				Address:        "0x1d7e57aa55817448",
				FqAddress:      "A.0x1d7e57aa55817448.NonFungibleToken",
				Contract:       "NonFungibleToken",
				Pin:            "83c9e3d61d3b5ebf24356a9f17b5b57b12d6d56547abc73e05f820a0ae7d9cf5",
				PinBlockHeight: 47509012,
			},
			"testnet": {
				Address:        "0x631e88ae7f1d7c20",
				FqAddress:      "A.0x631e88ae7f1d7c20.NonFungibleToken",
				Contract:       "NonFungibleToken",
				Pin:            "83c9e3d61d3b5ebf24356a9f17b5b57b12d6d56547abc73e05f820a0ae7d9cf5",
				PinBlockHeight: 95808228,
			},
		},
		"MetadataViews": {
			"mainnet": {
				Address:        "0x1d7e57aa55817448",
				FqAddress:      "A.0x1d7e57aa55817448.NonFungibleToken",
				Contract:       "MetadataViews",
				Pin:            "ba061d95016d5506e9f5d1afda15d82eb066aa8b0552e8b26dc7950fa5714d51",
				PinBlockHeight: 47487348,
			},
			"testnet": {
				Address:        "0x631e88ae7f1d7c20",
				FqAddress:      "A.0x631e88ae7f1d7c20.NonFungibleToken",
				Contract:       "MetadataViews",
				Pin:            "ba061d95016d5506e9f5d1afda15d82eb066aa8b0552e8b26dc7950fa5714d51",
				PinBlockHeight: 95782517,
			},
		},
		"FlowToken": {
			"mainnet": {
				Address:        "0x1654653399040a61",
				FqAddress:      "A.0x1654653399040a61.FlowToken",
				Contract:       "FlowToken",
				Pin:            "0326c320322c4e8dde768ba2975c384184fb7e41765c2c87e79a2040bfc71be8",
				PinBlockHeight: 47509023,
			},
			"testnet": {
				Address:        "0x7e60df042a9c0868",
				FqAddress:      "A.0x7e60df042a9c0868.FlowToken",
				Contract:       "FlowToken",
				Pin:            "0326c320322c4e8dde768ba2975c384184fb7e41765c2c87e79a2040bfc71be8",
				PinBlockHeight: 95808240,
			},
		},
	}

	// To lookup a specific entry:
	_, exists := contracts[contractName]
	if exists {
		return contracts[contractName]
	}

	return nil

}
