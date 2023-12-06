package flixkitv1_1

import (
	"context"
	"encoding/json"
	"fmt"
	"strings"

	"github.com/onflow/cadence/runtime/ast"
	"github.com/onflow/cadence/runtime/cmd"
	"github.com/onflow/cadence/runtime/common"
	"github.com/onflow/cadence/runtime/parser"
	"github.com/onflow/flixkit-go"
	v1_1 "github.com/onflow/flixkit-go/flixkitv1_1"
	"github.com/onflow/flixkit-go/generator"
	"github.com/onflow/flow-cli/flowkit"
	"github.com/onflow/flow-cli/flowkit/config"
	"github.com/onflow/flow-cli/flowkit/gateway"
	"github.com/onflow/flow-cli/flowkit/output"
	"github.com/onflow/flow-go-sdk"
	"github.com/onflow/flow-go-sdk/crypto"
	"github.com/spf13/afero"
)

type Generator struct {
	deployedContracts []v1_1.Contract
	coreContracts     []v1_1.Contract
	testnetClient     *flowkit.Flowkit
	mainnetClient     *flowkit.Flowkit
	template          *v1_1.InteractionTemplate
}

// stubb to pass in parameters
func NewGenerator(deployedContracts []v1_1.Contract, coreContracts []v1_1.Contract, logger output.Logger) (*Generator, error) {
	loader := afero.Afero{Fs: afero.NewOsFs()}

	gwt, err := gateway.NewGrpcGateway(config.TestnetNetwork)
	if err != nil {
		return nil, fmt.Errorf("could not create grpc gateway for testnet %w", err)
	}

	gwm, err := gateway.NewGrpcGateway(config.MainnetNetwork)
	if err != nil {
		return nil, fmt.Errorf("could not create grpc gateway for mainnet %w", err)
	}

	state, err := flowkit.Init(loader, crypto.ECDSA_P256, crypto.SHA3_256)
	if err != nil {
		return nil, fmt.Errorf("could not initialize flowkit state %w", err)
	}
	testnetClient := flowkit.NewFlowkit(state, config.TestnetNetwork, gwt, logger)
	mainnetClient := flowkit.NewFlowkit(state, config.MainnetNetwork, gwm, logger)

	if coreContracts == nil {
		// get default core contracts
		//coreContracts = generator.GetDefaultCoreContracts()
		coreContracts = []v1_1.Contract{
			{
				Contract: "FungibleToken",
				Networks: []v1_1.Network{
					{
						Network: "emulator",
						Address: "0xee82856bf20e2aa6",
					},
				},
			},
			{
				Contract: "NonFungibleToken",
				Networks: []v1_1.Network{
					{
						Network: "emulator",
						Address: "0x01cf0e2f2f715450",
					},
				},
			},
		}
	}

	j, err := json.MarshalIndent(deployedContracts, " ", "     ")
	if err != nil {
		return nil, err
	}
	fmt.Println(string(j))

	return &Generator{
		deployedContracts: deployedContracts,
		coreContracts:     coreContracts,
		testnetClient:     testnetClient,
		mainnetClient:     mainnetClient,
		template:          &v1_1.InteractionTemplate{},
	}, nil
}

func (g Generator) Generate(ctx context.Context, code string, preFill string) (string, error) {
	g.template = &v1_1.InteractionTemplate{}
	g.template.Init()
	if preFill != "" {
		t, err := v1_1.ParseFlix(preFill)
		if err != nil {
			return "", err
		}
		g.template = t
	}

	// make sure imports use new import syntax "string import"
	g.template.ProcessImports(code)
	program, err := parser.ParseProgram(nil, []byte(g.template.Data.Cadence.Body), parser.Config{})
	if err != nil {
		return "", err
	}

	err = g.template.DetermineCadenceType(program)
	if err != nil {
		return "", err
	}

	err = g.template.ParsePragma(program)
	if err != nil {
		return "", err
	}

	err = g.template.ProcessParameters(program)
	if err != nil {
		return "", err
	}

	err = g.processDependencies(ctx, program)
	if err != nil {
		return "", err
	}

	// need to process dependencies before calculating network pins
	// not all networks are supported, depends on what clients are passed in
	_ = g.calculateNetworkPins(program)

	id, _ := v1_1.GenerateFlixID(g.template)
	g.template.ID = id

	templateJson, err := json.MarshalIndent(g.template, "", "    ")
	fmt.Println(string(templateJson))

	return string(templateJson), err
}

