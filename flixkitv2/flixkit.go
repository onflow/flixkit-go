package flixkitv2

import (
	"context"

	internal "github.com/onflow/flixkit-go/internal/flixkitv2"
)

// FlixService is the interface for the flix service
type FlixService interface {
	// GetTemplate returns the raw flix template
	GetTemplate(ctx context.Context, templateName string) (string, error)
	// GetAndReplaceCadenceImports returns the raw flix template with cadence imports replaced
	GetTemplateAndReplaceImports(ctx context.Context, templateName string, network string) (*FlowInteractionTemplateExecution, error)
}

type FlowInteractionTemplateExecution = internal.FlowInteractionTemplateExecution

type FlixGenerator interface {
	// GenerateTemplate returns the generated template 
	GenerateTemplate(ctx context.Context, code string, preFill string) (string, error)
	// GenerateBinding returns the generated binding given the language
	GenerateBinding(ctx context.Context, flixString string, templateLocation string, lang string) (string, error)

}

type FlixServiceConfig = internal.FlixServiceConfig 

type FlixGeneratorConfig = internal.FlixServiceConfig

func NewFlixService(config *FlixServiceConfig) FlixService {
	return internal.NewFlixService(config)
}

func NewFlixGenaarator(config *FlixGeneratorConfig) FlixGenerator {
	return internal.NewFlixGenerator(config)
}