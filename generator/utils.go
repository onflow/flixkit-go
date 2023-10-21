package generator

import (
	"context"
	"errors"
	"regexp"
	"strings"

	"github.com/onflow/flixkit-go"
	"github.com/onflow/flow-cli/flowkit"
)

// todo: remove this and update tests
func ExtractImports(cadenceCode string) []string {
	// Regex pattern to match Cadence import lines
	pattern := `import [\w\s\"\.]+(?:from 0x[\w]+)?`
	r := regexp.MustCompile(pattern)

	// Find all matches in the given code
	matches := r.FindAllString(cadenceCode, -1)

	return matches
}

func getDependencyContractCode(contractName string, flowkit flowkit.State, flow flowkit.Services) (string, error) {
	// use flow client to get contract code
	return "", nil
}

func findImportDetails(contractName string, flowkit flowkit.State) (string, string, error) {
	// look up core contracts if not found in flow.json state

	return "", "", nil
}

func ParseImport(ctx context.Context, line string, flowkit flowkit.State) (flixkit.Contracts, error) {
	// Define regex patterns
	importSyntax := `import "(?P<contract>[^"]+)"`
	oldImportSyntax := `import (?P<contract>\w+) from (?P<address>0x[\w]+)`

	contractInfo := flixkit.Contracts{}
	// Use regex to extract relevant information
	// structure for flix is dependency -> import placeholder -> contract -> network
	// if old import syntax and uses address then use the address as the placeholder
	// if new import syntax then generate a placeholder by 0xContractNameAddress
	if matches, _ := regexpMatch(importSyntax, line); matches != nil {
		// new import syntax need to find the contract deployment to get address
		contraceName := matches["contract"]
		info := GetContractInformation(contraceName)
		contractInfo = flixkit.Contracts{
			contraceName: info,
		}
		// if contract info is nil then need to look up in flow.json
	} else if matches, _ := regexpMatch(oldImportSyntax, line); matches != nil {
		contraceName := matches["contract"]
		info := GetContractInformation(contraceName)
		contractInfo = flixkit.Contracts{
			contraceName: info,
		}
		// if contract info is nil then not core contract then
		// determine if contract has been deployed in flow.json
		//
	}

	return contractInfo, nil
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
	// TODO: determine if code is an interface, flix can be interfaces
	interfacePattern := regexp.MustCompile(`pub\s+interface`)
	if strings.Contains(code, "transaction(") {
		return "transaction", nil
	} else if strings.Contains(code, "pub fun main()") {
		return "script", nil
	} else if interfacePattern.MatchString(code) {
		return "interface", nil
	}
	return "", errors.New("could not determine if code is transaction or script")
}
