package parser

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

// TestParseServerConfig tests parsing a server configuration with nested TLS settings
func TestParseServerConfig(t *testing.T) {
	input := `server {
  host "0.0.0.0"
  port 8080

  tls enabled=#true {
    cert-file #"C:\certs\server.crt"#
    key-file #"C:\certs\server.key"#
  }
}`

	doc, err := New().Parse(input)
	want := &Document{
		Nodes: []Node{
			{
				Name:       "server",
				Arguments:  []Value{},
				Properties: []Property{},
				Children: []Node{
					{
						Name: "host",
						Arguments: []Value{
							{Type: ValueTypeString, Value: "0.0.0.0"},
						},
						Properties: []Property{},
						Children:   []Node{},
					},
					{
						Name: "port",
						Arguments: []Value{
							{Type: ValueTypeNumber, Value: "8080"},
						},
						Properties: []Property{},
						Children:   []Node{},
					},
					{
						Name:      "tls",
						Arguments: []Value{},
						Properties: []Property{
							{Key: "enabled", Value: Value{Type: ValueTypeBoolean, Value: "true"}},
						},
						Children: []Node{
							{
								Name: "cert-file",
								Arguments: []Value{
									{Type: ValueTypeString, Value: `C:\certs\server.crt`},
								},
								Properties: []Property{},
								Children:   []Node{},
							},
							{
								Name: "key-file",
								Arguments: []Value{
									{Type: ValueTypeString, Value: `C:\certs\server.key`},
								},
								Properties: []Property{},
								Children:   []Node{},
							},
						},
					},
				},
			},
		},
	}

	assert.NoError(t, err)
	assert.Equal(t, want, doc)
}

// TestParseDatabaseConfig tests database config with type annotations
func TestParseDatabaseConfig(t *testing.T) {
	input := `database {
  host "db.example.com"
  port 5432
  pool {
    min-size (int)5
    max-size (int)20
    timeout (duration)"30s"
  }
}`

	doc, err := New().Parse(input)
	want := &Document{
		Nodes: []Node{
			{
				Name:       "database",
				Arguments:  []Value{},
				Properties: []Property{},
				Children: []Node{
					{
						Name: "host",
						Arguments: []Value{
							{Type: ValueTypeString, Value: "db.example.com"},
						},
						Properties: []Property{},
						Children:   []Node{},
					},
					{
						Name: "port",
						Arguments: []Value{
							{Type: ValueTypeNumber, Value: "5432"},
						},
						Properties: []Property{},
						Children:   []Node{},
					},
					{
						Name:       "pool",
						Arguments:  []Value{},
						Properties: []Property{},
						Children: []Node{
							{
								Name: "min-size",
								Arguments: []Value{
									{Type: ValueTypeNumber, Value: "5", TypeAnnotation: "int"},
								},
								Properties: []Property{},
								Children:   []Node{},
							},
							{
								Name: "max-size",
								Arguments: []Value{
									{Type: ValueTypeNumber, Value: "20", TypeAnnotation: "int"},
								},
								Properties: []Property{},
								Children:   []Node{},
							},
							{
								Name: "timeout",
								Arguments: []Value{
									{Type: ValueTypeString, Value: "30s", TypeAnnotation: "duration"},
								},
								Properties: []Property{},
								Children:   []Node{},
							},
						},
					},
				},
			},
		},
	}

	assert.NoError(t, err)
	assert.Equal(t, want, doc)
}

// TestParsePackageWithDependencies tests package manifest with bare identifiers
func TestParsePackageWithDependencies(t *testing.T) {
	input := `package {
  name awesome-app
  version "1.2.3"
  license MIT

  dependencies {
    dep http-client version="^2.0.0" optional=#false
    dep logger version="~1.5.0"
  }
}`

	doc, err := New().Parse(input)
	want := &Document{
		Nodes: []Node{
			{
				Name:       "package",
				Arguments:  []Value{},
				Properties: []Property{},
				Children: []Node{
					{
						Name: "name",
						Arguments: []Value{
							{Type: ValueTypeString, Value: "awesome-app"},
						},
						Properties: []Property{},
						Children:   []Node{},
					},
					{
						Name: "version",
						Arguments: []Value{
							{Type: ValueTypeString, Value: "1.2.3"},
						},
						Properties: []Property{},
						Children:   []Node{},
					},
					{
						Name: "license",
						Arguments: []Value{
							{Type: ValueTypeString, Value: "MIT"},
						},
						Properties: []Property{},
						Children:   []Node{},
					},
					{
						Name:       "dependencies",
						Arguments:  []Value{},
						Properties: []Property{},
						Children: []Node{
							{
								Name: "dep",
								Arguments: []Value{
									{Type: ValueTypeString, Value: "http-client"},
								},
								Properties: []Property{
									{Key: "version", Value: Value{Type: ValueTypeString, Value: "^2.0.0"}},
									{Key: "optional", Value: Value{Type: ValueTypeBoolean, Value: "false"}},
								},
								Children: []Node{},
							},
							{
								Name: "dep",
								Arguments: []Value{
									{Type: ValueTypeString, Value: "logger"},
								},
								Properties: []Property{
									{Key: "version", Value: Value{Type: ValueTypeString, Value: "~1.5.0"}},
								},
								Children: []Node{},
							},
						},
					},
				},
			},
		},
	}

	assert.NoError(t, err)
	assert.Equal(t, want, doc)
}

