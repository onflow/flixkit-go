package bindings

import (
	"encoding/json"
	"testing"

	"github.com/hexops/autogold/v2"
	v1_1 "github.com/onflow/flixkit-go/flixkitv1_1"
	"github.com/stretchr/testify/assert"
)

var minimumNoParamTemplateTS_SCRIPT = &v1_1.InteractionTemplate{
	FType:    "InteractionTemplate",
	FVersion: "1.1.0",
	ID:       "290b6b6222b2a77b16db896a80ddf29ebd1fa3038c9e6625a933fa213fce51fa",
	Data: v1_1.Data{
		Type:      "script",
		Interface: "",
		Cadence: v1_1.Cadence{
			Body:        "pub fun main(): Int { return 1 }",
			NetworkPins: []v1_1.NetworkPin{},
		},
	},
}

var minimumNoParamTemplateTS_TX = &v1_1.InteractionTemplate{
	FType:    "InteractionTemplate",
	FVersion: "1.1.0",
	ID:       "290b6b6222b2a77b16db896a80ddf29ebd1fa3038c9e6625a933fa213fce51fa",
	Data: v1_1.Data{
		Type:      "transaction",
		Interface: "",
		Cadence: v1_1.Cadence{
			Body:        "import \"HelloWorld\"\n\n#interaction (\n  version: \"1.1.0\",\n\ttitle: \"Update Greeting\",\n\tdescription: \"Update the greeting on the HelloWorld contract\",\n\tlanguage: \"en-US\",\n)\ntransaction() {\n\n  prepare(acct: AuthAccount) {\n     }\n\n  execute {\n   \n  }\n}\n",
			NetworkPins: []v1_1.NetworkPin{},
		},
	},
}

func TestTSGenNoParamsScript(t *testing.T) {
	ttemp, err := json.Marshal(minimumNoParamTemplateTS_SCRIPT)
	assert.NoError(t, err, "marshal template to json should not return an error")

	generator := NewFclTSGenerator()
	assert := assert.New(t)

	out, err := generator.GenerateTS(string(ttemp), "./min.template.json")
	assert.NoError(err, "ParseTemplate should not return an error")
	autogold.ExpectFile(t, out)
}

func TestTSGenNoParamsTx(t *testing.T) {
	ttemp, err := json.Marshal(minimumNoParamTemplateTS_TX)
	assert.NoError(t, err, "marshal template to json should not return an error")

	generator := NewFclTSGenerator()
	assert := assert.New(t)

	out, err := generator.GenerateTS(string(ttemp), "./min.template.json")
	assert.NoError(err, "ParseTemplate should not return an error")
	autogold.ExpectFile(t, out)
}
