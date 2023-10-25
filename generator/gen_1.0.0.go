package generator

import (
	"context"
	"errors"
	"fmt"
	"regexp"

	"github.com/onflow/cadence/runtime/ast"
	"github.com/onflow/cadence/runtime/parser"
	"github.com/onflow/flixkit-go"
)

type Generator1_0_0 struct{}

// stubb if parameters are needed to be passed in
func NewGenerator() *Generator1_0_0 {
	return &Generator1_0_0{}
}

func (g Generator1_0_0) Generate(code string) (*flixkit.FlowInteractionTemplate, error) {
	template := &flixkit.FlowInteractionTemplate{}

	codeBytes := []byte(code)
	program, err := parser.ParseProgram(nil, codeBytes, parser.Config{})
	if err != nil {
		return nil, err
	}

	err = processParameters(program, code, template)
	if err != nil {
		return nil, err
	}

	err = processCadenceCommentBlock(code, template)
	if err != nil {
		return nil, err
	}

	err = processDependencies(program, code, template)
	if err != nil {
		return nil, err
	}

	err = processTemplateHashes(program, code, template)
	if err != nil {
		return nil, err
	}

	// TODO: coordinate with flix team to synch up on how to generate template hash id
	id, err := generateTemplateId(template)
	if err != nil {
		return nil, err
	}
	template.ID = id

	return template, nil
}

func processTemplateHashes(program *ast.Program, code string, template *flixkit.FlowInteractionTemplate) error {
	return nil
}

func processDependencies(program *ast.Program, code string, template *flixkit.FlowInteractionTemplate) error {
	ctx := context.Background()
	imports := program.ImportDeclarations()
	if len(imports) == 0 {
		return nil
	}
	// fill in dependence information
	deps := make(flixkit.Dependencies, len(imports))
	for _, imp := range imports {
		dep, err := parseImport(ctx, imp.String())
		if err != nil {
			return err
		}
		for contractName, contract := range dep {
			// todo: check if contract is already in template.Data.Dependencies
			// need Placeholder instead of contract name
			deps[contractName] = contract
		}
		template.Data.Dependencies = deps
	}

	// get dep contract
	// get dep pin

	// todo: process imports, use flowJson to fill in contract addresses for networks in template.Data.Dependencies

	return nil
}

func processParameters(program *ast.Program, code string, template *flixkit.FlowInteractionTemplate) error {

	template.Data.Arguments = make(flixkit.Arguments)
	if program.SoleTransactionDeclaration() != nil {
		if program.SoleTransactionDeclaration().ParameterList != nil {
			for i, param := range program.SoleTransactionDeclaration().ParameterList.Parameters {
				template.Data.Arguments[param.Identifier.String()] = flixkit.Argument{
					Type:  param.TypeAnnotation.Type.String(),
					Index: i,
				}
			}
		}
	}

	return nil
}

func processCadenceCommentBlock(cadenceCode string, template *flixkit.FlowInteractionTemplate) error {
	commentBlockPattern := regexp.MustCompile(`/\*\*[\s\S]*?@f_version[\s\S]*?\*/`)
	codeCommentBlock := commentBlockPattern.FindString(cadenceCode)

	versionRE := regexp.MustCompile(`@f_version (.+)`)
	template.FVersion = versionRE.FindStringSubmatch(codeCommentBlock)[1]
	// TODO: determine if there are other values for f_type
	template.FType = "InteractionTemplate"
	// branch logic for version 1.0.0 and future 1.1.0, currently 1.1.0 not supported
	if template.FVersion == "1.1.0" {
		return errors.New("version 1.1.0 not supported")
	}

	template.Data.Cadence = cadenceCode
	fType, err := determineCadenceType(cadenceCode)

	fmt.Println(fType)

	if err != nil {
		return err
	}
	template.Data.Type = fType
	if fType == "interface" {
		template.FType = "InteractionTemplateInterface"
	}

	// Regular expressions for various properties
	messageTitleRE := regexp.MustCompile(`@message title: (.+)`)
	messageDescRE := regexp.MustCompile(`@message description: (.+)`)

	langRE := regexp.MustCompile(`@lang (.+)`)
	paramTitleRE := regexp.MustCompile(`@parameter title (\w+): (.+)`)
	paramDescRE := regexp.MustCompile(`@parameter description (\w+): (.+)`)
	balanceRE := regexp.MustCompile(`@balance (\w+): (.+)`)

	// Populate the template with extracted data

	template.Data.Messages.Title = &flixkit.Title{I18N: map[string]string{langRE.FindStringSubmatch(codeCommentBlock)[1]: messageTitleRE.FindStringSubmatch(codeCommentBlock)[1]}}
	template.Data.Messages.Description = &flixkit.Description{I18N: map[string]string{langRE.FindStringSubmatch(codeCommentBlock)[1]: messageDescRE.FindStringSubmatch(codeCommentBlock)[1]}}

	paramTitleMatches := paramTitleRE.FindAllStringSubmatch(codeCommentBlock, -1)
	paramDescMatches := paramDescRE.FindAllStringSubmatch(codeCommentBlock, -1)
	balanceMatches := balanceRE.FindAllStringSubmatch(codeCommentBlock, -1)

	for _, match := range paramTitleMatches {
		argName := match[1]
		argTitle := match[2]

		if arg, exists := template.Data.Arguments[argName]; exists {
			arg.Messages.Title = &flixkit.Title{I18N: map[string]string{langRE.FindStringSubmatch(codeCommentBlock)[1]: argTitle}}
			template.Data.Arguments[argName] = arg
		}
	}

	for _, match := range paramDescMatches {
		argName := match[1]
		argDesc := match[2]

		if arg, exists := template.Data.Arguments[argName]; exists {
			arg.Messages.Description = &flixkit.Description{I18N: map[string]string{langRE.FindStringSubmatch(codeCommentBlock)[1]: argDesc}}
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
