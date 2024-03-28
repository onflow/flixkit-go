package flixkit

import (
	"github.com/onflow/flixkit-go/internal"
	filereader "github.com/onflow/flixkit-go/internal/file-reader"
	"github.com/onflow/flixkit-go/types"
)

func NewFlixService(config *types.FlixServiceConfig) types.FlixService {
	if config.FlixServerUrl == "" {
		config.FlixServerUrl = internal.DefaultFlixServerUrl
	}

	if config.FileReader == nil {
		config.FileReader = filereader.GetDefaultFileReader()
	}

	return internal.NewFlixService(config)
}
