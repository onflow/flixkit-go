package v1_1

import (
	"encoding/json"
	"testing"

	"github.com/hexops/autogold/v2"
	"github.com/onflow/cadence/parser"
	"github.com/stretchr/testify/assert"
)

var templateWithDepsMissingLeading0x = `
{
	"f_type": "InteractionTemplate",
	"f_version": "1.1.0",
	"id": "",
	"data": {
		"type": "transaction",
		"interface": "",
		"messages": [],
		"cadence": {
			"body": "import \"HelloWorld\"\n\n#interaction (\n  version: \"1.1.0\",\n\ttitle: \"Update Greeting\",\n\tdescription: \"Update the greeting on the HelloWorld contract\",\n\tlanguage: \"en-US\",\n\tparameters: [\n\t\tParameter(\n\t\t\tname: \"greeting\", \n\t\t\ttitle: \"Greeting\", \n\t\t\tdescription: \"The greeting to set on the HelloWorld contract\"\n\t\t)\n\t],\n)\ntransaction(greeting: String) {\n\n  prepare(acct: &Account) {\n    log(acct.address)\n  }\n\n  execute {\n    HelloWorld.updateGreeting(newGreeting: greeting)\n  }\n}\n",
			"network_pins": []
		},
		"dependencies": [
			{
				"contracts": [
					{
						"contract": "HelloWorld",
						"networks": [
							{
								"network": "mainnet",
								"address": "e15193734357cf5c"
							}
						]
					}
				]
			}
		],
		"parameters": []
	}
}
`

var pragmaWithParameters = `
#interaction(
	version: "1.1.0",
	title: "Update Greeting",
	description: "Update the greeting on the HelloWorld contract",
	language: "en-US",
	parameters: [
		Parameter(
			name: "greeting",
			title: "Greeting",
			description: "The greeting to set on the HelloWorld contract"
		),
		Parameter(
			name: "amount",
			title: "Amount",
			description: "The amount parameter to Test"
		)
	],
)

import "HelloWorld"

transaction(greeting: String, amount: UFix64) {

	prepare(acct: &Account) {
		log(acct.address)
	}

	execute {
		HelloWorld.updateGreeting(newGreeting: greeting)
	}
}
`

var pragmaWithoutParameters = `
#interaction(
	version: "1.1.0",
	title: "Update Greeting",
	description: "Update the greeting on the HelloWorld contract",
	language: "en-US",
)

import "HelloWorld"

transaction(greeting: String) {

	prepare(acct: &Account) {
		log(acct.address)
	}

	execute {
		HelloWorld.updateGreeting(newGreeting: greeting)
	}
}
`

var pragmaMinimum = `
#interaction(
	version: "1.1.0",
)

import "HelloWorld"

transaction(greeting: String) {

	prepare(acct: &Account) {
		log(acct.address)
	}

	execute {
		HelloWorld.updateGreeting(newGreeting: greeting)
	}
}
`

var PragmaEmpty = `
import "HelloWorld"
transaction(greeting: String) {

	prepare(acct: &Account) {
		log(acct.address)
	}

	execute {
		HelloWorld.updateGreeting(newGreeting: greeting)
	}
}
`

func TestParsePragma(t *testing.T) {
	assert := assert.New(t)

	tests := []struct {
		name    string
		wantErr bool
		code    string
	}{
		{
			name:    "WithParameters",
			wantErr: false,
			code:    pragmaWithParameters,
		},
		{
			name:    "WithoutParameters",
			wantErr: false,
			code:    pragmaWithoutParameters,
		},
		{
			name:    "Minimum",
			wantErr: false,
			code:    pragmaMinimum,
		},
		{
			name:    "Empty",
			wantErr: false,
			code:    PragmaEmpty,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			codeBytes := []byte(tt.code)
			program, err := parser.ParseProgram(nil, codeBytes, parser.Config{})
			if err != nil {
				t.Fatal(err)
			}
			var template InteractionTemplate
			if err != nil {
				t.Fatal(err)
			}

			err = template.ParsePragma(program)
			if err != nil && !tt.wantErr {
				t.Fatal(err)
			}
			if err != nil {
				t.Fatal(err)
			}
			prettyJSON, err := json.MarshalIndent(template, "", "    ")
			assert.NoError(err, "marshal template to json should not return an error")
			autogold.ExpectFile(t, string(prettyJSON))

		})
	}
}

