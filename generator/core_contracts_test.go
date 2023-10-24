package generator

import "testing"

func TestContractExist(t *testing.T) {
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
			got := getContractInformation(tt.input)
			if tt.output && got == nil {
				t.Errorf("ContractExist() got = %v, want %v", got, tt.output)
			}
		})
	}
}
