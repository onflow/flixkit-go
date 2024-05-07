package v1_1

import (
	"context"
	"encoding/json"
	"fmt"
	"strings"

	"github.com/onflow/cadence/runtime/ast"
	"github.com/onflow/cadence/runtime/cmd"
	"github.com/onflow/cadence/runtime/common"
	"github.com/onflow/cadence/runtime/parser"
	"github.com/onflow/flixkit-go/internal/contracts"
	"github.com/onflow/flow-go-sdk"
	"github.com/onflow/flow-go-sdk/crypto"
	"github.com/onflow/flowkit"
	"github.com/onflow/flowkit/config"
	"github.com/onflow/flowkit/gateway"
	"github.com/onflow/flowkit/output"
	"github.com/spf13/afero"
)

/*
Same structure as core contracts, using config network names
*/
type NetworkAddressMap map[string]string

/*
Same structure as core contracts, keyed by contract name
*/
type ContractInfos map[string]NetworkAddressMap

type Generator struct {
	deployedContracts []Contract
	testnetClient     *flowkit.Flowkit
	mainnetClient     *flowkit.Flowkit
	template          *InteractionTemplate
}

func NewTemplateGenerator(contractInfos ContractInfos, logger output.Logger) (*Generator, error) {
	loader := afero.Afero{Fs: afero.NewOsFs()}

	gwt, err := gateway.NewGrpcGateway(config.TestnetNetwork)
	if err != nil {
		return nil, fmt.Errorf("could not create grpc gateway for testnet %w", err)
	}

	gwm, err := gateway.NewGrpcGateway(config.MainnetNetwork)
	if err != nil {
		return nil, fmt.Errorf("could not create grpc gateway for mainnet %w", err)
	}

	state, err := flowkit.Init(loader, crypto.ECDSA_P256, crypto.SHA3_256, "")
	if err != nil {
		return nil, fmt.Errorf("could not initialize flowkit state %w", err)
	}
	testnetClient := flowkit.NewFlowkit(state, config.TestnetNetwork, gwt, logger)
	mainnetClient := flowkit.NewFlowkit(state, config.MainnetNetwork, gwm, logger)
	// add core contracts to deployed contracts
	cc := contracts.GetCoreContracts()
	deployedContracts := make([]Contract, 0)
	for contractName, c := range cc {
		contract := Contract{
			Contract: contractName,
			Networks: []Network{
				{Network: config.MainnetNetwork.Name, Address: c[config.MainnetNetwork.Name]},
				{Network: config.TestnetNetwork.Name, Address: c[config.TestnetNetwork.Name]},
				{Network: config.EmulatorNetwork.Name, Address: c[config.EmulatorNetwork.Name]},
			},
		}
		deployedContracts = append(deployedContracts, contract)
	}
	// allow user contracts to override core contracts
	for contractInfo, networks := range contractInfos {
		contract := Contract{
			Contract: contractInfo,
			Networks: make([]Network, 0),
		}
		for network, address := range networks {
			addr := flow.HexToAddress(address)
			contract.Networks = append(contract.Networks, Network{
				Network: network,
				Address: "0x" + addr.Hex(),
			})
		}
		deployedContracts = append(deployedContracts, contract)
	}

	return &Generator{
		deployedContracts: deployedContracts,
		testnetClient:     testnetClient,
		mainnetClient:     mainnetClient,
		template:          &InteractionTemplate{},
	}, nil
}

