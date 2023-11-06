package generator

import (
	"context"
	"errors"
	"fmt"
	"regexp"
	"strings"

	"github.com/onflow/cadence/runtime/ast"
	"github.com/onflow/cadence/runtime/cmd"
	"github.com/onflow/cadence/runtime/common"
	"github.com/onflow/flixkit-go"
	"github.com/onflow/flow-cli/flowkit"
	"github.com/onflow/flow-go-sdk"
)

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

// Support FLIP - Interaction Template Cadence Doc
// # Interaction Template Cadence Doc (v1.0.0)
// https://github.com/onflow/flips/pull/80
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
	template.FType = "InteractionTemplate"
	// branch logic for version 1.0.0 and future 1.1.0, currently 1.1.0 not supported
	if template.FVersion != "1.0.0" {
		return errors.New("only version 1.0.0 is supported at this time")
	}

	// Regular expressions for various properties
	messageTitleRE := regexp.MustCompile(`\s*@message title: (.+)`)
	messageDescRE := regexp.MustCompile(`\s*@message description: (.+)`)

	langRE := regexp.MustCompile(`\s*@lang (.+)`)
	paramTitleRE := regexp.MustCompile(`\s*@parameter title (\w+): (.+)`)
	paramDescRE := regexp.MustCompile(`\s*@parameter description (\w+): (.+)`)
	balanceRE := regexp.MustCompile(`\s*@balance (\w+): (.+)`)

	langMatch := langRE.FindStringSubmatch(codeCommentBlock)
	if langMatch == nil {
		langMatch = []string{"", "en-US"}
	}
	messageTitleMatch := messageTitleRE.FindStringSubmatch(codeCommentBlock)
	if len(messageTitleMatch) > 0 {
		// Populate the template with extracted data
		if template.Data.Messages.Title == nil {
			template.Data.Messages.Title = &flixkit.Title{
				I18N: map[string]string{
					"en-US": "",
				},
			}
		}
		template.Data.Messages.Title = &flixkit.Title{
			I18N: map[string]string{
				langMatch[1]: messageTitleMatch[1],
			},
		}
	}

	messageDescMatch := messageDescRE.FindStringSubmatch(codeCommentBlock)
	if len(messageDescMatch) > 0 {
		if template.Data.Messages.Description == nil {
			template.Data.Messages.Description = &flixkit.Description{
				I18N: map[string]string{
					"en-US": "",
				},
			}
		}
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

func regexpMatch(pattern, text string) (map[string]string, error) {
	r := regexp.MustCompile(pattern)
	names := r.SubexpNames()
	match := r.FindStringSubmatch(text)
	if match == nil {
		return nil, nil
	}

	m := map[string]string{}
	for i, n := range match {
		m[names[i]] = n
	}

	return m, nil
}

func determineCadenceType(code string) (string, error) {
	// Use regex to match only occurrences not in comments or strings.
	transactionRegex := regexp.MustCompile(`(?s)\btransaction\s*(?:\([^)]*\))?\s*{.*`)
	scriptRegex := regexp.MustCompile(`(?m)^\s*pub\s+fun\s+main\(`)
	interfaceRegex := regexp.MustCompile(`(?m)^\s*(pub|priv)\s+(resource|struct|contract)\s+interface`)

	if transactionRegex.MatchString(code) {
		return "transaction", nil
	} else if scriptRegex.MatchString(code) {
		return "script", nil
	} else if interfaceRegex.MatchString(code) {
		// future support for interface, continue to process code
		return "interface", nil
	}

	return "", errors.New("could not determine if code is transaction or script")
}

func stripComments(cadenceCode string) string {
	// Strip block comments
	blockCommentRe := regexp.MustCompile(`(?s)/\*.*?\*/`)
	cadenceCode = blockCommentRe.ReplaceAllString(cadenceCode, "")

	// Strip single line comments
	singleCommentRe := regexp.MustCompile(`//.*\n?`)
	cadenceCode = singleCommentRe.ReplaceAllString(cadenceCode, "")

	return cadenceCode
}

func stripImports(cadenceCode string) string {
	// Match lines starting with optional leading whitespaces followed by the word "import"
	re := regexp.MustCompile(`(?m)^\s*import.*$\n?`)
	return re.ReplaceAllString(cadenceCode, "")
}

func normalizeImports(cadenceCode string) string {
	// replace new import syntax with old import syntax to be used in templates
	// import "0xNonFungibleTokenAddress" -> import NonFungibleToken from 0xNonFungibleTokenAddress
	// Use a regex pattern to match the new import syntax
	pattern := regexp.MustCompile(`import "(.+?)"`)

	// Replace the matched pattern with the old syntax
	replaced := pattern.ReplaceAllStringFunc(cadenceCode, func(match string) string {
		submatch := pattern.FindStringSubmatch(match)
		if len(submatch) > 1 {
			return fmt.Sprintf(`import %s from 0x%s`, submatch[1], submatch[1])
		}
		return match
	})

	return replaced

}

/*
 Thanks to Overflow, https://github.com/bjartek/overflow/ for all the contract pinning code
*/
// https://github.com/onflow/fcl-js/blob/master/packages/fcl/src/interaction-template-utils/generate-dependency-pin.js
func generateDependentPin(ctx context.Context, flowkit flowkit.Flowkit, address string, name string, cache map[string][]string) ([]string, error) {

	identifier := fmt.Sprintf("A.%s.%s", strings.ReplaceAll(address, "0x", ""), name)
	existingHash, ok := cache[identifier]
	if ok {
		return existingHash, nil
	}

	account, err := flowkit.GetAccount(ctx, flow.HexToAddress(address))
	if err != nil {
		return nil, err
	}
	code := account.Contracts[name]
	imports := getAddressImports(code, name)
	hashes := []string{flixkit.ShaHex(code, "")}

	for _, imp := range imports {
		split := strings.Split(imp, ".")
		address, name := split[0], split[1]
		dep, err := generateDependentPin(ctx, flowkit, address, name, cache)
		if err != nil {
			return nil, err
		}
		hashes = append(hashes, dep...)
	}
	cache[identifier] = hashes
	return hashes, nil
}

func getAddressImports(code []byte, name string) []string {
	deps := []string{}
	codes := map[common.Location][]byte{}
	location := common.StringLocation(name)
	program, _ := cmd.PrepareProgram(code, location, codes)
	for _, imp := range program.ImportDeclarations() {
		address, isAddressImport := imp.Location.(common.AddressLocation)
		if isAddressImport {
			adr := address.Address.Hex()
			impName := imp.Identifiers[0].Identifier
			deps = append(deps, fmt.Sprintf("%s.%s", adr, impName))
		}
	}
	return deps
}

// future: use service to get deployed contracts hashes
func generatePinDebthFirst(ctx context.Context, flowkit flowkit.Flowkit, address string, name string) (string, uint64, error) {

	memoize := map[string][]string{}
	pin, err := generateDependentPin(ctx, flowkit, address, name, memoize)

	if err != nil {
		return "", 0, err
	}

	block, _ := flowkit.Gateway().GetLatestBlock()
	height := block.Height
	hash := flixkit.ShaHex(strings.Join(pin, ""), "")

	return hash, height, nil
}
