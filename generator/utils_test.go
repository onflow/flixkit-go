package generator

import (
	"context"
	"encoding/json"
	"strings"
	"testing"

	"github.com/hexops/autogold/v2"
	"github.com/onflow/flixkit-go"
)

func TestParseImport(t *testing.T) {
	fungi := flixkit.Contracts{
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
	}
	tests := []struct {
		cadence string
		want    flixkit.Contracts
	}{
		{
			cadence: `import FungibleToken from 0xFungibleTokenAddress`,
			want:    fungi,
		},
		{
			cadence: `import "FungibleToken"`,
			want:    fungi,
		},
		{
			cadence: `import FungibleToken from 0x9a0766d93b6608b7`,
			want:    fungi,
		},
		{
			cadence: `import "FungibleToken"`,
			want:    fungi,
		},
	}

	for _, tt := range tests {
		t.Run(tt.cadence, func(t *testing.T) {
			got, err := parseImport(context.Background(), tt.cadence, nil)
			if err != nil {
				t.Errorf("parseImport() err %v", err)
			}
			if got == nil {
				t.Errorf("parseImport() got = %v, want %v", got, tt.want)
			}
			prettyJSON, err := json.MarshalIndent(got, "", "    ")
			if err != nil {
				t.Errorf("parseImport() err %v", err)
			}
			autogold.ExpectFile(t, string(prettyJSON))
		})
	}
}
func TestDepCheck(t *testing.T) {
	tests := []struct {
		cadence     string
		cadenceType string
	}{
		{
			cadence: `import FungibleToken from 0xFungibleTokenAddress

			pub fun main(accountAddress: Address): UFix64 {
				let account = getAccount(accountAddress)
				let balanceRef = account.borrow<&FungibleToken.Vault{FungibleToken.Balance}>(from: /public/fungibleTokenBalance)
					?? panic("Could not borrow reference to the account's Vault")
			
				return balanceRef.balance
			}
			`,
			cadenceType: "script",
		},
		{
			cadence: `import FungibleToken from 0xFungibleTokenAddress

			transaction(amount: UFix64, recipient: Address) {
				
				// Reference to the sender's Vault
				let tokenVault: @FungibleToken.Vault
			
				prepare(signer: AuthAccount) {
					self.tokenVault = signer.borrow<&FungibleToken.Vault>(from: /storage/fungibleTokenVault)
						?? panic("Could not borrow reference to the owner's Vault")
				}
			
				execute {
					let recipient = getAccount(recipient)
					let recipientVault = recipient.getCapability(/public/fungibleTokenReceiver)
						.borrow<&{FungibleToken.Receiver}>()
						?? panic("Could not borrow receiver reference from the recipient")
			
					recipientVault.deposit(from: <-self.tokenVault.withdraw(amount: amount))
				}
			}
			`,
			cadenceType: "transaction",
		}, {
			cadence: `
/*
Here are some comments and transaction is on start of a line
*/
transaction(amount: UFix64, recipient: Address) {
// More comments
// Reference to the sender's Vault
prepare(signer: AuthAccount) {}
execute {}
			`,
			cadenceType: "transaction",
		},
		{
			cadence: `import NonFungibleToken from 0xNonFungibleTokenAddress

			pub fun main(accountAddress: Address, tokenId: UInt64): Bool {
				let account = getAccount(accountAddress)
				let collectionRef = account.borrow<&NonFungibleToken.Collection{NonFungibleToken.CollectionPublic}>(from: /public/nftCollection)
					?? panic("Could not borrow reference to the NFT Collection")
			
				return collectionRef.borrowNFT(id: tokenId) != nil
			}
			`,
			cadenceType: "script",
		},
		{
			cadence: `pub contract interface TokenContract {

				// Returns the total supply of tokens
				pub fun totalSupply(): UFix64
			
				// Returns the balance of the specified address
				pub fun balanceOf(address: Address): UFix64
			
				// Transfers tokens from one address to another
				pub fun transfer(from: Address, to: Address, amount: UFix64): Bool
			}
			`,
			cadenceType: "interface",
		},
	}

	for _, tt := range tests {
		t.Run(tt.cadence, func(t *testing.T) {
			got, err := determineCadenceType(tt.cadence)
			if err != nil {
				t.Errorf("determineCadenceType() err %v", err)
			}
			if got != tt.cadenceType {
				t.Errorf("determineCadenceType() got = %v, want %v", got, tt.cadenceType)
			}
		})
	}

}

func TestStripComments(t *testing.T) {
	removedComment := "Here is a comment"
	tests := []struct {
		cadence string
	}{
		{
			cadence: `import FungibleToken from 0xFungibleTokenAddress
			/* Here is a comment */
			pub fun main(accountAddress: Address): UFix64 {
				return balanceRef.balance
			}
			`,
		},
		{
			cadence: `import FungibleToken from 0xFungibleTokenAddress

			transaction(amount: UFix64, recipient: Address) {
				/* 
				Here is a comment
				*/
				execute {

				}
			}
			`,
		}, {
			cadence: `
/*
Here are some comments and transaction is on start of a line
*/
transaction(amount: UFix64, recipient: Address) {
// Here is a comment
// Here is a comment
prepare(signer: AuthAccount) {}
execute {}
			`,
		},
		{
			cadence: `import NonFungibleToken from 0xNonFungibleTokenAddress
			/*
			Here is a comment
			Her is a comment
			*/
			pub fun main(accountAddress: Address, tokenId: UInt64): Bool {
				return collectionRef.borrowNFT(id: tokenId) != nil
			}
			`,
		},
		{
			cadence: `pub contract interface TokenContract {
				/* Here is a comment */
				// Here is a comment
				pub fun totalSupply(): UFix64
			
				// Returns the balance of the specified address
				pub fun balanceOf(address: Address): UFix64
			
				// Transfers tokens from one address to another
				pub fun transfer(from: Address, to: Address, amount: UFix64): Bool
			}
			`,
		},
	}

	for _, tt := range tests {
		t.Run(tt.cadence, func(t *testing.T) {
			got := stripComments(tt.cadence)

			if strings.Contains(got, removedComment) {
				t.Errorf("stripComments got = %v, want no comments", got)
			}
		})
	}
}
