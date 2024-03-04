package flixkit

import (
	"context"
	"fmt"
)

func NewFlixService(config *FlixServiceConfig) *FlixService {
	if config.FlixServerUrl == "" {
		config.FlixServerUrl = DefaultFlixServerUrl
	}

	return &FlixService{
		config: config,
	}
}

func (s *FlixService) GetTemplate(ctx context.Context, flixQuery string) (FlixInterface, string, error) {
	// Determine if the template string is flix id, flix path, url, or file path
	templateType := getType(flixQuery, s.config.FileReader)

	// Get template using proper method for string type above
	switch templateType {
	case FlixPath:
		// TODO:
	case FlixFilePath:
		if s.config.FileReader == nil {
			return nil, flixQuery, fmt.Errorf("file reader not provided")
		}

		return getFlixByFilePath(flixQuery, s.config.FileReader)
	case FlixId:
		// TODO:
	case FlixJson:
		// TODO:
	case FlixName:
		// TODO:
	case FlixUrl:
		// TODO:
	}

	return nil, flixQuery, fmt.Errorf("invalid template type: %v", templateType)
}

func (s *FlixService) CreateTemplate(ctx context.Context) (FlixInterface, error) {
	return nil, nil
}
