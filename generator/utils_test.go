package generator

import (
	"testing"

	"github.com/onflow/cadence/runtime/parser"
	"github.com/onflow/flixkit-go"
)

func TestDepCheck(t *testing.T) {
	tests := []struct {
		cadence     string
		cadenceType string
	}{
		{
			cadence: `import "FungibleToken"

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
			cadence: `import "FungibleToken"

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
execute {}}
			`,
			cadenceType: "transaction",
		},
		{
			cadence: `import "NonFungibleToken"

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
			codeBytes := []byte(tt.cadence)
			program, _ := parser.ParseProgram(nil, codeBytes, parser.Config{})
			got := DetermineCadenceType(program)
			if got != tt.cadenceType {
				t.Errorf("DetermineCadenceType() got = %v, want %v", got, tt.cadenceType)
			}
		})
	}

}

func TestGenerateTemplateId(t *testing.T) {
	templateId := "bd10ab0bf472e6b58ecc0398e9b3d1bd58a4205f14a7099c52c0640d9589295f"
	code := `
	{
		"f_type": "InteractionTemplate",
		"f_version": "1.0.0",
		"id": "",
		"data": {
		  "type": "script",
		  "interface": "",
		  "messages": {
			"title": {
			  "i18n": {
				"en-US": "Multiply Two Integers"
			  }
			},
			"description": {
			  "i18n": {
				"en-US": "Multiplies two integer arguments together and returns the result."
			  }
			}
		  },
		  "cadence": "pub fun main(x: Int, y: Int): Int { return x * y }",
		  "dependencies": {},
		  "arguments": {
			"x": {
			  "index": 0,
			  "type": "Int",
			  "messages": {
				"title": {
				  "i18n": {
					"en-US": "Int 1"
				  }
				}
			  }
			},
			"y": {
			  "index": 1,
			  "type": "Int",
			  "messages": {
				"title": {
				  "i18n": {
					"en-US": "Int 2"
				  }
				}
			  }
			}
		  }
		}
	  }`

	template, err := flixkit.ParseFlix(code)
	if err != nil {
		t.Errorf("ParseFlix() err %v", err)
	}
	id, err := flixkit.GenerateFlixID(template)
	if err != nil {
		t.Errorf("call GenerateFlixID err %v", err)
	}

	if id != templateId {
		t.Errorf("GenerateFlixID got = %v, want %v", id, templateId)
	}

}

func TestGenerateTemplateIdWithDeps(t *testing.T) {
	templateId := "290b6b6222b2a77b16db896a80ddf29ebd1fa3038c9e6625a933fa213fce51fa"
	code := `
	{
		"f_type": "InteractionTemplate",
		"f_version": "1.0.0",
		"id": "290b6b6222b2a77b16db896a80ddf29ebd1fa3038c9e6625a933fa213fce51fa",
		"data": {
		  "type": "transaction",
		  "interface": "",
		  "messages": {
			"title": {
			  "i18n": {
				"en-US": "Transfer Tokens"
			  }
			},
			"description": {
			  "i18n": {
				"en-US": "Transfer tokens from one account to another"
			  }
			}
		  },
		  "cadence": "import FungibleToken from 0xFUNGIBLETOKENADDRESS\ntransaction(amount: UFix64, to: Address) {\nlet vault: @FungibleToken.Vault\nprepare(signer: AuthAccount) {\nself.vault <- signer\n.borrow<&{FungibleToken.Provider}>(from: /storage/flowTokenVault)!\n.withdraw(amount: amount)\n}\nexecute {\ngetAccount(to)\n.getCapability(/public/flowTokenReceiver)!\n.borrow<&{FungibleToken.Receiver}>()!\n.deposit(from: <-self.vault)\n}\n}",
		  "dependencies": {
			"0xFUNGIBLETOKENADDRESS": {
			  "FungibleToken": {
				"mainnet": {
				  "address": "0xf233dcee88fe0abe",
				  "fq_address": "A.0xf233dcee88fe0abe.FungibleToken",
				  "contract": "FungibleToken",
				  "pin": "83c9e3d61d3b5ebf24356a9f17b5b57b12d6d56547abc73e05f820a0ae7d9cf5",
				  "pin_block_height": 34166296
				},
				"testnet": {
				  "address": "0x9a0766d93b6608b7",
				  "fq_address": "A.0x9a0766d93b6608b7.FungibleToken",
				  "contract": "FungibleToken",
				  "pin": "83c9e3d61d3b5ebf24356a9f17b5b57b12d6d56547abc73e05f820a0ae7d9cf5",
				  "pin_block_height": 74776482
				}
			  }
			}
		  },
		  "arguments": {
			"amount": {
			  "index": 0,
			  "type": "UFix64",
			  "messages": {
				"title": {
				  "i18n": {
					"en-US": "The amount of FLOW tokens to send"
				  }
				}
			  },
			  "balance": ""
			},
			"to": {
			  "index": 1,
			  "type": "Address",
			  "messages": {
				"title": {
				  "i18n": {
					"en-US": "The Flow account the tokens will go to"
				  }
				}
			  },
			  "balance": ""
			}
		  }
		}
	  }`

	template, err := flixkit.ParseFlix(code)
	if err != nil {
		t.Errorf("ParseFlix() err %v", err)
	}
	id, err := flixkit.GenerateFlixID(template)
	if err != nil {
		t.Errorf("GenerateFlixID err %v", err)
	}

	if id != templateId {
		t.Errorf("GenerateFlixID got = %v, want %v", id, templateId)
	}

}

