package generator

import (
	"fmt"
	"reflect"
	"testing"
)

func TestExtractImports(t *testing.T) {
	tests := []struct {
		name   string
		input  string
		output []string
	}{
		{
			name:   "basic imports",
			input:  `import "ContractA"`,
			output: []string{`import "ContractA"`},
		},
		{
			name:   "multiple imports",
			input:  `import "ContractA"\nimport ContractB from 0x123456`,
			output: []string{`import "ContractA"`, `import ContractB from 0x123456`},
		},
		{
			name:   "mixed imports",
			input:  `import "ContractA"\nimport ContractB from 0x123456\nimport ContractC from 0xSomePlaceHolderAddress`,
			output: []string{`import "ContractA"`, `import ContractB from 0x123456`, `import ContractC from 0xSomePlaceHolderAddress`},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := ExtractImports(tt.input)
			fmt.Println(got, tt.output)
			if !reflect.DeepEqual(got, tt.output) {
				t.Errorf("ExtractImports() got = %v, want %v", got, tt.output)
			}
		})
	}
}
