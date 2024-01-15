package internal

import (
	"context"

	v1_1 "github.com/onflow/flixkit-go/internal/v1_1"
	"github.com/onflow/flow-cli/flowkit/output"
)

type flixGenerator struct {
	generator  v1_1.Generator
	fileReader FileReader
}

type FlixTemplateGeneratorConfig struct {
	FileReader        FileReader
	deployedContracts v1_1.ContractInfos
	logger            output.Logger
}

func NewFlixTemplateGenerator(conf FlixTemplateGeneratorConfig) flixGenerator {
	var gen, err = v1_1.NewTemplateGenerator(conf.deployedContracts, conf.logger)
	if err != nil {
		panic(err)
	}

	return flixGenerator{
		generator:  *gen,
		fileReader: conf.FileReader,
	}
}

func (s flixGenerator) CreateTemplate(ctx context.Context, code string, preFill string) (string, error) {
	return s.generator.CreateTemplate(ctx, code, preFill)
}
