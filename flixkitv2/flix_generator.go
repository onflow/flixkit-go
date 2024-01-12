package flixkitv2

import (
	"context"

	internal "github.com/onflow/flixkit-go/internal/flixkitv2"
)

type FlixGenerator interface {
	// GenerateTemplate returns the generated raw template 
	GenerateTemplate(ctx context.Context, code string, preFill string) (string, error)
	// GenerateBinding returns the generated binding given the language
	GenerateBinding(ctx context.Context, flixString string, templateLocation string, lang string) (string, error)
}

type FlixGeneratorConfig = internal.FlixGeneratorConfig

// NewFlixGenerator returns a new FlixGenerator
func NewFlixGenerator(conf FlixGeneratorConfig) FlixGenerator {
	return internal.NewFlixGenerator(conf)
}