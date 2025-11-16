# KDLy

A Go parser for [KDL v2](https://kdl.dev) (KDL Document Language).

> [!WARNING] The parser is mostly implmented using guided vibe-coding approach. There might be bugs.

> [!WARNING] Public API might change before 1.0 release.

> [!WARNING] The parser might not yet fully compliant with the KDL v2 specification. This library is not tested in production.

## Installation

```bash
go get github.com/dector/kdly
```

## Usage

### Basic Parsing

```go
package main

import (
    "fmt"
    "log"
    "github.com/dector/kdly"
)

func main() {
    input := `
		  server {
		    host "localhost"
		    port 8080
		    env "production"
		
		    database {
		      driver "postgres"
		      connection-string "postgres://localhost/mydb"
		      pool-size 10
		    }
		  }
	  `

    doc, err := kdly.Parse(input)
    if err != nil {
        log.Fatal(err)
    }

    // Access server configuration
    server := doc.Nodes[0]
    fmt.Println(server.Name) // "server"

    // Access properties
    for _, arg := range server.Children[0].Arguments {
        fmt.Printf("%s: %s\n", server.Children[0].Name, arg.Value)
    }
}
```

### Serialization

```go
// Convert document back to KDL text
output := kdly.ToKDL(doc)
fmt.Println(output)
```

### Parser Options

```go
// Keep duplicate properties instead of using rightmost value
parser := kdly.NewParser().WithAllowDuplicateProperties()

// Reject type annotations
parser := kdly.NewParser().WithNoTypeAnnotations()
```

## Data Types

The parser exposes these key types:

- `Document` - Top-level KDL document containing nodes
- `Node` - A KDL node with name, arguments, properties, and children
- `Value` - Any value (string, number, boolean, null)
- `Property` - Key-value pair

## Features

- Almost full KDL v2 spec support
- Type annotations
- Raw strings
- Multiline strings with dedentation
- Hex, binary, octal numbers
- Scientific notation
- Comments (line, multiline, slashdash)
- Configurable duplicate property handling

## License

See LICENSE file for details.
