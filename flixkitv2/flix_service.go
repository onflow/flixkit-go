package flixkitv2

import (
	"context"

	internal "github.com/onflow/flixkit-go/internal/flixkitv2"
)

type FlixService interface {
	// GetTemplate returns the raw flix template
	GetTemplate(ctx context.Context, templateName string) (string, error)
	// GetAndReplaceImports returns the raw flix template with cadence imports replaced
	GetTemplateAndReplaceImports(ctx context.Context, templateName string, network string) (*FlowInteractionTemplateExecution, error)
}

type FlowInteractionTemplateExecution = internal.FlowInteractionTemplateExecution

type FlixServiceConfig = internal.FlixServiceConfig 

func NewFlixService(config *FlixServiceConfig) FlixService {
	return internal.NewFlixService(config)
}