func (g Generator) calculateNetworkPins(program *ast.Program) error {
	// only interested in public networks
	networksOfInterest := []string{
		config.MainnetNetwork.Name,
		config.TestnetNetwork.Name,
	}
	networkPins := make([]v1_1.NetworkPin, 0)
	for _, netName := range networksOfInterest {
		cad, err := g.template.GetAndReplaceCadenceImports(netName)
		if err != nil {
			continue
		}
		networkPins = append(networkPins, v1_1.NetworkPin{
			Network: netName,
			PinSelf: flixkit.ShaHex(cad, ""),
		})
	}
	g.template.Data.Cadence.NetworkPins = networkPins
	return nil
}

func (g Generator) processDependencies(ctx context.Context, program *ast.Program) error {
	imports := program.ImportDeclarations()

	if len(imports) == 0 {
		return nil
	}

	// fill in dependence information
	g.template.Data.Dependencies = make([]v1_1.Dependency, 0)
	for _, imp := range imports {
		contractName, err := generator.ExtractContractName(imp.String())
		if err != nil {
			return err
		}
		networks, err := g.generateDependenceInfo(ctx, contractName)
		if err != nil {
			return err
		}
		c := v1_1.Contract{
			Contract: contractName,
			Networks: networks,
		}
		dep := v1_1.Dependency{
			Contracts: []v1_1.Contract{c},
		}
		g.template.Data.Dependencies = append(g.template.Data.Dependencies, dep)
	}

	return nil
}

func (g *Generator) generateDependenceInfo(ctx context.Context, contractName string) ([]v1_1.Network, error) {
	// only support string import syntax
	contractNetworks := g.LookupImportContractInfo(contractName)
	if len(contractNetworks) == 0 {
		return nil, fmt.Errorf("could not find contract dependency %s", contractName)
	}
	var networks []v1_1.Network
	for _, network := range contractNetworks {
		var flowkit *flowkit.Flowkit
		if network.Network == config.MainnetNetwork.Name && g.mainnetClient != nil {
			flowkit = g.mainnetClient
		} else if network.Network == config.TestnetNetwork.Name && g.testnetClient != nil {
			flowkit = g.testnetClient
		}
		if network.DependencyPinBlockHeight == 0 && flowkit != nil {
			block, _ := flowkit.Gateway().GetLatestBlock()
			details, err := g.GenerateDepPinDepthFirst(ctx, *flowkit, network.Address, contractName, block.Height)
			if err != nil {
				return nil, err
			}
			network.DependencyPinBlockHeight = block.Height
			network.DependencyPin = details
		}

		networks = append(networks, network)
	}

	return networks, nil
}

// basic contract lookup, (address, network) using information given to generator, deployed contracts in flow.json or core contracts
func (g *Generator) LookupImportContractInfo(contractName string) []v1_1.Network {
	contracts := make([]v1_1.Contract, len(g.coreContracts)+len(g.deployedContracts))
	contracts = append(contracts, g.coreContracts...)
	contracts = append(contracts, g.deployedContracts...)

	// add deployed contracts to contracts map
	for _, contract := range contracts {
		if contractName == contract.Contract {
			return contract.Networks
		}
	}
	return nil
}

func (g *Generator) GenerateDepPinDepthFirst(ctx context.Context, flowkit flowkit.Flowkit, address string, name string, height uint64) (details v1_1.PinDetail, err error) {
	memoize := make(map[string]v1_1.PinDetail)
	networkPinDetail, err := generateDependencyNetworks(ctx, flowkit, address, name, memoize, height)
	if err != nil {
		return networkPinDetail, err
	}

	return networkPinDetail, nil
}

func generateDependencyNetworks(ctx context.Context, flowkit flowkit.Flowkit, address string, name string, cache map[string]v1_1.PinDetail, height uint64) (v1_1.PinDetail, error) {
	identifier := fmt.Sprintf("A.%s.%s", strings.ReplaceAll(address, "0x", ""), name)
	pinDetail, ok := cache[identifier]
	if ok {
		return pinDetail, nil
	}

	account, err := flowkit.GetAccount(ctx, flow.HexToAddress(address))
	if err != nil {
		return pinDetail, err
	}
	code := account.Contracts[name]
	depend := v1_1.PinDetail{
		PinContractName:    name,
		PinContractAddress: address,
		PinSelf:            flixkit.ShaHex(code, ""),
	}
	depend.CalculatePin(height)
	imports := getAddressImports(code, name)
	detailImports := make([]v1_1.PinDetail, 0)
	for _, imp := range imports {
		split := strings.Split(imp, ".")
		address, name := split[0], split[1]
		dep, err := generateDependencyNetworks(ctx, flowkit, address, name, cache, height)
		if err != nil {
			return depend, err
		}
		detailImports = append(detailImports, dep)
		cache[identifier] = dep
	}
	depend.Imports = detailImports
	return depend, nil
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
