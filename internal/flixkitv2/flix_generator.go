package flixkitv2

import (
	"context"

	v1_1 "github.com/onflow/flixkit-go/internal/flixkitv2/v1_1"
)

type flixGenerator struct {
	generator  v1_1.Generator
	bindingCreator FclCreator
	fileReader    FileReader
}

type FlixGeneratorConfig struct {
	FileReader    FileReader
}

func NewFlixGenerator(conf FlixGeneratorConfig) flixGenerator{
	return flixGenerator{
		fileReader: conf.FileReader,
	}
}

func (s flixGenerator) GenerateBinding(ctx context.Context, flixString string, templateLocation string, lang string) (string, error) {
	return s.bindingCreator.Generate(flixString, templateLocation)
}

func (s flixGenerator) GenerateTemplate(ctx context.Context, code string, preFill string) (string, error) {
	return s.generator.GenerateTemplate(ctx, code, preFill)
}
