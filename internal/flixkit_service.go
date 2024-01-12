package internal

import (
	"context"
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
}

func (s flixService) GetTemplate(ctx context.Context, templateName string) (string, error) {
	return "", nil
}

func (s flixService) GetAndReplaceCadenceImports(ctx context.Context, templateName string, network string) (string, error) {
	return "", nil
}