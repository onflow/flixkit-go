package generator

import (
	"encoding/json"
	"testing"

	"github.com/hexops/autogold/v2"
	"github.com/stretchr/testify/assert"
)

func TestGenerateCommentBlock(t *testing.T) {
	generator := Generator1_0_0{}
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

	template, err := generator.Generate(code)
	assert.NoError(err, "Generate should not return an error")
	prettyJSON, err := json.MarshalIndent(template, "", "    ")
	assert.NoError(err, "marshal template to json should not return an error")
	autogold.ExpectFile(t, string(prettyJSON))
}
