package js

import (
	"strings"
	"testing"

	"github.com/onflow/flixkit-go/common"
	"github.com/stretchr/testify/assert"
)


var parsedTemplateTX = &common.FlowInteractionTemplate{
	FType:    "InteractionTemplate",
	FVersion: "1.0.0",
	ID:       "290b6b6222b2a77b16db896a80ddf29ebd1fa3038c9e6625a933fa213fce51fa",
	Data: common.Data{
		Type:      "transaction",
		Interface: "",
		Messages: common.Messages{
			Title: &common.Title{
				I18N: map[string]string{
					"en-US": "Transfer Tokens",
				},
			},
			Description: &common.Description{
				I18N: map[string]string{
					"en-US": "Transfer tokens from one account to another",
				},
			},
		},
		Cadence: "import FungibleToken from 0xFUNGIBLETOKENADDRESS\ntransaction(amount: UFix64, to: Address) {\nlet vault: @FungibleToken.Vault\nprepare(signer: AuthAccount) {\nself.vault <- signer\n.borrow<&{FungibleToken.Provider}>(from: /storage/flowTokenVault)!\n.withdraw(amount: amount)\n}\nexecute {\ngetAccount(to)\n.getCapability(/public/flowTokenReceiver)!\n.borrow<&{FungibleToken.Receiver}>()!\n.deposit(from: <-self.vault)\n}\n}",
		Dependencies: common.Dependencies{
			"0xFUNGIBLETOKENADDRESS": common.Contracts{
				"FungibleToken": common.Networks{
					"mainnet": common.Network{
						Address:        "0xf233dcee88fe0abe",
						FqAddress:      "A.0xf233dcee88fe0abe.FungibleToken",
						Contract:       "FungibleToken",
						Pin:            "83c9e3d61d3b5ebf24356a9f17b5b57b12d6d56547abc73e05f820a0ae7d9cf5",
						PinBlockHeight: 34166296,
					},
					"testnet": common.Network{
						Address:        "0x9a0766d93b6608b7",
						FqAddress:      "A.0x9a0766d93b6608b7.FungibleToken",
						Contract:       "FungibleToken",
						Pin:            "83c9e3d61d3b5ebf24356a9f17b5b57b12d6d56547abc73e05f820a0ae7d9cf5",
						PinBlockHeight: 74776482,
					},
				},
			},
		},
		Arguments: common.Arguments{
			"amount": common.Argument{
				Index: 0,
				Type:  "UFix64",
				Messages: common.Messages{
					Title: &common.Title{
						I18N: map[string]string{
							"en-US": "The amount of FLOW tokens to send",
						},
					},
				},
				Balance: "",
			},
			"to": common.Argument{
				Index: 1,
				Type:  "Address",
				Messages: common.Messages{
					Title: &common.Title{
						I18N: map[string]string{
							"en-US": "The Flow account the tokens will go to",
						},
					},
				},
				Balance: "",
			},
		},
	},
}


var parsedTemplateScript = &common.FlowInteractionTemplate{
	FType:    "InteractionTemplate",
	FVersion: "1.0.0",
	ID:       "290b6b6222b2a77b16db896a80ddf29ebd1fa3038c9e6625a933fa213fce51fa",
	Data: common.Data{
		Type:      "script",
		Interface: "",
		Messages: common.Messages{
			Title: &common.Title{
				I18N: map[string]string{
					"en-US": "Transfer Tokens",
				},
			},
			Description: &common.Description{
				I18N: map[string]string{
					"en-US": "Transfer tokens from one account to another",
				},
			},
		},
		Cadence: "pub fun main(x: Int, y: Int): Int { return x * y }",
		Arguments: common.Arguments{
			"x": common.Argument{
				Index: 0,
				Type:  "Int",
				Messages: common.Messages{
					Title: &common.Title{
						I18N: map[string]string{
							"en-US": "number to be multiplied",
						},
					},
				},
				Balance: "",
			},
			"y": common.Argument{
				Index: 1,
				Type:  "Int",
				Messages: common.Messages{
					Title: &common.Title{
						I18N: map[string]string{
							"en-US": "second number to be multiplied",
						},
					},
				},
				Balance: "",
			},
		},
	},
}


var ArrayTypeScript = &common.FlowInteractionTemplate{
	FType:    "InteractionTemplate",
	FVersion: "1.0.0",
	ID:       "290b6b6222b2a77b16db896a80ddf29ebd1fa3038c9e6625a933fa213fce51fa",
	Data: common.Data{
		Type:      "script",
		Interface: "",
		Messages: common.Messages{
			Title: &common.Title{
				I18N: map[string]string{
					"en-US": "Multiply Numbers",
				},
			},
			Description: &common.Description{
				I18N: map[string]string{
					"en-US": "Multiply numbers in an array",
				},
			},
		},
		Cadence: "pub fun main(numbers: [Int]): Int { var total = 1; for x in numbers { total = total * x }; return total }",
		Arguments: common.Arguments{
			"numbers": common.Argument{
				Index: 0,
				Type:  "[Int]",
				Messages: common.Messages{
					Title: &common.Title{
						I18N: map[string]string{
							"en-US": "Array of numbers to be multiplied",
						},
					},
				},
				Balance: "",
			},
		},
	},
}
func TestJSGenTransaction(t *testing.T) {
	assert := assert.New(t)

	contents, err := GenerateJavaScript(parsedTemplateTX, "./transfer_token.json")
	assert.NoError(err, "ParseTemplate should not return an error")
	assert.NotNil(contents, "Parsed template should not be nil")
	assert.True(strings.Contains(contents, "await fcl.mutate("), "Expected '%s'", "await fcl.mutate(")
	println(contents)

}


func TestJSGenScript(t *testing.T) {
	assert := assert.New(t)

	contents, err := GenerateJavaScript(parsedTemplateScript, "./transfer_token.json")
	assert.NoError(err, "ParseTemplate should not return an error")
	assert.NotNil(contents, "Parsed template should not be nil")
	assert.True(strings.Contains(contents, "await fcl.query("), "Expected '%s'", "await fcl.query(")
	println(contents)

}



func TestJSGenArrayScript(t *testing.T) {
	assert := assert.New(t)

	contents, err := GenerateJavaScript(ArrayTypeScript, "./multiply-numbers.template.json")
	assert.NoError(err, "ParseTemplate should not return an error")
	assert.NotNil(contents, "Parsed template should not be nil")
	assert.True(strings.Contains(contents, "await fcl.query("), "Expected '%s'", "await fcl.query(")
	assert.True(strings.Contains(contents, `args: (arg, t) => ([ args(numbers, t.Array(t.Int))])`), "Expected '%s'", `args([ fcl.arg(numbers, t.Array(t.Int)) ])`)
	println(contents)

}