package flixkit

import (
	"github.com/stretchr/testify/assert"
	"net/http"
	"net/http/httptest"
	"testing"
)

var template = `{
	  "f_type": "InteractionTemplate",
	  "f_version": "1.0.0",
	  "id": "290b6b6222b2a77b16db896a80ddf29ebd1fa3038c9e6625a933fa213fce51fa",
	  "data": {
		"type": "transaction",
		"interface": "",
		"messages": {
		  "title": {
			"i18n": {
			  "en-US": "Transfer Tokens"
			}
		  },
		  "description": {
			"i18n": {
			  "en-US": "Transfer tokens from one account to another"
			}
		  }
		},
		"cadence": "import FungibleToken from 0xFUNGIBLETOKENADDRESS\ntransaction(amount: UFix64, to: Address) {\nlet vault: @FungibleToken.Vault\nprepare(signer: AuthAccount) {\nself.vault <- signer\n.borrow<&{FungibleToken.Provider}>(from: /storage/flowTokenVault)!\n.withdraw(amount: amount)\n}\nexecute {\ngetAccount(to)\n.getCapability(/public/flowTokenReceiver)!\n.borrow<&{FungibleToken.Receiver}>()!\n.deposit(from: <-self.vault)\n}\n}",
		"dependencies": {
		  "0xFUNGIBLETOKENADDRESS": {
			"FungibleToken": {
			  "mainnet": {
				"address": "0xf233dcee88fe0abe",
				"fq_address": "A.0xf233dcee88fe0abe.FungibleToken",
				"contract": "FungibleToken",
				"pin": "83c9e3d61d3b5ebf24356a9f17b5b57b12d6d56547abc73e05f820a0ae7d9cf5",
				"pin_block_height": 34166296
			  },
			  "testnet": {
				"address": "0x9a0766d93b6608b7",
				"fq_address": "A.0x9a0766d93b6608b7.FungibleToken",
				"contract": "FungibleToken",
				"pin": "83c9e3d61d3b5ebf24356a9f17b5b57b12d6d56547abc73e05f820a0ae7d9cf5",
				"pin_block_height": 74776482
			  }
			}
		  }
		},
		"arguments": {
		  "amount": {
			"index": 0,
			"type": "UFix64",
			"messages": {
			  "title": {
				"i18n": {
				  "en-US": "The amount of FLOW tokens to send"
				}
			  }
			},
			"balance": ""
		  },
		  "to": {
			"index": 1,
			"type": "Address",
			"messages": {
			  "title": {
				"i18n": {
				  "en-US": "The Flow account the tokens will go to"
				}
			  }
			},
			"balance": ""
		  }
		}
	  }
	}`

var parsedTemplate = &FlowInteractionTemplate{
	FType:    "InteractionTemplate",
	FVersion: "1.0.0",
	ID:       "290b6b6222b2a77b16db896a80ddf29ebd1fa3038c9e6625a933fa213fce51fa",
	Data: Data{
		Type:      "transaction",
		Interface: "",
		Messages: Messages{
			Title: &Title{
				I18N: map[string]string{
					"en-US": "Transfer Tokens",
				},
			},
			Description: &Description{
				I18N: map[string]string{
					"en-US": "Transfer tokens from one account to another",
				},
			},
		},
		Cadence: "import FungibleToken from 0xFUNGIBLETOKENADDRESS\ntransaction(amount: UFix64, to: Address) {\nlet vault: @FungibleToken.Vault\nprepare(signer: AuthAccount) {\nself.vault <- signer\n.borrow<&{FungibleToken.Provider}>(from: /storage/flowTokenVault)!\n.withdraw(amount: amount)\n}\nexecute {\ngetAccount(to)\n.getCapability(/public/flowTokenReceiver)!\n.borrow<&{FungibleToken.Receiver}>()!\n.deposit(from: <-self.vault)\n}\n}",
		Dependencies: Dependencies{
			"0xFUNGIBLETOKENADDRESS": Contracts{
				"FungibleToken": Networks{
					"mainnet": Network{
						Address:        "0xf233dcee88fe0abe",
						FqAddress:      "A.0xf233dcee88fe0abe.FungibleToken",
						Contract:       "FungibleToken",
						Pin:            "83c9e3d61d3b5ebf24356a9f17b5b57b12d6d56547abc73e05f820a0ae7d9cf5",
						PinBlockHeight: 34166296,
					},
					"testnet": Network{
						Address:        "0x9a0766d93b6608b7",
						FqAddress:      "A.0x9a0766d93b6608b7.FungibleToken",
						Contract:       "FungibleToken",
						Pin:            "83c9e3d61d3b5ebf24356a9f17b5b57b12d6d56547abc73e05f820a0ae7d9cf5",
						PinBlockHeight: 74776482,
					},
				},
			},
		},
		Arguments: Arguments{
			"amount": Argument{
				Index: 0,
				Type:  "UFix64",
				Messages: Messages{
					Title: &Title{
						I18N: map[string]string{
							"en-US": "The amount of FLOW tokens to send",
						},
					},
				},
				Balance: "",
			},
			"to": Argument{
				Index: 1,
				Type:  "Address",
				Messages: Messages{
					Title: &Title{
						I18N: map[string]string{
							"en-US": "The Flow account the tokens will go to",
						},
					},
				},
				Balance: "",
			},
		},
	},
}

