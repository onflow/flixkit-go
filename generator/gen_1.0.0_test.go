package generator

import (
	"encoding/json"
	"testing"

	"github.com/hexops/autogold/v2"
	"github.com/onflow/flixkit-go"
	"github.com/stretchr/testify/assert"
)

func TestGenerateCommentBlock(t *testing.T) {
	assert := assert.New(t)

	code := `
	import FungibleToken from 0xFungibleTokenAddress
	import FlowToken from 0xFlowTokenAddress
	
	/**
	@f_version 1.0.0
	@lang en-US

	@message title: Transfer Tokens
	@message description: Transfer tokens from one account to another
	
	@parameter title amount: Amount
	@parameter title to: To
	@parameter description amount: The amount of FLOW tokens to send
	@parameter description to: The Flow account the tokens will go to
	
	@balance amount: FlowToken
	*/

	transaction(amount: UFix64, to: Address) {
	  let vault: @FungibleToken.Vault
	
	  prepare(signer: AuthAccount) {
		self.vault <- signer
		  .borrow<&{FungibleToken.Provider}>(from: /storage/flowTokenVault)!
		  .withdraw(amount: amount)
	  }
	
	  execute {
		getAccount(to)
		  .getCapability(/public/flowTokenReceiver)!
		  .borrow<&{FungibleToken.Receiver}>()!
		  .deposit(from: <-self.vault)
	  }
	}
`
	gen := NewGenerator([]flixkit.Contracts{})
	template, err := gen.Generate(code)

	assert.NoError(err, "Generate should not return an error")
	prettyJSON, err := json.MarshalIndent(template, "", "    ")
	assert.NoError(err, "marshal template to json should not return an error")
	autogold.ExpectFile(t, string(prettyJSON))
}

func TestScriptGenCommentBlock(t *testing.T) {
	contracts := []flixkit.Contracts{
		{
			"HelloWorld": flixkit.Networks{
				"emulator": flixkit.Network{
					Address:   "0x01cf0e2f2f715450",
					FqAddress: "A.01cf0e2f2f715450.HelloWorld",
					Contract:  "HelloWorld",
				},
			},
		},
	}
	generator := NewGenerator(contracts)
	assert := assert.New(t)

	code := `
	/*
		@f_version 1.0.0
		@lang en-US

		@message title: read greeting
		@message description: read greeting of the HelloWorld smart contract

		@return title greeting: Greeting
	*/
	import "HelloWorld"

	pub fun main(): String {
	return HelloWorld.greeting
	}
`
	template, err := generator.Generate(code)
	assert.NoError(err, "Generate should not return an error")
	prettyJSON, err := json.MarshalIndent(template, "", "    ")
	assert.NoError(err, "marshal template to json should not return an error")
	autogold.ExpectFile(t, string(prettyJSON))

}
