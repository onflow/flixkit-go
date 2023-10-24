package generator

import (
	"context"
	"errors"
	"fmt"
	"regexp"

	"github.com/onflow/flixkit-go"
)

func getDependencyContractCode(contractName string) (string, error) {
	// use flow client to get contract code
	return "", nil
}

func findImportDetails(contractName string) (string, string, error) {
	// look up core contracts if not found in flow.json state

	return "", "", nil
}

func parseImport(ctx context.Context, line string) (map[string]flixkit.Contracts, error) {
	// Define regex patterns
	importSyntax := `import "(?P<contract>[^"]+)"`
	oldImportSyntax := `import (?P<contract>\w+) from (?P<address>[\w]+)`

	contractInfo := flixkit.Contracts{}
	placeholder := ""
	// Use regex to extract relevant information
	// structure for flix is dependency -> import placeholder -> contract -> network
	// if old import syntax and uses address then use the address as the placeholder
	// if new import syntax then generate a placeholder by 0xContractNameAddress
	if matches, _ := regexpMatch(importSyntax, line); matches != nil {
		// new import syntax need to find the contract deployment to get address
		contractName := matches["contract"]
		placeholder = "0x" + contractName
		info := getContractInformation(contractName)
		contractInfo = flixkit.Contracts{
			contractName: info,
		}
		// if contract info is nil then need to look up in flow.json
	} else if matches, _ := regexpMatch(oldImportSyntax, line); matches != nil {
		contraceName := matches["contract"]
		placeholder = matches["address"]
		info := getContractInformation(contraceName)
		contractInfo = flixkit.Contracts{
			contraceName: info,
		}
		// if contract info is nil then no core contract then
		// determine if contract has been deployed in flow.json
		//

	}

	return map[string]flixkit.Contracts{
		placeholder: contractInfo,
	}, nil
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
		return "interface", nil
	}

	fmt.Println(code, transactionRegex.MatchString(code), scriptRegex.MatchString(code), interfaceRegex.MatchString(code))

	return "", errors.New("could not determine if code is transaction or script")
}
