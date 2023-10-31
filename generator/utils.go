package generator

import (
	"bytes"
	"context"
	"encoding/hex"
	"errors"
	"fmt"
	"regexp"
	"sort"

	"github.com/ethereum/go-ethereum/rlp"
	"github.com/onflow/flixkit-go"
	"golang.org/x/crypto/sha3"
)

func getDependencyContractCode(contractName string) (string, error) {
	// use flow client to get contract code
	return "", nil
}

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

func genHash(utf8String string) string {
	hasher := sha3.New256()          // Create a new SHA3 256 hasher
	hasher.Write([]byte(utf8String)) // Write the utf8 string to the hasher
	hash := hasher.Sum(nil)          // Get the hash result

	return hex.EncodeToString(hash) // Convert the hash result to hex
}

func generateTemplateId(template *flixkit.FlowInteractionTemplate) (string, error) {
	// Your normalization function
	// template = normalizeInteractionTemplate(template)

	var buffer bytes.Buffer
	// Mimicking the hashing order in the JS code
	if template.FType != "" {
		buffer.WriteString(genHash(template.FType))
	}
	if template.FVersion != "" {
		buffer.WriteString(genHash(template.FVersion))
	}
	if template.Data.Type != "" {
		buffer.WriteString(genHash(template.Data.Type))
	}
	if template.Data.Interface != "" {
		buffer.WriteString(genHash(template.Data.Interface))
	}

	if template.Data.Messages.Title != nil {
		for _, i18nKey := range template.Data.Messages.Title.I18N {
			buffer.WriteString(genHash(i18nKey))
			buffer.WriteString(genHash(template.Data.Messages.Title.I18N[i18nKey]))
		}
	}

	if template.Data.Messages.Description != nil {
		for _, i18nKey := range template.Data.Messages.Description.I18N {
			buffer.WriteString(genHash(i18nKey))
			buffer.WriteString(genHash(template.Data.Messages.Description.I18N[i18nKey]))
		}
	}
	if template.Data.Cadence != "" {
		buffer.WriteString(genHash(template.Data.Cadence))
	}

	// Continue for dependencies and arguments in a similar fashion...

	encoded, err := rlp.EncodeToBytes(buffer.String())
	if err != nil {
		return "", err
	}
	encodedHex := hex.EncodeToString(encoded)
	return genHash(encodedHex), nil
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
