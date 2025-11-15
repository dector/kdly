package serializer

import (
	"testing"

	"github.com/dector/kdly/pkg/parser"
	"github.com/stretchr/testify/assert"
)

// TestSerializeServerConfig tests serialization of server configuration with nested TLS settings
func TestSerializeServerConfig(t *testing.T) {
	doc := &parser.Document{
		Nodes: []parser.Node{
			{
				Name:       "server",
				Arguments:  []parser.Value{},
				Properties: []parser.Property{},
				Children: []parser.Node{
					{
						Name: "host",
						Arguments: []parser.Value{
							{Type: parser.ValueTypeString, Value: "0.0.0.0"},
						},
						Properties: []parser.Property{},
						Children:   []parser.Node{},
					},
					{
						Name: "port",
						Arguments: []parser.Value{
							{Type: parser.ValueTypeNumber, Value: "8080"},
						},
						Properties: []parser.Property{},
						Children:   []parser.Node{},
					},
					{
						Name:      "tls",
						Arguments: []parser.Value{},
						Properties: []parser.Property{
							{Key: "enabled", Value: parser.Value{Type: parser.ValueTypeBoolean, Value: "true"}},
						},
						Children: []parser.Node{
							{
								Name: "cert-file",
								Arguments: []parser.Value{
									{Type: parser.ValueTypeString, Value: `C:\certs\server.crt`},
								},
								Properties: []parser.Property{},
								Children:   []parser.Node{},
							},
							{
								Name: "key-file",
								Arguments: []parser.Value{
									{Type: parser.ValueTypeString, Value: `C:\certs\server.key`},
								},
								Properties: []parser.Property{},
								Children:   []parser.Node{},
							},
						},
					},
				},
			},
		},
	}

	result := ToKDL(doc)

	// Expected output with proper formatting
	expected := `server {
  host "0.0.0.0"
  port 8080
  tls enabled=#true {
    cert-file #"C:\certs\server.crt"#
    key-file #"C:\certs\server.key"#
  }
}`

	assert.Equal(t, expected, result)

	// Round-trip test: parse the serialized output
	roundTrip, err := parser.New().Parse(result)
	assert.NoError(t, err)
	assert.Equal(t, doc, roundTrip)
}

// TestSerializeDatabaseConfig tests database config with type annotations
func TestSerializeDatabaseConfig(t *testing.T) {
	doc := &parser.Document{
		Nodes: []parser.Node{
			{
				Name:       "database",
				Arguments:  []parser.Value{},
				Properties: []parser.Property{},
				Children: []parser.Node{
					{
						Name: "host",
						Arguments: []parser.Value{
							{Type: parser.ValueTypeString, Value: "db.example.com"},
						},
						Properties: []parser.Property{},
						Children:   []parser.Node{},
					},
					{
						Name: "port",
						Arguments: []parser.Value{
							{Type: parser.ValueTypeNumber, Value: "5432"},
						},
						Properties: []parser.Property{},
						Children:   []parser.Node{},
					},
					{
						Name:       "pool",
						Arguments:  []parser.Value{},
						Properties: []parser.Property{},
						Children: []parser.Node{
							{
								Name: "min-size",
								Arguments: []parser.Value{
									{Type: parser.ValueTypeNumber, Value: "5", TypeAnnotation: "int"},
								},
								Properties: []parser.Property{},
								Children:   []parser.Node{},
							},
							{
								Name: "max-size",
								Arguments: []parser.Value{
									{Type: parser.ValueTypeNumber, Value: "20", TypeAnnotation: "int"},
								},
								Properties: []parser.Property{},
								Children:   []parser.Node{},
							},
							{
								Name: "timeout",
								Arguments: []parser.Value{
									{Type: parser.ValueTypeString, Value: "30s", TypeAnnotation: "duration"},
								},
								Properties: []parser.Property{},
								Children:   []parser.Node{},
							},
						},
					},
				},
			},
		},
	}

	result := ToKDL(doc)

	expected := `database {
  host db.example.com
  port 5432
  pool {
    min-size (int)5
    max-size (int)20
    timeout (duration)"30s"
  }
}`

	assert.Equal(t, expected, result)

	// Round-trip test
	roundTrip, err := parser.New().Parse(result)
	assert.NoError(t, err)
	assert.Equal(t, doc, roundTrip)
}

