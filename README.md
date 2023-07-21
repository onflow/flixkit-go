# FlixKit

> FlixKit is a Go package that provides functionalities for interacting with Flow Interaction Templates (aka FLIX). Please note that this package is currently in alpha and may undergo significant changes.

The `flixkit` package is a Go module designed to interact with Flow Interaction Templates (FLIX). It allows users to fetch, parse, and handle Flow interaction templates.

## Structures

The package provides a range of structs to represent data fetched from FLIX service:

- `Network`: Contains information about a specific network like address, contract and pin.
- `Argument`: Represents the arguments that can be given to the contracts.
- `Title`, `Description`: Used for i18n (internationalization) purposes in messages.
- `Messages`: Contains a title and a description.
- `Data`: Provides detailed information about the Flow interaction template like type, interface, messages, dependencies and arguments.
- `FlowInteractionTemplate`: The main struct that contains all the details of a flow interaction template.

The package also defines the following interfaces:

- `FlixService`: This interface defines methods to interact with the FLIX service such as fetching raw data or parsed data by template name or template ID.

## Usage

The package provides a `FlixService` interface with a constructor function `NewFlixService(config *Config)`. `Config` contains `FlixServerURL` which should be provided. If no URL is provided, it defaults to `"https://flix.flow.com/v1/templates"`.

The `FlixService` interface provides the following methods:

- `GetFlixRaw(ctx context.Context, templateName string) (string, error)`: Fetches a raw Flix template by its name.
- `GetFlix(ctx context.Context, templateName string) (*FlowInteractionTemplate, error)`: Fetches and parses a Flix template by its name.
- `GetFlixByIDRaw(ctx context.Context, templateID string) (string, error)`: Fetches a raw Flix template by its ID.
- `GetFlixByID(ctx context.Context, templateID string) (*FlowInteractionTemplate, error)`: Fetches and parses a Flix template by its ID.

Each `FlowInteractionTemplate` instance also provides the following methods:

- `IsScript() bool`: Checks if the template is of type "script".
- `IsTransaction() bool`: Checks if the template is of type "transaction".
- `GetAndReplaceCadenceImports(networkName string) (string, error)`: Replaces cadence import statements in the cadence script with their respective network addresses.

## Examples

Here is a simple example of creating a new FlixService and fetching a template:

```go
flixService := flixkit.NewFlixService(&flixkit.Config{})

template, err := flixService.GetFlix(context.Background(), "transfer-flow")
if err != nil {
    log.Fatal(err)
}

fmt.Println(template)
```

Note: Remember to replace "transfer-flow" with the actual name of the template you wish to fetch.

To read more about Flow Interaction Templates, [see the docs](https://developers.flow.com/tooling/fcl-js/interaction-templates).