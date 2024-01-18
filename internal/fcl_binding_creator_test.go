package internal

import (
	"encoding/json"
	"testing"

	"github.com/hexops/autogold/v2"
	"github.com/onflow/flixkit-go/internal/templates"
	v1 "github.com/onflow/flixkit-go/internal/v1"
	v1_1 "github.com/onflow/flixkit-go/internal/v1_1"
	"github.com/stretchr/testify/assert"
)

var parsedTemplateTX = &v1.FlowInteractionTemplate{
	FType:    "InteractionTemplate",
	FVersion: "1.0.0",
	ID:       "290b6b6222b2a77b16db896a80ddf29ebd1fa3038c9e6625a933fa213fce51fa",
	Data: v1.Data{
		Type:      "transaction",
		Interface: "",
		Messages: v1.Messages{
			Title: &v1.Title{
				I18N: map[string]string{
					"en-US": "Transfer Tokens",
				},
			},
			Description: &v1.Description{
				I18N: map[string]string{
					"en-US": "Transfer tokens from one account to another",
				},
			},
		},
		Cadence: "import FungibleToken from 0xFUNGIBLETOKENADDRESS\ntransaction(amount: UFix64, to: Address) {\nlet vault: @FungibleToken.Vault\nprepare(signer: AuthAccount) {\nself.vault <- signer\n.borrow<&{FungibleToken.Provider}>(from: /storage/flowTokenVault)!\n.withdraw(amount: amount)\n}\nexecute {\ngetAccount(to)\n.getCapability(/public/flowTokenReceiver)!\n.borrow<&{FungibleToken.Receiver}>()!\n.deposit(from: <-self.vault)\n}\n}",
		Dependencies: v1.Dependencies{
			"0xFUNGIBLETOKENADDRESS": v1.Contracts{
				"FungibleToken": v1.Networks{
					"mainnet": v1.Network{
						Address:        "0xf233dcee88fe0abe",
						FqAddress:      "A.0xf233dcee88fe0abe.FungibleToken",
						Contract:       "FungibleToken",
						Pin:            "83c9e3d61d3b5ebf24356a9f17b5b57b12d6d56547abc73e05f820a0ae7d9cf5",
						PinBlockHeight: 34166296,
					},
					"testnet": v1.Network{
						Address:        "0x9a0766d93b6608b7",
						FqAddress:      "A.0x9a0766d93b6608b7.FungibleToken",
						Contract:       "FungibleToken",
						Pin:            "83c9e3d61d3b5ebf24356a9f17b5b57b12d6d56547abc73e05f820a0ae7d9cf5",
						PinBlockHeight: 74776482,
					},
				},
			},
		},
		Arguments: v1.Arguments{
			"amount": v1.Argument{
				Index: 0,
				Type:  "UFix64",
				Messages: v1.Messages{
					Title: &v1.Title{
						I18N: map[string]string{
							"en-US": "The amount of FLOW tokens to send",
						},
					},
				},
				Balance: "",
			},
			"to": v1.Argument{
				Index: 1,
				Type:  "Address",
				Messages: v1.Messages{
					Title: &v1.Title{
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

var parsedTemplateScript = &v1.FlowInteractionTemplate{
	FType:    "InteractionTemplate",
	FVersion: "1.0.0",
	ID:       "290b6b6222b2a77b16db896a80ddf29ebd1fa3038c9e6625a933fa213fce51fa",
	Data: v1.Data{
		Type:      "script",
		Interface: "",
		Messages: v1.Messages{
			Title: &v1.Title{
				I18N: map[string]string{
					"en-US": "Multiply Two Integers",
				},
			},
			Description: &v1.Description{
				I18N: map[string]string{
					"en-US": "Multiply two numbers to another",
				},
			},
		},
		Cadence: "pub fun main(x: Int, y: Int): Int { return x * y }",
		Arguments: v1.Arguments{
			"x": v1.Argument{
				Index: 0,
				Type:  "Int",
				Messages: v1.Messages{
					Title: &v1.Title{
						I18N: map[string]string{
							"en-US": "number to be multiplied",
						},
					},
				},
				Balance: "",
			},
			"y": v1.Argument{
				Index: 1,
				Type:  "Int",
				Messages: v1.Messages{
					Title: &v1.Title{
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

var ArrayTypeScript = &v1.FlowInteractionTemplate{
	FType:    "InteractionTemplate",
	FVersion: "1.0.0",
	ID:       "290b6b6222b2a77b16db896a80ddf29ebd1fa3038c9e6625a933fa213fce51fa",
	Data: v1.Data{
		Type:      "script",
		Interface: "",
		Messages: v1.Messages{
			Title: &v1.Title{
				I18N: map[string]string{
					"en-US": "Multiply Numbers",
				},
			},
			Description: &v1.Description{
				I18N: map[string]string{
					"en-US": "Multiply numbers in an array",
				},
			},
		},
		Cadence: "pub fun main(numbers: [Int]): Int { var total = 1; for x in numbers { total = total * x }; return total }",
		Arguments: v1.Arguments{
			"numbers": v1.Argument{
				Index: 0,
				Type:  "[Int]",
				Messages: v1.Messages{
					Title: &v1.Title{
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

var minimumTemplate = &v1.FlowInteractionTemplate{
	FType:    "InteractionTemplate",
	FVersion: "1.0.0",
	ID:       "290b6b6222b2a77b16db896a80ddf29ebd1fa3038c9e6625a933fa213fce51fa",
	Data: v1.Data{
		Type:      "script",
		Interface: "",
		Cadence:   "pub fun main(numbers: [Int]): Int { var total = 1; for x in numbers { total = total * x }; return total }",
		Arguments: v1.Arguments{
			"numbers": v1.Argument{
				Index: 0,
				Type:  "[Int]",
			},
		},
	},
}

var minimumNoParamTemplate = &v1.FlowInteractionTemplate{
	FType:    "InteractionTemplate",
	FVersion: "1.0.0",
	ID:       "290b6b6222b2a77b16db896a80ddf29ebd1fa3038c9e6625a933fa213fce51fa",
	Data: v1.Data{
		Type:      "script",
		Interface: "",
		Cadence:   "pub fun main(): Int { return 1 }",
	},
}

func TestJSGenTransaction(t *testing.T) {
	ttemp, err := json.Marshal(parsedTemplateTX)
	assert.NoError(t, err, "marshal template to json should not return an error")
	generator := FclCreator{
		templates: []string{
			templates.GetJsFclMainTemplate(),
			templates.GetJsFclScriptTemplate(),
			templates.GetJsFclTxTemplate(),
			templates.GetJsFclParamsTemplate(),
		},
	}
	got, _ := generator.Create(string(ttemp), "./transfer_token.json")
	autogold.ExpectFile(t, got)
}

func TestJSGenScript(t *testing.T) {
	ttemp, err := json.Marshal(parsedTemplateScript)
	assert.NoError(t, err, "marshal template to json should not return an error")

	generator := FclCreator{
		templates: []string{
			templates.GetJsFclMainTemplate(),
			templates.GetJsFclScriptTemplate(),
			templates.GetJsFclTxTemplate(),
			templates.GetJsFclParamsTemplate(),
		},
	}
	assert := assert.New(t)
	got, err := generator.Create(string(ttemp), "./multiply_two_integers.template.json")
	assert.NoError(err, "ParseTemplate should not return an error")
	autogold.ExpectFile(t, got)
}

func TestJSGenArrayScript(t *testing.T) {
	ttemp, err := json.Marshal(ArrayTypeScript)
	assert.NoError(t, err, "marshal template to json should not return an error")

	generator := FclCreator{
		templates: []string{
			templates.GetJsFclMainTemplate(),
			templates.GetJsFclScriptTemplate(),
			templates.GetJsFclTxTemplate(),
			templates.GetJsFclParamsTemplate(),
		},
	}
	assert := assert.New(t)

	out, err := generator.Create(string(ttemp), "./multiply-numbers.template.json")
	assert.NoError(err, "ParseTemplate should not return an error")
	autogold.ExpectFile(t, out)
}

func TestJSGenMinScript(t *testing.T) {
	ttemp, err := json.Marshal(minimumTemplate)
	assert.NoError(t, err, "marshal template to json should not return an error")

	generator := NewFclJSCreator()
	assert := assert.New(t)

	out, err := generator.Create(string(ttemp), "./min.template.json")
	assert.NoError(err, "ParseTemplate should not return an error")
	autogold.ExpectFile(t, out)
}
func TestJSGenNoParamsScript(t *testing.T) {
	ttemp, err := json.Marshal(minimumNoParamTemplate)
	assert.NoError(t, err, "marshal template to json should not return an error")

	generator := NewFclJSCreator()
	assert := assert.New(t)

	out, err := generator.Create(string(ttemp), "./min.template.json")
	assert.NoError(err, "ParseTemplate should not return an error")
	autogold.ExpectFile(t, out)
}

var minimumNoParamTemplateTS_SCRIPT = &v1_1.InteractionTemplate{
	FType:    "InteractionTemplate",
	FVersion: "1.1.0",
	ID:       "290b6b6222b2a77b16db896a80ddf29ebd1fa3038c9e6625a933fa213fce51fa",
	Data: v1_1.Data{
		Type:      "script",
		Interface: "",
		Cadence: v1_1.Cadence{
			Body:        "pub fun main(): Int { return 1 }",
			NetworkPins: []v1_1.NetworkPin{},
		},
		Output: &v1_1.Parameter{
			Label: "result",
			Type:  "Int",
			Messages: []v1_1.Message{
				{
					Key: "description",
					I18n: []v1_1.I18n{
						{
							Tag:         "en-US",
							Translation: "Result of some number plus one",
						},
					},
				},
			},
		},
	},
}

var minimumNoParamTemplateTS_TX = &v1_1.InteractionTemplate{
	FType:    "InteractionTemplate",
	FVersion: "1.1.0",
	ID:       "290b6b6222b2a77b16db896a80ddf29ebd1fa3038c9e6625a933fa213fce51fa",
	Data: v1_1.Data{
		Type:      "transaction",
		Interface: "",
		Cadence: v1_1.Cadence{
			Body:        "import \"HelloWorld\"\n\n#interaction (\n  version: \"1.1.0\",\n\ttitle: \"Update Greeting\",\n\tdescription: \"Update the greeting on the HelloWorld contract\",\n\tlanguage: \"en-US\",\n)\ntransaction() {\n\n  prepare(acct: AuthAccount) {\n     }\n\n  execute {\n   \n  }\n}\n",
			NetworkPins: []v1_1.NetworkPin{},
		},
	},
}

var minimumParamTemplateTS_SCRIPT = &v1_1.InteractionTemplate{
	FType:    "InteractionTemplate",
	FVersion: "1.1.0",
	ID:       "290b6b6222b2a77b16db896a80ddf29ebd1fa3038c9e6625a933fa213fce51fa",
	Data: v1_1.Data{
		Type:      "script",
		Interface: "",
		Cadence: v1_1.Cadence{
			Body:        "pub fun main(someNumber Int): Int { return 1 + someNumber }",
			NetworkPins: []v1_1.NetworkPin{},
		},
		Parameters: []v1_1.Parameter{
			{
				Label: "someNumber",
				Index: 0,
				Type:  "Int",
				Messages: []v1_1.Message{
					{
						Key: "title",
						I18n: []v1_1.I18n{
							{
								Tag:         "en-US",
								Translation: "Some Number",
							},
						},
					},
				},
			},
		},
		Output: &v1_1.Parameter{
			Label: "result",
			Type:  "Int",
			Messages: []v1_1.Message{
				{
					Key: "description",
					I18n: []v1_1.I18n{
						{
							Tag:         "en-US",
							Translation: "Result of some number plus one",
						},
					},
				},
			},
		},
	},
}

var minimumParamTemplateTS_TX = &v1_1.InteractionTemplate{
	FType:    "InteractionTemplate",
	FVersion: "1.1.0",
	ID:       "290b6b6222b2a77b16db896a80ddf29ebd1fa3038c9e6625a933fa213fce51fa",
	Data: v1_1.Data{
		Type:      "transaction",
		Interface: "",
		Cadence: v1_1.Cadence{
			Body:        "import \"HelloWorld\"\n\n#interaction (\n  version: \"1.1.0\",\n\ttitle: \"Update Greeting\",\n\tdescription: \"Update the greeting on the HelloWorld contract\",\n\tlanguage: \"en-US\",\n\tparameters: [\n\t\tParameter(\n\t\t\tname: \"greeting\", \n\t\t\ttitle: \"Greeting\", \n\t\t\tdescription: \"The greeting to set on the HelloWorld contract\"\n\t\t)\n\t],\n)\ntransaction(greeting: String) {\n\n  prepare(acct: AuthAccount) {\n    log(acct.address)\n  }\n\n  execute {\n    HelloWorld.updateGreeting(newGreeting: greeting)\n  }\n}\n",
			NetworkPins: []v1_1.NetworkPin{},
		},
		Messages: []v1_1.Message{
			{
				Key: "title",
				I18n: []v1_1.I18n{
					{
						Tag:         "en-US",
						Translation: "Update Greeting",
					},
				},
			},
			{
				Key: "description",
				I18n: []v1_1.I18n{
					{
						Tag:         "en-US",
						Translation: "Update HelloWorld Greeting",
					},
				},
			},
		},
		Parameters: []v1_1.Parameter{
			{
				Label: "greeting",
				Index: 0,
				Type:  "String",
				Messages: []v1_1.Message{
					{
						Key: "title",
						I18n: []v1_1.I18n{
							{
								Tag:         "en-US",
								Translation: "Greeting",
							},
						},
					},
				},
			},
		},
	},
}

func TestTSGenNoParamsScript(t *testing.T) {
	ttemp, err := json.Marshal(minimumNoParamTemplateTS_SCRIPT)
	assert.NoError(t, err, "marshal template to json should not return an error")

	generator := NewFclTSCreator()
	assert := assert.New(t)

	out, err := generator.Create(string(ttemp), "./min.template.json")
	assert.NoError(err, "ParseTemplate should not return an error")
	autogold.ExpectFile(t, out)
}

func TestTSGenNoParamsTx(t *testing.T) {
	ttemp, err := json.Marshal(minimumNoParamTemplateTS_TX)
	assert.NoError(t, err, "marshal template to json should not return an error")

	generator := NewFclTSCreator()
	assert := assert.New(t)

	out, err := generator.Create(string(ttemp), "./min.template.json")
	assert.NoError(err, "ParseTemplate should not return an error")
	autogold.ExpectFile(t, out)
}

func TestTSGenParamsScript(t *testing.T) {
	ttemp, err := json.Marshal(minimumParamTemplateTS_SCRIPT)
	assert.NoError(t, err, "marshal template to json should not return an error")

	generator := NewFclTSCreator()
	assert := assert.New(t)

	out, err := generator.Create(string(ttemp), "./min.template.json")
	assert.NoError(err, "ParseTemplate should not return an error")
	autogold.ExpectFile(t, out)
}

func TestTSGenParamsTx(t *testing.T) {
	ttemp, err := json.Marshal(minimumParamTemplateTS_TX)
	assert.NoError(t, err, "marshal template to json should not return an error")

	generator := NewFclTSCreator()
	assert := assert.New(t)

	out, err := generator.Create(string(ttemp), "./min.template.json")
	assert.NoError(err, "ParseTemplate should not return an error")
	autogold.ExpectFile(t, out)
}

const ReadTokenScript = `
{
    "f_type": "InteractionTemplate",
    "f_version": "1.1.0",
    "id": "29d03aafbbb5a02e0d5f4ffee685c12494915410812305c2858008d3e2902b72",
    "data": {
        "type": "script",
        "interface": "",
        "messages": null,
        "cadence": {
            "body": "import \"FungibleToken\"\nimport \"FlowToken\"\n\npub fun main(address: Address): UFix64 {\n    let account = getAccount(address)\n\n    let vaultRef = account\n        .getCapability(/public/flowTokenBalance)\n        .borrow\u003c\u0026FlowToken.Vault{FungibleToken.Balance}\u003e()\n        ?? panic(\"Could not borrow balance reference to the Vault\")\n\n    return vaultRef.balance\n}\n",
            "network_pins": [
                {
                    "network": "mainnet",
                    "pin_self": "e0a1c0443b724d1238410c4a05c48441ee974160cad8cf1103c63b6999f81dd5"
                },
                {
                    "network": "testnet",
                    "pin_self": "6fee459b35d7013a83070c9ac42ea43ee04a3925deca445c34614c1bd6dc4cb8"
                }
            ]
        },
        "dependencies": [
            {
                "contracts": [
                    {
                        "contract": "FungibleToken",
                        "networks": [
                            {
                                "network": "mainnet",
                                "address": "0xf233dcee88fe0abe",
                                "dependency_pin_block_height": 69539302,
                                "dependency_pin": {
                                    "pin": "ac0208f93d07829ec96584d618ddbec6af3cf4e2866bd5071249e8ec93c7e0dc",
                                    "pin_self": "cdadd5b5897f2dfe35d8b25f4e41fea9f8fca8f40f8a8b506b33701ef5033076",
                                    "pin_contract_name": "FungibleToken",
                                    "pin_contract_address": "0xf233dcee88fe0abe",
                                    "imports": []
                                }
                            },
                            {
                                "network": "testnet",
                                "address": "0x9a0766d93b6608b7",
                                "dependency_pin_block_height": 146201102,
                                "dependency_pin": {
                                    "pin": "ac0208f93d07829ec96584d618ddbec6af3cf4e2866bd5071249e8ec93c7e0dc",
                                    "pin_self": "cdadd5b5897f2dfe35d8b25f4e41fea9f8fca8f40f8a8b506b33701ef5033076",
                                    "pin_contract_name": "FungibleToken",
                                    "pin_contract_address": "0x9a0766d93b6608b7",
                                    "imports": []
                                }
                            },
                            {
                                "network": "emulator",
                                "address": "0xee82856bf20e2aa6",
                                "dependency_pin_block_height": 0
                            }
                        ]
                    }
                ]
            },
            {
                "contracts": [
                    {
                        "contract": "FlowToken",
                        "networks": [
                            {
                                "network": "mainnet",
                                "address": "0x1654653399040a61",
                                "dependency_pin_block_height": 69539302,
                                "dependency_pin": {
                                    "pin": "a341e772da413bfbcf43b0fc167bd50a20c9f40baf10e12d3dbc2f5181526de9",
                                    "pin_self": "0e932728b73bff3c09dd58922f2529fc7b7fe7477f1dcc61169bc8f46948ad91",
                                    "pin_contract_name": "FlowToken",
                                    "pin_contract_address": "0x1654653399040a61",
                                    "imports": [
                                        {
                                            "pin": "ac0208f93d07829ec96584d618ddbec6af3cf4e2866bd5071249e8ec93c7e0dc",
                                            "pin_self": "cdadd5b5897f2dfe35d8b25f4e41fea9f8fca8f40f8a8b506b33701ef5033076",
                                            "pin_contract_name": "FungibleToken",
                                            "pin_contract_address": "0xf233dcee88fe0abe",
                                            "imports": []
                                        }
                                    ]
                                }
                            },
                            {
                                "network": "testnet",
                                "address": "0x7e60df042a9c0868",
                                "dependency_pin_block_height": 146201102,
                                "dependency_pin": {
                                    "pin": "9cc21a34a01486ebd6f044e99dbcdd58671850f81fcc345d071181c19f61aaa4",
                                    "pin_self": "6f01c7001e2d6635b667a170d3ccbc13659c40d01bb35e56979fcc7fa2d18646",
                                    "pin_contract_name": "FlowToken",
                                    "pin_contract_address": "0x7e60df042a9c0868",
                                    "imports": [
                                        {
                                            "pin": "ac0208f93d07829ec96584d618ddbec6af3cf4e2866bd5071249e8ec93c7e0dc",
                                            "pin_self": "cdadd5b5897f2dfe35d8b25f4e41fea9f8fca8f40f8a8b506b33701ef5033076",
                                            "pin_contract_name": "FungibleToken",
                                            "pin_contract_address": "0x9a0766d93b6608b7",
                                            "imports": []
                                        }
                                    ]
                                }
                            },
                            {
                                "network": "emulator",
                                "address": "0x0ae53cb6e3f42a79",
                                "dependency_pin_block_height": 0
                            }
                        ]
                    }
                ]
            }
        ],
        "parameters": [
            {
                "label": "address",
                "index": 0,
                "type": "Address",
                "messages": []
            }
        ]
    }
}
`

func TestBindingReadTokenBalance(t *testing.T) {
	generator := NewFclTSCreator()
	assert := assert.New(t)

	out, err := generator.Create(ReadTokenScript, "./read-token-balance.template.json")
	assert.NoError(err, "ParseTemplate should not return an error")
	autogold.ExpectFile(t, out)
}