func TestGenerateParametersScripts(t *testing.T) {
	templateString := `
	{
		"f_type": "InteractionTemplate",
		"f_version": "1.1.0",
		"id": "",
		"data":
		{
			"type": "script",
			"interface": "",
			"messages": [
			{
				"key": "title",
				"i18n": [
					{
						"tag": "en-US",
						"translation": "User Balance"
					}
				]
			},
			{
				"key": "description",
				"i18n": [
					{
						"tag": "en-US",
						"translation": "Get User Balance"
					}
				]
			}],
			"cadence": {},
			"dependencies":	[],
			"parameters": [
            {
                "label": "address",
                "index": 0,
                "type": "Address",
                "messages": [
                    {
                        "key": "title",
                        "i18n": [
                            {
                                "tag": "en-US",
                                "translation": "User Address"
                            }
                        ]
                    },
                    {
                        "key": "description",
                        "i18n": [
                            {
                                "tag": "en-US",
                                "translation": "Get user balance"
                            }
                        ]
                    }
                ]
            }
        ]
		}
	}`

	cadence := `
		import "FungibleToken"
		import "FlowToken"

		access(all)
		fun main(address: Address): UFix64 {
			let account = getAccount(address)
			let vaultRef = account.capabilities
				.borrow<&FlowToken.Vault>(/public/flowTokenBalance)
				?? panic("Could not borrow balance reference to the Vault")

			return vaultRef.balance
		}
	`

	codeBytes := []byte(cadence)
	program, err := parser.ParseProgram(nil, codeBytes, parser.Config{})
	if err != nil {
		t.Errorf("ParseProgram() err %v", err)
	}

	template, err := ParseFlix(templateString)
	if err != nil {
		t.Errorf("ParseFlix() err %v", err)
	}
	err = template.ProcessParameters(program)
	if err != nil {
		t.Errorf("process parameters err %v", err)
	}
	prettyJSON, err := json.MarshalIndent(template, "", "    ")
	if err != nil {
		t.Errorf("process parameters err %v", err)
	}

	autogold.ExpectFile(t, string(prettyJSON))
}

func TestGenerateTemplateIdWithDeps(t *testing.T) {
	templateId := "2ced1d914c5ffdcfb6e772e57e2c19d9ce467e0aa1f1c56fe4302f722896f22b"
	code := `
	{
		"f_type": "InteractionTemplate",
		"f_version": "1.1.0",
		"id": "",
		"data": {
			"type": "transaction",
			"interface": "",
			"messages": [
				{
					"key": "title",
					"i18n": [
						{
							"tag": "en-US",
							"translation": "Update Greeting"
						}
					]
				},
				{
					"key": "description",
					"i18n": [
						{
							"tag": "en-US",
							"translation": "Update the greeting on the HelloWorld contract"
						}
					]
				}
			],
			"cadence": {
				"body": "import \"HelloWorld\"\n\n#interaction (\n  version: \"1.1.0\",\n\ttitle: \"Update Greeting\",\n\tdescription: \"Update the greeting on the HelloWorld contract\",\n\tlanguage: \"en-US\",\n\tparameters: [\n\t\tParameter(\n\t\t\tname: \"greeting\", \n\t\t\ttitle: \"Greeting\", \n\t\t\tdescription: \"The greeting to set on the HelloWorld contract\"\n\t\t)\n\t],\n)\ntransaction(greeting: String) {\n\n  prepare(acct: &Account) {\n    log(acct.address)\n  }\n\n  execute {\n    HelloWorld.updateGreeting(newGreeting: greeting)\n  }\n}\n",
				"network_pins": [
					{
						"network": "testnet",
						"pin_self": "f61e68b5ba6987aaee393401889d5410b01ffa603a66952307319ea09fd505e7"
					}
				]
			},
			"dependencies": [
				{
					"contracts": [
						{
							"contract": "HelloWorld",
							"networks": [
								{
									"network": "testnet",
									"address": "0xe15193734357cf5c",
									"dependency_pin_block_height": 139331034,
									"dependency_pin": {
										"pin": "38b038a23c5975f90a797d6a821f9a8c4e4325a661f92513aedd73fda0e3300c",
										"pin_self": "a06b3cd29330a3c22df3ac2383653e89c249c5e773fd4bbee73c45ea10294b97",
										"pin_contract_name": "HelloWorld",
										"pin_contract_address": "0xe15193734357cf5c",
										"imports": [
											{
												"pin": "3efc62adadbb1dedab0716ac031066a431cd7d627bc1b9260dd08a5a67b26b55",
												"pin_self": "403cd82df774d247bc1fd7471e5ef1fdb7e2e0cb8ec44dce3af5473627179f9a",
												"pin_contract_name": "GiveNumber",
												"pin_contract_address": "0xe15193734357cf5c",
												"imports": []
											}
										]
									}
								}
							]
						}
					]
				}
			],
			"parameters": [
				{
					"label": "greeting",
					"index": 0,
					"type": "String",
					"messages": [
						{
							"key": "title",
							"i18n": [
								{
									"tag": "en-US",
									"translation": "Greeting"
								}
							]
						},
						{
							"key": "description",
							"i18n": [
								{
									"tag": "en-US",
									"translation": "The greeting to set on the HelloWorld contract"
								}
							]
						}
					]
				}
			]
		}
	}`

	template, err := ParseFlix(code)
	if err != nil {
		t.Errorf("ParseFlix() err %v", err)
	}
	id, err := GenerateFlixID(template)
	if err != nil {
		t.Errorf("GenerateFlixID err %v", err)
	}

	if id != templateId {
		t.Errorf("GenerateFlixID got = %v, want %v", id, templateId)
	}

}

