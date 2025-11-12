package kdly

import (
	"encoding/json"
	"encoding/xml"
	"fmt"
	"os"
	"path/filepath"
	"sort"
	"strings"
	"testing"
)

// TestSpec represents the structure of a test specification file
type TestSpec struct {
	Input  string            `json:"input"`
	Output map[string]string `json:"output"` // format -> filename (e.g., "json": "test001-out.json")
	Title  string            `json:"title"`
	Errors []TestError       `json:"errors"`
}

// TestError represents an expected error
type TestError struct {
	Type string `json:"type"`
	Line int    `json:"line"`
	Pos  int    `json:"pos"`
}

// MinimalNode represents a node in the minimalist JSON format
type MinimalNode struct {
	Name           string                 `json:"name"`
	Args           []MinimalValue         `json:"args,omitempty"`
	Props          map[string]MinimalValue `json:"props,omitempty"`
	Children       []MinimalNode          `json:"children,omitempty"`
	TypeAnnotation *string                `json:"typeAnnotation,omitempty"`
}

// MinimalValue represents a value in the minimalist JSON format
type MinimalValue struct {
	Type           string      `json:"type"`
	Value          interface{} `json:"value"`
	TypeAnnotation *string     `json:"typeAnnotation,omitempty"`
}

// MinimalDocument represents a document in the minimalist JSON format
type MinimalDocument struct {
	Nodes []MinimalNode `json:"nodes"`
}

// toMinimalValue converts a Value to MinimalValue
func toMinimalValue(v *Value) MinimalValue {
	typeStr := ""
	switch v.Type {
	case ValueTypeString:
		typeStr = "string"
	case ValueTypeNumber:
		typeStr = "number"
	case ValueTypeBoolean:
		typeStr = "boolean"
	case ValueTypeNull:
		typeStr = "null"
	}

	return MinimalValue{
		Type:           typeStr,
		Value:          v.Value,
		TypeAnnotation: v.TypeAnnotation,
	}
}

// toMinimalNode converts a Node to MinimalNode
func toMinimalNode(n *Node) MinimalNode {
	mn := MinimalNode{
		Name:           n.Name,
		TypeAnnotation: n.TypeAnnotation,
	}

	// Convert arguments
	if len(n.Arguments) > 0 {
		mn.Args = make([]MinimalValue, len(n.Arguments))
		for i, arg := range n.Arguments {
			mn.Args[i] = toMinimalValue(arg)
		}
	}

	// Convert properties
	if len(n.Properties) > 0 {
		mn.Props = make(map[string]MinimalValue)
		for key, val := range n.Properties {
			mn.Props[key] = toMinimalValue(val)
		}
	}

	// Convert children
	if len(n.Children) > 0 {
		mn.Children = make([]MinimalNode, len(n.Children))
		for i, child := range n.Children {
			mn.Children[i] = toMinimalNode(child)
		}
	}

	return mn
}

// toMinimalDocument converts a Document to MinimalDocument
func toMinimalDocument(doc *Document) MinimalDocument {
	md := MinimalDocument{
		Nodes: make([]MinimalNode, len(doc.Nodes)),
	}
	for i, node := range doc.Nodes {
		md.Nodes[i] = toMinimalNode(node)
	}
	return md
}

// XMLValue represents a value in XML format
type XMLValue struct {
	Type           string `xml:"type,attr"`
	Value          string `xml:",chardata"`
	TypeAnnotation string `xml:"typeAnnotation,attr,omitempty"`
}

// XMLProp represents a property in XML format
type XMLProp struct {
	Name           string `xml:"name,attr"`
	Type           string `xml:"type,attr"`
	Value          string `xml:",chardata"`
	TypeAnnotation string `xml:"typeAnnotation,attr,omitempty"`
}

// XMLNode represents a node in XML format
type XMLNode struct {
	XMLName        xml.Name   `xml:"node"`
	Name           string     `xml:"name,attr"`
	TypeAnnotation string     `xml:"typeAnnotation,attr,omitempty"`
	Args           []XMLValue `xml:"arg,omitempty"`
	Props          []XMLProp  `xml:"prop,omitempty"`
	Children       []XMLNode  `xml:"node,omitempty"`
}

// XMLDocument represents a document in XML format
type XMLDocument struct {
	XMLName xml.Name  `xml:"document"`
	Nodes   []XMLNode `xml:"node"`
}

