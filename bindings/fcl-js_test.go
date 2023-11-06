package bindings

import (
	"testing"

	"github.com/hexops/autogold/v2"
	"github.com/onflow/flixkit-go"
	"github.com/stretchr/testify/assert"
)

var parsedTemplateTX = &flixkit.FlowInteractionTemplate{
	FType:    "InteractionTemplate",
	FVersion: "1.0.0",
	ID:       "290b6b6222b2a77b16db896a80ddf29ebd1fa3038c9e6625a933fa213fce51fa",
	Data: flixkit.Data{
		Type:      "transaction",
		Interface: "",
		Messages: flixkit.Messages{
			Title: &flixkit.Title{
				I18N: map[string]string{
					"en-US": "Transfer Tokens",
				},
			},
			Description: &flixkit.Description{
				I18N: map[string]string{
					"en-US": "Transfer tokens from one account to another",
				},
			},
		},
		Cadence: "import FungibleToken from 0xFUNGIBLETOKENADDRESS\ntransaction(amount: UFix64, to: Address) {\nlet vault: @FungibleToken.Vault\nprepare(signer: AuthAccount) {\nself.vault <- signer\n.borrow<&{FungibleToken.Provider}>(from: /storage/flowTokenVault)!\n.withdraw(amount: amount)\n}\nexecute {\ngetAccount(to)\n.getCapability(/public/flowTokenReceiver)!\n.borrow<&{FungibleToken.Receiver}>()!\n.deposit(from: <-self.vault)\n}\n}",
		Dependencies: flixkit.Dependencies{
			"0xFUNGIBLETOKENADDRESS": flixkit.Contracts{
				"FungibleToken": flixkit.Networks{
					"mainnet": flixkit.Network{
						Address:        "0xf233dcee88fe0abe",
						FqAddress:      "A.0xf233dcee88fe0abe.FungibleToken",
						Contract:       "FungibleToken",
						Pin:            "83c9e3d61d3b5ebf24356a9f17b5b57b12d6d56547abc73e05f820a0ae7d9cf5",
						PinBlockHeight: 34166296,
					},
					"testnet": flixkit.Network{
						Address:        "0x9a0766d93b6608b7",
						FqAddress:      "A.0x9a0766d93b6608b7.FungibleToken",
						Contract:       "FungibleToken",
						Pin:            "83c9e3d61d3b5ebf24356a9f17b5b57b12d6d56547abc73e05f820a0ae7d9cf5",
						PinBlockHeight: 74776482,
					},
				},
			},
		},
		Arguments: flixkit.Arguments{
			"amount": flixkit.Argument{
				Index: 0,
				Type:  "UFix64",
				Messages: flixkit.Messages{
					Title: &flixkit.Title{
						I18N: map[string]string{
							"en-US": "The amount of FLOW tokens to send",
						},
					},
				},
				Balance: "",
			},
			"to": flixkit.Argument{
				Index: 1,
				Type:  "Address",
				Messages: flixkit.Messages{
					Title: &flixkit.Title{
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

var parsedTemplateScript = &flixkit.FlowInteractionTemplate{
	FType:    "InteractionTemplate",
	FVersion: "1.0.0",
	ID:       "290b6b6222b2a77b16db896a80ddf29ebd1fa3038c9e6625a933fa213fce51fa",
	Data: flixkit.Data{
		Type:      "script",
		Interface: "",
		Messages: flixkit.Messages{
			Title: &flixkit.Title{
				I18N: map[string]string{
					"en-US": "Multiply Two Integers",
				},
			},
			Description: &flixkit.Description{
				I18N: map[string]string{
					"en-US": "Multiply two numbers to another",
				},
			},
		},
		Cadence: "pub fun main(x: Int, y: Int): Int { return x * y }",
		Arguments: flixkit.Arguments{
			"x": flixkit.Argument{
				Index: 0,
				Type:  "Int",
				Messages: flixkit.Messages{
					Title: &flixkit.Title{
						I18N: map[string]string{
							"en-US": "number to be multiplied",
						},
					},
				},
				Balance: "",
			},
			"y": flixkit.Argument{
				Index: 1,
				Type:  "Int",
				Messages: flixkit.Messages{
					Title: &flixkit.Title{
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

var ArrayTypeScript = &flixkit.FlowInteractionTemplate{
	FType:    "InteractionTemplate",
	FVersion: "1.0.0",
	ID:       "290b6b6222b2a77b16db896a80ddf29ebd1fa3038c9e6625a933fa213fce51fa",
	Data: flixkit.Data{
		Type:      "script",
		Interface: "",
		Messages: flixkit.Messages{
			Title: &flixkit.Title{
				I18N: map[string]string{
					"en-US": "Multiply Numbers",
				},
			},
			Description: &flixkit.Description{
				I18N: map[string]string{
					"en-US": "Multiply numbers in an array",
				},
			},
		},
		Cadence: "pub fun main(numbers: [Int]): Int { var total = 1; for x in numbers { total = total * x }; return total }",
		Arguments: flixkit.Arguments{
			"numbers": flixkit.Argument{
				Index: 0,
				Type:  "[Int]",
				Messages: flixkit.Messages{
					Title: &flixkit.Title{
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

var minimumTemplate = &flixkit.FlowInteractionTemplate{
	FType:    "InteractionTemplate",
	FVersion: "1.0.0",
	ID:       "290b6b6222b2a77b16db896a80ddf29ebd1fa3038c9e6625a933fa213fce51fa",
	Data: flixkit.Data{
		Type:      "script",
		Interface: "",
		Cadence:   "pub fun main(numbers: [Int]): Int { var total = 1; for x in numbers { total = total * x }; return total }",
		Arguments: flixkit.Arguments{
			"numbers": flixkit.Argument{
				Index: 0,
				Type:  "[Int]",
			},
		},
	},
}

var minimumNoParamTemplate = &flixkit.FlowInteractionTemplate{
	FType:    "InteractionTemplate",
	FVersion: "1.0.0",
	ID:       "290b6b6222b2a77b16db896a80ddf29ebd1fa3038c9e6625a933fa213fce51fa",
	Data: flixkit.Data{
		Type:      "script",
		Interface: "",
		Cadence:   "pub fun main(): Int { return 1 }",
	},
}

func TestJSGenTransaction(t *testing.T) {
	generator := FclJSGenerator{
		TemplateDir: "./templates",
	}
	got, _ := generator.Generate(parsedTemplateTX, "./transfer_token.json", true)
	autogold.ExpectFile(t, got)
}

func TestJSGenScript(t *testing.T) {
	generator := FclJSGenerator{
		TemplateDir: "./templates",
	}
	assert := assert.New(t)
	got, err := generator.Generate(parsedTemplateScript, "./multiply_two_integers.template.json", true)
	assert.NoError(err, "ParseTemplate should not return an error")
	autogold.ExpectFile(t, got)
}

func TestJSGenArrayScript(t *testing.T) {
	generator := FclJSGenerator{
		TemplateDir: "./templates",
	}
	assert := assert.New(t)

	out, err := generator.Generate(ArrayTypeScript, "./multiply-numbers.template.json", true)
	assert.NoError(err, "ParseTemplate should not return an error")
	autogold.ExpectFile(t, out)
}

func TestJSGenMinScript(t *testing.T) {
	generator := NewFclJSGenerator()
	assert := assert.New(t)

	out, err := generator.Generate(minimumTemplate, "./min.template.json", true)
	assert.NoError(err, "ParseTemplate should not return an error")
	autogold.ExpectFile(t, out)
}
func TestJSGenNoParamsScript(t *testing.T) {
	generator := NewFclJSGenerator()
	assert := assert.New(t)

	out, err := generator.Generate(minimumNoParamTemplate, "./min.template.json", true)
	assert.NoError(err, "ParseTemplate should not return an error")
	autogold.ExpectFile(t, out)
}