func TestEmptyPragmaWithParameter(t *testing.T) {
	codeBytes := []byte(PragmaEmpty)
	program, err := parser.ParseProgram(nil, codeBytes, parser.Config{})
	if err != nil {
		t.Fatal(err)
	}
	var template InteractionTemplate
	err = template.ProcessParameters(program)
	if err != nil {
		t.Errorf("process parameters err %v", err)
	}
	prettyJSON, err := json.MarshalIndent(template, "", "    ")
	if err != nil {
		t.Errorf("process parameters err %v", err)
	}

	autogold.ExpectFile(t, string(prettyJSON))

}

const template = `
{
	"f_type": "InteractionTemplate",
	"f_version": "1.1.0",
	"id": "a2b2d73def...aabc5472d2",
	"data": {
	  "type": "transaction",
	  "interface": "asadf23234...fas234234",
	  "messages": [
		{
		  "key": "title",
		  "i18n": [
			{
			  "tag": "en-US",
			  "translation": "Transfer FLOW"
			},
			{
			  "tag": "fr-FR",
			  "translation": "FLOW de transfert"
			},
			{
			  "tag": "zh-CN",
			  "translation": "转移流程"
			}
		  ]
		},
		{
		  "key": "description",
		  "i18n": [
			{
			  "tag": "en-US",
			  "translation": "Transfer {amount} FLOW to {to}"
			},
			{
			  "tag": "fr-FR",
			  "translation": "Transférez {amount} FLOW à {to}"
			},
			{
			  "tag": "zh-CN",
			  "translation": "将 {amount} FLOW 转移到 {to}"
			}
		  ]
		},
		{
		  "key": "signer",
		  "i18n": [
			{
			  "tag": "en-US",
			  "translation": "Sign this message to transfer FLOW"
			},
			{
			  "tag": "fr-FR",
			  "translation": "Signez ce message pour transférer FLOW."
			},
			{
			  "tag": "zh-CN",
			  "translation": "签署此消息以转移FLOW。"
			}
		  ]
		}
	  ],
	  "cadence": {
		"body": "import \"FlowToken\"\n        transaction(amount: UFix64, to: Address) {\n            let vault: @FungibleToken.Vault\n            prepare(signer: auth(Storage) &Account) {\n                %%self.vault <- signer.storage\n                .borrow<&{FungibleToken.Provider}>(from: /storage/flowTokenVault)!\n                .withdraw(amount: amount)\n                self.vault <- FungibleToken.getVault(signer)\n            }\n            execute {\n                getAccount(to).capabilities\n                .borrow<&{FungibleToken.Receiver}>(/public/flowTokenReceiver)!\n                .deposit(from: <-self.vault)\n            }\n        }",
		"network_pins": [
		  {
			"network": "mainnet",
			"pin_self": "186e262ce6fe06b5075ec6569a0e5482a79c471881182612d8e4a665c2977f3e"
		  },
		  {
			"network": "testnet",
			"pin_self": "f93977d7a297f559e97259cb2a95fed0f87cfeec46c5257a26adc26a260d6c4c"
		  }
		]
	  },
	  "dependencies": [
		{
		  "contracts": [
			{
			  "contract": "FlowToken",
			  "networks": [
				{
				  "network": "mainnet",
				  "address": "0x1654653399040a61",
				  "dependency_pin_block_height": 10123123123,
				  "dependency_pin": {
					"pin": "c8cb7cc7a1c2a329de65d83455016bc3a9b53f9668c74ef555032804bac0b25b",
					"pin_self": "38d0cca4b74c4e88213df636b4cfc2eb6e86fd8b2b84579d3b9bffab3e0b1fcb",
					"pin_contract_name": "FlowToken",
					"imports": [
					  {
						"pin": "b8a3ed26c222ed67016a28021d8fee5603b948533cbc992b3c90f71a61b2b312",
						"pin_self": "7bc3056ba5d39d130f45411c2c05bb549db8ce727c11a1cb821136a621be27fb",
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
				  "dependency_pin_block_height": 10123123123,
				  "dependency_pin": {
					"pin": "c8cb7cc7a1c2a329de65d83455016bc3a9b53f9668c74ef555032804bac0b25b",
					"pin_self": "38d0cca4b74c4e88213df636b4cfc2eb6e86fd8b2b84579d3b9bffab3e0b1fcb",
					"pin_contract_name": "FlowToken",
					"pin_contract_address": "0x7e60df042a9c0868",
					"imports": [
					  {
						"pin": "b8a3ed26c222ed67016a28021d8fee5603b948533cbc992b3c90f71a61b2b312",
						"pin_self": "7bc3056ba5d39d130f45411c2c05bb549db8ce727c11a1cb821136a621be27fb",
						"pin_contract_name": "FungibleToken",
						"pin_contract_address": "0x9a0766d93b6608b7",
						"imports": []
					  }
					]
				  }
				}
			  ]
			}
		  ]
		}
	  ],
	  "parameters": [
		{
		  "label": "amount",
		  "index": 0,
		  "type": "UFix64",
		  "messages": [
			{
			  "key": "title",
			  "i18n": [
				{
				  "tag": "en-US",
				  "translation": "Amount"
				},
				{
				  "tag": "fr-FR",
				  "translation": "Montant"
				},
				{
				  "tag": "zh-CN",
				  "translation": "数量"
				}
			  ]
			},
			{
			  "key": "description",
			  "i18n": [
				{
				  "tag": "en-US",
				  "translation": "Amount of FLOW token to transfer"
				},
				{
				  "tag": "fr-FR",
				  "translation": "Quantité de token FLOW à transférer"
				},
				{
				  "tag": "zh-CN",
				  "translation": "要转移的 FLOW 代币数量"
				}
			  ]
			}
		  ],
		  "balance": "FlowToken"
		},
		{
		  "label": "to",
		  "index": 1,
		  "type": "Address",
		  "messages": [
			{
			  "key": "title",
			  "i18n": [
				{
				  "tag": "en-US",
				  "translation": "To"
				},
				{
				  "tag": "fr-FR",
				  "translation": "Pour"
				},
				{
				  "tag": "zh-CN",
				  "translation": "到"
				}
			  ]
			},
			{
			  "key": "description",
			  "i18n": [
				{
				  "tag": "en-US",
				  "translation": "Amount of FLOW token to transfer"
				},
				{
				  "tag": "fr-FR",
				  "translation": "Le compte vers lequel transférer les jetons FLOW"
				},
				{
				  "tag": "zh-CN",
				  "translation": "将 FLOW 代币转移到的帐户"
				}
			  ]
			}
		  ]
		}
	  ]
	}
  }
`
const templateMultipleImports = `
{
    "f_type": "InteractionTemplate",
    "f_version": "1.1.0",
    "id": "29d03aafbbb5a02e0d5f4ffee685c12494915410812305c2858008d3e2902b72",
    "data": {
        "type": "script",
        "interface": "",
        "messages": null,
        "cadence": {
            "body": "import \"FungibleToken\"\nimport \"FlowToken\"\n\naccess(all)\nfun main(address: Address): UFix64 {\n    let account = getAccount(address)\n\n    let vaultRef = account.capabilities\n        .borrow<&FlowToken.Vault>(/public/flowTokenBalance)\n        ?? panic(\"Could not borrow balance reference to the Vault\")\n\n    return vaultRef.balance\n}\n",
            "network_pins": [
                {
                    "network": "mainnet",
                    "pin_self": "c9aef2c441b2ff0e1a724fcd72f7a48ae7fbbba3c6e72c530607a90ea0fdf93a"
                },
                {
                    "network": "testnet",
                    "pin_self": "74331585cf3df9cd60e6570566d079f97b3e28b0e2156a06731e73e492fe120e"
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
                                "dependency_pin_block_height": 67669170,
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
                                "dependency_pin_block_height": 139547221,
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
                                "dependency_pin_block_height": 67669170,
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
                                "dependency_pin_block_height": 139547221,
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
const templateMissing = `
{
    "f_type": "InteractionTemplate",
    "f_version": "1.1.0",
    "id": "a2b2d73def...aabc5472d2",
    "data": {
        "type": "transaction",
        "interface": "asadf23234...fas234234",
        "messages": [],
        "cadence": {
            "body": "import \"FlowTokenAA\"\n        transaction(amount: UFix64, to: Address) {\n            let vault: @FungibleToken.Vault\n            prepare(signer: auth(Storage) &Account) {\n                %%self.vault <- signer.storage\n                .borrow<&{FungibleToken.Provider}>(from: /storage/flowTokenVault)!\n                .withdraw(amount: amount)\n                self.vault <- FungibleToken.getVault(signer)\n            }\n            execute {\n                getAccount(to).capabilities\n                .borrow<&{FungibleToken.Receiver}>(/public/flowTokenReceiver)!\n                .deposit(from: <-self.vault)\n            }\n        }",
            "network_pins": []
        },
        "dependencies": [
            {
                "contracts": [],
                "parameters": []
            }
        ]
    }
}
`

func TestGetAndReplaceCadenceImports(t *testing.T) {
	assert := assert.New(t)

	tests := []struct {
		name       string
		network    string
		wantErr    bool
		wantImport string
		template   string
	}{
		{
			name:       "Mainnet",
			network:    "mainnet",
			wantErr:    false,
			wantImport: "import FlowToken from 0x1654653399040a61",
			template:   template,
		},
		{
			name:       "Testnet",
			network:    "testnet",
			wantErr:    false,
			wantImport: "import FlowToken from 0x7e60df042a9c0868",
			template:   template,
		},
		{
			name:     "MissingNetwork",
			network:  "missing",
			wantErr:  true,
			template: template,
		},
		{
			name:     "MissingCadence",
			network:  "mainnet",
			wantErr:  true,
			template: templateMissing,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			parsedTemplate, err := ParseFlix(tt.template)
			if err != nil {
				t.Fatal(err)
			}

			cadenceCode, err := parsedTemplate.ReplaceCadenceImports(tt.network)
			if tt.wantErr {
				assert.Error(err, tt.name, "GetCadenceWithReplacedImports should return an error")
			} else {
				assert.NoError(err, "GetCadenceWithReplacedImports should not return an error")
				assert.NotEmpty(cadenceCode, tt.name, "Cadence should not be empty")

				assert.Contains(cadenceCode, tt.wantImport, "Cadence should contain the expected import")
			}
		})
	}
}

func TestGetAndReplaceCadenceImportsMultipleImports(t *testing.T) {
	template, err := ParseFlix(templateMultipleImports)
	if err != nil {
		t.Fatal(err)
	}
	cadenceCode, err := template.ReplaceCadenceImports("mainnet")
	if err != nil {
		t.Fatal(err)
	}
	assert.Contains(t, cadenceCode, "import FungibleToken from 0xf233dcee88fe0abe", "Cadence should contain the expected FungibleToken import")
	assert.Contains(t, cadenceCode, "import FlowToken from 0x1654653399040a61", "Cadence should contain the expected FlowTokenimport")

}

func TestReplaceCadenceImports(t *testing.T) {
	// test replacing import with address missing leading 0x
	template, err := ParseFlix(templateWithDepsMissingLeading0x)
	if err != nil {
		t.Fatal(err)
	}
	cadenceCode, err := template.ReplaceCadenceImports("mainnet")
	if err != nil {
		t.Fatal(err)
	}
	assert.Contains(t, cadenceCode, "import HelloWorld from 0xe15193734357cf5c", "Cadence should contain the expected HelloWorld import with address missing leading 0x")
}
