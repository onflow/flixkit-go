package flixkit

import (
	"context"

	"github.com/onflow/flixkit-go/internal"
)

type FlixService interface {
	// GetTemplate returns the raw flix template
	GetTemplate(ctx context.Context, templateName string) (string, string, error)
	// GetAndReplaceImports returns the raw flix template with cadence imports replaced
	GetTemplateAndReplaceImports(ctx context.Context, templateName string, network string) (*FlowInteractionTemplateExecution, error)
	// GenerateBinding returns the generated binding given the language
	GetTemplateAndCreateBinding(ctx context.Context, templateName string, lang string, destFile string) (string, error)
	// GenerateTemplate returns the generated raw template
	CreateTemplate(ctx context.Context, contractInfos ContractInfos, code string, preFill string) (string, error)
}

type FlowInteractionTemplateExecution = internal.FlowInteractionTemplateExecution
type ContractInfos = internal.ContractInfos
type NetworkAddressMap = internal.NetworkAddressMap
type FlixServiceConfig = internal.FlixServiceConfig

func NewFlixService(config *FlixServiceConfig) FlixService {
	return internal.NewFlixService(config)
}
