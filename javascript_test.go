package flixkit

import (
	"testing"

	"github.com/hexops/autogold/v2"
	"github.com/stretchr/testify/assert"
)


var parsedTemplateTX = &FlowInteractionTemplate{
	FType:    "InteractionTemplate",
	FVersion: "1.0.0",
	ID:       "290b6b6222b2a77b16db896a80ddf29ebd1fa3038c9e6625a933fa213fce51fa",
	Data: Data{
		Type:      "transaction",
		Interface: "",
		Messages: Messages{
			Title: &Title{
				I18N: map[string]string{
					"en-US": "Transfer Tokens",
				},
			},
			Description: &Description{
				I18N: map[string]string{
					"en-US": "Transfer tokens from one account to another",
				},
			},
		},
		Cadence: "import FungibleToken from 0xFUNGIBLETOKENADDRESS\ntransaction(amount: UFix64, to: Address) {\nlet vault: @FungibleToken.Vault\nprepare(signer: AuthAccount) {\nself.vault <- signer\n.borrow<&{FungibleToken.Provider}>(from: /storage/flowTokenVault)!\n.withdraw(amount: amount)\n}\nexecute {\ngetAccount(to)\n.getCapability(/public/flowTokenReceiver)!\n.borrow<&{FungibleToken.Receiver}>()!\n.deposit(from: <-self.vault)\n}\n}",
		Dependencies: Dependencies{
			"0xFUNGIBLETOKENADDRESS": Contracts{
				"FungibleToken": Networks{
					"mainnet": Network{
						Address:        "0xf233dcee88fe0abe",
						FqAddress:      "A.0xf233dcee88fe0abe.FungibleToken",
						Contract:       "FungibleToken",
						Pin:            "83c9e3d61d3b5ebf24356a9f17b5b57b12d6d56547abc73e05f820a0ae7d9cf5",
						PinBlockHeight: 34166296,
					},
					"testnet": Network{
						Address:        "0x9a0766d93b6608b7",
						FqAddress:      "A.0x9a0766d93b6608b7.FungibleToken",
						Contract:       "FungibleToken",
						Pin:            "83c9e3d61d3b5ebf24356a9f17b5b57b12d6d56547abc73e05f820a0ae7d9cf5",
						PinBlockHeight: 74776482,
					},
				},
			},
		},
		Arguments: Arguments{
			"amount": Argument{
				Index: 0,
				Type:  "UFix64",
				Messages: Messages{
					Title: &Title{
						I18N: map[string]string{
							"en-US": "The amount of FLOW tokens to send",
						},
					},
				},
				Balance: "",
			},
			"to": Argument{
				Index: 1,
				Type:  "Address",
				Messages: Messages{
					Title: &Title{
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


var parsedTemplateScript = &FlowInteractionTemplate{
	FType:    "InteractionTemplate",
	FVersion: "1.0.0",
	ID:       "290b6b6222b2a77b16db896a80ddf29ebd1fa3038c9e6625a933fa213fce51fa",
	Data: Data{
		Type:      "script",
		Interface: "",
		Messages: Messages{
			Title: &Title{
				I18N: map[string]string{
					"en-US": "Multiply Two Integers",
				},
			},
			Description: &Description{
				I18N: map[string]string{
					"en-US": "Multiply two numbers to another",
				},
			},
		},
		Cadence: "pub fun main(x: Int, y: Int): Int { return x * y }",
		Arguments: Arguments{
			"x": Argument{
				Index: 0,
				Type:  "Int",
				Messages: Messages{
					Title: &Title{
						I18N: map[string]string{
							"en-US": "number to be multiplied",
						},
					},
				},
				Balance: "",
			},
			"y": Argument{
				Index: 1,
				Type:  "Int",
				Messages: Messages{
					Title: &Title{
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


var ArrayTypeScript = &FlowInteractionTemplate{
	FType:    "InteractionTemplate",
	FVersion: "1.0.0",
	ID:       "290b6b6222b2a77b16db896a80ddf29ebd1fa3038c9e6625a933fa213fce51fa",
	Data: Data{
		Type:      "script",
		Interface: "",
		Messages: Messages{
			Title: &Title{
				I18N: map[string]string{
					"en-US": "Multiply Numbers",
				},
			},
			Description: &Description{
				I18N: map[string]string{
					"en-US": "Multiply numbers in an array",
				},
			},
		},
		Cadence: "pub fun main(numbers: [Int]): Int { var total = 1; for x in numbers { total = total * x }; return total }",
		Arguments: Arguments{
			"numbers": Argument{
				Index: 0,
				Type:  "[Int]",
				Messages: Messages{
					Title: &Title{
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

var minimumTemplate = &FlowInteractionTemplate{
	FType:    "InteractionTemplate",
	FVersion: "1.0.0",
	ID:       "290b6b6222b2a77b16db896a80ddf29ebd1fa3038c9e6625a933fa213fce51fa",
	Data: Data{
		Type:      "script",
		Interface: "",
		Cadence: "pub fun main(numbers: [Int]): Int { var total = 1; for x in numbers { total = total * x }; return total }",
		Arguments: Arguments{
			"numbers": Argument{
				Index: 0,
				Type:  "[Int]",
			},
		},
	},
}
func TestJSGenTransaction(t *testing.T) {
	got, _ := GenerateJavaScript(parsedTemplateTX, "./transfer_token.json", true)
	autogold.ExpectFile(t, got)
}

func TestJSGenScript(t *testing.T) {
	assert := assert.New(t)

	got, err:= GenerateJavaScript(parsedTemplateScript, "./multiply_two_integers.template.json", true)
	assert.NoError(err, "ParseTemplate should not return an error")
	autogold.ExpectFile(t, got)
}

func TestJSGenArrayScript(t *testing.T) {
	assert := assert.New(t)

	out, err := GenerateJavaScript(ArrayTypeScript, "./multiply-numbers.template.json", true)
	assert.NoError(err, "ParseTemplate should not return an error")
	autogold.ExpectFile(t, out)
}

func TestJSGenMinScript(t *testing.T) {
	assert := assert.New(t)

	out, err := GenerateJavaScript(minimumTemplate, "./min.template.json", true)
	assert.NoError(err, "ParseTemplate should not return an error")
	autogold.ExpectFile(t, out)
}