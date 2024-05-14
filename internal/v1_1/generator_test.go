package v1_1

import (
	"bytes"
	"context"
	_ "embed"
	"encoding/hex"
	"testing"

	"github.com/hexops/autogold/v2"
	"github.com/onflow/flixkit-go/internal/contracts"
	"github.com/onflow/flow-go-sdk"
	"github.com/onflow/flow-go-sdk/crypto"
	"github.com/onflow/flowkit"
	"github.com/onflow/flowkit/accounts"
	"github.com/onflow/flowkit/config"
	"github.com/onflow/flowkit/gateway/mocks"
	"github.com/onflow/flowkit/output"
	"github.com/onflow/flowkit/tests"
	"github.com/onflow/go-ethereum/rlp"
	"github.com/spf13/afero"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

func TestHelloScript(t *testing.T) {
	contracts := []Contract{
		{
			Contract: "HelloWorld",
			Networks: []Network{
				{
					Network: "testnet",
					Address: "0xee82856bf20e2aa6",
				},
				{
					Network: "mainnet",
					Address: "0xee82856bf20e2aa6",
				},
				{
					Network: "emulator",
					Address: "0xee82856bf20e2aa6",
				},
			},
		},
	}

	generator := Generator{
		deployedContracts: contracts,
		testnetClient:     nil,
		mainnetClient:     nil,
	}

	assert := assert.New(t)
	code := `
	#interaction(
		version: "1.1.0",
		title: "Say Hello",
		description: "Read the greeting from the HelloWorld contract",
		language: "en-US",
		parameters: [],
	)
	
	import "HelloWorld"

	pub fun main(): String {
	return HelloWorld.greeting
	}
`
	ctx := context.Background()
	template, err := generator.CreateTemplate(ctx, code, "")
	assert.NoError(err, "Generate should not return an error")
	autogold.ExpectFile(t, template)

}

func TestTransactionValue(t *testing.T) {
	contracts := []Contract{
		{
			Contract: "HelloWorld",
			Networks: []Network{
				{
					Network: "testnet",
					Address: "0xee82856bf20e2aa6",
				},
				{
					Network: "mainnet",
					Address: "0xee82856bf20e2aa6",
				},
				{
					Network: "emulator",
					Address: "0xee82856bf20e2aa6",
				},
			},
		},
	}

	generator := Generator{
		deployedContracts: contracts,
		testnetClient:     nil,
		mainnetClient:     nil,
	}
	assert := assert.New(t)
	code := `
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
			)
		],
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
	ctx := context.Background()
	template, err := generator.CreateTemplate(ctx, code, "")
	assert.NoError(err, "Generate should not return an error")
	autogold.ExpectFile(t, template)

}

func TestTransferFlowTransaction(t *testing.T) {
	cs := []Contract{}
	cc := contracts.GetCoreContracts()
	for contractName, c := range cc {
		contract := Contract{
			Contract: contractName,
			Networks: []Network{
				{Network: config.MainnetNetwork.Name, Address: c[config.MainnetNetwork.Name]},
				{Network: config.TestnetNetwork.Name, Address: c[config.TestnetNetwork.Name]},
				{Network: config.EmulatorNetwork.Name, Address: c[config.EmulatorNetwork.Name]},
			},
		}
		cs = append(cs, contract)
	}
	// fill in top level dependencies for the generator
	generator := &Generator{
		deployedContracts: cs,
		testnetClient:     nil,
		mainnetClient:     nil,
	}
	assert := assert.New(t)
	code := `
	#interaction(
		version: "1.1.0",
		title: "Transfer Flow",
		description: "Transfer Flow to account",
		language: "en-US",
		parameters: [
			Parameter(
				name: "amount", 
				title: "Amount", 
				description: "Amount of Flow to transfer"
			),
			Parameter(
				name: "to", 
				title: "Reciever", 
				description: "Destination address to receive Flow Tokens"
			)
		],
	)
	
	import "FlowToken"
	transaction(amount: UFix64, to: Address) {
		let vault: @FlowToken.Vault
		prepare(signer: AuthAccount) {
		
		}
	}
`
	ctx := context.Background()
	template, err := generator.CreateTemplate(ctx, code, "")
	assert.NoError(err, "Generate should not return an error")
	autogold.ExpectFile(t, template)

}

var templateHashedTester = `
	{
		"f_type": "InteractionTemplate",
		"f_version": "1.1.0",
		"id": "3accd8c0bf4c7b543a80287d6c158043b4c2e737c2205dba6e009abbbf1328a4",
		"data": {
		  "type": "transaction",
		  "interface": "",
		  "messages": [
			{
			  "key": "title",
			  "i18n": [
				{
				  "tag": "en-US",
				  "translation": "Transfer Tokens"
				}
			  ]
			},
			{
			  "key": "description",
			  "i18n": [
				{
				  "tag": "en-US",
				  "translation": "Transfer Flow to account"
				}
			  ]
			}
		  ],
		  "cadence": {
			"body": "import \"FungibleToken\"\n\n#interaction(\n    version: \"1.1.0\",\n    title: \"Transfer Flow\",\n    description: \"Transfer Flow to account\",\n    language: \"en-US\",\n    parameters: [\n        Parameter(\n            name: \"amount\", \n            title: \"Amount\", \n            description: \"The amount of FLOW tokens to send\"\n        ),\n        Parameter(\n            name: \"to\", \n            title: \"To\",\n            description: \"The Flow account the tokens will go to\"\n        )\n    ],\n)\n\ntransaction(amount: UFix64, to: Address) {\n    let vault: @FungibleToken.Vault\n    \n    prepare(signer: AuthAccount) {\n        self.vault <- signer\n            .borrow<&{FungibleToken.Provider}>(from: /storage/flowTokenVault)!\n            .withdraw(amount: amount)\n    }\n\n    execute {\n        getAccount(to)\n            .getCapability(/public/flowTokenReceiver)!\n            .borrow<&{FungibleToken.Receiver}>()!\n            .deposit(from: <-self.vault)\n    }\n}",
			"network_pins": [
			  {
				"network": "mainnet",
				"pin_self": "dd046de8ef442e4d708124d5710cb78962eb884a4387df1f0b1daf374bd28278"
			  },
			  {
				"network": "testnet",
				"pin_self": "4089786f5e19fe66b39e347634ca28229851f4de1fd469bd8f327d79510e771f"
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
					  "dependency_pin_block_height": 70493190,
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
					  "dependency_pin_block_height": 149595558,
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
					}
				  ]
				},
				{
				  "key": "description",
				  "i18n": [
					{
					  "tag": "en-US",
					  "translation": "The amount of FLOW tokens to send"
					}
				  ]
				}
			  ]
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
					}
				  ]
				},
				{
				  "key": "description",
				  "i18n": [
					{
					  "tag": "en-US",
					  "translation": "The Flow account the tokens will go to"
					}
				  ]
				}
			  ]
			}
		  ]
		}
	  }
