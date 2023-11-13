package v1_0_0

import (
	"context"
	"encoding/json"
	"testing"

	"github.com/hexops/autogold/v2"
	"github.com/onflow/flixkit-go"
	"github.com/onflow/flixkit-go/generator"
	"github.com/stretchr/testify/assert"
)

func TestGenerateWithPrefill(t *testing.T) {
	assert := assert.New(t)
	templatePreFill := `{
		"f_type": "InteractionTemplate",
		"f_version": "1.0.0",
		"id": "",
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
			"cadence": "",
			"dependencies": {},
			"arguments": {
				"amount": {
					"index": 0,
					"type": "UFix64",
					"messages": {
						"title": {
							"i18n": {
								"en-US": "Amount "
							}
						},
						"description": {
							"i18n": {
								"en-US": "Number of tokens to transfer"
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
								"en-US": "Account"
							}
						},
						"description": {
							"i18n": {
								"en-US": "Destination Account"
							}
						}
					},
					"balance": ""
				}
			}
		}
	}`
	prefill, _ := flixkit.ParseFlix(templatePreFill)

	code := `
	import FungibleToken from 0xFungibleTokenAddress
	import FlowToken from 0xFlowTokenAddress
	
	transaction(amount: UFix64, to: Address) {
	  let vault: @FungibleToken.Vault
	
	  prepare(signer: AuthAccount) {
		self.vault <- signer
		  .borrow<&{FungibleToken.Provider}>(from: /storage/flowTokenVault)!
		  .withdraw(amount: amount)
	  }
	
	  execute {
		getAccount(to)
		  .getCapability(/public/flowTokenReceiver)!
		  .borrow<&{FungibleToken.Receiver}>()!
		  .deposit(from: <-self.vault)
	  }
	}
`
	gen := GeneratorV1_0_0{
		deployedContracts: []flixkit.Contracts{},
		coreContracts:     generator.GetDefaultCoreContracts(),
		testnetClient:     nil,
		mainnetClient:     nil,
	}
	ctx := context.Background()
	template, err := gen.Generate(ctx, code, prefill)
	assert.NoError(err, "Generate should not return an error")
	prettyJSON, err := json.MarshalIndent(template, "", "    ")
	assert.NoError(err, "marshal template to json should not return an error")
	autogold.ExpectFile(t, string(prettyJSON))
}

func TestSimpleScriptGen(t *testing.T) {
	templatePreFill := `{
		"f_type": "InteractionTemplate",
		"f_version": "1.0.0",
		"id": "",
		"data": {
			"type": "script",
			"interface": "",
			"messages": {
				"title": {
					"i18n": {
						"en-US": "read Greeting"
					}
				},
				"description": {
					"i18n": {
						"en-US": "read greeting of the HelloWorld smart contract"
					}
				}
			},
			"cadence": "",
			"dependencies": {},
			"arguments": {}
		}
	}`
	prefill, _ := flixkit.ParseFlix(templatePreFill)
	contracts := []flixkit.Contracts{
		{
			"HelloWorld": flixkit.Networks{
				"emulator": flixkit.Network{
					Address:   "0x01cf0e2f2f715450",
					FqAddress: "A.01cf0e2f2f715450.HelloWorld",
					Contract:  "HelloWorld",
				},
			},
		},
	}
	generator := GeneratorV1_0_0{
		deployedContracts: contracts,
		testnetClient:     nil,
		mainnetClient:     nil,
	}
	assert := assert.New(t)
	code := `
	import "HelloWorld"

	pub fun main(): String {
	return HelloWorld.greeting
	}
`
	ctx := context.Background()
	template, err := generator.Generate(ctx, code, prefill)
	assert.NoError(err, "Generate should not return an error")
	prettyJSON, err := json.MarshalIndent(template, "", "    ")
	assert.NoError(err, "marshal template to json should not return an error")
	autogold.ExpectFile(t, string(prettyJSON))

}

func TestMinimumValues(t *testing.T) {
	contracts := []flixkit.Contracts{
		{
			"HelloWorld": flixkit.Networks{
				"emulator": flixkit.Network{
					Address:        "0x01cf0e2f2f715450",
					FqAddress:      "A.01cf0e2f2f715450.HelloWorld",
					Contract:       "HelloWorld",
					Pin:            "xxxxxmake-up-data",
					PinBlockHeight: 10,
				},
			},
		},
	}
	generator := GeneratorV1_0_0{
		deployedContracts: contracts,
		testnetClient:     nil,
		mainnetClient:     nil,
	}
	assert := assert.New(t)
	code := `
	import "HelloWorld"

	pub fun main(): String {
	return HelloWorld.greeting
	}
`
	ctx := context.Background()
	template, err := generator.Generate(ctx, code, nil)
	assert.NoError(err, "Generate should not return an error")
	prettyJSON, err := json.MarshalIndent(template, "", "    ")
	assert.NoError(err, "marshal template to json should not return an error")
	autogold.ExpectFile(t, string(prettyJSON))

}

func TestGetDependenceContract(t *testing.T) {
	fungi := flixkit.Contracts{
		"LocalContract": {
			"mainnet": {
				Address:        "0xf233dcee88fe0abe",
				FqAddress:      "A.0xf233dcee88fe0abe.LocalContract",
				Contract:       "LocalContract",
				Pin:            "83c9e3d61d3b5ebf24356a9f17b5b57b12d6d56547abc73e05f820a0ae7d9cf5",
				PinBlockHeight: 34166296,
			},
			"testnet": {
				Address:        "0x9a0766d93b6608b7",
				FqAddress:      "A.0x9a0766d93b6608b7.LocalContract",
				Contract:       "LocalContract",
				Pin:            "83c9e3d61d3b5ebf24356a9f17b5b57b12d6d56547abc73e05f820a0ae7d9cf5",
				PinBlockHeight: 74776482,
			},
		},
	}
	genV1 := GeneratorV1_0_0{
		deployedContracts: []flixkit.Contracts{fungi},
		coreContracts:     generator.GetDefaultCoreContracts(),
		testnetClient:     nil,
		mainnetClient:     nil,
	}

	tests := []struct {
		contractName string
		want         flixkit.Contracts
	}{
		{
			contractName: `FungibleToken`,
			want:         fungi,
		},
		{
			contractName: `LocalContract`,
			want:         fungi,
		},
	}

	for _, tt := range tests {
		t.Run(tt.contractName, func(t *testing.T) {
			got, err := genV1.generateDependenceInfo(context.Background(), tt.contractName)
			if err != nil {
				t.Errorf("generateDependenceInfo() err %v", err)
			}
			if got == nil {
				t.Errorf("generateDependenceInfo() got = %v, want %v", got, tt.want)
			}
			prettyJSON, err := json.MarshalIndent(got, "", "    ")
			if err != nil {
				t.Errorf("generateDependenceInfo() err %v", err)
			}
			autogold.ExpectFile(t, string(prettyJSON))
		})
	}
}
