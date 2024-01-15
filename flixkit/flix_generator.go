package flixkit

import (
	"context"

	"github.com/onflow/flixkit-go/internal"
)

type FlixGenerator interface {
	// GenerateTemplate returns the generated raw template
	CreateTemplate(ctx context.Context, code string, preFill string) (string, error)
}

type FlixTemplateGeneratorConfig = internal.FlixTemplateGeneratorConfig

// NewFlixGenerator returns a new FlixGenerator
func NewFlixTemplateGenerator(conf FlixTemplateGeneratorConfig) FlixGenerator {
	return internal.NewFlixTemplateGenerator(conf)
}
