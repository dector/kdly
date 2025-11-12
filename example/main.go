package main

import (
	"fmt"
	"github.com/dector/kdly"
)

func main() {
	// Example KDL document
	input := `
// Configuration file
database {
	host "localhost"
	port 5432
	credentials {
		username "admin"
		password "secret"
	}
}

// Server configuration
server {
	(production)instance "prod-1" region="us-east" enabled=#true
	(development)instance "dev-1" region="us-west" enabled=#false
}

// Feature flags
features analytics=#true caching=#false max-connections=100
`

	fmt.Println("Parsing KDL document...")
	fmt.Println("======================================================")

	parser := kdly.NewParser(input)
	doc, err := parser.Parse()
	if err != nil {
		fmt.Printf("Error parsing: %v\n", err)
		return
	}

	fmt.Printf("Successfully parsed %d top-level nodes\n\n", len(doc.Nodes))

	// Display parsed structure
	for i, node := range doc.Nodes {
		fmt.Printf("%d. Node: %s\n", i+1, node.Name)

		if len(node.Arguments) > 0 {
			fmt.Println("   Arguments:")
			for j, arg := range node.Arguments {
				typeAnnotation := ""
				if arg.TypeAnnotation != nil {
					typeAnnotation = fmt.Sprintf(" (%s)", *arg.TypeAnnotation)
				}
				fmt.Printf("     %d. %v [%s]%s\n", j+1, arg.Value, arg.Type, typeAnnotation)
			}
		}

		if len(node.Properties) > 0 {
			fmt.Println("   Properties:")
			for key, value := range node.Properties {
				typeAnnotation := ""
				if value.TypeAnnotation != nil {
					typeAnnotation = fmt.Sprintf(" (%s)", *value.TypeAnnotation)
				}
				fmt.Printf("     %s = %v [%s]%s\n", key, value.Value, value.Type, typeAnnotation)
			}
		}

		if len(node.Children) > 0 {
			fmt.Printf("   Children: %d node(s)\n", len(node.Children))
			for _, child := range node.Children {
				fmt.Printf("     - %s", child.Name)
				if len(child.Arguments) > 0 {
					fmt.Printf(" (with %d arg(s))", len(child.Arguments))
				}
				if len(child.Properties) > 0 {
					fmt.Printf(" (with %d property/ies)", len(child.Properties))
				}
				if child.TypeAnnotation != nil {
					fmt.Printf(" [type: %s]", *child.TypeAnnotation)
				}
				fmt.Println()
			}
		}

		fmt.Println()
	}
}