// TestSerializePackageWithDependencies tests package manifest with bare identifiers
func TestSerializePackageWithDependencies(t *testing.T) {
	doc := &parser.Document{
		Nodes: []parser.Node{
			{
				Name:       "package",
				Arguments:  []parser.Value{},
				Properties: []parser.Property{},
				Children: []parser.Node{
					{
						Name: "name",
						Arguments: []parser.Value{
							{Type: parser.ValueTypeString, Value: "awesome-app"},
						},
						Properties: []parser.Property{},
						Children:   []parser.Node{},
					},
					{
						Name: "version",
						Arguments: []parser.Value{
							{Type: parser.ValueTypeString, Value: "1.2.3"},
						},
						Properties: []parser.Property{},
						Children:   []parser.Node{},
					},
					{
						Name: "license",
						Arguments: []parser.Value{
							{Type: parser.ValueTypeString, Value: "MIT"},
						},
						Properties: []parser.Property{},
						Children:   []parser.Node{},
					},
					{
						Name:       "dependencies",
						Arguments:  []parser.Value{},
						Properties: []parser.Property{},
						Children: []parser.Node{
							{
								Name: "dep",
								Arguments: []parser.Value{
									{Type: parser.ValueTypeString, Value: "http-client"},
								},
								Properties: []parser.Property{
									{Key: "version", Value: parser.Value{Type: parser.ValueTypeString, Value: "^2.0.0"}},
									{Key: "optional", Value: parser.Value{Type: parser.ValueTypeBoolean, Value: "false"}},
								},
								Children: []parser.Node{},
							},
							{
								Name: "dep",
								Arguments: []parser.Value{
									{Type: parser.ValueTypeString, Value: "logger"},
								},
								Properties: []parser.Property{
									{Key: "version", Value: parser.Value{Type: parser.ValueTypeString, Value: "~1.5.0"}},
								},
								Children: []parser.Node{},
							},
						},
					},
				},
			},
		},
	}

	result := ToKDL(doc)

	expected := `package {
  name awesome-app
  version "1.2.3"
  license MIT
  dependencies {
    dep http-client version=^2.0.0 optional=#false
    dep logger version=~1.5.0
  }
}`

	assert.Equal(t, expected, result)

	// Round-trip test
	roundTrip, err := parser.New().Parse(result)
	assert.NoError(t, err)
	assert.Equal(t, doc, roundTrip)
}

// TestSerializeCIPipelineStage tests CI config with properties
func TestSerializeCIPipelineStage(t *testing.T) {
	doc := &parser.Document{
		Nodes: []parser.Node{
			{
				Name:       "pipeline",
				Arguments:  []parser.Value{},
				Properties: []parser.Property{},
				Children: []parser.Node{
					{
						Name: "stage",
						Arguments: []parser.Value{
							{Type: parser.ValueTypeString, Value: "build"},
						},
						Properties: []parser.Property{},
						Children: []parser.Node{
							{
								Name: "image",
								Arguments: []parser.Value{
									{Type: parser.ValueTypeString, Value: "golang:1.21"},
								},
								Properties: []parser.Property{},
								Children:   []parser.Node{},
							},
							{
								Name:       "script",
								Arguments:  []parser.Value{},
								Properties: []parser.Property{},
								Children: []parser.Node{
									{
										Name: "run",
										Arguments: []parser.Value{
											{Type: parser.ValueTypeString, Value: "go build ./..."},
										},
										Properties: []parser.Property{},
										Children:   []parser.Node{},
									},
									{
										Name: "run",
										Arguments: []parser.Value{
											{Type: parser.ValueTypeString, Value: "go test ./..."},
										},
										Properties: []parser.Property{},
										Children:   []parser.Node{},
									},
								},
							},
						},
					},
					{
						Name: "stage",
						Arguments: []parser.Value{
							{Type: parser.ValueTypeString, Value: "deploy"},
						},
						Properties: []parser.Property{
							{Key: "depends-on", Value: parser.Value{Type: parser.ValueTypeString, Value: "build"}},
							{Key: "enabled", Value: parser.Value{Type: parser.ValueTypeBoolean, Value: "true"}},
						},
						Children: []parser.Node{},
					},
				},
			},
		},
	}

	result := ToKDL(doc)

	expected := `pipeline {
  stage build {
    image golang:1.21
    script {
      run "go build ./..."
      run "go test ./..."
    }
  }
  stage deploy depends-on=build enabled=#true
}`

	assert.Equal(t, expected, result)

	// Round-trip test
	roundTrip, err := parser.New().Parse(result)
	assert.NoError(t, err)
	assert.Equal(t, doc, roundTrip)
}

