package v1

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestReplaceCadenceImports(t *testing.T) {
	tests := []struct {
		name        string
		template    *FlowInteractionTemplate
		network     string
		want        string
		wantErr     bool
		errContains string
	}{
		{
			name: "replaces contract addresses from dependencies",
			template: &FlowInteractionTemplate{
				Data: Data{
					Cadence: `
						import FungibleToken from 0x1234
						import FlowToken from 0x5678
					`,
					Dependencies: Dependencies{
						"0x1234": Contracts{
							"FungibleToken": Networks{
								"testnet": Network{
									Address: "0x9a0766d93b6608b7",
								},
							},
						},
						"0x5678": Contracts{
							"FlowToken": Networks{
								"testnet": Network{
									Address: "0x7e60df042a9c0868",
								},
							},
						},
					},
				},
			},
			network: "testnet",
			want: `
						import FungibleToken from 0x9a0766d93b6608b7
						import FlowToken from 0x7e60df042a9c0868
					`,
		},
		{
			name: "handles missing network in dependencies",
			template: &FlowInteractionTemplate{
				Data: Data{
					Cadence: `import MyContract from 0xPLACEHOLDER`,
					Dependencies: Dependencies{
						"0xPLACEHOLDER": Contracts{
							"MyContract": Networks{},
						},
					},
				},
			},
			network:     "testnet",
			wantErr:     true,
			errContains: "network testnet not found",
		},
		{
			name: "handles missing contract in dependencies",
			template: &FlowInteractionTemplate{
				Data: Data{
					Cadence:      `import FungibleToken from 0x1234`,
					Dependencies: Dependencies{},
				},
			},
			network:     "testnet",
			wantErr:     true,
			errContains: "network testnet not found",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := tt.template.ReplaceCadenceImports(tt.network)
			if tt.wantErr {
				assert.Error(t, err)
				if tt.errContains != "" {
					assert.Contains(t, err.Error(), tt.errContains)
				}
				return
			}
			assert.NoError(t, err)
			assert.Equal(t, tt.want, got)
		})
	}
}
