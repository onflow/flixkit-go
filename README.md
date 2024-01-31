# FlixKit

> FlixKit is a Go package that provides functionalities for interacting with Flow Interaction Templates (aka FLIX). This package supports generating v1.1 FLIX template json, creating binding files for v1.0 and v1.1 and replacing import for v1.0 and v1.1. More information about FLIX [FLIX FLIP](https://github.com/onflow/flips/blob/main/application/20230330-interaction-templates-1.1.0.md)

The `flixkit` package is a Go module designed to interact with Flow Interaction Templates (FLIX). It allows users to fetch, parse, generate and create binding files for Flow interaction templates aka FLIX, aka Verified Interactions. 

## Structures

See FLIP that describes json structure, [Here](https://github.com/onflow/flips/blob/main/application/20230330-interaction-templates-1.1.0.md) current version is v1.1.0

This package provides three functionalities. 
 - Getting network specific Cadence from FLIX
 - Generate FLIX from Cadence
 - Create binding files based on FLIX (javascript, typescript)

The package also defines the following interfaces:

- `FlixService`: This interface defines methods to fetch templates using template name or template id (from flix service), URL or local file path. 

### Methods
```go
// GetTemplate returns the raw flix template
GetTemplate(ctx context.Context, templateName string) (string, string, error)
// GetAndReplaceImports returns the raw flix template with cadence imports replaced
GetTemplateAndReplaceImports(ctx context.Context, templateName string, network string) (*FlowInteractionTemplateExecution, error)
// GenerateBinding returns the generated binding given the language
GetTemplateAndCreateBinding(ctx context.Context, templateName string, lang string, destFile string) (string, error)
// GenerateTemplate returns the generated raw template
CreateTemplate(ctx context.Context, contractInfos ContractInfos, code string, preFill string) (string, error)
```

## Usage

The package provides a `FlixService` interface with a constructor function `NewFlixService(config *FlixServiceConfig)`. `FlixServiceConfig`
contains 
 - `FlixServerURL` which is defaulted to `"https://flix.flow.com/v1/templates"`. User can provide their own service url endpoint
 - `FileReader` which is used to read local FLIX json template files
 - `Logger` which is used in creating `flowkit.NewFlowkit` for FLIX template generation

The `FlixService` interface provides the following methods:

- `GetTemplate`: Fetches template and returns as a string.
- `GetTemplateAndReplaceImports` returns `FlowInteractionTemplateExecution`: Fetches and parses a Flix template and provides the cadence for the network provided. There are two helper methods to assist in determining if the Cadence is a transaction or a script.

- Note: `templateName` parameter can be the id or name of a template from the interactive template service. A local file or url to the FLIX json file or the template string itself.

Result form GetAndReplaceCadenceImports is a `FlowInteractionTemplateExecution` instance also provides the following methods:

- `IsScript() bool`: Checks if the template is of type "script".
- `IsTransaction() bool`: Checks if the template is of type "transaction".
- `Cadence`: Replaced cadence with respective network addresses.
- `Network`: Name of network used to get import addresses

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


## Binding Files

> Binding files are client code files used to call Cadence contracts using the scripts or transactions in a FLIX. These client files can be created given a FLIX, currently TypeScript and JavaScript are supported.

### Usage

The `bindings` module has two public methods `Generate` and `NewFclJSGenerator`. `FclJSGenerator` takes a template directory. `bindings` has fcl-js templates.

```go
flixService := flixkit.NewFlixService(&flixkit.Config{
	FileReader: myFileReader
})

binding, err := flixService.GetTemplateAndCreateBinding(context.Background(), "transfer-flow", "js", "./bindingFiles/transferFlow.js")
if err != nil {
    log.Fatal(err)
}

fmt.Println(binding)
```

```go
GetTemplateAndCreateBinding(ctx context.Context, templateName string, lang string, destFile string) (string, error)
```

 - `templateName` value can be template name, template id, url or local file. 
 - `lang` values supported are "js", "javascript", "ts", "typescript" 
 - `destFile` is the location of the destination binding file, this is used to create the relative path if the template is local. If the template is a template name, template id or url `destFile` isn't used

## Generate Templates

> CreateTemplate creates the newest ratified version of FLIX, as of this update, see link to FLIP Flip above for more information. 
```go
	CreateTemplate(ctx context.Context, contractInfos ContractInfos, code string, preFill string) (string, error)
```

### Usage
```go
flixService := flixkit.NewFlixService(&flixkit.Config{
	FileReader: myFileReader,
	Logger: myLogger,
})

prettyJSON, err := flixService.CreateTemplate(ctx, depContracts, string(code), preFilled)

fmt.Println(prettyJSON)
```
- `contractInfos` is an array of v1_1.Contract struct. This provides the network information about the deployed contracts that are dependencies in the FLIX Cadence code.
- `code` is the actual Cadence code the template is based on
- `preFilled` is a partially filled out FLIX template. This can be a template name, template id, url or local file. Alternatively to using a prefilled template, the Cadence itself can provide metadata using a FLIX specific Cadence pragma, more on that below, [See Cadence Doc Flip](https://github.com/onflow/flips/blob/main/application/20230406-interaction-template-cadence-doc.md)


### Cadence docs pragma

> Using Cadence pragma the metadata can exist along with the Cadence code. Therefore a prefilled template isn't necessary

### Example

```go
#interaction(
		version: "1.1.0",
		title: "Transfer Flow",
		description: "Transfer Flow to account",
		language: "en-US",
	)
	
	import "FlowToken"
	transaction(amount: UFix64, to: Address) {
		let vault: @FlowToken.Vault
        prepare(signer: &Account) {
		...
		}
	}
`

```

A `pragma` gives special instruction or processors, in this case FLIX specific information that describes the transaction or script.

Notice: Nested structures in Cadence pragma will be supported in future, this will allow describing parameters
```go
...
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
...		
```