// TestSerializeLoggingConfig tests logging with hex numbers and type annotations
func TestSerializeLoggingConfig(t *testing.T) {
	doc := &parser.Document{
		Nodes: []parser.Node{
			{
				Name:       "logging",
				Arguments:  []parser.Value{},
				Properties: []parser.Property{},
				Children: []parser.Node{
					{
						Name: "level",
						Arguments: []parser.Value{
							{Type: parser.ValueTypeString, Value: "info"},
						},
						Properties: []parser.Property{},
						Children:   []parser.Node{},
					},
					{
						Name: "format",
						Arguments: []parser.Value{
							{Type: parser.ValueTypeString, Value: "json"},
						},
						Properties: []parser.Property{},
						Children:   []parser.Node{},
					},
					{
						Name:       "rotation",
						Arguments:  []parser.Value{},
						Properties: []parser.Property{},
						Children: []parser.Node{
							{
								Name: "max-size",
								Arguments: []parser.Value{
									{Type: parser.ValueTypeNumber, Value: "0x6400000", TypeAnnotation: "bytes"},
								},
								Properties: []parser.Property{},
								Children:   []parser.Node{},
							},
							{
								Name: "compress",
								Arguments: []parser.Value{
									{Type: parser.ValueTypeBoolean, Value: "true"},
								},
								Properties: []parser.Property{},
								Children:   []parser.Node{},
							},
						},
					},
					{
						Name:       "loggers",
						Arguments:  []parser.Value{},
						Properties: []parser.Property{},
						Children: []parser.Node{
							{
								Name: "logger",
								Arguments: []parser.Value{
									{Type: parser.ValueTypeString, Value: "http"},
								},
								Properties: []parser.Property{
									{Key: "level", Value: parser.Value{Type: parser.ValueTypeString, Value: "debug"}},
								},
								Children: []parser.Node{},
							},
							{
								Name: "logger",
								Arguments: []parser.Value{
									{Type: parser.ValueTypeString, Value: "database"},
								},
								Properties: []parser.Property{
									{Key: "level", Value: parser.Value{Type: parser.ValueTypeString, Value: "warn"}},
								},
								Children: []parser.Node{},
							},
						},
					},
				},
			},
		},
	}

	result := ToKDL(doc)

	expected := `logging {
  level info
  format json
  rotation {
    max-size (bytes)0x6400000
    compress #true
  }
  loggers {
    logger http level=debug
    logger database level=warn
  }
}`

	assert.Equal(t, expected, result)

	// Round-trip test
	roundTrip, err := parser.New().Parse(result)
	assert.NoError(t, err)
	assert.Equal(t, doc, roundTrip)
}

// TestSerializeMetricsWithMultipleArguments tests multiple arguments
func TestSerializeMetricsWithMultipleArguments(t *testing.T) {
	doc := &parser.Document{
		Nodes: []parser.Node{
			{
				Name:       "metrics",
				Arguments:  []parser.Value{},
				Properties: []parser.Property{},
				Children: []parser.Node{
					{
						Name: "enabled",
						Arguments: []parser.Value{
							{Type: parser.ValueTypeBoolean, Value: "true"},
						},
						Properties: []parser.Property{},
						Children:   []parser.Node{},
					},
					{
						Name: "provider",
						Arguments: []parser.Value{
							{Type: parser.ValueTypeString, Value: "prometheus"},
						},
						Properties: []parser.Property{},
						Children:   []parser.Node{},
					},
					{
						Name: "port",
						Arguments: []parser.Value{
							{Type: parser.ValueTypeNumber, Value: "9090"},
						},
						Properties: []parser.Property{},
						Children:   []parser.Node{},
					},
					{
						Name: "histogram",
						Arguments: []parser.Value{
							{Type: parser.ValueTypeString, Value: "http_duration"},
							{Type: parser.ValueTypeNumber, Value: "0.005"},
							{Type: parser.ValueTypeNumber, Value: "0.01"},
							{Type: parser.ValueTypeNumber, Value: "0.025"},
						},
						Properties: []parser.Property{},
						Children:   []parser.Node{},
					},
				},
			},
		},
	}

	result := ToKDL(doc)

	expected := `metrics {
  enabled #true
  provider prometheus
  port 9090
  histogram http_duration 0.005 0.01 0.025
}`

	assert.Equal(t, expected, result)

	// Round-trip test
	roundTrip, err := parser.New().Parse(result)
	assert.NoError(t, err)
	assert.Equal(t, doc, roundTrip)
}

