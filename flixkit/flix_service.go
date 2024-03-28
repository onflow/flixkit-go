package flixkit

import (
	"context"

	"github.com/onflow/flixkit-go/internal"
	"github.com/onflow/flowkit/v2/config"
)

type FlixService interface {
	// GetTemplate returns the raw flix template
	GetTemplate(ctx context.Context, templateName string) (string, string, error)
	// GetAndReplaceImports returns the raw flix template with cadence imports replaced
	GetTemplateAndReplaceImports(ctx context.Context, templateName string, network string) (*FlowInteractionTemplateExecution, error)
	// GenerateBinding returns the generated binding given the language
	GetTemplateAndCreateBinding(ctx context.Context, templateName string, lang string, destFile string) (string, error)
	// GenerateTemplate returns the generated raw template
	CreateTemplate(ctx context.Context, contractInfos ContractInfos, code string, preFill string, networks []config.Network) (string, error)
}

// FlowInteractionTemplateCadence is the interface returned from Replacing imports, it provides helper methods to assist in executing the resulting Cadence.
type FlowInteractionTemplateExecution = internal.FlowInteractionTemplateExecution

// ContractInfos is an input into generating a template, it is a map of contract name to network information of deployed contracts of the source Cadence code.
type ContractInfos = internal.ContractInfos
type NetworkAddressMap = internal.NetworkAddressMap

// FlixServiceConfig is the configuration for the FlixService that provides a override for FlixServerURL and default values for FileReader and Logger.
type FlixServiceConfig = internal.FlixServiceConfig

// NewFlixService returns a new FlixService given a FlixServiceConfig
func NewFlixService(config *FlixServiceConfig) FlixService {
	return internal.NewFlixService(config)
}
