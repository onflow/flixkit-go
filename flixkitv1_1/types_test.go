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

			err = ParsePragma(program.PragmaDeclarations(), &template)
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
