# FlixKit

> FlixKit is a Go package that provides functionalities for interacting with Flow Interaction Templates (aka FLIX). Please note that this package is currently in alpha and may undergo significant changes.

The `flixkit` package is a Go module designed to interact with Flow Interaction Templates (FLIX). It allows users to fetch, parse, generate and create binding files for Flow interaction templates aka FLIX, aka Verified Interactions. 

## Structures

See FLIP that descibes json structure, [Here](https://github.com/onflow/flips/blob/main/application/20230330-interaction-templates-1.1.0.md) current version is v1.1.0

This package provides three functionalities. 
 - Getting network specific Cadence from FLIX
 - Generate FLIX from Cadence
 - Create binding files based on FLIX (javascript, typescript)

The package also defines the following interfaces:

- `FlixService`: This interface defines methods to interact with the FLIX service such as fetching raw data or parsed data by template name or template ID.
- `Generator`: This interface generates FLIX json given Cadence, metadata can be provided in two ways:
   - prefilled out json 
   - Cadence docs in the form of a pragma
- `Bindings`: This interface has two implementations for javascript and typescript using fcl

## Usage

The package provides a `FlixService` interface with a constructor function `NewFlixService(config *Config)`. `Config` contains `FlixServerURL` which should be provided. If no URL is provided, it defaults to `"https://flix.flow.com/v1/templates"`.

The `FlixService` interface provides the following methods:

- `GetTemplate(ctx context.Context, templateName string) (string, error)`: Fetches template and returns as a string.
- `GetAndReplaceCadenceImports(ctx context.Context, templateName string, network string) (*FlowInteractionTemplateExecution, error)`: Fetches and parses a Flix template and provides the cadence for the network provided.

- Note: `templateName` parameter can be the id or name of a template from the interactive template service. A local file or url to the FLIX json file.

Result form GetAndReplaceCadenceImports is a `FlowInteractionTemplateExecution` instance also provides the following methods:

- `IsScript() bool`: Checks if the template is of type "script".
- `IsTransaction() bool`: Checks if the template is of type "transaction".
- `GetAndReplaceCadenceImports(networkName string) (string, error)`: Replaces cadence import statements in the cadence script or transaction with their respective network addresses.

## Examples

Here is a simple example of creating a new FlixService and fetching a template:

```go
flixService := flixkit.NewFlixService(&flixkit.Config{})

template, err := flixService.GetTemplate(context.Background(), "transfer-flow")
if err != nil {
    log.Fatal(err)
}

fmt.Println(template)
```

Note: Remember to replace "transfer-flow" with the actual name of the template you wish to fetch.

To read more about Flow Interaction Templates, [see the docs](https://developers.flow.com/tooling/fcl-js/interaction-templates).


## Bindings

> Binding files are code files that bind consuming code with FLIX. The `bindings` module in Flixkit generates code that calls the FLIX cadence code. FLIX cadence is primarily transactions and scripts. 

### Usage

The `bindings` module has two public methods `Generate` and `NewFclJSGenerator`. `FclJSGenerator` takes a template directory. `bindings` has fcl-js templates.


 - `NewFclJSGenerator() *FclJSGenerator` // uses default fcl-js vanilla javascript
 - `Generate(template string, templateLocation string) (string, error)` // flix is the hydrated template struct, templateLocation is the file location of the flix json file, isLocal is a flag that indicates if the template is local or on remote server

### Example

```go

// uses default fcl-js templates
fclJsGen := flixkit.NewFclJSGenerator() 

output, err := fclJsGen.Generate(template, flixQuery, isLocal)
if err != nil {
    log.Fatal(err)
}

// output is the javascript binding code
fmt.Println(output])

```

## Generate

> Generate creates the newest ratified version of FLIX, as of this update, v1.1 has been passed. Version 1.0.0 will be supported with `FlixService` and `bindings`. 

- `deployedContracts` is an array of v1_1.Contract structs of the contract dependencies the Cadence code depends on, Core contracts are already configured, look in `internal/contracts/core.go` for details

### Example
```go
generator, err := flixkit.NewGenerator(deployedContracts, logger output.Logger)
// preFilledTemplate is a partially populated flix template with human readable messages
// see FLIX flip for details
prettyJSON, err := generator.Generate(ctx, string(code), preFilledTemplate)

fmt.Println(prettyJSON)
```
### Cadence docs pragma

> Using Cadence pragma the metadata can live with the Cadence code. Therefore a prefilled template isn't necessary. More information [Cadence Doc FLIP](https://github.com/onflow/flips/blob/main/application/20230406-interaction-template-cadence-doc.md)

### Example

```go
#interaction(
		version: "1.1.0",
		title: "Transfer Flow",
		description: "Transfer Flow to account",
		language: "en-US",	
	)
#interaction-param-amount(
        title: "Amount", 
        description: "Amount of Flow to transfer",
		language: "en-US",
)	
#interaction-param-to(
		title: "Reciever", 
		description: "Destination address to receive Flow Tokens",
		language: "en-US",
)
	import "FlowToken"
	transaction(amount: UFix64, to: Address) {
		let vault: @FlowToken.Vault
		prepare(signer: AuthAccount) {
		...
		}
	}
`

```

The pragma describes the transaction parameters and reason for the transaction.