// toXMLValue converts a Value to XMLValue
func toXMLValue(v *Value) XMLValue {
	typeStr := ""
	valueStr := ""

	switch v.Type {
	case ValueTypeString:
		typeStr = "string"
		if v.Value != nil {
			valueStr = v.Value.(string)
		}
	case ValueTypeNumber:
		typeStr = "number"
		if v.Value != nil {
			valueStr = fmt.Sprintf("%v", v.Value)
		}
	case ValueTypeBoolean:
		typeStr = "boolean"
		if v.Value != nil {
			valueStr = fmt.Sprintf("%v", v.Value)
		}
	case ValueTypeNull:
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

// toXMLNode converts a Node to XMLNode
func toXMLNode(n *Node) XMLNode {
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
			case ValueTypeString:
				prop.Type = "string"
				if val.Value != nil {
					prop.Value = val.Value.(string)
				}
			case ValueTypeNumber:
				prop.Type = "number"
				if val.Value != nil {
					prop.Value = fmt.Sprintf("%v", val.Value)
				}
			case ValueTypeBoolean:
				prop.Type = "boolean"
				if val.Value != nil {
					prop.Value = fmt.Sprintf("%v", val.Value)
				}
			case ValueTypeNull:
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

// toXMLDocument converts a Document to XMLDocument
func toXMLDocument(doc *Document) XMLDocument {
	xd := XMLDocument{
		Nodes: make([]XMLNode, len(doc.Nodes)),
	}
	for i, node := range doc.Nodes {
		xd.Nodes[i] = toXMLNode(node)
	}
	return xd
}

func TestGoldenTests(t *testing.T) {
	// Get all test spec files
	testFiles, err := filepath.Glob("tests/test*.json")
	if err != nil {
		t.Fatalf("Failed to list test files: %v", err)
	}

	// Filter only test spec files (not input or output files)
	var specFiles []string
	for _, file := range testFiles {
		base := filepath.Base(file)
		if !strings.Contains(base, "-in") && !strings.Contains(base, "-out") {
			specFiles = append(specFiles, file)
		}
	}

	// Sort files to ensure consistent test order
	sort.Strings(specFiles)

	if len(specFiles) == 0 {
		t.Fatal("No test spec files found in tests/ directory")
	}

	for _, specFile := range specFiles {
		// Read test spec
		specData, err := os.ReadFile(specFile)
		if err != nil {
			t.Fatalf("Failed to read test spec %s: %v", specFile, err)
		}

		var spec TestSpec
		if err := json.Unmarshal(specData, &spec); err != nil {
			t.Fatalf("Failed to parse test spec %s: %v", specFile, err)
		}

		// Run test as a subtest
		t.Run(spec.Title, func(t *testing.T) {
			runGoldenTest(t, specFile, spec)
		})
	}
}

func runGoldenTest(t *testing.T, specFile string, spec TestSpec) {
	testDir := filepath.Dir(specFile)

	// Read input file
	inputPath := filepath.Join(testDir, spec.Input)
	inputData, err := os.ReadFile(inputPath)
	if err != nil {
		t.Fatalf("Failed to read input file %s: %v", inputPath, err)
	}

	// Parse input
	parser := NewParser(string(inputData))
	doc, parseErr := parser.Parse()

	// Check for expected errors
	if len(spec.Errors) > 0 {
		// Test expects errors
		if parseErr == nil {
			t.Fatalf("Expected parse errors but got none")
		}

		// Convert parseErr to error list if needed
		var errorList []error
		switch e := parseErr.(type) {
		case ErrorList:
			errorList = []error(e)
		default:
			errorList = []error{parseErr}
		}

		// Check if we have the right number of errors
		if len(errorList) != len(spec.Errors) {
			t.Errorf("Expected %d errors, got %d", len(spec.Errors), len(errorList))
		}

		// Check each expected error
		for i, expectedErr := range spec.Errors {
			if i >= len(errorList) {
				break
			}

			actualErr := errorList[i]

			// Check error position
			var actualLine, actualPos int
			switch e := actualErr.(type) {
			case *ParseError:
				actualLine = e.Position.Line
				actualPos = e.Position.Column
			case *LexError:
				actualLine = e.Position.Line
				actualPos = e.Position.Column
			}

			if actualLine != expectedErr.Line {
				t.Errorf("Error %d: expected line %d, got %d", i, expectedErr.Line, actualLine)
			}
			if actualPos != expectedErr.Pos {
				t.Errorf("Error %d: expected pos %d, got %d", i, expectedErr.Pos, actualPos)
			}
		}

		return
	}

	// Test expects successful parse
	if parseErr != nil {
		t.Fatalf("Parse error: %v", parseErr)
	}

	if doc == nil {
		t.Fatal("Document is nil")
	}

	// Test each output format
	for format, outputFile := range spec.Output {
		t.Run(format, func(t *testing.T) {
			outputPath := filepath.Join(testDir, outputFile)

			switch format {
			case "json":
				compareJSONOutput(t, doc, outputPath)
			case "xml":
				compareXMLOutput(t, doc, outputPath)
			default:
				t.Errorf("Unknown output format: %s", format)
			}
		})
	}
}

func compareJSONOutput(t *testing.T, doc *Document, expectedPath string) {
	// Convert to minimal format
	minimalDoc := toMinimalDocument(doc)

	// Serialize to JSON
	actualJSON, err := json.MarshalIndent(minimalDoc, "", "  ")
	if err != nil {
		t.Fatalf("Failed to marshal actual output: %v", err)
	}

	// Read expected output
	expectedJSON, err := os.ReadFile(expectedPath)
	if err != nil {
		t.Fatalf("Failed to read expected output %s: %v", expectedPath, err)
	}

	// Normalize JSON for comparison (parse and re-serialize both)
	var expectedObj, actualObj interface{}
	if err := json.Unmarshal(expectedJSON, &expectedObj); err != nil {
		t.Fatalf("Failed to parse expected JSON: %v", err)
	}
	if err := json.Unmarshal(actualJSON, &actualObj); err != nil {
		t.Fatalf("Failed to parse actual JSON: %v", err)
	}

	expectedNormalized, _ := json.MarshalIndent(expectedObj, "", "  ")
	actualNormalized, _ := json.MarshalIndent(actualObj, "", "  ")

	// Compare
	if string(expectedNormalized) != string(actualNormalized) {
		t.Errorf("JSON output mismatch:\n\nExpected:\n%s\n\nActual:\n%s",
			string(expectedNormalized), string(actualNormalized))

		// Show diff
		showDiff(t, string(expectedNormalized), string(actualNormalized))
	}
}

func compareXMLOutput(t *testing.T, doc *Document, expectedPath string) {
	// Convert to XML format
	xmlDoc := toXMLDocument(doc)

	// Serialize to XML
	actualXML, err := xml.MarshalIndent(xmlDoc, "", "  ")
	if err != nil {
		t.Fatalf("Failed to marshal actual XML output: %v", err)
	}

	// Add XML declaration
	actualXML = []byte(xml.Header + string(actualXML))

	// Read expected output
	expectedXML, err := os.ReadFile(expectedPath)
	if err != nil {
		t.Fatalf("Failed to read expected XML output %s: %v", expectedPath, err)
	}

	// Normalize XML for comparison (parse and re-serialize both)
	var expectedDoc, actualDoc XMLDocument
	if err := xml.Unmarshal(expectedXML, &expectedDoc); err != nil {
		t.Fatalf("Failed to parse expected XML: %v", err)
	}
	if err := xml.Unmarshal(actualXML, &actualDoc); err != nil {
		t.Fatalf("Failed to parse actual XML: %v", err)
	}

	expectedNormalized, _ := xml.MarshalIndent(expectedDoc, "", "  ")
	actualNormalized, _ := xml.MarshalIndent(actualDoc, "", "  ")

	expectedNormalized = []byte(xml.Header + string(expectedNormalized))
	actualNormalized = []byte(xml.Header + string(actualNormalized))

	// Compare
	if string(expectedNormalized) != string(actualNormalized) {
		t.Errorf("XML output mismatch:\n\nExpected:\n%s\n\nActual:\n%s",
			string(expectedNormalized), string(actualNormalized))

		// Show diff
		showDiff(t, string(expectedNormalized), string(actualNormalized))
	}
}

func showDiff(t *testing.T, expected, actual string) {
	expectedLines := strings.Split(expected, "\n")
	actualLines := strings.Split(actual, "\n")

	maxLines := len(expectedLines)
	if len(actualLines) > maxLines {
		maxLines = len(actualLines)
	}

	t.Log("\nLine-by-line comparison:")
	for i := 0; i < maxLines; i++ {
		var expLine, actLine string
		if i < len(expectedLines) {
			expLine = expectedLines[i]
		}
		if i < len(actualLines) {
			actLine = actualLines[i]
		}

		if expLine != actLine {
			t.Logf("Line %d:", i+1)
			t.Logf("  Expected: %q", expLine)
			t.Logf("  Actual:   %q", actLine)
		}
	}
}

// Helper function to update golden files (for development)
func updateGoldenFile(t *testing.T, outputPath string, doc *Document) {
	minimalDoc := toMinimalDocument(doc)
	data, err := json.MarshalIndent(minimalDoc, "", "  ")
	if err != nil {
		t.Fatalf("Failed to marshal output: %v", err)
	}

	if err := os.WriteFile(outputPath, data, 0644); err != nil {
		t.Fatalf("Failed to write golden file: %v", err)
	}

	fmt.Printf("Updated golden file: %s\n", outputPath)
}
