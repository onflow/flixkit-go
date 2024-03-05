package flixkit

import (
	"context"
	"fmt"

	"github.com/onflow/flixkit-go/internal"
	"github.com/onflow/flixkit-go/types"
)

type FlixService struct {
	config *types.FlixServiceConfig
}

func NewFlixService(config *types.FlixServiceConfig) *FlixService {
	if config.FlixServerUrl == "" {
		config.FlixServerUrl = internal.DefaultFlixServerUrl
	}

	return &FlixService{
		config: config,
	}
}

func (s *FlixService) GetTemplate(ctx context.Context, flixQuery string) (types.FlixInterface, string, error) {
	// Determine if the template string is flix id, flix path, url, or file path
	templateType := internal.GetType(flixQuery, s.config.FileReader)

	// Get template using proper method for string type above
	switch templateType {
	case internal.FlixPath:
		// TODO:
	case internal.FlixFilePath:
		if s.config.FileReader == nil {
			return nil, flixQuery, fmt.Errorf("file reader not provided")
		}

		return internal.GetFlixByFilePath(flixQuery, s.config.FileReader)
	case internal.FlixId:
		// TODO:
	case internal.FlixJson:
		// TODO:
	case internal.FlixName:
		// TODO:
	case internal.FlixUrl:
		// TODO:
	}

	return nil, flixQuery, fmt.Errorf("invalid template type: %v", templateType)
}

func (s *FlixService) CreateTemplate(ctx context.Context) (types.FlixInterface, error) {
	return nil, nil
}
