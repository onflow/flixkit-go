package v1

import "github.com/onflow/cadence"

type FLIX struct{}

func ParseJSON(flixJSON []byte) (FLIX, error) {
	return FLIX{}, nil
}

func (f FLIX) AsCadance() (cadence.Value, error) {
	return nil, nil
}

func (f FLIX) AsJSON() string {
	return ""
}

func (f FLIX) ReplaceImports() {

}

func (f FLIX) CreateBindings() (string, error) {
	return "", nil
}