// TestSerializeEdgeCases tests various edge cases
func TestSerializeEdgeCases(t *testing.T) {
	t.Run("empty document", func(t *testing.T) {
		doc := &parser.Document{Nodes: []parser.Node{}}
		result := ToKDL(doc)
		assert.Equal(t, "", result)
	})

	t.Run("node with no arguments or properties", func(t *testing.T) {
		doc := &parser.Document{
			Nodes: []parser.Node{
				{Name: "simple", Arguments: []parser.Value{}, Properties: []parser.Property{}, Children: []parser.Node{}},
			},
		}
		result := ToKDL(doc)
		assert.Equal(t, "simple", result)
	})

	t.Run("string with special characters", func(t *testing.T) {
		doc := &parser.Document{
			Nodes: []parser.Node{
				{
					Name: "node",
					Arguments: []parser.Value{
						{Type: parser.ValueTypeString, Value: "hello\nworld\ttab"},
					},
					Properties: []parser.Property{},
					Children:   []parser.Node{},
				},
			},
		}
		result := ToKDL(doc)
		assert.Equal(t, `node "hello\nworld\ttab"`, result)
	})

	t.Run("node with type annotation", func(t *testing.T) {
		doc := &parser.Document{
			Nodes: []parser.Node{
				{
					Name:           "person",
					TypeAnnotation: "author",
					Arguments:      []parser.Value{},
					Properties:     []parser.Property{},
					Children:       []parser.Node{},
				},
			},
		}
		result := ToKDL(doc)
		assert.Equal(t, "(author)person", result)
	})

	t.Run("null value", func(t *testing.T) {
		doc := &parser.Document{
			Nodes: []parser.Node{
				{
					Name: "value",
					Arguments: []parser.Value{
						{Type: parser.ValueTypeNull},
					},
					Properties: []parser.Property{},
					Children:   []parser.Node{},
				},
			},
		}
		result := ToKDL(doc)
		assert.Equal(t, "value #null", result)
	})

	t.Run("multiple top-level nodes", func(t *testing.T) {
		doc := &parser.Document{
			Nodes: []parser.Node{
				{Name: "node1", Arguments: []parser.Value{}, Properties: []parser.Property{}, Children: []parser.Node{}},
				{Name: "node2", Arguments: []parser.Value{}, Properties: []parser.Property{}, Children: []parser.Node{}},
			},
		}
		result := ToKDL(doc)
		assert.Equal(t, "node1\nnode2", result)
	})
}

// TestRoundTrip tests that parse -> serialize -> parse produces identical results
func TestRoundTrip(t *testing.T) {
	inputs := []string{
		`server {
  host "0.0.0.0"
  port 8080
}`,
		`node arg1 arg2 key1=val1 key2=val2`,
		`(type)node (int)42`,
		`parent {
  child1
  child2 {
    grandchild
  }
}`,
	}

	for _, input := range inputs {
		t.Run(input, func(t *testing.T) {
			// Parse
			doc1, err := parser.New().Parse(input)
			assert.NoError(t, err)

			// Serialize
			serialized := ToKDL(doc1)

			// Parse again
			doc2, err := parser.New().Parse(serialized)
			assert.NoError(t, err)

			// Compare
			assert.Equal(t, doc1, doc2)
		})
	}
}
