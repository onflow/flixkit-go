package internal

import (
	"context"
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net/http"
	"strings"

	"github.com/bjartek/underflow"
	"github.com/onflow/cadence"
	v1 "github.com/onflow/flixkit-go/internal/v1"
	v1_1 "github.com/onflow/flixkit-go/internal/v1_1"
	"github.com/onflow/flow-cli/flowkit/output"
)

type FileReader interface {
	ReadFile(path string) ([]byte, error)
}

type FlixServiceConfig struct {
	FlixServerURL string
	FileReader    FileReader
	Logger        output.Logger
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
	config *FlixServiceConfig
}

type FlowInteractionTemplateExecution struct {
	Network       string
	Cadence       string
	IsTransaciton bool
	IsScript      bool
}

type VerCheck struct {
	FVersion string `json:"f_version"`
}

/*
Deployed contracts to network addresses
*/
type NetworkAddressMap = v1_1.NetworkAddressMap

/*
contract name associated with network information
*/
type ContractInfos = v1_1.ContractInfos

type flowInteractionTemplateCadence interface {
	ReplaceCadenceImports(templateName string) (string, error)
	IsTransaction() bool
	IsScript() bool
}

func (s flixService) GetTemplate(ctx context.Context, flixQuery string) (string, string, error) {
	var template string
	source := flixQuery
	var err error

	if flixQuery == "" {
		return "", source, fmt.Errorf("flix query cannot be empty")
	}

	switch getType(flixQuery, s.config.FileReader) {
	case flixId:
		template, source, err = s.getFlixByID(ctx, flixQuery)
		if err != nil {
			return "", source, fmt.Errorf("could not find flix with id %s: %w", flixQuery, err)
		}

	case flixName:
		template, source, err = s.getFlix(ctx, flixQuery)
		if err != nil {
			return "", source, fmt.Errorf("could not find flix with name %s: %w", flixQuery, err)
		}

	case flixPath:
		source = flixQuery
		if s.config.FileReader == nil {
			return "", source, fmt.Errorf("file reader not provided")
		}
		file, err := s.config.FileReader.ReadFile(flixQuery)
		if err != nil {
			return "", source, fmt.Errorf("could not read flix file %s: %w", flixQuery, err)
		}
		template = string(file)
		if err != nil {
			return "", source, fmt.Errorf("could not parse flix from file %s: %w", flixQuery, err)
		}

	case flixUrl:
		template, source, err = fetchFlixWithContext(ctx, flixQuery)
		if err != nil {
			return "", source, fmt.Errorf("could not parse flix from url %s: %w", flixQuery, err)
		}
	case flixJson:
		template = flixQuery
		source = "json"
	default:
		return "", source, fmt.Errorf("invalid flix query type: %s", flixQuery)
	}

	return template, source, nil
}

func (s flixService) GetTemplateAndReplaceImports(ctx context.Context, templateName string, network string) (*FlowInteractionTemplateExecution, error) {
	template, _, err := s.GetTemplate(ctx, templateName)
	if err != nil {
		return nil, err
	}
	var execution FlowInteractionTemplateExecution
	var cadenceCode string
	ver, err := getTemplateVersion(template)
	if err != nil {
		return nil, fmt.Errorf("invalid flix template version, %w", err)
	}
	var replaceableCadence flowInteractionTemplateCadence
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

func (s flixService) GetTemplateAndCreateBinding(ctx context.Context, templateName string, lang string, destFileLocation string) (string, error) {
	template, source, err := s.GetTemplate(ctx, templateName)
	if err != nil {
		return "", err
	}
	language := strings.ToLower(lang)
	var gen *FclCreator

	switch language {
	case "js", "javascript":
		gen = NewFclJSCreator()
	case "ts", "typescript":
		gen = NewFclTSCreator()
	default:
		return "", fmt.Errorf("language %s not supported", lang)
	}

	relativeTemplateLocation := source
	flixType := getType(source, s.config.FileReader)
	if flixType == flixPath && destFileLocation != "" {
		relativeTemplateLocation, err = getRelativePath(templateName, destFileLocation)
		if err != nil {
			return "", err
		}
	}

	return gen.Create(template, relativeTemplateLocation)
}

func (s flixService) GetTemplateAsCadanceValue(ctx context.Context, templateName string, network string) (cadence.Value, error) {
	tmplString, _, err := s.GetTemplate(ctx, templateName)
	if err != nil {
		return nil, err
	}

	// Check FLIX version
	var ver VerCheck
	err = json.Unmarshal([]byte(tmplString), &ver)
	if err != nil {
		return nil, err
	}

	switch ver.FVersion {
	case "1.0.0":
		var v1Template v1.FlowInteractionTemplate
		return unmarshalToCadenceValue(tmplString, v1Template, network)
	case "1.1.0":
		var v1_1Template v1_1.InteractionTemplate
		return unmarshalToCadenceValue(tmplString, v1_1Template, network)
	}

	return nil, nil
}

func (s flixService) CreateTemplate(ctx context.Context, deployedContracts ContractInfos, code string, preFill string) (string, error) {
	template, _, _ := s.GetTemplate(ctx, preFill)
	var gen *v1_1.Generator
	var err2 error
	gen, err2 = v1_1.NewTemplateGenerator(deployedContracts, s.config.Logger)
	if err2 != nil {
		return "", err2
	}
	return gen.CreateTemplate(ctx, code, template)
}

func (s flixService) getFlixRaw(ctx context.Context, templateName string) (string, string, error) {
	url := fmt.Sprintf("%s?name=%s", s.config.FlixServerURL, templateName)
	return fetchFlixWithContext(ctx, url)
}

func (s flixService) getFlix(ctx context.Context, templateName string) (string, string, error) {
	template, url, err := s.getFlixRaw(ctx, templateName)
	if err != nil {
		return "", url, err
	}

	return template, url, nil
}

func (s flixService) getFlixByIDRaw(ctx context.Context, templateID string) (string, string, error) {
	url := fmt.Sprintf("%s/%s", s.config.FlixServerURL, templateID)
	return fetchFlixWithContext(ctx, url)
}

func (s flixService) getFlixByID(ctx context.Context, templateID string) (string, string, error) {
	template, url, err := s.getFlixByIDRaw(ctx, templateID)
	if err != nil {
		return "", url, err
	}
	return template, url, nil
}

func fetchFlixWithContext(ctx context.Context, url string) (template string, templateUrl string, err error) {
	req, err := http.NewRequestWithContext(ctx, http.MethodGet, url, nil)
	if err != nil {
		return "", url, err
	}

	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		return "", url, err
	}
	defer func() {
		if err := resp.Body.Close(); err != nil {
			log.Printf("Warning: error while closing the response body: %v", err)
		}
	}()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return "", url, err
	}
	template = string(body)
	return template, url, err
}

func unmarshalToCadenceValue(template string, v interface{}, network string) (cadence.Value, error) {
	err := json.Unmarshal([]byte(template), &v)
	if err != nil {
		return nil, err
	}

	switch v.(type) {
	case v1.FlowInteractionTemplate:
		return underflow.InputToCadence(v, func(s string) (string, error) {
			fmt.Println(s)
			return "A.123.Foo.Bar", nil
		})
	case v1_1.InteractionTemplate:
		return underflow.InputToCadence(v, func(s string) (string, error) {
			fmt.Println(s)
			return "A.123.Foo.Bar", nil
		})
	}

	return nil, nil
}
