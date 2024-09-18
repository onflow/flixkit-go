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
	"github.com/onflow/flow-go-sdk"
	"github.com/onflow/flowkit/v2"
	"github.com/onflow/flowkit/v2/config"
	"github.com/onflow/flowkit/v2/gateway"
	"github.com/onflow/flowkit/v2/output"
	"github.com/spf13/afero"

	"github.com/onflow/flixkit-go/internal/contracts"
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
	clients           []*flowkit.Flowkit
	template          *InteractionTemplate
}

func NewTemplateGenerator(contractInfos ContractInfos, logger output.Logger, networks []config.Network) (*Generator, error) {
	loader := afero.Afero{Fs: afero.NewOsFs()}

	var clients []*flowkit.Flowkit
	for _, network := range networks {
		gw, err := gateway.NewGrpcGateway(network)
		if err != nil {
			return nil, fmt.Errorf("could not create grpc gateway for %s %w", network.Name, err)
		}
		state, err := flowkit.Init(loader)
		if err != nil {
			return nil, fmt.Errorf("could not initialize flowkit state %w", err)
		}

		client := flowkit.NewFlowkit(state, network, gw, logger)
		clients = append(clients, client)
	}
	networkNames := make([]string, 0)
	for _, network := range networks {
		networkNames = append(networkNames, network.Name)
	}
	// add core contracts to deployed contracts
	cc := contracts.GetCoreContracts()
	deployedContracts := make([]Contract, 0)
	for contractName, c := range cc {
		var nets []Network
		for network, address := range c {
			// if network is in user defined networks then add to deployed contracts
			if isItemInArray(network, networkNames) {
				addr := flow.HexToAddress(address)
				nets = append(nets, Network{
					Network: network,
					Address: addr.HexWithPrefix(),
				})
			}
		}
		if len(nets) > 0 {
			contract := Contract{
				Contract: contractName,
				Networks: nets,
			}
			deployedContracts = append(deployedContracts, contract)
		}
	}

	deployedContracts = mergeContractsAndInfos(deployedContracts, contractInfos)

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
	networksOfInterest := []string{}
	// only interested in the client networks
	for _, client := range g.clients {
		networksOfInterest = append(networksOfInterest, client.Network().Name)
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

		// Built-in contracts imports are represented with identifier location
		_, isBuiltInContract := imp.Location.(common.IdentifierLocation)
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

func getNetworkClient(networkName string, clients []*flowkit.Flowkit) *flowkit.Flowkit {
	for _, client := range clients {
		if client.Network().Name == networkName {
			return client
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
		flowkit := getNetworkClient(n.Network, g.clients)
		if n.DependencyPinBlockHeight == 0 && flowkit != nil {
			block, err := flowkit.Gateway().GetLatestBlock(ctx)
			if err != nil {
				return nil, err
			}
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

// Helper function to merge ContractInfos into []Contract and add missing contracts
func mergeContractsAndInfos(contracts []Contract, infos ContractInfos) []Contract {
	// Track existing contracts for quick lookup
	existingContracts := make(map[string]int)
	for i, contract := range contracts {
		existingContracts[contract.Contract] = i

		if info, exists := infos[contract.Contract]; exists {
			// Create a map to track existing networks for duplicate check
			existingNetworks := make(map[string]bool)
			for _, network := range contract.Networks {
				existingNetworks[network.Network] = true
			}

			// Iterate over the networks in the ContractInfos
			for network, address := range info {
				if !existingNetworks[network] {
					// If the network doesn't exist in the contract, add it
					contracts[i].Networks = append(contracts[i].Networks, Network{
						Network: network,
						Address: address,
					})
				}
			}
		}
	}

	// Add contracts from infos that don't exist in the current contracts array
	for contractName, networks := range infos {
		if _, exists := existingContracts[contractName]; !exists {
			newContract := Contract{Contract: contractName}
			for network, address := range networks {
				newContract.Networks = append(newContract.Networks, Network{
					Network: network,
					Address: address,
				})
			}
			contracts = append(contracts, newContract)
		}
	}

	return contracts
}
