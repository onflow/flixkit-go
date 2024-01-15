package internal

import (
	"context"
	"fmt"
	"io"
	"log"
	"net/http"

	v1 "github.com/onflow/flixkit-go/internal/v1"
	v1_1 "github.com/onflow/flixkit-go/internal/v1_1"
)

type FileReader interface {
	ReadFile(path string) ([]byte, error)
}

type FlixServiceConfig struct {
	FlixServerURL string
	FileReader    FileReader
}

func NewFlixService(config *FlixServiceConfig) flixService {
	if config.FlixServerURL == "" {
		config.FlixServerURL = "https://flix.flow.com/v1/templates"
	}

	return flixService{
		config: config,
	}
}

type flixService struct {
	config         *FlixServiceConfig
	bindingCreator FclCreator
}

type FlowInteractionTemplateExecution struct {
	Network       string
	Cadence       string
	IsTransaciton bool
	IsScript      bool
}

type FlowInteractionTemplateCadence interface {
	ReplaceCadenceImports(templateName string) (string, error)
	IsTransaction() bool
	IsScript() bool
}

func (s flixService) GetTemplate(ctx context.Context, flixQuery string) (string, error) {
	var template string
	var err error

	switch getType(flixQuery) {
	case flixId:
		template, err = s.getFlixByID(ctx, flixQuery)
		if err != nil {
			return "", fmt.Errorf("could not find flix with id %s: %w", flixQuery, err)
		}

	case flixName:
		template, err = s.getFlix(ctx, flixQuery)
		if err != nil {
			return "", fmt.Errorf("could not find flix with name %s: %w", flixQuery, err)
		}

	case flixPath:
		if s.config.FileReader == nil {
			return "", fmt.Errorf("file reader not provided")
		}
		file, err := s.config.FileReader.ReadFile(flixQuery)
		if err != nil {
			return "", fmt.Errorf("could not read flix file %s: %w", flixQuery, err)
		}
		template = string(file)
		if err != nil {
			return "", fmt.Errorf("could not parse flix from file %s: %w", flixQuery, err)
		}

	case flixUrl:
		template, err = fetchFlixWithContext(ctx, flixQuery)
		if err != nil {
			return "", fmt.Errorf("could not parse flix from url %s: %w", flixQuery, err)
		}

	default:
		return "", fmt.Errorf("invalid flix query type: %s", flixQuery)
	}

	return template, nil
}

func (s flixService) GetTemplateAndReplaceImports(ctx context.Context, templateName string, network string) (*FlowInteractionTemplateExecution, error) {
	template, err := s.GetTemplate(ctx, templateName)
	if err != nil {
		return nil, err
	}
	var execution FlowInteractionTemplateExecution
	var cadenceCode string
	ver, err := getTemplateVersion(template)
	if err != nil {
		return nil, fmt.Errorf("invalid flix template version, %w", err)
	}
	var replaceableCadence FlowInteractionTemplateCadence
	switch ver {
	case "1.1.0":
		replaceableCadence, err = v1_1.ParseFlix(template)
		if err != nil {
			return nil, err
		}
		cadenceCode, err = replaceableCadence.ReplaceCadenceImports(network)
		if err != nil {
			return nil, err
		}
		execution.Cadence = cadenceCode
		execution.IsScript = replaceableCadence.IsScript()
		execution.IsTransaciton = replaceableCadence.IsTransaction()
	case "1.0.0":
		replaceableCadence, err = v1.ParseFlix(template)
		if err != nil {
			return nil, err
		}
		cadenceCode, err = replaceableCadence.ReplaceCadenceImports(network)
		if err != nil {
			return nil, err
		}
		execution.Cadence = cadenceCode
		execution.IsScript = replaceableCadence.IsScript()
		execution.IsTransaciton = replaceableCadence.IsTransaction()
	default:
		return nil, fmt.Errorf("flix template version: %s not supported", ver)
	}

	if execution.Cadence == "" {
		return nil, fmt.Errorf("could not parse template, invalid flix template")
	}

	return &execution, nil
}

func (s flixService) GetTemplateAndCreateBinding(ctx context.Context, templateName string, lang string) (string, error) {
	template, err := s.GetTemplate(ctx, templateName)
	if err != nil {
		return "", err
	}

	return s.bindingCreator.Generate(template, templateName)
}

func (s flixService) GenerateBinding(ctx context.Context, flixString string, templateLocation string, lang string) (string, error) {
	return s.bindingCreator.Generate(flixString, templateLocation)
}

func (s flixService) getFlixRaw(ctx context.Context, templateName string) (string, error) {
	url := fmt.Sprintf("%s?name=%s", s.config.FlixServerURL, templateName)
	return fetchFlixWithContext(ctx, url)
}

func (s flixService) getFlix(ctx context.Context, templateName string) (string, error) {
	template, err := s.getFlixRaw(ctx, templateName)
	if err != nil {
		return "", err
	}

	return template, nil
}

func (s flixService) getFlixByIDRaw(ctx context.Context, templateID string) (string, error) {
	url := fmt.Sprintf("%s/%s", s.config.FlixServerURL, templateID)
	return fetchFlixWithContext(ctx, url)
}

func (s flixService) getFlixByID(ctx context.Context, templateID string) (string, error) {
	template, err := s.getFlixByIDRaw(ctx, templateID)
	if err != nil {
		return "", err
	}
	return template, nil
}

func fetchFlixWithContext(ctx context.Context, url string) (string, error) {
	req, err := http.NewRequestWithContext(ctx, http.MethodGet, url, nil)
	if err != nil {
		return "", err
	}

	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		return "", err
	}
	defer func() {
		if err := resp.Body.Close(); err != nil {
			log.Printf("Warning: error while closing the response body: %v", err)
		}
	}()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return "", err
	}

	return string(body), nil
}