func TestUnNormalizeCode(t *testing.T) {
	tests := []struct {
		cadence      string
		unNormalized string
	}{
		{
			cadence: `import "FungibleToken"
			/* Here is a comment */
			pub fun main(accountAddress: Address): UFix64 {
				return balanceRef.balance
			}
			`,
			unNormalized: `import FungibleToken from 0xFungibleToken
			/* Here is a comment */
			pub fun main(accountAddress: Address): UFix64 {
				return balanceRef.balance
			}
			`,
		},
		{
			cadence: `import "NonFungibleToken"
			import "FungibleToken"

			transaction(amount: UFix64, recipient: Address) {
				/* 
				Here is a comment
				*/
				execute {

				}
			}
			`,
			unNormalized: `import NonFungibleToken from 0xNonFungibleToken
			import FungibleToken from 0xFungibleToken

			transaction(amount: UFix64, recipient: Address) {
				/* 
				Here is a comment
				*/
				execute {

				}
			}
			`,
		},
	}

	for _, tt := range tests {
		t.Run(tt.cadence, func(t *testing.T) {
			got := UnNormalizeImports(tt.cadence)

			if got != tt.unNormalized {
				t.Errorf("UnNormalizedImports got = %v, want %v", got, tt.unNormalized)
			}
		})
	}

}

func TestNormalizeCode(t *testing.T) {
	tests := []struct {
		cadence    string
		normalized string
	}{
		{
			cadence: `import FungibleToken from 0xFungibleToken
			/* Here is a comment */
			pub fun main(accountAddress: Address): UFix64 {
				return balanceRef.balance
			}
			`,
			normalized: `import "FungibleToken"
			/* Here is a comment */
			pub fun main(accountAddress: Address): UFix64 {
				return balanceRef.balance
			}
			`,
		},
		{
			cadence: `import NonFungibleToken from 0xNonFungibleToken
			import FungibleToken from 0xFungibleToken

			transaction(amount: UFix64, recipient: Address) {
				/* 
				Here is a comment
				*/
				execute {

				}
			}
			`,
			normalized: `import "NonFungibleToken"
			import "FungibleToken"

			transaction(amount: UFix64, recipient: Address) {
				/* 
				Here is a comment
				*/
				execute {

				}
			}
			`,
		},
	}

	for _, tt := range tests {
		t.Run(tt.cadence, func(t *testing.T) {
			got := NormalizeImports(tt.cadence)

			if got != tt.normalized {
				t.Errorf("NormalizedImports got = %v, want %v", got, tt.normalized)
			}
		})
	}

}
