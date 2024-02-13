package internal

import (
	"context"
	"net/http"
	"net/http/httptest"
	"testing"

	v1 "github.com/onflow/flixkit-go/internal/v1"
	"github.com/onflow/flixkit-go/internal/v1_1"
	"github.com/stretchr/testify/assert"
)

var flix_template = `{
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

var parsedTemplate = &v1.FlowInteractionTemplate{
	FType:    "InteractionTemplate",
	FVersion: "1.0.0",
	ID:       "290b6b6222b2a77b16db896a80ddf29ebd1fa3038c9e6625a933fa213fce51fa",
	Data: v1.Data{
		Type:      "transaction",
		Interface: "",
		Messages: v1.Messages{
			Title: &v1.Title{
				I18N: map[string]string{
					"en-US": "Transfer Tokens",
				},
			},
			Description: &v1.Description{
				I18N: map[string]string{
					"en-US": "Transfer tokens from one account to another",
				},
			},
		},
		Cadence: "import FungibleToken from 0xFUNGIBLETOKENADDRESS\ntransaction(amount: UFix64, to: Address) {\nlet vault: @FungibleToken.Vault\nprepare(signer: AuthAccount) {\nself.vault <- signer\n.borrow<&{FungibleToken.Provider}>(from: /storage/flowTokenVault)!\n.withdraw(amount: amount)\n}\nexecute {\ngetAccount(to)\n.getCapability(/public/flowTokenReceiver)!\n.borrow<&{FungibleToken.Receiver}>()!\n.deposit(from: <-self.vault)\n}\n}",
		Dependencies: v1.Dependencies{
			"0xFUNGIBLETOKENADDRESS": v1.Contracts{
				"FungibleToken": v1.Networks{
					"mainnet": v1.Network{
						Address:        "0xf233dcee88fe0abe",
						FqAddress:      "A.0xf233dcee88fe0abe.FungibleToken",
						Contract:       "FungibleToken",
						Pin:            "83c9e3d61d3b5ebf24356a9f17b5b57b12d6d56547abc73e05f820a0ae7d9cf5",
						PinBlockHeight: 34166296,
					},
					"testnet": v1.Network{
						Address:        "0x9a0766d93b6608b7",
						FqAddress:      "A.0x9a0766d93b6608b7.FungibleToken",
						Contract:       "FungibleToken",
						Pin:            "83c9e3d61d3b5ebf24356a9f17b5b57b12d6d56547abc73e05f820a0ae7d9cf5",
						PinBlockHeight: 74776482,
					},
				},
			},
		},
		Arguments: v1.Arguments{
			"amount": v1.Argument{
				Index: 0,
				Type:  "UFix64",
				Messages: v1.Messages{
					Title: &v1.Title{
						I18N: map[string]string{
							"en-US": "The amount of FLOW tokens to send",
						},
					},
				},
				Balance: "",
			},
			"to": v1.Argument{
				Index: 1,
				Type:  "Address",
				Messages: v1.Messages{
					Title: &v1.Title{
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

	parsedTemplate, err := v1.ParseFlix(flix_template)
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
			cadence, err := parsedTemplate.ReplaceCadenceImports(tt.network)
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

func TestGetAndReplaceCadenceSimpleImportsV1(t *testing.T) {
	assert := assert.New(t)

	tests := []struct {
		name    string
		network string
		wantErr bool
		source  *v1.FlowInteractionTemplate
		result  string
	}{
		{
			name:    "Mainnet",
			network: "mainnet",
			wantErr: false,
			source: &v1.FlowInteractionTemplate{
				Data: v1.Data{
					Cadence: "access(all) fun main(x: Int, y: Int): Int { return x * y }",
				},
			},
			result: "access(all) fun main(x: Int, y: Int): Int { return x * y }",
		},
		{
			name:    "Testnet",
			network: "testnet",
			wantErr: false,
			source: &v1.FlowInteractionTemplate{
				Data: v1.Data{
					Cadence: "import FungibleToken from 0xFUNGIBLETOKENADDRESS",
					Dependencies: v1.Dependencies{
						"0xFUNGIBLETOKENADDRESS": v1.Contracts{
							"FungibleToken": v1.Networks{
								"testnet": v1.Network{
									Address: "0x9a0766d93b6608b7",
								},
							},
						},
					},
				},
			},
			result: "import FungibleToken from 0x9a0766d93b6608b7",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			cadence, err := tt.source.ReplaceCadenceImports(tt.network)
			if tt.wantErr {
				assert.Error(err, "GetCadenceWithReplacedImports should return an error")
			} else {
				assert.NoError(err, "GetCadenceWithReplacedImports should not return an error")
				assert.NotEmpty(cadence, "Cadence should not be empty")

				assert.Contains(cadence, tt.result, "Cadence should contain the expected import")
			}
		})
	}
}

func TestGetAndReplaceCadenceSimpleImportsV11(t *testing.T) {
	assert := assert.New(t)

	tests := []struct {
		name    string
		network string
		wantErr bool
		source  *v1_1.InteractionTemplate
		result  string
	}{
		{
			name:    "Mainnet",
			network: "mainnet",
			wantErr: false,
			source: &v1_1.InteractionTemplate{
				FType:    "InteractionTemplate",
				FVersion: "1.1.0",
				Data: v1_1.Data{
					Cadence: v1_1.Cadence{
						Body: "access(all) fun main(x: Int, y: Int): Int { return x * y }",
					},
				},
			},
			result: "access(all) fun main(x: Int, y: Int): Int { return x * y }",
		},
		{
			name:    "Testnet",
			network: "testnet",
			wantErr: false,
			source: &v1_1.InteractionTemplate{
				Data: v1_1.Data{
					Cadence: v1_1.Cadence{
						Body: "import \"FungibleToken\"",
					},
					Dependencies: []v1_1.Dependency{
						{
							Contracts: []v1_1.Contract{
								{
									Contract: "FungibleToken",
									Networks: []v1_1.Network{
										{
											Network: "testnet",
											Address: "0x9a0766d93b6608b7",
										},
									},
								},
							},
						},
					},
				},
			},
			result: "import FungibleToken from 0x9a0766d93b6608b7",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			cadence, err := tt.source.ReplaceCadenceImports(tt.network)
			if tt.wantErr {
				assert.Error(err, "GetCadenceWithReplacedImports should return an error")
			} else {
				assert.NoError(err, "GetCadenceWithReplacedImports should not return an error")
				assert.NotEmpty(cadence, "Cadence should not be empty")

				assert.Contains(cadence, tt.result, "Cadence should contain the expected import")
			}
		})
	}
}

func TestIsScript(t *testing.T) {
	assert := assert.New(t)

	scriptTemplate := &v1.FlowInteractionTemplate{
		Data: v1.Data{
			Type: "script",
		},
	}
	assert.True(scriptTemplate.IsScript(), "IsScript() should return true")

	transactionTemplate := &v1.FlowInteractionTemplate{
		Data: v1.Data{
			Type: "transaction",
		},
	}
	assert.False(transactionTemplate.IsScript(), "IsScript() should return false")
}

func TestIsTransaction(t *testing.T) {
	assert := assert.New(t)

	scriptTemplate := &v1.FlowInteractionTemplate{
		Data: v1.Data{
			Type: "script",
		},
	}
	assert.False(scriptTemplate.IsTransaction(), "IsTransaction() should return false")

	transactionTemplate := &v1.FlowInteractionTemplate{
		Data: v1.Data{
			Type: "transaction",
		},
	}
	assert.True(transactionTemplate.IsTransaction(), "IsTransaction() should return true")
}

func TestFetchFlix(t *testing.T) {
	assert := assert.New(t)

	server := httptest.NewServer(http.HandlerFunc(func(rw http.ResponseWriter, req *http.Request) {
		rw.Write([]byte("Hello World"))
	}))
	defer server.Close()

	ctx := context.Background()
	body, source, err := fetchFlixWithContext(ctx, server.URL)
	assert.NoError(err, "GetFlix should not return an error")
	assert.Equal("Hello World", body, "GetFlix should return the correct body")
	assert.Equal(server.URL, source, "GetFlix should return the correct source")
}

type DefaultReader struct{}

func (d DefaultReader) ReadFile(path string) ([]byte, error) {
	return []byte(flix_template), nil
}

func TestGetFlixRaw(t *testing.T) {
	assert := assert.New(t)

	server := httptest.NewServer(http.HandlerFunc(func(rw http.ResponseWriter, req *http.Request) {
		assert.Equal("/?name=templateName", req.URL.String(), "GetFlixByName should request the correct query string")
		rw.Write([]byte("Hello World"))
	}))
	defer server.Close()

	flixService := NewFlixService(&FlixServiceConfig{FlixServerURL: server.URL, FileReader: nil})
	ctx := context.Background()
	body, source, err := flixService.GetTemplate(ctx, "templateName")
	assert.NoError(err, "GetFlixByName should not return an error")
	assert.Equal("Hello World", body, "GetFlixByName should return the correct body")
	assert.Equal(server.URL+"?name=templateName", source, "GetFlixByName should return the correct source")
}

func TestGetFlixFilename(t *testing.T) {
	assert := assert.New(t)

	server := httptest.NewServer(http.HandlerFunc(func(rw http.ResponseWriter, req *http.Request) {
		rw.Write([]byte(flix_template))
	}))
	defer server.Close()
	flixService := NewFlixService(&FlixServiceConfig{FlixServerURL: server.URL, FileReader: DefaultReader{}})
	ctx := context.Background()
	flix, source, err := flixService.GetTemplate(ctx, "./templateFileName")
	assert.NoError(err, "GetParsedFlixByName should not return an error")
	assert.NotNil(flix, "GetParsedFlixByName should not return a nil Flix")
	assert.Equal(flix_template, flix, "GetParsedFlixByName should return the correct Flix")
	assert.Equal("./templateFileName", source, "GetParsedFlixByName should return the correct source")
}

func TestGetFlixByIDRaw(t *testing.T) {
	assert := assert.New(t)

	server := httptest.NewServer(http.HandlerFunc(func(rw http.ResponseWriter, req *http.Request) {
		assert.Equal("/1234567890ABCDEF1234567890ABCDEF1234567890ABCDEF1234567890ABCDEF", req.URL.String(), "GetFlixByID should request the correct path")
		rw.Write([]byte("Hello World"))
	}))
	defer server.Close()

	flixService := NewFlixService(&FlixServiceConfig{FlixServerURL: server.URL})
	ctx := context.Background()
	body, source, err := flixService.GetTemplate(ctx, "1234567890ABCDEF1234567890ABCDEF1234567890ABCDEF1234567890ABCDEF")
	assert.NoError(err, "GetFlixByID should not return an error")
	assert.Equal("Hello World", body, "GetFlixByID should return the correct body")
	assert.Equal(server.URL+"/1234567890ABCDEF1234567890ABCDEF1234567890ABCDEF1234567890ABCDEF", source, "GetFlixByID should return the correct source")
}

func TestGetFlixByID(t *testing.T) {
	assert := assert.New(t)

	server := httptest.NewServer(http.HandlerFunc(func(rw http.ResponseWriter, req *http.Request) {
		rw.Write([]byte(flix_template))
	}))
	defer server.Close()

	flixService := NewFlixService(&FlixServiceConfig{FlixServerURL: server.URL})
	ctx := context.Background()
	flix, source, err := flixService.GetTemplate(ctx, "templateID")
	assert.NoError(err, "GetParsedFlixByID should not return an error")
	assert.NotNil(flix, "GetParsedFlixByID should not return a nil Flix")
	assert.Equal(flix_template, flix, "GetParsedFlixByID should return the correct Flix")
	assert.Equal(server.URL+"?name=templateID", source, "GetParsedFlixByID should return the correct source")
}

func TestTemplateVersion(t *testing.T) {
	assert := assert.New(t)

	tests := []struct {
		templateStr string
		version     string
		wantErr     bool
	}{
		{
			templateStr: `{
				"f_version": "1.0.0"
			}`,
			version: "1.0.0",
			wantErr: false,
		},
		{
			templateStr: `{
				"f_version": "1.1.0"
			}`,
			version: "1.1.0",
			wantErr: false,
		},
		{
			templateStr: `{
				"f_ver": "1.x"
			}`,
			version: "",
			wantErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.templateStr, func(t *testing.T) {
			ver, err := getTemplateVersion(tt.templateStr)
			if tt.wantErr {
				assert.Error(err, "TemplateVersion should return an error")
			} else {
				assert.NoError(err, "TemplateVersion should not return an error")
				assert.Equal(tt.version, ver, "TemplateVersion should return the correct version")
			}
		})
	}
}

func TestParseFlixV1(t *testing.T) {
	temp := `{
		"f_type": "InteractionTemplate",
		"f_version": "1.0.0",
		"id": "bd10ab0bf472e6b58ecc0398e9b3d1bd58a4205f14a7099c52c0640d9589295f",
		"data": {
		  "type": "script",
		  "interface": "",
		  "messages": {
			"title": {
			  "i18n": {
				"en-US": "Multiply Two Integers"
			  }
			},
			"description": {
			  "i18n": {
				"en-US": "Multiplies two integer arguments together and returns the result."
			  }
			}
		  },
		  "cadence": "pub fun main(x: Int, y: Int): Int { return x * y }",
		  "dependencies": {},
		  "arguments": {
			"x": {
			  "index": 0,
			  "type": "Int",
			  "messages": {
				"title": {
				  "i18n": {
					"en-US": "Int 1"
				  }
				}
			  }
			},
			"y": {
			  "index": 1,
			  "type": "Int",
			  "messages": {
				"title": {
				  "i18n": {
					"en-US": "Int 2"
				  }
				}
			  }
			}
		  }
		}
	  }`
	assert := assert.New(t)

	parsedTemplate, err := v1.ParseFlix(temp)
	assert.NoError(err, "ParseTemplate should not return an error")
	assert.NotNil(parsedTemplate, "Parsed template should not be nil")

	expectedType := "script"
	assert.Equal(expectedType, parsedTemplate.Data.Type, "Parsed template should have the correct type")
	v, err := parsedTemplate.ReplaceCadenceImports("mainnet")
	assert.NoError(err, "ReplaceCadenceImports should not return an error")
	assert.Equal("pub fun main(x: Int, y: Int): Int { return x * y }", v, "ReplaceCadenceImports should return the correct cadence")
}
