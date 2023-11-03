package generator

import (
	"context"
	"errors"
	"fmt"
	"regexp"
	"sort"

	"github.com/onflow/flixkit-go"
)

func parseImport(ctx context.Context, line string, deployedContracts []flixkit.Contracts) (map[string]flixkit.Contracts, error) {
	// Define regex patterns
	importSyntax := `import "(?P<contract>[^"]+)"`
	oldImportSyntax := `import (?P<contract>\w+) from (?P<address>[\w]+)`

	placeholder := ""
	// Use regex to extract relevant information
	// structure for flix is dependency -> import placeholder -> contract -> network
	// if old import syntax and uses address then use the address as the placeholder
	// if new import syntax then placeholder is ""
	var contractName string
	var info flixkit.Networks
	if matches, _ := regexpMatch(importSyntax, line); matches != nil {
		// new import syntax need to find the contract deployment to get address
		contractName := matches["contract"]
		placeholder = "0x" + contractName
		info = getContractInformation(contractName, deployedContracts)
		// if contract info is nil then need to look up in deployed contracts
		// need to change the import statement to use the placeholder

	} else if matches, _ := regexpMatch(oldImportSyntax, line); matches != nil {
		contractName = matches["contract"]
		placeholder = matches["address"]
		info = getContractInformation(contractName, deployedContracts)
		// if contract info is nil then no core contract then
		// determine if contract has been deployed in deployed contracts
		//
	}

	if info == nil {
		return nil, fmt.Errorf("contract %s not found", contractName)
	}

	return map[string]flixkit.Contracts{
		placeholder: {
			contractName: info,
		},
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

// TODO: make sure message types are sorted when there is user created types
func sortedKeys(m map[string]interface{}) []string {
	keys := make([]string, 0, len(m))
	for k := range m {
		keys = append(keys, k)
	}
	sort.Strings(keys)
	return keys
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
