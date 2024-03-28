package internal

import (
	"context"
	"fmt"

	"github.com/onflow/flixkit-go/types"
)

type FlixService struct {
	config *types.FlixServiceConfig
}

func NewFlixService(config *types.FlixServiceConfig) FlixService {
	return FlixService{
		config: config,
	}
}

func (s FlixService) GetTemplate(ctx context.Context, flixQuery string) (types.FlixInterface, string, error) {
	// Determine if the template string is flix id, flix path, url, or file path
	templateType := GetType(flixQuery, s.config.FileReader)

	// Get template using proper method for string type above
	switch templateType {
	case FlixPath:
		return getFlixByFlixPath(flixQuery)
	case FlixFilePath:
		return getFlixByFilePath(flixQuery, s.config.FileReader)
	case FlixId:
		return getFlixById(flixQuery)
	case FlixJson:
		return parseFlixJSON(flixQuery)
	case FlixName:
		return getFlixByName(flixQuery)
	case FlixUrl:
		return getFlixByUrl(flixQuery)
	}

	return nil, flixQuery, fmt.Errorf("invalid template type: %v", templateType)
}

func (s FlixService) VerifyTemplateID(template types.FlixInterface) bool {
	return false
}

func (s FlixService) CreateTemplate(ctx context.Context) (types.FlixInterface, error) {
	return nil, nil
}
