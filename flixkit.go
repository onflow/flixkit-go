package flixkit

import (
	"context"
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net/http"
	"net/url"
	"os"
	"strings"

	"github.com/onflow/flixkit-go/bindings"
	"github.com/onflow/flixkit-go/common"
)


type FlixService interface {
	GetFlixRaw(ctx context.Context, templateName string) (string, error)
	GetFlix(ctx context.Context, templateName string) (*common.FlowInteractionTemplate, error)
	GetFlixByIDRaw(ctx context.Context, templateID string) (string, error)
	GetFlixByID(ctx context.Context, templateID string) (*common.FlowInteractionTemplate, error)
	GenFlixBinding(ctx context.Context, templateID string, lang string, tmplReferencePath string) (string, error)
}

type flixServiceImpl struct {
	config *Config
}

type Config struct {
	FlixServerURL string
}

func NewFlixService(config *Config) FlixService {
	if config.FlixServerURL == "" {
		config.FlixServerURL = "https://flix.flow.com/v1/templates"
	}

	return &flixServiceImpl{
		config: config,
	}
}

func (s *flixServiceImpl) GetFlixRaw(ctx context.Context, templateName string) (string, error) {
	url := fmt.Sprintf("%s?name=%s", s.config.FlixServerURL, templateName)
	return FetchFlix(ctx, url)
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
	return FetchFlix(ctx, url)
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


func (s *flixServiceImpl) GenFlixBinding(ctx context.Context, templateLocation string, lang string, tmplReferencePath string) (string, error) {
	template, err := FetchFlix(ctx, templateLocation)
	if err != nil {
		return "", err
	}

	parsedTemplate, err := ParseFlix(template)
	if err != nil {
		return "", err
	}

	contents, bindingErr := bindings.Generate(lang, parsedTemplate, tmplReferencePath);

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

func FetchFlix(ctx context.Context, fileUrl string) (string, error) {
	u, err := url.Parse(fileUrl)

	if err != nil {
		return "", err
	}

	switch u.Scheme {
	case "file":
		return FetchFlixWithContextFromFile(ctx, fileUrl)
	case "http", "https":
		return FetchFlixWithContext(ctx, fileUrl)
	default:
		return "", fmt.Errorf("Unsupported URL scheme", u.Scheme)
	}

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

func FetchFlixWithContextFromFile(ctx context.Context, url string) (string, error) {
	localFilePath := strings.TrimPrefix(url, "file://")

	// Read the file
	body, err := os.ReadFile(localFilePath)
	if err != nil {
		log.Fatalf("Failed to read file: %v", err)
	}
	return string(body), nil
}
