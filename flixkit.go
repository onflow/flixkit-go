package flixkit

import (
	"context"
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net/http"
	"regexp"
)

type Network struct {
	Address        string `json:"address"`
	FqAddress      string `json:"fq_address"`
	Contract       string `json:"contract"`
	Pin            string `json:"pin"`
	PinBlockHeight uint64 `json:"pin_block_height"`
}

type Argument struct {
	Index    int      `json:"index"`
	Type     string   `json:"type"`
	Messages Messages `json:"messages"`
	Balance  string   `json:"balance"`
}

type Title struct {
	I18N map[string]string `json:"i18n"`
}

type Description struct {
	I18N map[string]string `json:"i18n"`
}

type Messages struct {
	Title       *Title       `json:"title,omitempty"`
	Description *Description `json:"description,omitempty"`
}

type Dependencies map[string]Contracts
type Contracts map[string]Networks
type Networks map[string]Network
type Arguments map[string]Argument

type Data struct {
	Type         string       `json:"type"`
	Interface    string       `json:"interface"`
	Messages     Messages     `json:"messages"`
	Cadence      string       `json:"cadence"`
	Dependencies Dependencies `json:"dependencies"`
	Arguments    Arguments    `json:"arguments"`
}

type FlowInteractionTemplate struct {
	FType    string `json:"f_type"`
	FVersion string `json:"f_version"`
	ID       string `json:"id"`
	Data     Data   `json:"data"`
}

type FlixService interface {
	GetFlixRaw(ctx context.Context, templateName string) (string, error)
	GetFlix(ctx context.Context, templateName string) (*FlowInteractionTemplate, error)
	GetFlixByIDRaw(ctx context.Context, templateID string) (string, error)
	GetFlixByID(ctx context.Context, templateID string) (*FlowInteractionTemplate, error)
}

type flixServiceImpl struct {
	config *Config
}

var _ FlixService = (*flixServiceImpl)(nil)

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
	return FetchFlixWithContext(ctx, url)
}

func (s *flixServiceImpl) GetFlix(ctx context.Context, templateName string) (*FlowInteractionTemplate, error) {
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

func (s *flixServiceImpl) GetFlixByID(ctx context.Context, templateID string) (*FlowInteractionTemplate, error) {
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

func (t *FlowInteractionTemplate) IsScript() bool {
	return t.Data.Type == "script"
}

func (t *FlowInteractionTemplate) IsTransaction() bool {
	return t.Data.Type == "transaction"
}

func (t *FlowInteractionTemplate) GetAndReplaceCadenceImports(networkName string) (string, error) {
	cadence := t.Data.Cadence

	for dependencyAddress, contracts := range t.Data.Dependencies {
		for contractName, networks := range contracts {
			network, ok := networks[networkName]
			if !ok {
				return "", fmt.Errorf("network %s not found for contract %s", networkName, contractName)
			}

			pattern := fmt.Sprintf(`import\s*%s\s*from\s*%s`, contractName, dependencyAddress)
			re, err := regexp.Compile(pattern)
			if err != nil {
				return "", fmt.Errorf("invalid regex pattern: %v", err)
			}

			replacement := fmt.Sprintf("import %s from %s", contractName, network.Address)
			cadence = re.ReplaceAllString(cadence, replacement)
		}
	}

	return cadence, nil
}

func ParseFlix(template string) (*FlowInteractionTemplate, error) {
	var flowTemplate FlowInteractionTemplate

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