func (g Generator) CreateTemplate(ctx context.Context, code string, preFill string) (string, error) {
	g.template = &InteractionTemplate{}
	g.template.Init()
	if preFill != "" {
		t, err := ParseFlix(preFill)
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
	_ = g.calculateNetworkPins(program)
	id, _ := GenerateFlixID(g.template)
	g.template.ID = id
	templateJson, err := json.MarshalIndent(g.template, "", "    ")

	return string(templateJson), err

}

func (g Generator) calculateNetworkPins(program *ast.Program) error {
	networksOfInterest := []string{
		config.MainnetNetwork.Name,
		config.TestnetNetwork.Name,
	}
	networkPins := make([]NetworkPin, 0)
	for _, netName := range networksOfInterest {
		cad, err := g.template.ReplaceCadenceImports(netName)
		if err != nil {
			continue
		}
		networkPins = append(networkPins, NetworkPin{
			Network: netName,
			PinSelf: ShaHex(cad, ""),
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
	g.template.Data.Dependencies = make([]Dependency, 0)
	for _, imp := range imports {
		contractName, err := ExtractContractName(imp.String())
		if err != nil {
			return err
		}
		networks, err := g.generateDependenceInfo(ctx, contractName)
		if err != nil {
			return err
		}
		c := Contract{
			Contract: contractName,
			Networks: networks,
		}
		dep := Dependency{
			Contracts: []Contract{c},
		}
		g.template.Data.Dependencies = append(g.template.Data.Dependencies, dep)
	}

	return nil
}

func (g *Generator) generateDependenceInfo(ctx context.Context, contractName string) ([]Network, error) {
	// only support string import syntax
	contractNetworks := g.LookupImportContractInfo(contractName)
	if len(contractNetworks) == 0 {
		return nil, fmt.Errorf("could not find contract dependency %s", contractName)
	}
	var networks []Network
	for _, n := range contractNetworks {
		network := Network{
			Network: n.Network,
			Address: n.Address,
		}
		var flowkit *flowkit.Flowkit
		if n.Network == config.MainnetNetwork.Name && g.mainnetClient != nil {
			flowkit = g.mainnetClient
		} else if n.Network == config.TestnetNetwork.Name && g.testnetClient != nil {
			flowkit = g.testnetClient
		}
		if n.DependencyPinBlockHeight == 0 && flowkit != nil {
			block, _ := flowkit.Gateway().GetLatestBlock(ctx)
			height := block.Height

			details, err := g.GenerateDepPinDepthFirst(ctx, flowkit, n.Address, contractName, height)
			if err != nil {
				return nil, err
			}
			network.DependencyPinBlockHeight = height
			network.DependencyPin = details
		}
		networks = append(networks, network)
	}

	return networks, nil
}

func (g *Generator) LookupImportContractInfo(contractName string) []Network {
	for _, contract := range g.deployedContracts {
		if contractName == contract.Contract {
			return contract.Networks
		}
	}
	return nil
}

func (g *Generator) GenerateDepPinDepthFirst(ctx context.Context, flowkit *flowkit.Flowkit, address string, name string, height uint64) (details *PinDetail, err error) {
	memoize := make(map[string]PinDetail)
	networkPinDetail, err := generateDependencyNetworks(ctx, flowkit, address, name, memoize, height)
	if err != nil {
		return nil, err
	}

	return networkPinDetail, nil
}

func generateDependencyNetworks(ctx context.Context, flowkit *flowkit.Flowkit, address string, name string, cache map[string]PinDetail, height uint64) (*PinDetail, error) {
	addr := flow.HexToAddress(address)
	identifier := fmt.Sprintf("A.%s.%s", addr.Hex(), name)
	pinDetail, ok := cache[identifier]
	if ok {
		return &pinDetail, nil
	}

	account, err := flowkit.GetAccount(ctx, addr)
	if err != nil {
		return nil, err
	}

	code := account.Contracts[name]
	depend := PinDetail{
		PinContractName:    name,
		PinContractAddress: "0x" + addr.Hex(),
		PinSelf:            ShaHex(code, ""),
	}
	depend.CalculatePin(height)
	pins := []string{depend.PinSelf}
	imports := getAddressImports(code, name)
	detailImports := make([]PinDetail, 0)
	for _, imp := range imports {
		split := strings.Split(imp, ".")
		address, name := split[0], split[1]
		dep, err := generateDependencyNetworks(ctx, flowkit, address, name, cache, height)
		if err != nil {
			return nil, err
		}
		if dep != nil {
			detailImports = append(detailImports, *dep)
			cache[identifier] = *dep
		}
		pins = append(pins, dep.PinSelf)
	}

	depend.Imports = detailImports
	depend.Pin = ShaHex(strings.Join(pins, ""), "")
	return &depend, nil
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
