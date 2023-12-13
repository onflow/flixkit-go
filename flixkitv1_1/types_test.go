package v1_1

import (
	"encoding/json"
	"testing"

	"github.com/hexops/autogold/v2"
	"github.com/onflow/cadence/runtime/parser"
	"github.com/stretchr/testify/assert"
)

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

	prepare(acct: AuthAccount) {
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

	prepare(acct: AuthAccount) {
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

	prepare(acct: AuthAccount) {
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

	prepare(acct: AuthAccount) {
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

	   pub fun main(address: Address): UFix64 {
	   	let account = getAccount(address)
	   	let vaultRef = account.getCapability(/public/flowTokenBalance)
	   						.borrow<&FlowToken.Vault{FungibleToken.Balance}>()
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
	templateId := "fcada4d7a654a0386a4bb048ac4c851ad7de3945e6e835dc4593581b8c8113da"
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
				"body": "import \"HelloWorld\"\n\n#interaction (\n  version: \"1.1.0\",\n\ttitle: \"Update Greeting\",\n\tdescription: \"Update the greeting on the HelloWorld contract\",\n\tlanguage: \"en-US\",\n\tparameters: [\n\t\tParameter(\n\t\t\tname: \"greeting\", \n\t\t\ttitle: \"Greeting\", \n\t\t\tdescription: \"The greeting to set on the HelloWorld contract\"\n\t\t)\n\t],\n)\ntransaction(greeting: String) {\n\n  prepare(acct: AuthAccount) {\n    log(acct.address)\n  }\n\n  execute {\n    HelloWorld.updateGreeting(newGreeting: greeting)\n  }\n}\n",
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
