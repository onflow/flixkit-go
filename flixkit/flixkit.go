package flixkit

import (
	"context"
	"encoding/hex"
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net/http"
	"net/url"
	"os"

	v1 "github.com/onflow/flixkit-go/flixkit/v1"
	v1_1 "github.com/onflow/flixkit-go/flixkit/v1_1"
)

type FlowInteractionTemplateExecution struct {
	Network       string
	Cadence       string
	IsTransaciton bool
	IsScript      bool
}

type FlowInteractionTemplateVersion struct {
	FVersion string `json:"f_version"`
}

type GenerateTemplate interface {
	Generate(ctx context.Context, code string, preFill string) (string, error)
}

type GenerateBinding interface {
	Generate(flixString string, templateLocation string) (string, error)
}

type FlowInteractionTemplateCadence interface {
	GetAndReplaceCadenceImports(templateName string) (string, error)
	IsTransaction() bool
	IsScript() bool
}

type FlixService interface {
	GetTemplate(ctx context.Context, templateName string) (string, error)
	GetAndReplaceCadenceImports(ctx context.Context, templateName string, network string) (*FlowInteractionTemplateExecution, error)
}


type FileReader interface {
	ReadFile(path string) ([]byte, error)
}

type Config struct {
	FlixServerURL string
	FileReader    FileReader
}

func NewFlixService(config *Config) FlixService {
	if config.FlixServerURL == "" {
		config.FlixServerURL = "https://flix.flow.com/v1/templates"
	}

	return &flixServiceImpl{
		config: config,
	}
}


func GetTemplateVersion(template string) (string, error) {
	var flowTemplate FlowInteractionTemplateVersion

	err := json.Unmarshal([]byte(template), &flowTemplate)
	if err != nil {
		return "", err
	}

	if flowTemplate.FVersion == "" {
		return "", fmt.Errorf("version not found")
	}

	return flowTemplate.FVersion, nil
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