`

func TestGenerateTemplateIdTransaction(t *testing.T) {

	flix, err := ParseFlix(templateHashedTester)
	assert.NoError(t, err)
	assert.Equal(t, "1.1.0", flix.FVersion)

	//	program, err := parser.ParseProgram(nil, []byte(flix.Data.Cadence.Body), parser.Config{})
	assert.NoError(t, err)
	prlp := parametersToRlp(flix.Data.Parameters)

	fullyHashed := getHashedValue(prlp)
	assert.Equal(t, "b5838adae6019bef5455241281621d45543fc27f48815476cddc9b1939de0d76", fullyHashed)

	m := messagesToRlp(flix.Data.Messages)
	hashed := getHashedValue(m)
	assert.Equal(t, "064d5a1903b9ebb4224c3345432cae5aa01650b71408f2518a6458e744afaf6d", hashed)

}

func TestGenerateFlixNetworkRLP(t *testing.T) {

	flix, err := ParseFlix(templateHashedTester)
	assert.NoError(t, err)
	assert.Equal(t, "1.1.0", flix.FVersion)

	//	program, err := parser.ParseProgram(nil, []byte(flix.Data.Cadence.Body), parser.Config{})
	assert.NoError(t, err)
	networkHahes := networksToRlp(flix.Data.Dependencies[0].Contracts[0].Networks)

	fullyHashed := getHashedValue(networkHahes)
	assert.Equal(t, "2d4a4c4f90aab978e297d14ebd9713e37173cefd84ba55fd4d4a6c5bb018ec63", fullyHashed)

	contractHahes := contractsToRlp(flix.Data.Dependencies[0].Contracts)
	contractCollectionHashed := getHashedValue(contractHahes)
	assert.Equal(t, "e12a2186edfc2e0ea06c29195576371591fc09d2c05bdea2048f5a3d674f17c2", contractCollectionHashed)

	idValue, err := flix.EncodeRLP()
	assert.NoError(t, err)
	assert.Equal(t, "3accd8c0bf4c7b543a80287d6c158043b4c2e737c2205dba6e009abbbf1328a4", idValue)
}

//go:embed FungibleToken.cdc
var fungibleTokenContract string

func TestNetworkHashingIds(t *testing.T) {
	contractName := "FungibleToken"
	flix, err := ParseFlix(templateHashedTester)
	assert.NoError(t, err)
	assert.NotNil(t, flix)
	assert.Equal(t, "1.1.0", flix.FVersion)
	var mockFS = afero.NewMemMapFs()
	var rw = afero.Afero{Fs: mockFS}

	_, fKit, gw := setup(rw)
	a := bobMainnet()

	gw.GetAccount.Run(func(args mock.Arguments) {
		addr := args.Get(1).(flow.Address)
		racc := tests.NewAccountWithAddress(addr.String())
		racc.Contracts = map[string][]byte{
			contractName: []byte(fungibleTokenContract),
		}

		gw.GetAccount.Return(racc, nil)
	})

	memoize := make(map[string]PinDetail)
	ctx := context.Background()

	details, err := generateDependencyNetworks(ctx, fKit, a.Address.String(), contractName, memoize, 0)
	assert.NoError(t, err)
	assert.Equal(t, "ac0208f93d07829ec96584d618ddbec6af3cf4e2866bd5071249e8ec93c7e0dc", details.Pin)
	assert.NotNil(t, details)

}

func bobMainnet() *accounts.Account {
	return newAccount("BobMainnet", "0xf233dcee88fe0abe", "seedseedseedseedseedseedseedseedseedseedseedseedAlice")
}

func newAccount(name string, address string, seed string) *accounts.Account {
	privateKey, _ := crypto.GeneratePrivateKey(crypto.ECDSA_P256, []byte(seed))

	return &accounts.Account{
		Name:    name,
		Address: flow.HexToAddress(address),
		Key:     accounts.NewHexKeyFromPrivateKey(0, crypto.SHA3_256, privateKey),
	}
}

func setup(rw flowkit.ReaderWriter) (*flowkit.State, *flowkit.Flowkit, *mocks.TestGateway) {
	state, err := flowkit.Init(rw)
	if err != nil {
		panic(err)
	}
	emulatorServiceAccount, _ := accounts.NewEmulatorAccount(rw, crypto.ECDSA_P256, crypto.SHA3_256, "")
	state.Accounts().AddOrUpdate(emulatorServiceAccount)
	gw := mocks.DefaultMockGateway()
	logger := output.NewStdoutLogger(output.NoneLog)
	flowkit := flowkit.NewFlowkit(state, config.TestnetNetwork, gw.Mock, logger)

	return state, flowkit, gw
}
func getHashedValue(values []interface{}) string {
	input := []interface{}{values}
	var buffer bytes.Buffer
	rlp.Encode(&buffer, input)
	hexString := hex.EncodeToString(buffer.Bytes())

	return ShaHex(hexString, "parameters")

}
