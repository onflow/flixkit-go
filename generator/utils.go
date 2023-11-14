package generator

import (
	"context"
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

func ProcessParameters(program *ast.Program, template *flixkit.FlowInteractionTemplate) error {
	if program != nil && program.SoleTransactionDeclaration() != nil && program.SoleTransactionDeclaration().ParameterList != nil {
		if template.Data.Arguments == nil {
			template.Data.Arguments = flixkit.Arguments{}
		}

		for i, param := range program.SoleTransactionDeclaration().ParameterList.Parameters {
			argMessages := flixkit.Messages{}
			if template.Data.Arguments != nil && template.Data.Arguments[param.Identifier.String()] != (flixkit.Argument{}) {
				argMessages = template.Data.Arguments[param.Identifier.String()].Messages
				fmt.Println("found argMessages", argMessages)
			}
			template.Data.Arguments[param.Identifier.String()] = flixkit.Argument{
				Type:     param.TypeAnnotation.Type.String(),
				Index:    i,
				Messages: argMessages,
			}
		}
	}

	return nil
}

func RegexpMatch(pattern, text string) (map[string]string, error) {
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

func DetermineCadenceType(program *ast.Program) string {
	funcs := program.FunctionDeclarations()
	trans := program.TransactionDeclarations()

	if len(funcs) > 0 {
		return "script"
	} else if len(trans) > 0 {
		return "transaction"
	}
	return "interface"
}

func NormalizeImports(cadenceCode string) string {
	// Define a regular expression to match the "import ContractName from 0xContractName" pattern
	pattern := regexp.MustCompile(`import\s+(\w+)\s+from\s+0x\w+`)
	// Replace the matched pattern with "import \"ContractName\""
	replaced := pattern.ReplaceAllString(cadenceCode, `import "$1"`)
	return replaced
}

func UnNormalizeImports(cadenceCode string) string {
	// Define a regular expression to match the import "ContractName" pattern
	pattern := regexp.MustCompile(`import "(.+?)"`)
	// Replace the matched pattern with "import ContractName from 0xContractName"
	replaced := pattern.ReplaceAllString(cadenceCode, `import $1 from 0x$1`)
	return replaced
}

func ExtractContractName(importStr string) (string, error) {
	// Create a regex pattern to find the contract name inside the quotes
	pattern := regexp.MustCompile(`import "([^"]+)"`)
	matches := pattern.FindStringSubmatch(importStr)

	if len(matches) >= 2 {
		return matches[1], nil
	}

	return "", fmt.Errorf("no contract name found in string")
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
func GeneratePinDebthFirst(ctx context.Context, flowkit flowkit.Flowkit, address string, name string) (string, uint64, error) {

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
