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

type flixServiceImpl struct {
	config *Config
}

var _ FlixService = (*flixServiceImpl)(nil)

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

func (s *flixServiceImpl) GetFlixRaw(ctx context.Context, templateName string) (string, error) {
	url := fmt.Sprintf("%s?name=%s", s.config.FlixServerURL, templateName)
	return FetchFlixWithContext(ctx, url)
}

func (s *flixServiceImpl) GetFlix(ctx context.Context, templateName string) (string, error) {
	template, err := s.GetFlixRaw(ctx, templateName)
	if err != nil {
		return "", err
	}

	return template, nil
}

func (s *flixServiceImpl) GetFlixByIDRaw(ctx context.Context, templateID string) (string, error) {
	url := fmt.Sprintf("%s/%s", s.config.FlixServerURL, templateID)
	return FetchFlixWithContext(ctx, url)
}

func (s *flixServiceImpl) GetFlixByID(ctx context.Context, templateID string) (string, error) {
	template, err := s.GetFlixByIDRaw(ctx, templateID)
	if err != nil {
		return "", err
	}
	return template, nil
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

func (s *flixServiceImpl) GetAndReplaceCadenceImports(ctx context.Context, templateName string, network string) (*FlowInteractionTemplateExecution, error) {
	template, err := s.GetTemplate(ctx, templateName)
	if err != nil {
		return nil, err
	}
	var execution FlowInteractionTemplateExecution
	var cadenceCode string
	ver, err := GetTemplateVersion(template)
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
		cadenceCode, err = replaceableCadence.GetAndReplaceCadenceImports(network)
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
		cadenceCode, err = replaceableCadence.GetAndReplaceCadenceImports(network)
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

type flixQueryTypes string

const (
	flixName flixQueryTypes = "name"
	flixPath flixQueryTypes = "path"
	flixId   flixQueryTypes = "id"
	flixUrl  flixQueryTypes = "url"
	flixJson flixQueryTypes = "json"
)

func isHex(str string) bool {
	if len(str) != 64 {
		return false
	}
	_, err := hex.DecodeString(str)
	return err == nil
}

func isPath(path string) bool {
	_, err := os.Stat(path)
	return err == nil
}

func isUrl(str string) bool {
	u, err := url.Parse(str)
	return err == nil && u.Scheme != "" && u.Host != ""
}

func isJson(str string) bool {
	var js json.RawMessage
	return json.Unmarshal([]byte(str), &js) == nil
}

func getType(s string) flixQueryTypes {
	switch {
	case isPath(s):
		return flixPath
	case isHex(s):
		return flixId
	case isUrl(s):
		return flixUrl
	case isJson(s):
		return flixJson
	default:
		return flixName
	}
}

func (s *flixServiceImpl) GetTemplate(ctx context.Context, flixQuery string) (string, error) {
	var template string
	var err error

	switch getType(flixQuery) {
	case flixId:
		template, err = s.GetFlixByID(ctx, flixQuery)
		if err != nil {
			return "", fmt.Errorf("could not find flix with id %s: %w", flixQuery, err)
		}

	case flixName:
		template, err = s.GetFlix(ctx, flixQuery)
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
		template, err = FetchFlixWithContext(ctx, flixQuery)
		if err != nil {
			return "", fmt.Errorf("could not parse flix from url %s: %w", flixQuery, err)
		}

	default:
		return "", fmt.Errorf("invalid flix query type: %s", flixQuery)
	}

	return template, nil
}
