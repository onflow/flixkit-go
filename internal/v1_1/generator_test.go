package v1_1

import (
	"context"
	_ "embed"
	"testing"

	"github.com/hexops/autogold/v2"
	"github.com/stretchr/testify/assert"
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

	access(all)
	fun main(): String {
		return HelloWorld.greeting
	}
`
	ctx := context.Background()
	template, err := generator.CreateTemplate(ctx, code, "")
	assert.NoError(err, "Generate should not return an error")
	autogold.ExpectFile(t, template)

}

func TestValidImports(t *testing.T) {
	contracts := []Contract{
		{
			Contract: "Alice",
			Networks: []Network{
				{
					Network: "testnet",
					Address: "0x0000000000000001",
				},
				{
					Network: "mainnet",
					Address: "0x0000000000000001",
				},
				{
					Network: "emulator",
					Address: "0x0000000000000001",
				},
			},
		},
		{
			Contract: "Bob",
			Networks: []Network{
				{
					Network: "testnet",
					Address: "0x0000000000000002",
				},
				{
					Network: "mainnet",
					Address: "0x0000000000000002",
				},
				{
					Network: "emulator",
					Address: "0x0000000000000002",
				},
			},
		},
	}

	generator := Generator{
		deployedContracts: contracts,
	}

	assert := assert.New(t)
	code := `
	import "Alice"
	import Bob from 0x0000000000000002

	access(all)
	fun main(): Void {}
`
	ctx := context.Background()
	template, err := generator.CreateTemplate(ctx, code, "")
	assert.NoError(err, "Generate should not return an error")
	autogold.ExpectFile(t, template)

}
func TestInValidImports(t *testing.T) {
	contracts := []Contract{
		{
			Contract: "Alice",
			Networks: []Network{
				{
					Network: "testnet",
					Address: "0x0000000000000001",
				},
				{
					Network: "mainnet",
					Address: "0x0000000000000001",
				},
				{
					Network: "emulator",
					Address: "0x0000000000000001",
				},
			},
		},
		{
			Contract: "Bob",
			Networks: []Network{
				{
					Network: "testnet",
					Address: "0x0000000000000002",
				},
				{
					Network: "mainnet",
					Address: "0x0000000000000002",
				},
				{
					Network: "emulator",
					Address: "0x0000000000000002",
				},
			},
		},
	}

	generator := Generator{
		deployedContracts: contracts,
	}

	assert := assert.New(t)
	code := `
	import "Alice"
	import Bob from 0x0000000000000002
	import "Joe"

	access(all)
	fun main(): Void {}
`
	ctx := context.Background()
	_, err := generator.CreateTemplate(ctx, code, "")
	assert.Error(err, "Generate should not return an error")

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

		prepare(acct: &Account) {
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
	cs := []Contract{
		{
			Contract: "FlowToken",
			Networks: []Network{
				{Network: "mainnet", Address: "0x1654653399040a61"},
				{Network: "testnet", Address: "0x7e60df042a9c0868"},
				{Network: "emulator", Address: "0x0ae53cb6e3f42a79"},
			},
		},
	}
	// fill in top level dependencies for the generator
	generator := &Generator{
		deployedContracts: cs,
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
				title: "Receiver",
				description: "Destination address to receive Flow Tokens"
			)
		],
	)

	import "FlowToken"

	transaction(amount: UFix64, to: Address) {
		let vault: @FlowToken.Vault

		prepare(signer: &Account) {}
	}
`
	ctx := context.Background()
	template, err := generator.CreateTemplate(ctx, code, "")
	assert.NoError(err, "Generate should not return an error")
	autogold.ExpectFile(t, template)

}
