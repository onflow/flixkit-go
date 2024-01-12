package flixkitv2

import (
	"context"

	"github.com/onflow/flixkit-go/internal"
)

// FlixService is the interface for the flix service
type FlixService interface {
	// GetTemplate returns the raw flix template
	GetTemplate(ctx context.Context, templateName string) (string, error)
	// GetAndReplaceCadenceImports returns the raw flix template with cadence imports replaced
	GetAndReplaceCadenceImports(ctx context.Context, templateName string, network string) (string, error)
	// GenerateTemplate returns the generated template 
	GenerateTemplate(ctx context.Context, code string, preFill string) (string, error)
	// GenerateBinding returns the generated binding
	GenerateBinding(ctx context.Context, flixString string, templateLocation string) (string, error)
}

type Config = internal.FlixServiceConfig 

func NewFlixService(config *Config) FlixService {
	return internal.NewFlixService(config)
}