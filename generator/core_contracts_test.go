package generator

import "testing"

func TestExistContract(t *testing.T) {
	tests := []struct {
		name   string
		input  string
		output bool
	}{
		{
			name:   "existing contract",
			input:  "FungibleToken",
			output: true,
		},
		{
			name:   "non-existing contract",
			input:  "NonExistingContract",
			output: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := GetContractInformation(tt.input, nil, GetDefaultCoreContracts())
			if tt.output && got == nil {
				t.Errorf("ContractExist() got = %v, want %v", got, tt.output)
			}
		})
	}
}
