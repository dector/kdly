package kdly

import (
	"encoding/json"
	"encoding/xml"
	"fmt"
	"os"
	"path/filepath"
	"sort"
	"strconv"
	"strings"
	"testing"
)

// TestSpec represents the structure of a test specification file
type TestSpec struct {
	Input   string            `json:"input"`
	Output  map[string]string `json:"output"` // format -> filename (e.g., "json": "test001-out.json")
	Title   string            `json:"title"`
	Errors  []TestError       `json:"errors"`
	Source  *string           `json:"source,omitempty"`  // Optional source URL for test origin
	Enabled *bool             `json:"enabled,omitempty"` // Defaults to true if not specified
}

// TestError represents an expected error
type TestError struct {
	Type string `json:"type"`
	Line int    `json:"line"`
	Pos  int    `json:"pos"`
}

// MinimalNode represents a node in the minimalist JSON format
type MinimalNode struct {
	Name           string                  `json:"name"`
	Args           []MinimalValue          `json:"args,omitempty"`
	Props          map[string]MinimalValue `json:"props,omitempty"`
	Children       []MinimalNode           `json:"children,omitempty"`
	TypeAnnotation *string                 `json:"typeAnnotation,omitempty"`
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

// toKDL converts a Document to canonical KDL format
func toKDL(doc *Document) string {
	if len(doc.Nodes) == 0 {
		return "\n"
	}
	var sb strings.Builder
	for i, node := range doc.Nodes {
		if i > 0 {
			sb.WriteString("\n")
		}
		formatNode(&sb, node, 0)
	}
	// formatNode already adds a newline after each node, so no need to add another
	return sb.String()
}

// formatNode formats a single node to KDL
func formatNode(sb *strings.Builder, node *Node, indent int) {
	// Write indentation
	for i := 0; i < indent; i++ {
		sb.WriteString("    ")
	}

	// Write type annotation if present
	if node.TypeAnnotation != nil {
		sb.WriteString("(")
		sb.WriteString(*node.TypeAnnotation)
		sb.WriteString(")")
	}

	// Write node name
	sb.WriteString(formatIdentifier(node.Name))

	// Write arguments
	for _, arg := range node.Arguments {
		sb.WriteString(" ")
		formatValue(sb, arg)
	}

	// Write properties (sorted by key)
	if len(node.Properties) > 0 {
		keys := make([]string, 0, len(node.Properties))
		for k := range node.Properties {
			keys = append(keys, k)
		}
		sort.Strings(keys)

		for _, key := range keys {
			sb.WriteString(" ")
			sb.WriteString(formatIdentifier(key))
			sb.WriteString("=")
			formatValue(sb, node.Properties[key])
		}
	}

	// Write children
	if len(node.Children) > 0 {
		sb.WriteString(" {\n")
		for _, child := range node.Children {
			formatNode(sb, child, indent+1)
		}
		for i := 0; i < indent; i++ {
			sb.WriteString("    ")
		}
		sb.WriteString("}")
	}

	sb.WriteString("\n")
}

// formatValue formats a value to KDL
func formatValue(sb *strings.Builder, val *Value) {
	// Write type annotation if present
	if val.TypeAnnotation != nil {
		sb.WriteString("(")
		sb.WriteString(*val.TypeAnnotation)
		sb.WriteString(")")
	}

	switch val.Type {
	case ValueTypeString:
		if val.Value != nil {
			s := val.Value.(string)
			// Only quote if necessary
			if needsQuoting(s) {
				sb.WriteString(formatString(s))
			} else {
				sb.WriteString(s)
			}
		} else {
			sb.WriteString(`""`)
		}
	case ValueTypeNumber:
		if val.Value != nil {
			sb.WriteString(formatNumber(val.Value))
		} else {
			sb.WriteString("0")
		}
	case ValueTypeBoolean:
		if val.Value != nil && val.Value.(bool) {
			sb.WriteString("#true")
		} else {
			sb.WriteString("#false")
		}
	case ValueTypeNull:
		sb.WriteString("#null")
	}
}

// formatString formats a string value with proper escaping
func formatString(s string) string {
	var sb strings.Builder
	sb.WriteString(`"`)
	for _, r := range s {
		switch r {
		case '\b':
			sb.WriteString(`\b`)
		case '\f':
			sb.WriteString(`\f`)
		case '\n':
			sb.WriteString(`\n`)
		case '\r':
			sb.WriteString(`\r`)
		case '\t':
			sb.WriteString(`\t`)
		case '\\':
			sb.WriteString(`\\`)
		case '"':
			sb.WriteString(`\"`)
		default:
			sb.WriteRune(r)
		}
	}
	sb.WriteString(`"`)
	return sb.String()
}

// formatNumber formats a number value to its simplest decimal representation
func formatNumber(v interface{}) string {
	switch n := v.(type) {
	case int64:
		return strconv.FormatInt(n, 10)
	case float64:
		// Check if it's an integer value
		if n == float64(int64(n)) {
			return strconv.FormatFloat(n, 'f', 0, 64)
		}
		// Use scientific notation for very large or small numbers
		if n >= 1e15 || n <= -1e15 || (n != 0 && (n < 1e-4 && n > -1e-4)) {
			return strconv.FormatFloat(n, 'E', -1, 64)
		}
		return strconv.FormatFloat(n, 'f', -1, 64)
	default:
		return fmt.Sprintf("%v", v)
	}
}

// formatIdentifier formats an identifier, quoting if necessary
func formatIdentifier(s string) string {
	if needsQuoting(s) {
		return formatString(s)
	}
	return s
}

// needsQuoting checks if an identifier needs to be quoted
func needsQuoting(s string) bool {
	if len(s) == 0 {
		return true
	}

	// Check for keywords
	switch s {
	case "true", "false", "null":
		return true
	}

	// Check if it looks like a number
	if len(s) > 0 {
		first := s[0]
		if first >= '0' && first <= '9' {
			return true
		}
		if (first == '+' || first == '-') && len(s) > 1 {
			second := s[1]
			if second >= '0' && second <= '9' {
				return true
			}
		}
	}

	// Check for special characters that require quoting
	for i, r := range s {
		// First character rules
		if i == 0 {
			if !isInitialIdentChar(r) {
				return true
			}
		} else {
			if !isIdentChar(r) {
				return true
			}
		}
	}

	return false
}

// isInitialIdentChar checks if a rune can start an identifier
func isInitialIdentChar(r rune) bool {
	return (r >= 'a' && r <= 'z') ||
		(r >= 'A' && r <= 'Z') ||
		r == '_' ||
		r == '-' ||
		r == '.' ||
		r == '?' ||
		r == '!' ||
		r >= 0x80 // Unicode
}

// isIdentChar checks if a rune can be in an identifier
func isIdentChar(r rune) bool {
	return isInitialIdentChar(r) ||
		(r >= '0' && r <= '9')
}

func TestGoldenTests(t *testing.T) {
	// Get all test spec files
	testFiles, err := filepath.Glob("tests/*.json")
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

		// Run test unless explicitly disabled (enabled defaults to true)
		if spec.Enabled != nil && !*spec.Enabled {
			// Create a subtest that skips itself so it shows in test output
			t.Run(spec.Title, func(t *testing.T) {
				t.Skip("Test explicitly disabled")
			})
		} else {
			t.Run(spec.Title, func(t *testing.T) {
				runGoldenTest(t, specFile, spec)
			})
		}
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

	// Check if test expects errors or success
	// Tests with no output and no errors specified expect parse to fail (kdl-org convention)
	expectsFailure := len(spec.Output) == 0 && len(spec.Errors) == 0

	// Check for expected errors
	if len(spec.Errors) > 0 {
		// Test expects specific errors
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
	} else if expectsFailure {
		// Test expects failure but no specific errors (kdl-org convention)
		if parseErr == nil {
			t.Fatalf("Expected parse to fail but it succeeded")
		}
		// Parse failed as expected, test passes
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
			case "kdl":
				compareKDLOutput(t, doc, outputPath)
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

func compareKDLOutput(t *testing.T, doc *Document, expectedPath string) {
	// Convert to canonical KDL format
	actualKDL := toKDL(doc)

	// Read expected output
	expectedKDL, err := os.ReadFile(expectedPath)
	if err != nil {
		t.Fatalf("Failed to read expected KDL output %s: %v", expectedPath, err)
	}

	// Normalize line endings for comparison
	expectedStr := strings.ReplaceAll(string(expectedKDL), "\r\n", "\n")
	actualStr := strings.ReplaceAll(actualKDL, "\r\n", "\n")

	// Compare
	if expectedStr != actualStr {
		t.Errorf("KDL output mismatch:\n\nExpected:\n%s\n\nActual:\n%s",
			expectedStr, actualStr)

		// Show diff
		showDiff(t, expectedStr, actualStr)
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
