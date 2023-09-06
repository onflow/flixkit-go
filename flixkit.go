package flixkit

import (
	"context"
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net/http"
	"os"

	"github.com/onflow/flixkit-go/bindings"
	"github.com/onflow/flixkit-go/common"
)


type FlixService interface {
	GetFlixRaw(ctx context.Context, templateName string) (string, error)
	GetFlix(ctx context.Context, templateName string) (*common.FlowInteractionTemplate, error)
	GetFlixByIDRaw(ctx context.Context, templateID string) (string, error)
	GetFlixByID(ctx context.Context, templateID string) (*common.FlowInteractionTemplate, error)
	GenFlixBinding(ctx context.Context, templateID string, lang string) (string, error)
}

// OsFileReader is a real implementation that calls os.ReadFile.
type OsFileReader struct{}

func (o OsFileReader) ReadFile(filename string) ([]byte, error) {
	return os.ReadFile(filename)
}

type FileReader interface {
	ReadFile(filename string) ([]byte, error)
}

type flixServiceImpl struct {
	config *Config
}

type Config struct {
	FlixServerURL string
	FileReader   FileReader
}

func NewFlixService(config *Config) FlixService {
	if config.FlixServerURL == "" {
		config.FlixServerURL = "https://flix.flow.com/v1/templates"
	}

	if config.FileReader == nil {
		config.FileReader = OsFileReader{}
	}

	return &flixServiceImpl{
		config: config,
	}
}

func (s *flixServiceImpl) GetFlixRaw(ctx context.Context, templateName string) (string, error) {
	url := fmt.Sprintf("%s?name=%s", s.config.FlixServerURL, templateName)
	return FetchFlixWithContext(ctx, url)
}

func (s *flixServiceImpl) GetFlix(ctx context.Context, templateName string) (*common.FlowInteractionTemplate, error) {
	template, err := s.GetFlixRaw(ctx, templateName)
	if err != nil {
		return nil, err
	}

	parsedTemplate, err := ParseFlix(template)
	if err != nil {
		return nil, err
	}

	return parsedTemplate, nil
}

func (s *flixServiceImpl) GetFlixByIDRaw(ctx context.Context, templateID string) (string, error) {
	url := fmt.Sprintf("%s/%s", s.config.FlixServerURL, templateID)
	return FetchFlixWithContext(ctx, url)
}

func (s *flixServiceImpl) GetFlixByID(ctx context.Context, templateID string) (*common.FlowInteractionTemplate, error) {
	template, err := s.GetFlixByIDRaw(ctx, templateID)
	if err != nil {
		return nil, err
	}

	parsedTemplate, err := ParseFlix(template)
	if err != nil {
		return nil, err
	}

	return parsedTemplate, nil
}


func (s *flixServiceImpl) GenFlixBinding(ctx context.Context, templateLocation string, lang string) (string, error) {
	var template string
	var err error
	isLocal := common.IsLocalTemplate(templateLocation)

	if isLocal {
		template, err = FetchFlixWithContextFromFile(s.config.FileReader, ctx, templateLocation)
	} else {
		template, err = FetchFlixWithContext(ctx, templateLocation)
	}

	if err != nil {
		fmt.Println("can not get flix:", err)
		return "", err
	}

	parsedTemplate, err := ParseFlix(template)
	if err != nil {
		return "", err
	}

	contents, bindingErr := bindings.Generate(lang, parsedTemplate, templateLocation)

	return contents, bindingErr
}


func ParseFlix(template string) (*common.FlowInteractionTemplate, error) {
	var flowTemplate common.FlowInteractionTemplate

	err := json.Unmarshal([]byte(template), &flowTemplate)
	if err != nil {
		return nil, err
	}

	return &flowTemplate, nil
}

func FetchFlixWithContext(ctx context.Context, url string) (string, error) {
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

func FetchFlixWithContextFromFile(reader FileReader, ctx context.Context, url string) (string, error) {
	body, err := reader.ReadFile(url)
	if err != nil {
		log.Fatalf("Failed to read file: %v", err)
	}
	return string(body), nil
}
