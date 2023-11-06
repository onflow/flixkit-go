package generator

import (
	"context"
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
	gen := Generator1_0_0{
		deployedContracts: []flixkit.Contracts{},
		testnetClient:     nil,
		mainnetClient:     nil,
	}
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
	generator := Generator1_0_0{
		deployedContracts: contracts,
		testnetClient:     nil,
		mainnetClient:     nil,
	}
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

func TestParseImport(t *testing.T) {
	generator := Generator1_0_0{
		deployedContracts: []flixkit.Contracts{},
		testnetClient:     nil,
		mainnetClient:     nil,
	}
	fungi := flixkit.Contracts{
		"FungibleToken": {
			"mainnet": {
				Address:        "0xf233dcee88fe0abe",
				FqAddress:      "A.0xf233dcee88fe0abe.FungibleToken",
				Contract:       "FungibleToken",
				Pin:            "83c9e3d61d3b5ebf24356a9f17b5b57b12d6d56547abc73e05f820a0ae7d9cf5",
				PinBlockHeight: 34166296,
			},
			"testnet": {
				Address:        "0x9a0766d93b6608b7",
				FqAddress:      "A.0x9a0766d93b6608b7.FungibleToken",
				Contract:       "FungibleToken",
				Pin:            "83c9e3d61d3b5ebf24356a9f17b5b57b12d6d56547abc73e05f820a0ae7d9cf5",
				PinBlockHeight: 74776482,
			},
		},
	}
	tests := []struct {
		cadence string
		want    flixkit.Contracts
	}{
		{
			cadence: `import FungibleToken from 0xFungibleTokenAddress`,
			want:    fungi,
		},
		{
			cadence: `import "FungibleToken"`,
			want:    fungi,
		},
		{
			cadence: `import FungibleToken from 0x9a0766d93b6608b7`,
			want:    fungi,
		},
		{
			cadence: `import "FungibleToken"`,
			want:    fungi,
		},
	}

	for _, tt := range tests {
		t.Run(tt.cadence, func(t *testing.T) {
			got, err := generator.parseImport(context.Background(), tt.cadence, nil)
			if err != nil {
				t.Errorf("parseImport() err %v", err)
			}
			if got == nil {
				t.Errorf("parseImport() got = %v, want %v", got, tt.want)
			}
			prettyJSON, err := json.MarshalIndent(got, "", "    ")
			if err != nil {
				t.Errorf("parseImport() err %v", err)
			}
			autogold.ExpectFile(t, string(prettyJSON))
		})
	}
}
