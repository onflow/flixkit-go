package flixkitv1_1

import (
	"context"
	"encoding/json"
	"fmt"

	"github.com/onflow/cadence/runtime/parser"
	v1_1 "github.com/onflow/flixkit-go/flixkitv1_1"
	"github.com/onflow/flow-cli/flowkit"
	"github.com/onflow/flow-cli/flowkit/config"
	"github.com/onflow/flow-cli/flowkit/gateway"
	"github.com/onflow/flow-cli/flowkit/output"
	"github.com/onflow/flow-go-sdk/crypto"
	"github.com/spf13/afero"
)

type Generator struct {
	deployedContracts []v1_1.Contract
	coreContracts     []v1_1.Contract
	testnetClient     *flowkit.Flowkit
	mainnetClient     *flowkit.Flowkit
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

	return &Generator{
		deployedContracts: deployedContracts,
		coreContracts:     coreContracts,
		testnetClient:     testnetClient,
		mainnetClient:     mainnetClient,
	}, nil
}

func (g Generator) Generate(ctx context.Context, code string, preFill string) (string, error) {
	template := &v1_1.InteractionTemplate{}
	if preFill != "" {
		t, err := v1_1.ParseFlix(preFill)
		if err != nil {
			return "", err
		}
		template = t
	}

	codeBytes := []byte(code)
	program, err := parser.ParseProgram(nil, codeBytes, parser.Config{})
	if err != nil {
		return "", err
	}

	prags := program.PragmaDeclarations()
	err = v1_1.ParsePragma(prags, template)
	if err != nil {
		return "", err
	}

	templateJson, err := json.MarshalIndent(template, "", "    ")
	fmt.Println(string(templateJson))
	return string(templateJson), err
}
