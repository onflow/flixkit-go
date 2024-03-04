package flixkit

import (
	"github.com/onflow/cadence"
	"github.com/onflow/flowkit/output"
)

const (
	DefaultFlixServerUrl = "https://flix.flow.com/v1/templates"
)

type FlixServiceConfig struct {
	FileReader    FileReader
	FlixServerUrl string
	Logger        output.Logger
}

type FlixService struct {
	config *FlixServiceConfig
}

type FlixInterface interface {
	AsCadance() (cadence.Value, error)
	AsJSON() string
	ReplaceImports()
	CreateBindings() (string, error)
}

type FlixQueryTypes string

const (
	FlixName     FlixQueryTypes = "name"
	FlixFilePath FlixQueryTypes = "filePath"
	FlixPath     FlixQueryTypes = "path"
	FlixId       FlixQueryTypes = "id"
	FlixUrl      FlixQueryTypes = "url"
	FlixJson     FlixQueryTypes = "json"
)

type FileReader interface {
	ReadFile(path string) ([]byte, error)
}

type VerCheck struct {
	FVersion string `json:"f_version"`
}
