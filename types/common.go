package types

import (
	"github.com/onflow/cadence"
	"github.com/onflow/flowkit/output"
)

type FlixServiceConfig struct {
	FileReader    FileReader
	FlixServerUrl string
	Logger        output.Logger
}

type FlixInterface interface {
	AsCadance(status string, network string) (cadence.Value, error)
	AsJSON() ([]byte, error)
	ReplaceImports()
	CreateBindings() (string, error)
}

type FileReader interface {
	ReadFile(path string) ([]byte, error)
}
