package kdly

import (
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"
	"sort"
	"strings"
	"testing"
)

// TestSpec represents the structure of a test specification file
type TestSpec struct {
	Input  string      `json:"input"`
	Output string      `json:"output"`
	Title  string      `json:"title"`
	Errors []TestError `json:"errors"`
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
		if !strings.Contains(base, "-input") && !strings.Contains(base, "-output") {
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

	// Convert to minimal format
	minimalDoc := toMinimalDocument(doc)

	// Serialize to JSON
	actualJSON, err := json.MarshalIndent(minimalDoc, "", "  ")
	if err != nil {
		t.Fatalf("Failed to marshal actual output: %v", err)
	}

	// Read expected output
	outputPath := filepath.Join(testDir, spec.Output)
	expectedJSON, err := os.ReadFile(outputPath)
	if err != nil {
		t.Fatalf("Failed to read expected output %s: %v", outputPath, err)
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
		t.Errorf("Output mismatch:\n\nExpected:\n%s\n\nActual:\n%s",
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
