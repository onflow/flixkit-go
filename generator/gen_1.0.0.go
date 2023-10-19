package generator

import (
	"errors"
	"regexp"

	"github.com/onflow/cadence/runtime/parser"
	"github.com/onflow/flixkit-go"
	"github.com/onflow/flow-cli/flowkit"
)

type Generator1_0_0 struct{}

func (g Generator1_0_0) Generate(code string, flowJson flowkit.State) (*flixkit.FlowInteractionTemplate, error) {
	template := &flixkit.FlowInteractionTemplate{}

	err := processCadenceCommentBlock(code, template)
	if err != nil {
		return nil, err
	}
	err = processParameters(code, template)
	if err != nil {
		return nil, err
	}
	/*
		err = processDependencies(code, template, flowJson)
		if err != nil {
			return nil, err
		}
	*/
	return template, nil
}

func processDependencies(code string, template *flixkit.FlowInteractionTemplate, flowJson flowkit.State) error {
	// todo: process imports, need flowkit to fill in address and determine if address is placeholder
	contractImports := ExtractImports(code)
	if len(contractImports) > 0 {
		return nil
	}

	return nil
}

func processParameters(code string, template *flixkit.FlowInteractionTemplate) error {
	codeBytes := []byte(code)
	program, err := parser.ParseProgram(nil, codeBytes, parser.Config{})
	if err != nil {
		return errors.New("failed to parse cadence code")
	}

	if program.SoleTransactionDeclaration() != nil {
		if program.SoleTransactionDeclaration().ParameterList != nil {
			for i, param := range program.SoleTransactionDeclaration().ParameterList.Parameters {
				template.Data.Arguments[param.Identifier.String()] = flixkit.Argument{
					Type:  param.TypeAnnotation.String(),
					Index: i,
				}
			}
		}
	}

	return nil
}

func processCadenceCommentBlock(code string, template *flixkit.FlowInteractionTemplate) error {
	template.Data.Arguments = make(flixkit.Arguments)

	commentBlockPattern := regexp.MustCompile(`/\*\*[\s\S]*?@f_version[\s\S]*?\*/`)

	// todo: process imports
	code = commentBlockPattern.FindString(code)

	template.Data.Cadence = code
	// Regular expressions for various properties
	messageTitleRE := regexp.MustCompile(`@message title: (.+)`)
	messageDescRE := regexp.MustCompile(`@message description: (.+)`)
	versionRE := regexp.MustCompile(`@f_version (.+)`)
	langRE := regexp.MustCompile(`@lang (.+)`)
	paramTitleRE := regexp.MustCompile(`@parameter title (\w+): (.+)`)
	paramDescRE := regexp.MustCompile(`@parameter description (\w+): (.+)`)
	balanceRE := regexp.MustCompile(`@balance (\w+): (.+)`)

	// Populate the template with extracted data
	template.FVersion = versionRE.FindStringSubmatch(code)[1]
	template.Data.Messages.Title = &flixkit.Title{I18N: map[string]string{langRE.FindStringSubmatch(code)[1]: messageTitleRE.FindStringSubmatch(code)[1]}}
	template.Data.Messages.Description = &flixkit.Description{I18N: map[string]string{langRE.FindStringSubmatch(code)[1]: messageDescRE.FindStringSubmatch(code)[1]}}

	paramTitleMatches := paramTitleRE.FindAllStringSubmatch(code, -1)
	paramDescMatches := paramDescRE.FindAllStringSubmatch(code, -1)
	balanceMatches := balanceRE.FindAllStringSubmatch(code, -1)

	for i, match := range paramTitleMatches {
		argName := match[1]
		argTitle := match[2]

		template.Data.Arguments[argName] = flixkit.Argument{
			Index: i,
			Messages: flixkit.Messages{
				Title: &flixkit.Title{I18N: map[string]string{langRE.FindStringSubmatch(code)[1]: argTitle}},
			},
		}
	}

	for _, match := range paramDescMatches {
		argName := match[1]
		argDesc := match[2]

		if arg, exists := template.Data.Arguments[argName]; exists {
			arg.Messages.Description = &flixkit.Description{I18N: map[string]string{langRE.FindStringSubmatch(code)[1]: argDesc}}
			template.Data.Arguments[argName] = arg
		}
	}

	for _, match := range balanceMatches {
		argName := match[1]
		balance := match[2]

		if arg, exists := template.Data.Arguments[argName]; exists {
			arg.Balance = balance
			template.Data.Arguments[argName] = arg
		}
	}

	return nil
}