func TestParseFlix(t *testing.T) {
	assert := assert.New(t)

	parsedTemplate, err := ParseFlix(template)
	assert.NoError(err, "ParseTemplate should not return an error")
	assert.NotNil(parsedTemplate, "Parsed template should not be nil")

	expectedType := "transaction"
	assert.Equal(expectedType, parsedTemplate.Data.Type, "Parsed template should have the correct type")
}

func TestGetAndReplaceCadenceImports(t *testing.T) {
	assert := assert.New(t)

	tests := []struct {
		name       string
		network    string
		wantErr    bool
		wantImport string
	}{
		{
			name:       "Mainnet",
			network:    "mainnet",
			wantErr:    false,
			wantImport: "import FungibleToken from 0xf233dcee88fe0abe",
		},
		{
			name:       "Testnet",
			network:    "testnet",
			wantErr:    false,
			wantImport: "import FungibleToken from 0x9a0766d93b6608b7",
		},
		{
			name:    "MissingNetwork",
			network: "missing",
			wantErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			cadence, err := parsedTemplate.GetAndReplaceCadenceImports(tt.network)
			if tt.wantErr {
				assert.Error(err, "GetCadenceWithReplacedImports should return an error")
			} else {
				assert.NoError(err, "GetCadenceWithReplacedImports should not return an error")
				assert.NotEmpty(cadence, "Cadence should not be empty")

				assert.Contains(cadence, tt.wantImport, "Cadence should contain the expected import")
			}
		})
	}
}

func TestIsScript(t *testing.T) {
	assert := assert.New(t)

	scriptTemplate := &FlowInteractionTemplate{
		Data: Data{
			Type: "script",
		},
	}
	assert.True(scriptTemplate.IsScript(), "IsScript() should return true")

	transactionTemplate := &FlowInteractionTemplate{
		Data: Data{
			Type: "transaction",
		},
	}
	assert.False(transactionTemplate.IsScript(), "IsScript() should return false")
}

func TestIsTransaction(t *testing.T) {
	assert := assert.New(t)

	scriptTemplate := &FlowInteractionTemplate{
		Data: Data{
			Type: "script",
		},
	}
	assert.False(scriptTemplate.IsTransaction(), "IsTransaction() should return false")

	transactionTemplate := &FlowInteractionTemplate{
		Data: Data{
			Type: "transaction",
		},
	}
	assert.True(transactionTemplate.IsTransaction(), "IsTransaction() should return true")
}

func TestGetFlix(t *testing.T) {
	assert := assert.New(t)

	server := httptest.NewServer(http.HandlerFunc(func(rw http.ResponseWriter, req *http.Request) {
		rw.Write([]byte("Hello World"))
	}))
	defer server.Close()

	body, err := GetFlix(server.URL)
	assert.NoError(err, "GetFlix should not return an error")
	assert.Equal("Hello World", body, "GetFlix should return the correct body")
}

func TestGetFlixByID(t *testing.T) {
	assert := assert.New(t)

	server := httptest.NewServer(http.HandlerFunc(func(rw http.ResponseWriter, req *http.Request) {
		assert.Equal("/templateID", req.URL.String(), "GetFlixByID should request the correct path")
		rw.Write([]byte("Hello World"))
	}))
	defer server.Close()

	body, err := GetFlixByID(server.URL, "templateID")
	assert.NoError(err, "GetFlixByID should not return an error")
	assert.Equal("Hello World", body, "GetFlixByID should return the correct body")
}

func TestGetFlixByName(t *testing.T) {
	assert := assert.New(t)

	server := httptest.NewServer(http.HandlerFunc(func(rw http.ResponseWriter, req *http.Request) {
		assert.Equal("/?name=templateName", req.URL.String(), "GetFlixByName should request the correct query string")
		rw.Write([]byte("Hello World"))
	}))
	defer server.Close()

	body, err := GetFlixByName(server.URL, "templateName")
	assert.NoError(err, "GetFlixByName should not return an error")
	assert.Equal("Hello World", body, "GetFlixByName should return the correct body")
}

func TestGetParsedFlixByID(t *testing.T) {
	assert := assert.New(t)

	server := httptest.NewServer(http.HandlerFunc(func(rw http.ResponseWriter, req *http.Request) {
		rw.Write([]byte(template))
	}))
	defer server.Close()

	flix, err := GetParsedFlixByID(server.URL, "templateID")
	assert.NoError(err, "GetParsedFlixByID should not return an error")
	assert.NotNil(flix, "GetParsedFlixByID should not return a nil Flix")
	assert.Equal(parsedTemplate, flix, "GetParsedFlixByID should return the correct Flix")
}

func TestGetParsedFlixByName(t *testing.T) {
	assert := assert.New(t)

	server := httptest.NewServer(http.HandlerFunc(func(rw http.ResponseWriter, req *http.Request) {
		rw.Write([]byte(template))
	}))
	defer server.Close()

	flix, err := GetParsedFlixByName(server.URL, "templateName")
	assert.NoError(err, "GetParsedFlixByName should not return an error")
	assert.NotNil(flix, "GetParsedFlixByName should not return a nil Flix")
	assert.Equal(parsedTemplate, flix, "GetParsedFlixByName should return the correct Flix")
}
