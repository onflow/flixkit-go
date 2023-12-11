package flixkitv1_1

import (
	"context"
	"testing"

	"github.com/hexops/autogold/v2"
	v1_1 "github.com/onflow/flixkit-go/flixkitv1_1"
	"github.com/stretchr/testify/assert"
)

func TestHelloScript(t *testing.T) {
	contracts := []v1_1.Contract{
		{
			Contract: "HelloWorld",
			Networks: []v1_1.Network{
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
	template, err := generator.Generate(ctx, code, "")
	assert.NoError(err, "Generate should not return an error")
	autogold.ExpectFile(t, template)

}

func TestTransactionValue(t *testing.T) {
	contracts := []v1_1.Contract{
		{
			Contract: "HelloWorld",
			Networks: []v1_1.Network{
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
	template, err := generator.Generate(ctx, code, "")
	assert.NoError(err, "Generate should not return an error")
	autogold.ExpectFile(t, template)

}

func TestTransferFlowTransaction(t *testing.T) {
	contracts := []v1_1.Contract{}

	generator, err := NewGenerator(contracts, nil)
	if err != nil {
		t.Fatal(err)
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
	template, err := generator.Generate(ctx, code, "")
	assert.NoError(err, "Generate should not return an error")
	autogold.ExpectFile(t, template)

}
