package flixkitv2

import (
	"context"

	v1_1 "github.com/onflow/flixkit-go/internal/flixkitv2/v1_1"
)

type FileReader interface {
	ReadFile(path string) ([]byte, error)
}

type FlixServiceConfig struct {
	FlixServerURL string
	FileReader    FileReader
}

func NewFlixService(config *FlixServiceConfig) flixService {
	return flixService{}
}

type flixService struct {
	generator v1_1.Generator
	bindingCreator FclCreator

}


type FlowInteractionTemplateExecution struct {
	Network       string
	Cadence       string
	IsTransaciton bool
	IsScript      bool
}

func (s flixService) GetTemplate(ctx context.Context, templateName string) (string, error) {
	panic("unimplemented")
}

func (flixService) GetTemplateAndReplaceImports(ctx context.Context, templateName string, network string) (*FlowInteractionTemplateExecution, error) {
	panic("unimplemented")
}


func (s flixService) GenerateBinding(ctx context.Context, flixString string, templateLocation string, lang string) (string, error) {
	return s.bindingCreator.Generate(flixString, templateLocation)
}

func (s flixService) GenerateTemplate(ctx context.Context, code string, preFill string) (string, error) {
	return s.generator.GenerateTemplate(ctx, code, preFill)
}