// TestParseCIPipelineStage tests CI config with comments and line continuations
func TestParseCIPipelineStage(t *testing.T) {
	input := `// CI Pipeline configuration
pipeline {
  stage build {
    image "golang:1.21"
    script {
      run "go build ./..."
      run "go test ./..."
    }
  }

  /-stage old-deploy {
    script {
      run "old-deploy.sh"
    }
  }

  stage deploy depends-on=build \
    enabled=#true
}`

	doc, err := New().Parse(input)
	want := &Document{
		Nodes: []Node{
			{
				Name:       "pipeline",
				Arguments:  []Value{},
				Properties: []Property{},
				Children: []Node{
					{
						Name: "stage",
						Arguments: []Value{
							{Type: ValueTypeString, Value: "build"},
						},
						Properties: []Property{},
						Children: []Node{
							{
								Name: "image",
								Arguments: []Value{
									{Type: ValueTypeString, Value: "golang:1.21"},
								},
								Properties: []Property{},
								Children:   []Node{},
							},
							{
								Name:       "script",
								Arguments:  []Value{},
								Properties: []Property{},
								Children: []Node{
									{
										Name: "run",
										Arguments: []Value{
											{Type: ValueTypeString, Value: "go build ./..."},
										},
										Properties: []Property{},
										Children:   []Node{},
									},
									{
										Name: "run",
										Arguments: []Value{
											{Type: ValueTypeString, Value: "go test ./..."},
										},
										Properties: []Property{},
										Children:   []Node{},
									},
								},
							},
						},
					},
					{
						Name: "stage",
						Arguments: []Value{
							{Type: ValueTypeString, Value: "deploy"},
						},
						Properties: []Property{
							{Key: "depends-on", Value: Value{Type: ValueTypeString, Value: "build"}},
							{Key: "enabled", Value: Value{Type: ValueTypeBoolean, Value: "true"}},
						},
						Children: []Node{},
					},
				},
			},
		},
	}

	assert.NoError(t, err)
	assert.Equal(t, want, doc)
}

// TestParseLoggingConfig tests logging with hex numbers and multiple properties
func TestParseLoggingConfig(t *testing.T) {
	input := `logging {
  level info
  format json

  rotation {
    max-size (bytes)0x6400000
    compress #true
  }

  /* Configure log levels */
  loggers {
    logger http level=debug
    logger database level=warn
  }
}`

	doc, err := New().Parse(input)
	want := &Document{
		Nodes: []Node{
			{
				Name:       "logging",
				Arguments:  []Value{},
				Properties: []Property{},
				Children: []Node{
					{
						Name: "level",
						Arguments: []Value{
							{Type: ValueTypeString, Value: "info"},
						},
						Properties: []Property{},
						Children:   []Node{},
					},
					{
						Name: "format",
						Arguments: []Value{
							{Type: ValueTypeString, Value: "json"},
						},
						Properties: []Property{},
						Children:   []Node{},
					},
					{
						Name:       "rotation",
						Arguments:  []Value{},
						Properties: []Property{},
						Children: []Node{
							{
								Name: "max-size",
								Arguments: []Value{
									{Type: ValueTypeNumber, Value: "0x6400000", TypeAnnotation: "bytes"},
								},
								Properties: []Property{},
								Children:   []Node{},
							},
							{
								Name: "compress",
								Arguments: []Value{
									{Type: ValueTypeBoolean, Value: "true"},
								},
								Properties: []Property{},
								Children:   []Node{},
							},
						},
					},
					{
						Name:       "loggers",
						Arguments:  []Value{},
						Properties: []Property{},
						Children: []Node{
							{
								Name: "logger",
								Arguments: []Value{
									{Type: ValueTypeString, Value: "http"},
								},
								Properties: []Property{
									{Key: "level", Value: Value{Type: ValueTypeString, Value: "debug"}},
								},
								Children: []Node{},
							},
							{
								Name: "logger",
								Arguments: []Value{
									{Type: ValueTypeString, Value: "database"},
								},
								Properties: []Property{
									{Key: "level", Value: Value{Type: ValueTypeString, Value: "warn"}},
								},
								Children: []Node{},
							},
						},
					},
				},
			},
		},
	}

	assert.NoError(t, err)
	assert.Equal(t, want, doc)
}

// TestParseMetricsWithMultipleArguments tests multiple arguments
func TestParseMetricsWithMultipleArguments(t *testing.T) {
	input := `metrics {
  enabled #true
  provider prometheus
  port 9090

  histogram http_duration 0.005 0.01 0.025
}`

	doc, err := New().Parse(input)
	want := &Document{
		Nodes: []Node{
			{
				Name:       "metrics",
				Arguments:  []Value{},
				Properties: []Property{},
				Children: []Node{
					{
						Name: "enabled",
						Arguments: []Value{
							{Type: ValueTypeBoolean, Value: "true"},
						},
						Properties: []Property{},
						Children:   []Node{},
					},
					{
						Name: "provider",
						Arguments: []Value{
							{Type: ValueTypeString, Value: "prometheus"},
						},
						Properties: []Property{},
						Children:   []Node{},
					},
					{
						Name: "port",
						Arguments: []Value{
							{Type: ValueTypeNumber, Value: "9090"},
						},
						Properties: []Property{},
						Children:   []Node{},
					},
					{
						Name: "histogram",
						Arguments: []Value{
							{Type: ValueTypeString, Value: "http_duration"},
							{Type: ValueTypeNumber, Value: "0.005"},
							{Type: ValueTypeNumber, Value: "0.01"},
							{Type: ValueTypeNumber, Value: "0.025"},
						},
						Properties: []Property{},
						Children:   []Node{},
					},
				},
			},
		},
	}

	assert.NoError(t, err)
	assert.Equal(t, want, doc)
}
