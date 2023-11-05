package generator

import (
	"context"
	"fmt"
	"regexp"
	"sort"

	"github.com/onflow/cadence/runtime/parser"
	"github.com/onflow/flixkit-go"
	"github.com/onflow/flow-cli/flowkit"
)

type Generator1_0_0 struct {
	deployedContracts []flixkit.Contracts
	testnetClient     *flowkit.Flowkit
	mainnetClient     *flowkit.Flowkit
}

// stubb to pass in parameters
func NewGenerator(deployedContracts []flixkit.Contracts, testnet *flowkit.Flowkit, mainnet *flowkit.Flowkit) *Generator1_0_0 {
	return &Generator1_0_0{
		deployedContracts: deployedContracts,
		testnetClient:     testnet,
		mainnetClient:     mainnet,
	}
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

	// parsing cadence manually cuz cadence parser does not like old import syntax statements "from 0xPLACEHOLDER"
	err = g.processDependencies(template)
	if err != nil {
		return nil, err
	}

	id, err := flixkit.GenerateFlixID(template)
	if err != nil {
		return nil, err
	}
	template.ID = id

	return template, nil
}

func (g Generator1_0_0) processDependencies(template *flixkit.FlowInteractionTemplate) error {
	ctx := context.Background()
	normalizedCode := normalizeImports(template.Data.Cadence)
	// update cadence code in template so that dependencies match
	template.Data.Cadence = normalizedCode

	noCommentsCode := stripComments(normalizedCode)
	re := regexp.MustCompile(`(?m)^\s*import.*$`)
	imports := re.FindAllString(noCommentsCode, -1)
	// sort imports so they are processed consistently
	sort.Strings(imports)

	if len(imports) == 0 {
		return nil
	}
	// fill in dependence information
	deps := make(flixkit.Dependencies, len(imports))
	for _, imp := range imports {
		dep, err := g.parseImport(ctx, imp, g.deployedContracts)
		if err != nil {
			return err
		}
		for contractName, contract := range dep {
			deps[contractName] = contract
		}
		template.Data.Dependencies = deps
	}

	return nil
}

func (g *Generator1_0_0) parseImport(ctx context.Context, line string, deployedContracts []flixkit.Contracts) (map[string]flixkit.Contracts, error) {
	// Define regex patterns
	importSyntax := `import "(?P<contract>[^"]+)"`
	oldImportSyntax := `import (?P<contract>\w+) from (?P<address>[\w]+)`

	var placeholder string
	var contractName string
	var info flixkit.Networks
	if matches, _ := regexpMatch(importSyntax, line); matches != nil {
		// new import syntax detected, convert to old import syntax, limitation of 1.0.0
		contractName := matches["contract"]
		placeholder = "0x" + contractName
		info = getContractInformation(contractName, deployedContracts)

	} else if matches, _ := regexpMatch(oldImportSyntax, line); matches != nil {
		contractName = matches["contract"]
		placeholder = matches["address"]
		info = getContractInformation(contractName, deployedContracts)
	}

	for name, network := range info {
		if name == "mainnet" || name == "testnet" {
			if g.mainnetClient != nil && g.testnetClient != nil {
				flowkit := *g.mainnetClient
				if name == "testnet" {
					flowkit = *g.testnetClient
				}
				if network.Pin == "" {
					hash, height, _ := generatePinDebthFirst(ctx, flowkit, network.Address, network.Contract)
					network.Pin = hash
					network.PinBlockHeight = height
				}
			}
			info[name] = network
		}
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
