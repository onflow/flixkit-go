package v1_1

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

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
		"body": "import \"FlowToken\"\n        transaction(amount: UFix64, to: Address) {\n            let vault: @FungibleToken.Vault\n            prepare(signer: AuthAccount) {\n                %%self.vault <- signer\n                .borrow<&{FungibleToken.Provider}>(from: /storage/flowTokenVault)!\n                .withdraw(amount: amount)\n                self.vault <- FungibleToken.getVault(signer)\n            }\n            execute {\n                getAccount(to)\n                .getCapability(/public/flowTokenReceiver)!\n                .borrow<&{FungibleToken.Receiver}>()!\n                .deposit(from: <-self.vault)\n            }\n        }",
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
            "body": "import \"FlowToken\"\n        transaction(amount: UFix64, to: Address) {\n            let vault: @FungibleToken.Vault\n            prepare(signer: AuthAccount) {\n                %%self.vault <- signer\n                .borrow<&{FungibleToken.Provider}>(from: /storage/flowTokenVault)!\n                .withdraw(amount: amount)\n                self.vault <- FungibleToken.getVault(signer)\n            }\n            execute {\n                getAccount(to)\n                .getCapability(/public/flowTokenReceiver)!\n                .borrow<&{FungibleToken.Receiver}>()!\n                .deposit(from: <-self.vault)\n            }\n        }",
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

			cadenceCode, err := parsedTemplate.GetAndReplaceCadenceImports(tt.network)
			if tt.wantErr {
				assert.Error(err, "GetCadenceWithReplacedImports should return an error")
			} else {
				assert.NoError(err, "GetCadenceWithReplacedImports should not return an error")
				assert.NotEmpty(cadenceCode, "Cadence should not be empty")

				assert.Contains(cadenceCode, tt.wantImport, "Cadence should contain the expected import")
			}
		})
	}
}
