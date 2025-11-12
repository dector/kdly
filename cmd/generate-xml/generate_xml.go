package main

import (
	"encoding/xml"
	"fmt"
	"os"
	"path/filepath"

	kdly "github.com/dector/kdly"
)

func main() {
	testDir := "tests"

	// Process each test input file
	for i := 1; i <= 20; i++ {
		testNum := fmt.Sprintf("%03d", i)
		inputFile := filepath.Join(testDir, fmt.Sprintf("test%s-in.kdl", testNum))
		outputFile := filepath.Join(testDir, fmt.Sprintf("test%s-out.xml", testNum))

		// Read input
		inputData, err := os.ReadFile(inputFile)
		if err != nil {
			fmt.Printf("Skipping test%s: %v\n", testNum, err)
			continue
		}

		// Parse
		parser := kdly.NewParser(string(inputData))
		doc, err := parser.Parse()
		if err != nil {
			fmt.Printf("Test%s has parse errors, skipping XML generation\n", testNum)
			continue
		}

		// Convert to XML
		xmlDoc := toXMLDocument(doc)

		// Write XML
		xmlData, err := xml.MarshalIndent(xmlDoc, "", "  ")
		if err != nil {
			fmt.Printf("Failed to marshal XML for test%s: %v\n", testNum, err)
			continue
		}

		// Add XML declaration
		xmlOutput := []byte(xml.Header + string(xmlData) + "\n")

		if err := os.WriteFile(outputFile, xmlOutput, 0644); err != nil {
			fmt.Printf("Failed to write XML for test%s: %v\n", testNum, err)
			continue
		}

		fmt.Printf("Generated %s\n", outputFile)
	}
}

// XML structures (copied from golden_test.go)
type XMLValue struct {
	Type           string `xml:"type,attr"`
	Value          string `xml:",chardata"`
	TypeAnnotation string `xml:"typeAnnotation,attr,omitempty"`
}

type XMLProp struct {
	Name           string `xml:"name,attr"`
	Type           string `xml:"type,attr"`
	Value          string `xml:",chardata"`
	TypeAnnotation string `xml:"typeAnnotation,attr,omitempty"`
}

type XMLNode struct {
	XMLName        xml.Name   `xml:"node"`
	Name           string     `xml:"name,attr"`
	TypeAnnotation string     `xml:"typeAnnotation,attr,omitempty"`
	Args           []XMLValue `xml:"arg,omitempty"`
	Props          []XMLProp  `xml:"prop,omitempty"`
	Children       []XMLNode  `xml:"node,omitempty"`
}

type XMLDocument struct {
	XMLName xml.Name  `xml:"document"`
	Nodes   []XMLNode `xml:"node"`
}

func toXMLValue(v *kdly.Value) XMLValue {
	typeStr := ""
	valueStr := ""

	switch v.Type {
	case kdly.ValueTypeString:
		typeStr = "string"
		if v.Value != nil {
			valueStr = v.Value.(string)
		}
	case kdly.ValueTypeNumber:
		typeStr = "number"
		if v.Value != nil {
			valueStr = fmt.Sprintf("%v", v.Value)
		}
	case kdly.ValueTypeBoolean:
		typeStr = "boolean"
		if v.Value != nil {
			valueStr = fmt.Sprintf("%v", v.Value)
		}
	case kdly.ValueTypeNull:
		typeStr = "null"
		valueStr = ""
	}

	xv := XMLValue{
		Type:  typeStr,
		Value: valueStr,
	}

	if v.TypeAnnotation != nil {
		xv.TypeAnnotation = *v.TypeAnnotation
	}

	return xv
}

func toXMLNode(n *kdly.Node) XMLNode {
	xn := XMLNode{
		Name: n.Name,
	}

	if n.TypeAnnotation != nil {
		xn.TypeAnnotation = *n.TypeAnnotation
	}

	// Convert arguments
	if len(n.Arguments) > 0 {
		xn.Args = make([]XMLValue, len(n.Arguments))
		for i, arg := range n.Arguments {
			xn.Args[i] = toXMLValue(arg)
		}
	}

	// Convert properties
	if len(n.Properties) > 0 {
		xn.Props = make([]XMLProp, 0, len(n.Properties))
		for key, val := range n.Properties {
			prop := XMLProp{
				Name: key,
				Type: "",
			}

			switch val.Type {
			case kdly.ValueTypeString:
				prop.Type = "string"
				if val.Value != nil {
					prop.Value = val.Value.(string)
				}
			case kdly.ValueTypeNumber:
				prop.Type = "number"
				if val.Value != nil {
					prop.Value = fmt.Sprintf("%v", val.Value)
				}
			case kdly.ValueTypeBoolean:
				prop.Type = "boolean"
				if val.Value != nil {
					prop.Value = fmt.Sprintf("%v", val.Value)
				}
			case kdly.ValueTypeNull:
				prop.Type = "null"
				prop.Value = ""
			}

			if val.TypeAnnotation != nil {
				prop.TypeAnnotation = *val.TypeAnnotation
			}

			xn.Props = append(xn.Props, prop)
		}
	}

	// Convert children
	if len(n.Children) > 0 {
		xn.Children = make([]XMLNode, len(n.Children))
		for i, child := range n.Children {
			xn.Children[i] = toXMLNode(child)
		}
	}

	return xn
}

func toXMLDocument(doc *kdly.Document) XMLDocument {
	xd := XMLDocument{
		Nodes: make([]XMLNode, len(doc.Nodes)),
	}
	for i, node := range doc.Nodes {
		xd.Nodes[i] = toXMLNode(node)
	}
	return xd
}
