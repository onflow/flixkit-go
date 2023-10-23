package generator

import (
	"testing"
)

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
