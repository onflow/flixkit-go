package v1_1

import (
	"context"
	"encoding/json"
	"fmt"
	"strings"

	"github.com/onflow/cadence/runtime/ast"
	"github.com/onflow/cadence/runtime/cmd"
	cadenceCommon "github.com/onflow/cadence/runtime/common"
	"github.com/onflow/cadence/runtime/parser"
	"github.com/onflow/flow-go-sdk"
	"github.com/onflow/flow-go-sdk/access/grpc"

	"github.com/onflow/flixkit-go/v2/internal/common"
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
	clients           []*grpc.Client
	template          *InteractionTemplate
}

func NewTemplateGenerator(contractInfos ContractInfos, logger common.Logger, networks []common.NetworkConfig) (*Generator, error) {
	var clients []*grpc.Client
	for _, network := range networks {
		client, err := grpc.NewClient(network.Host)
		if err != nil {
			return nil, fmt.Errorf("could not create client for %s: %w", network.Name, err)
		}
		clients = append(clients, client)
	}

	deployedContracts := contractInfosToContracts(contractInfos)

	return &Generator{
		deployedContracts: deployedContracts,
		clients:           clients,
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
	_ = g.calculateNetworkPins()
	id, _ := GenerateFlixID(g.template)
	g.template.ID = id
	templateJson, err := json.MarshalIndent(g.template, "", "    ")

	return string(templateJson), err

}

func (g Generator) calculateNetworkPins() error {
	networkPins := make([]NetworkPin, 0)
	// only interested in the client networks
	for _, c := range g.clients {
		params, err := c.GetNetworkParameters(context.Background())
		if err != nil {
			continue
		}
		// remove flow- prefix
		networkName := strings.TrimPrefix(params.ChainID.String(), "flow-")
		cad, err := g.template.ReplaceCadenceImports(networkName)
		if err != nil {
			continue
		}
		networkPins = append(networkPins, NetworkPin{
			Network: networkName,
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

		// Built-in contracts imports are represented with identifier location
		_, isBuiltInContract := imp.Location.(cadenceCommon.IdentifierLocation)
		if isBuiltInContract {
			continue
		}

		contractName := imp.Location.String()
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

func getNetworkClient(networkName string, clients []*grpc.Client) *grpc.Client {
	for _, c := range clients {
		netParams, err := c.GetNetworkParameters(context.Background())
		if err != nil {
			continue
		}
		// chain id contains network name like "flow-testnet" contains "testnet"
		if strings.Contains(netParams.ChainID.String(), networkName) {
			return c
		}
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
		c := getNetworkClient(n.Network, g.clients)
		if n.DependencyPinBlockHeight == 0 && c != nil {
			block, err := c.GetLatestBlockHeader(ctx, true)
			if err != nil {
				return nil, err
			}
			height := block.Height

			details, err := g.GenerateDepPinDepthFirst(ctx, c, n.Address, contractName, height)
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

func (g *Generator) GenerateDepPinDepthFirst(ctx context.Context, clnt *grpc.Client, address string, name string, height uint64) (details *PinDetail, err error) {
	memoize := make(map[string]PinDetail)
	networkPinDetail, err := generateDependencyNetworks(ctx, clnt, address, name, memoize, height)
	if err != nil {
		return nil, err
	}

	return networkPinDetail, nil
}

func generateDependencyNetworks(ctx context.Context, c *grpc.Client, address string, name string, cache map[string]PinDetail, height uint64) (*PinDetail, error) {
	addr := flow.HexToAddress(address)
	identifier := fmt.Sprintf("A.%s.%s", addr.Hex(), name)
	pinDetail, ok := cache[identifier]
	if ok {
		return &pinDetail, nil
	}

	account, err := c.GetAccount(ctx, addr)
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
		dep, err := generateDependencyNetworks(ctx, c, address, name, cache, height)
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
	codes := map[cadenceCommon.Location][]byte{}
	location := cadenceCommon.StringLocation(name)
	program, _ := cmd.PrepareProgram(code, location, codes)
	for _, imp := range program.ImportDeclarations() {
		address, isAddressImport := imp.Location.(cadenceCommon.AddressLocation)
		if isAddressImport {
			adr := address.Address.HexWithPrefix()
			impName := imp.Identifiers[0].Identifier
			deps = append(deps, fmt.Sprintf("%s.%s", adr, impName))
		}
	}
	return deps
}

// Add this helper function
func contractInfosToContracts(infos ContractInfos) []Contract {
	contracts := make([]Contract, 0)

	for contractName, networks := range infos {
		contract := Contract{
			Contract: contractName,
			Networks: make([]Network, 0),
		}

		for networkName, address := range networks {
			addr := flow.HexToAddress(address)
			network := Network{
				Network: networkName,
				Address: addr.HexWithPrefix(),
			}
			contract.Networks = append(contract.Networks, network)
		}

		contracts = append(contracts, contract)
	}

	return contracts
}
