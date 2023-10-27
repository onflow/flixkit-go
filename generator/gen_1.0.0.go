package generator

import (
	"context"
	"errors"
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

	withoutImports := stripImports(code)
	codeBytes := []byte(withoutImports)
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

	// need to address this
	// parsing cadence using cadence parser does not like import statements "from 0xPLACEHOLDER"
	err = processDependencies(code, template)
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

func processDependencies(code string, template *flixkit.FlowInteractionTemplate) error {
	ctx := context.Background()
	noCommentsCode := stripComments(code)
	re := regexp.MustCompile(`(?m)^\s*import.*$`)
	imports := re.FindAllString(noCommentsCode, -1)

	if len(imports) == 0 {
		return nil
	}
	// fill in dependence information
	deps := make(flixkit.Dependencies, len(imports))
	for _, imp := range imports {
		dep, err := parseImport(ctx, imp)
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
	commentBlockPattern := regexp.MustCompile(`/\*[\s\S]*?@f_version[\s\S]*?\*/`)
	codeCommentBlock := commentBlockPattern.FindString(cadenceCode)
	template.Data.Cadence = cadenceCode
	fType, err := determineCadenceType(cadenceCode)

	if err != nil {
		return err
	}
	template.Data.Type = fType
	if fType == "interface" {
		template.FType = "InteractionTemplateInterface"
	}

	// no comment block found
	if codeCommentBlock == "" {
		return nil
	}

	versionRE := regexp.MustCompile(`\s*@f_version (.+)`)
	template.FVersion = versionRE.FindStringSubmatch(codeCommentBlock)[1]
	// TODO: determine if there are other values for f_type
	template.FType = "InteractionTemplate"
	// branch logic for version 1.0.0 and future 1.1.0, currently 1.1.0 not supported
	if template.FVersion == "1.1.0" {
		return errors.New("version 1.1.0 not supported")
	}

	// Regular expressions for various properties
	messageTitleRE := regexp.MustCompile(`\s*@message title: (.+)`)
	messageDescRE := regexp.MustCompile(`\s*@message description: (.+)`)

	langRE := regexp.MustCompile(`\s*@lang (.+)`)
	paramTitleRE := regexp.MustCompile(`\s*@parameter title (\w+): (.+)`)
	paramDescRE := regexp.MustCompile(`\s*@parameter description (\w+): (.+)`)
	balanceRE := regexp.MustCompile(`\s*@balance (\w+): (.+)`)

	// Populate the template with extracted data
	if template.Data.Messages.Title == nil {
		template.Data.Messages.Title = &flixkit.Title{}
	}
	if template.Data.Messages.Title.I18N == nil {
		template.Data.Messages.Title.I18N = make(map[string]string)
	}
	if template.Data.Messages.Description == nil {
		template.Data.Messages.Description = &flixkit.Description{}
	}
	if template.Data.Messages.Title.I18N == nil {
		template.Data.Messages.Title.I18N = make(map[string]string)
	}

	langMatch := langRE.FindStringSubmatch(codeCommentBlock)
	if langMatch == nil {
		langMatch = []string{"", "en-US"}
	}
	messageTitleMatch := messageTitleRE.FindStringSubmatch(codeCommentBlock)
	if len(messageTitleMatch) > 0 {
		template.Data.Messages.Title = &flixkit.Title{
			I18N: map[string]string{
				langMatch[1]: messageTitleMatch[1],
			},
		}
	}

	messageDescMatch := messageDescRE.FindStringSubmatch(codeCommentBlock)
	if len(messageDescMatch) > 0 {
		template.Data.Messages.Description = &flixkit.Description{
			I18N: map[string]string{
				langMatch[1]: messageDescMatch[1],
			},
		}
	}

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
