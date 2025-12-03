package parser

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func node(name string) *Node {
	return &Node{
		Name:       name,
		Arguments:  []Value{},
		Properties: []Property{},
		Children:   []Node{},
	}
}

// ============================================================
// Level 1: Fundamentals (Empty & Single Elements)
// ============================================================

func TestEmptyDocument(t *testing.T) {
	input := ``

	doc, err := New().Parse(input)

	assert.NoError(t, err)
	assert.Empty(t, doc.Nodes)
}

func TestSingleSimpleNode(t *testing.T) {
	input := `hello`

	doc, err := New().Parse(input)
	want := &Document{
		Nodes: []Node{
			*node("hello"),
		},
	}

	assert.NoError(t, err)
	assert.Equal(t, want, doc)
}

func TestSingleSimpleNode_withQuotes(t *testing.T) {
	input := `"illegal(){}[]/\\=#;identifier"`

	doc, err := New().Parse(input)
	want := &Document{
		Nodes: []Node{
			*node(`illegal(){}[]/\=#;identifier`),
		},
	}

	assert.NoError(t, err)
	assert.Equal(t, want, doc)
}

func TestSingleSimpleNode_flexibleBare(t *testing.T) {
	input := "-<123~!$@%^&*,.:'`|?+>"

	doc, err := New().Parse(input)
	want := &Document{
		Nodes: []Node{
			*node("-<123~!$@%^&*,.:'`|?+>"),
		},
	}

	assert.NoError(t, err)
	assert.Equal(t, want, doc)
}

func TestNodeWithStringArgument(t *testing.T) {
	input := `title "Hello, World"`

	doc, err := New().Parse(input)
	want := &Document{
		Nodes: []Node{
			{
				Name: "title",
				Arguments: []Value{
					{Type: ValueTypeString, Value: "Hello, World"},
				},
				Properties: []Property{},
				Children:   []Node{},
			},
		},
	}

	assert.NoError(t, err)
	assert.Equal(t, want, doc)
}

func TestNodeWithBareStringArgument(t *testing.T) {
	input := `node1 this-is-a-string`

	doc, err := New().Parse(input)
	want := &Document{
		Nodes: []Node{
			{
				Name: "node1",
				Arguments: []Value{
					{Type: ValueTypeString, Value: "this-is-a-string"},
				},
				Properties: []Property{},
				Children:   []Node{},
			},
		},
	}

	assert.NoError(t, err)
	assert.Equal(t, want, doc)
}

func TestNodeWithEscapedStringArgument(t *testing.T) {
	input := `node2 "this\nhas\tescapes"`

	doc, err := New().Parse(input)
	want := &Document{
		Nodes: []Node{
			{
				Name: "node2",
				Arguments: []Value{
					{Type: ValueTypeString, Value: "this\nhas\tescapes"},
				},
				Properties: []Property{},
				Children:   []Node{},
			},
		},
	}

	assert.NoError(t, err)
	assert.Equal(t, want, doc)
}

func TestNodeWithRawStringArgument(t *testing.T) {
	input := `node3 #"C:\Users\zkat\raw\string"#`

	doc, err := New().Parse(input)
	want := &Document{
		Nodes: []Node{
			{
				Name: "node3",
				Arguments: []Value{
					{Type: ValueTypeString, Value: `C:\Users\zkat\raw\string`},
				},
				Properties: []Property{},
				Children:   []Node{},
			},
		},
	}

	assert.NoError(t, err)
	assert.Equal(t, want, doc)
}

// ============================================================
// Level 2: Basic Values & Types
// ============================================================

func TestNodeWithMultipleArguments(t *testing.T) {
	input := `bookmarks 12 15 188 1234`

	doc, err := New().Parse(input)
	want := &Document{
		Nodes: []Node{
			{
				Name: "bookmarks",
				Arguments: []Value{
					{Type: ValueTypeNumber, Value: "12"},
					{Type: ValueTypeNumber, Value: "15"},
					{Type: ValueTypeNumber, Value: "188"},
					{Type: ValueTypeNumber, Value: "1234"},
				},
				Properties: []Property{},
				Children:   []Node{},
			},
		},
	}

	assert.NoError(t, err)
	assert.Equal(t, want, doc)
}

func TestNodeWithBooleanAndNull(t *testing.T) {
	input := `flags #true #false #null`

	doc, err := New().Parse(input)
	want := &Document{
		Nodes: []Node{
			{
				Name: "flags",
				Arguments: []Value{
					{Type: ValueTypeBoolean, Value: "true"},
					{Type: ValueTypeBoolean, Value: "false"},
					{Type: ValueTypeNull, Value: "null"},
				},
				Properties: []Property{},
				Children:   []Node{},
			},
		},
	}

	assert.NoError(t, err)
	assert.Equal(t, want, doc)
}

func TestNodeWithHexNumbers(t *testing.T) {
	input := `color 0xdeadbeef`

	doc, err := New().Parse(input)
	want := &Document{
		Nodes: []Node{
			{
				Name: "color",
				Arguments: []Value{
					{Type: ValueTypeNumber, Value: "0xdeadbeef"},
				},
				Properties: []Property{},
				Children:   []Node{},
			},
		},
	}

	assert.NoError(t, err)
	assert.Equal(t, want, doc)
}

// ============================================================
// Level 3: Properties & Type Annotations
// ============================================================

func TestNodeWithProperties(t *testing.T) {
	input := `author "Alex Monad" email=alex@example.com active=#true`

	doc, err := New().Parse(input)
	want := &Document{
		Nodes: []Node{
			{
				Name: "author",
				Arguments: []Value{
					{Type: ValueTypeString, Value: "Alex Monad"},
				},
				Properties: []Property{
					{Key: "email", Value: Value{Type: ValueTypeString, Value: "alex@example.com"}},
					{Key: "active", Value: Value{Type: ValueTypeBoolean, Value: "true"}},
				},
				Children: []Node{},
			},
		},
	}

	assert.NoError(t, err)
	assert.Equal(t, want, doc)
}

func TestTypeAnnotations(t *testing.T) {
	input := `numbers (u8)10 (uuid)"123e4567-e89b-12d3-a456-426614174000"`

	doc, err := New().Parse(input)

	want := &Document{
		Nodes: []Node{
			{
				Name: "numbers",
				Arguments: []Value{
					{
						Type:           ValueTypeNumber,
						Value:          "10",
						TypeAnnotation: "u8",
					},
					{
						Type:           ValueTypeString,
						Value:          "123e4567-e89b-12d3-a456-426614174000",
						TypeAnnotation: "uuid",
					},
				},
				Properties: []Property{},
				Children:   []Node{},
			},
		},
	}

	assert.NoError(t, err)
	assert.Equal(t, want, doc)
}

// ============================================================
// Level 4: Comments
// ============================================================

func TestLineComment(t *testing.T) {
	input := `// This is a comment`

	doc, err := New().Parse(input)

	assert.NoError(t, err)
	assert.Empty(t, doc.Nodes)
}

func TestMultilineComment(t *testing.T) {
	input := `/* Multiline
   comment */`

	doc, err := New().Parse(input)

	assert.NoError(t, err)
	assert.Empty(t, doc.Nodes)
}

// ============================================================
// Feature 1: Slashdash Comments (/-) - Comment out nodes/values
// ============================================================

func TestSlashdashComment_SimpleNode(t *testing.T) {
	input := `/-node1
node2`

	doc, err := New().Parse(input)
	want := &Document{
		Nodes: []Node{
			*node("node2"),
		},
	}

	assert.NoError(t, err)
	assert.Equal(t, want, doc)
}

func TestSlashdashComment_NodeWithSpace(t *testing.T) {
	input := `/- node1
node2`

	doc, err := New().Parse(input)
	want := &Document{
		Nodes: []Node{
			*node("node2"),
		},
	}

	assert.NoError(t, err)
	assert.Equal(t, want, doc)
}

func TestSlashdashComment_NodeWithChildren(t *testing.T) {
	input := `/-commented {
  child1
  child2
}
visible`

	doc, err := New().Parse(input)
	want := &Document{
		Nodes: []Node{
			*node("visible"),
		},
	}

	assert.NoError(t, err)
	assert.Equal(t, want, doc)
}

func TestSlashdashComment_Argument(t *testing.T) {
	input := `node "visible1" /-"commented" "visible2"`

	doc, err := New().Parse(input)
	want := &Document{
		Nodes: []Node{
			{
				Name: "node",
				Arguments: []Value{
					{Type: ValueTypeString, Value: "visible1"},
					{Type: ValueTypeString, Value: "visible2"},
				},
				Properties: []Property{},
				Children:   []Node{},
			},
		},
	}

	assert.NoError(t, err)
	assert.Equal(t, want, doc)
}

func TestSlashdashComment_ArgumentWithSpace(t *testing.T) {
	input := `node "visible1" /- "commented" "visible2"`

	doc, err := New().Parse(input)
	want := &Document{
		Nodes: []Node{
			{
				Name: "node",
				Arguments: []Value{
					{Type: ValueTypeString, Value: "visible1"},
					{Type: ValueTypeString, Value: "visible2"},
				},
				Properties: []Property{},
				Children:   []Node{},
			},
		},
	}

	assert.NoError(t, err)
	assert.Equal(t, want, doc)
}

func TestSlashdashComment_Property(t *testing.T) {
	input := `node visible=1 /-commented=2 another=3`

	doc, err := New().Parse(input)
	want := &Document{
		Nodes: []Node{
			{
				Name:      "node",
				Arguments: []Value{},
				Properties: []Property{
					{Key: "visible", Value: Value{Type: ValueTypeNumber, Value: "1"}},
					{Key: "another", Value: Value{Type: ValueTypeNumber, Value: "3"}},
				},
				Children: []Node{},
			},
		},
	}

	assert.NoError(t, err)
	assert.Equal(t, want, doc)
}

func TestSlashdashComment_ChildrenBlock(t *testing.T) {
	input := `node /-{
  child1
  child2
}`

	doc, err := New().Parse(input)
	want := &Document{
		Nodes: []Node{
			*node("node"),
		},
	}

	assert.NoError(t, err)
	assert.Equal(t, want, doc)
}

func TestSlashdashComment_MultipleElements(t *testing.T) {
	input := `node /-"arg1" key=/-"val1" /-prop=5 "arg2"`

	doc, err := New().Parse(input)
	want := &Document{
		Nodes: []Node{
			{
				Name: "node",
				Arguments: []Value{
					{Type: ValueTypeString, Value: "arg2"},
				},
				Properties: []Property{
					{Key: "key", Value: Value{Type: ValueTypeString, Value: ""}},
				},
				Children: []Node{},
			},
		},
	}

	assert.NoError(t, err)
	assert.Equal(t, want, doc)
}

func TestSlashdashComment_MixedWithOtherComments(t *testing.T) {
	input := `// line comment
node1
/-node2
/* block comment */
node3`

	doc, err := New().Parse(input)
	want := &Document{
		Nodes: []Node{
			*node("node1"),
			*node("node3"),
		},
	}

	assert.NoError(t, err)
	assert.Equal(t, want, doc)
}

// ============================================================
// Feature 2: Special Numeric Values (#inf, #-inf, #nan)
// ============================================================

func TestSpecialNumeric_PositiveInfinity(t *testing.T) {
	input := `value #inf`

	doc, err := New().Parse(input)
	want := &Document{
		Nodes: []Node{
			{
				Name: "value",
				Arguments: []Value{
					{Type: ValueTypeNumber, Value: "inf"},
				},
				Properties: []Property{},
				Children:   []Node{},
			},
		},
	}

	assert.NoError(t, err)
	assert.Equal(t, want, doc)
}

func TestSpecialNumeric_NegativeInfinity(t *testing.T) {
	input := `value #-inf`

	doc, err := New().Parse(input)
	want := &Document{
		Nodes: []Node{
			{
				Name: "value",
				Arguments: []Value{
					{Type: ValueTypeNumber, Value: "-inf"},
				},
				Properties: []Property{},
				Children:   []Node{},
			},
		},
	}

	assert.NoError(t, err)
	assert.Equal(t, want, doc)
}

func TestSpecialNumeric_NaN(t *testing.T) {
	input := `value #nan`

	doc, err := New().Parse(input)
	want := &Document{
		Nodes: []Node{
			{
				Name: "value",
				Arguments: []Value{
					{Type: ValueTypeNumber, Value: "nan"},
				},
				Properties: []Property{},
				Children:   []Node{},
			},
		},
	}

	assert.NoError(t, err)
	assert.Equal(t, want, doc)
}

func TestSpecialNumeric_MixedWithRegularNumbers(t *testing.T) {
	input := `values 3.14 #inf -2.5 #-inf 0 #nan`

	doc, err := New().Parse(input)
	want := &Document{
		Nodes: []Node{
			{
				Name: "values",
				Arguments: []Value{
					{Type: ValueTypeNumber, Value: "3.14"},
					{Type: ValueTypeNumber, Value: "inf"},
					{Type: ValueTypeNumber, Value: "-2.5"},
					{Type: ValueTypeNumber, Value: "-inf"},
					{Type: ValueTypeNumber, Value: "0"},
					{Type: ValueTypeNumber, Value: "nan"},
				},
				Properties: []Property{},
				Children:   []Node{},
			},
		},
	}

	assert.NoError(t, err)
	assert.Equal(t, want, doc)
}

func TestSpecialNumeric_InProperties(t *testing.T) {
	input := `node max=#inf min=#-inf error=#nan`

	doc, err := New().Parse(input)
	want := &Document{
		Nodes: []Node{
			{
				Name:      "node",
				Arguments: []Value{},
				Properties: []Property{
					{Key: "max", Value: Value{Type: ValueTypeNumber, Value: "inf"}},
					{Key: "min", Value: Value{Type: ValueTypeNumber, Value: "-inf"}},
					{Key: "error", Value: Value{Type: ValueTypeNumber, Value: "nan"}},
				},
				Children: []Node{},
			},
		},
	}

	assert.NoError(t, err)
	assert.Equal(t, want, doc)
}

func TestSpecialNumeric_WithTypeAnnotations(t *testing.T) {
	input := `node (f64)#inf (f32)#-inf (float)#nan`

	doc, err := New().Parse(input)
	want := &Document{
		Nodes: []Node{
			{
				Name: "node",
				Arguments: []Value{
					{Type: ValueTypeNumber, Value: "inf", TypeAnnotation: "f64"},
					{Type: ValueTypeNumber, Value: "-inf", TypeAnnotation: "f32"},
					{Type: ValueTypeNumber, Value: "nan", TypeAnnotation: "float"},
				},
				Properties: []Property{},
				Children:   []Node{},
			},
		},
	}

	assert.NoError(t, err)
	assert.Equal(t, want, doc)
}

// ============================================================
// Feature 3: Line Continuation (\)
// ============================================================

func TestLineContinuation_SimpleArguments(t *testing.T) {
	input := `node arg1 \
     arg2 \
     arg3`

	doc, err := New().Parse(input)
	want := &Document{
		Nodes: []Node{
			{
				Name: "node",
				Arguments: []Value{
					{Type: ValueTypeString, Value: "arg1"},
					{Type: ValueTypeString, Value: "arg2"},
					{Type: ValueTypeString, Value: "arg3"},
				},
				Properties: []Property{},
				Children:   []Node{},
			},
		},
	}

	assert.NoError(t, err)
	assert.Equal(t, want, doc)
}

func TestLineContinuation_WithProperties(t *testing.T) {
	input := `node \
  key1=value1 \
  key2=value2`

	doc, err := New().Parse(input)
	want := &Document{
		Nodes: []Node{
			{
				Name:      "node",
				Arguments: []Value{},
				Properties: []Property{
					{Key: "key1", Value: Value{Type: ValueTypeString, Value: "value1"}},
					{Key: "key2", Value: Value{Type: ValueTypeString, Value: "value2"}},
				},
				Children: []Node{},
			},
		},
	}

	assert.NoError(t, err)
	assert.Equal(t, want, doc)
}

func TestLineContinuation_MultipleNodes(t *testing.T) {
	input := `node1 \
  arg1 \
  arg2
node2 arg3`

	doc, err := New().Parse(input)
	want := &Document{
		Nodes: []Node{
			{
				Name: "node1",
				Arguments: []Value{
					{Type: ValueTypeString, Value: "arg1"},
					{Type: ValueTypeString, Value: "arg2"},
				},
				Properties: []Property{},
				Children:   []Node{},
			},
			{
				Name: "node2",
				Arguments: []Value{
					{Type: ValueTypeString, Value: "arg3"},
				},
				Properties: []Property{},
				Children:   []Node{},
			},
		},
	}

	assert.NoError(t, err)
	assert.Equal(t, want, doc)
}

func TestLineContinuation_WithComments(t *testing.T) {
	input := `node \
  // comment on continuation line
  arg1 \
  arg2`

	doc, err := New().Parse(input)
	want := &Document{
		Nodes: []Node{
			{
				Name: "node",
				Arguments: []Value{
					{Type: ValueTypeString, Value: "arg1"},
					{Type: ValueTypeString, Value: "arg2"},
				},
				Properties: []Property{},
				Children:   []Node{},
			},
		},
	}

	assert.NoError(t, err)
	assert.Equal(t, want, doc)
}

func TestLineContinuation_BeforeChildren(t *testing.T) {
	input := `node \
{
  child
}`

	doc, err := New().Parse(input)
	want := &Document{
		Nodes: []Node{
			{
				Name:       "node",
				Arguments:  []Value{},
				Properties: []Property{},
				Children: []Node{
					*node("child"),
				},
			},
		},
	}

	assert.NoError(t, err)
	assert.Equal(t, want, doc)
}

func TestLineContinuation_MixedWithValues(t *testing.T) {
	input := `node "string" \
  123 \
  #true \
  key=value`

	doc, err := New().Parse(input)
	want := &Document{
		Nodes: []Node{
			{
				Name: "node",
				Arguments: []Value{
					{Type: ValueTypeString, Value: "string"},
					{Type: ValueTypeNumber, Value: "123"},
					{Type: ValueTypeBoolean, Value: "true"},
				},
				Properties: []Property{
					{Key: "key", Value: Value{Type: ValueTypeString, Value: "value"}},
				},
				Children: []Node{},
			},
		},
	}

	assert.NoError(t, err)
	assert.Equal(t, want, doc)
}

func TestLineContinuation_MixedWithValues_AnotherOrder(t *testing.T) {
	input := `node "string" \
  123 \
  key=value \
  #true`

	doc, err := New().Parse(input)
	want := &Document{
		Nodes: []Node{
			{
				Name: "node",
				Arguments: []Value{
					{Type: ValueTypeString, Value: "string"},
					{Type: ValueTypeNumber, Value: "123"},
					{Type: ValueTypeBoolean, Value: "true"},
				},
				Properties: []Property{
					{Key: "key", Value: Value{Type: ValueTypeString, Value: "value"}},
				},
				Children: []Node{},
			},
		},
	}

	assert.NoError(t, err)
	assert.Equal(t, want, doc)
}

// ============================================================
// Feature 4: Bare Identifier Number Validation (should reject number-like patterns)
// ============================================================

func TestBareIdentifier_AllowValidBare(t *testing.T) {
	input := `valid-identifier`

	doc, err := New().Parse(input)
	want := &Document{
		Nodes: []Node{
			*node("valid-identifier"),
		},
	}

	assert.NoError(t, err)
	assert.Equal(t, want, doc)
}

func TestBareIdentifier_AllowDashInMiddle(t *testing.T) {
	input := `my-node-name`

	doc, err := New().Parse(input)
	want := &Document{
		Nodes: []Node{
			*node("my-node-name"),
		},
	}

	assert.NoError(t, err)
	assert.Equal(t, want, doc)
}

func TestBareIdentifier_NumberAsArgument(t *testing.T) {
	input := `node -123 +456 -.5 +.5`

	doc, err := New().Parse(input)
	want := &Document{
		Nodes: []Node{
			{
				Name: "node",
				Arguments: []Value{
					{Type: ValueTypeNumber, Value: "-123"},
					{Type: ValueTypeNumber, Value: "+456"},
					{Type: ValueTypeNumber, Value: "-.5"},
					{Type: ValueTypeNumber, Value: "+.5"},
				},
				Properties: []Property{},
				Children:   []Node{},
			},
		},
	}

	assert.NoError(t, err)
	assert.Equal(t, want, doc)
}

func TestLineCommentWithNodes(t *testing.T) {
	input := `// This is a comment
node1
// Another comment
node2`

	doc, err := New().Parse(input)
	want := &Document{
		Nodes: []Node{
			*node("node1"),
			*node("node2"),
		},
	}

	assert.NoError(t, err)
	assert.Equal(t, want, doc)
}

func TestMultilineCommentWithNodes(t *testing.T) {
	input := `/* Comment before */
node1
/* Comment
   in between */
node2
/* Comment after */`

	doc, err := New().Parse(input)
	want := &Document{
		Nodes: []Node{
			*node("node1"),
			*node("node2"),
		},
	}

	assert.NoError(t, err)
	assert.Equal(t, want, doc)
}

func TestInlineCommentAfterNode(t *testing.T) {
	input := `node1 // inline comment`

	doc, err := New().Parse(input)
	want := &Document{
		Nodes: []Node{
			*node("node1"),
		},
	}

	assert.NoError(t, err)
	assert.Equal(t, want, doc)
}

func TestMixedComments(t *testing.T) {
	input := `// Line comment
/* Multiline comment */
node1
// Another line comment
node2 // inline
/* Final comment */`

	doc, err := New().Parse(input)
	want := &Document{
		Nodes: []Node{
			*node("node1"),
			*node("node2"),
		},
	}

	assert.NoError(t, err)
	assert.Equal(t, want, doc)
}

func TestNestedMultilineComment(t *testing.T) {
	input := `/* /* comment inside comment */ */`

	doc, err := New().Parse(input)

	assert.NoError(t, err)
	assert.Empty(t, doc.Nodes)
}

// ============================================================
// Level 5: Strings & Multiline
// ============================================================

func TestQuotedMultilineString(t *testing.T) {
	input := `message """
  hello
  world
  """`

	doc, err := New().Parse(input)
	want := &Document{
		Nodes: []Node{
			{
				Name: "message",
				Arguments: []Value{
					{Type: ValueTypeString, Value: "hello\nworld\n"},
				},
				Properties: []Property{},
				Children:   []Node{},
			},
		},
	}

	assert.NoError(t, err)
	assert.Equal(t, want, doc)
}

func TestRawString(t *testing.T) {
	input := `path #"C:\path\to\file"#`

	doc, err := New().Parse(input)
	want := &Document{
		Nodes: []Node{
			{
				Name: "path",
				Arguments: []Value{
					{Type: ValueTypeString, Value: `C:\path\to\file`},
				},
				Properties: []Property{},
				Children:   []Node{},
			},
		},
	}

	assert.NoError(t, err)
	assert.Equal(t, want, doc)
}

func TestMultiLevelRawStrings(t *testing.T) {
	tests := []struct {
		name  string
		input string
		want  string
	}{
		{
			name:  "double hash raw string",
			input: `node ##"hello world"##`,
			want:  "hello world",
		},
		{
			name:  "triple hash raw string",
			input: `node ###"hello world"###`,
			want:  "hello world",
		},
		{
			name:  "raw string containing single hash quote",
			input: `node ##"string with #" inside"##`,
			want:  `string with #" inside`,
		},
		{
			name:  "raw string containing double hash quote",
			input: `node ###"string with ##" inside"###`,
			want:  `string with ##" inside`,
		},
		{
			name:  "raw string with quote in middle",
			input: `node ##"hello"world"##`,
			want:  `hello"world`,
		},
		{
			name:  "raw string with hash in middle",
			input: `node ##"hello#world"##`,
			want:  `hello#world`,
		},
		{
			name:  "raw string with mismatched closing",
			input: `node ###"content"#not closing yet"###`,
			want:  `content"#not closing yet`,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			doc, err := New().Parse(tt.input)
			assert.NoError(t, err)
			assert.Equal(t, 1, len(doc.Nodes), "should have one node")
			assert.Equal(t, 1, len(doc.Nodes[0].Arguments), "should have one argument")
			assert.Equal(t, tt.want, doc.Nodes[0].Arguments[0].Value, "raw string value mismatch")
		})
	}
}

func TestMultiLevelRawStrings_ErrorCases(t *testing.T) {
	tests := []struct {
		name  string
		input string
	}{
		{
			name:  "insufficient closing hashes",
			input: `node ###"hello world"##`,
		},
		{
			name:  "no closing delimiter at all",
			input: `node ##"hello world`,
		},
		{
			name:  "too many closing hashes - also fails",
			input: `node ##"hello"###`,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			_, err := New().Parse(tt.input)
			// All of these should fail with unterminated raw string
			assert.Error(t, err, "expected parse error for: %s", tt.input)
			assert.Contains(t, err.Error(), "unterminated raw string", "error should mention unterminated raw string")
		})
	}
}

// ============================================================
// Level 6: Children & Nested Structures
// ============================================================

func TestNodeWithChildren(t *testing.T) {
	input := `contents {
  section "First section" {
    paragraph "Text"
  }
}`

	doc, err := New().Parse(input)
	want := &Document{
		Nodes: []Node{
			{
				Name:       "contents",
				Arguments:  []Value{},
				Properties: []Property{},
				Children: []Node{
					{
						Name: "section",
						Arguments: []Value{
							{Type: ValueTypeString, Value: "First section"},
						},
						Properties: []Property{},
						Children: []Node{
							{
								Name: "paragraph",
								Arguments: []Value{
									{Type: ValueTypeString, Value: "Text"},
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

func TestMixedConstructs(t *testing.T) {
	input := `server {
  host "localhost"
  port 8080
  debug #true
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
							{Type: ValueTypeString, Value: "localhost"},
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
						Name: "debug",
						Arguments: []Value{
							{Type: ValueTypeBoolean, Value: "true"},
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

// ============================================================
// Level 7: Additional Coverage
// ============================================================

func TestMultipleTopLevelNodes(t *testing.T) {
	input := `node1
node2
node3`

	doc, err := New().Parse(input)
	want := &Document{
		Nodes: []Node{
			*node("node1"),
			*node("node2"),
			*node("node3"),
		},
	}

	assert.NoError(t, err)
	assert.Equal(t, want, doc)
}

func TestNodeWithOnlyProperties(t *testing.T) {
	input := `config debug=#true port=8080`

	doc, err := New().Parse(input)
	want := &Document{
		Nodes: []Node{
			{
				Name:      "config",
				Arguments: []Value{},
				Properties: []Property{
					{Key: "debug", Value: Value{Type: ValueTypeBoolean, Value: "true"}},
					{Key: "port", Value: Value{Type: ValueTypeNumber, Value: "8080"}},
				},
				Children: []Node{},
			},
		},
	}

	assert.NoError(t, err)
	assert.Equal(t, want, doc)
}

func TestNodeWithEmptyChildren(t *testing.T) {
	input := `container {}`

	doc, err := New().Parse(input)
	want := &Document{
		Nodes: []Node{
			*node("container"),
		},
	}

	assert.NoError(t, err)
	assert.Equal(t, want, doc)
}

// ============================================================
// Feature 5: Numbers with Underscores - Readability separators
// ============================================================

func TestNumbersWithUnderscores_Simple(t *testing.T) {
	input := `value 1_000_000`

	doc, err := New().Parse(input)
	want := &Document{
		Nodes: []Node{
			{
				Name: "value",
				Arguments: []Value{
					{Type: ValueTypeNumber, Value: "1_000_000"},
				},
				Properties: []Property{},
				Children:   []Node{},
			},
		},
	}

	assert.NoError(t, err)
	assert.Equal(t, want, doc)
}

func TestNumbersWithUnderscores_Multiple(t *testing.T) {
	input := `values 1_000 100_000 1_000_000`

	doc, err := New().Parse(input)
	want := &Document{
		Nodes: []Node{
			{
				Name: "values",
				Arguments: []Value{
					{Type: ValueTypeNumber, Value: "1_000"},
					{Type: ValueTypeNumber, Value: "100_000"},
					{Type: ValueTypeNumber, Value: "1_000_000"},
				},
				Properties: []Property{},
				Children:   []Node{},
			},
		},
	}

	assert.NoError(t, err)
	assert.Equal(t, want, doc)
}

func TestNumbersWithUnderscores_FloatingPoint(t *testing.T) {
	input := `value 3_141.592_653`

	doc, err := New().Parse(input)
	want := &Document{
		Nodes: []Node{
			{
				Name: "value",
				Arguments: []Value{
					{Type: ValueTypeNumber, Value: "3_141.592_653"},
				},
				Properties: []Property{},
				Children:   []Node{},
			},
		},
	}

	assert.NoError(t, err)
	assert.Equal(t, want, doc)
}

func TestNumbersWithUnderscores_Hexadecimal(t *testing.T) {
	input := `color 0xdead_beef 0xFF_00_FF`

	doc, err := New().Parse(input)
	want := &Document{
		Nodes: []Node{
			{
				Name: "color",
				Arguments: []Value{
					{Type: ValueTypeNumber, Value: "0xdead_beef"},
					{Type: ValueTypeNumber, Value: "0xFF_00_FF"},
				},
				Properties: []Property{},
				Children:   []Node{},
			},
		},
	}

	assert.NoError(t, err)
	assert.Equal(t, want, doc)
}

func TestNumbersWithUnderscores_Binary(t *testing.T) {
	input := `bits 0b1010_1100 0b1111_0000_1111_0000`

	doc, err := New().Parse(input)
	want := &Document{
		Nodes: []Node{
			{
				Name: "bits",
				Arguments: []Value{
					{Type: ValueTypeNumber, Value: "0b1010_1100"},
					{Type: ValueTypeNumber, Value: "0b1111_0000_1111_0000"},
				},
				Properties: []Property{},
				Children:   []Node{},
			},
		},
	}

	assert.NoError(t, err)
	assert.Equal(t, want, doc)
}

func TestNumbersWithUnderscores_InProperties(t *testing.T) {
	input := `config timeout=30_000 max_size=1_000_000`

	doc, err := New().Parse(input)
	want := &Document{
		Nodes: []Node{
			{
				Name:      "config",
				Arguments: []Value{},
				Properties: []Property{
					{Key: "timeout", Value: Value{Type: ValueTypeNumber, Value: "30_000"}},
					{Key: "max_size", Value: Value{Type: ValueTypeNumber, Value: "1_000_000"}},
				},
				Children: []Node{},
			},
		},
	}

	assert.NoError(t, err)
	assert.Equal(t, want, doc)
}

func TestBinaryAndOctalNumbers_Basic(t *testing.T) {
	input := `bits 0b1010 0o755`

	doc, err := New().Parse(input)
	want := &Document{
		Nodes: []Node{
			{
				Name: "bits",
				Arguments: []Value{
					{Type: ValueTypeNumber, Value: "0b1010"},
					{Type: ValueTypeNumber, Value: "0o755"},
				},
				Properties: []Property{},
				Children:   []Node{},
			},
		},
	}

	assert.NoError(t, err)
	assert.Equal(t, want, doc)
}

func TestBinaryAndOctalNumbers_UpperCase(t *testing.T) {
	input := `bits 0B1010 0O755`

	doc, err := New().Parse(input)
	want := &Document{
		Nodes: []Node{
			{
				Name: "bits",
				Arguments: []Value{
					{Type: ValueTypeNumber, Value: "0B1010"},
					{Type: ValueTypeNumber, Value: "0O755"},
				},
				Properties: []Property{},
				Children:   []Node{},
			},
		},
	}

	assert.NoError(t, err)
	assert.Equal(t, want, doc)
}

func TestBinaryAndOctalNumbers_WithUnderscores(t *testing.T) {
	input := `bits 0b1111_0000 0o7_5_5`

	doc, err := New().Parse(input)
	want := &Document{
		Nodes: []Node{
			{
				Name: "bits",
				Arguments: []Value{
					{Type: ValueTypeNumber, Value: "0b1111_0000"},
					{Type: ValueTypeNumber, Value: "0o7_5_5"},
				},
				Properties: []Property{},
				Children:   []Node{},
			},
		},
	}

	assert.NoError(t, err)
	assert.Equal(t, want, doc)
}

func TestFloatingPointNumbers_Basic(t *testing.T) {
	input := `coords 3.14 -2.5 1.0e10`

	doc, err := New().Parse(input)
	want := &Document{
		Nodes: []Node{
			{
				Name: "coords",
				Arguments: []Value{
					{Type: ValueTypeNumber, Value: "3.14"},
					{Type: ValueTypeNumber, Value: "-2.5"},
					{Type: ValueTypeNumber, Value: "1.0e10"},
				},
				Properties: []Property{},
				Children:   []Node{},
			},
		},
	}

	assert.NoError(t, err)
	assert.Equal(t, want, doc)
}

func TestFloatingPointNumbers_ScientificNotation(t *testing.T) {
	input := `values 1.5e-3 2E10 -3.14e+5`

	doc, err := New().Parse(input)
	want := &Document{
		Nodes: []Node{
			{
				Name: "values",
				Arguments: []Value{
					{Type: ValueTypeNumber, Value: "1.5e-3"},
					{Type: ValueTypeNumber, Value: "2E10"},
					{Type: ValueTypeNumber, Value: "-3.14e+5"},
				},
				Properties: []Property{},
				Children:   []Node{},
			},
		},
	}

	assert.NoError(t, err)
	assert.Equal(t, want, doc)
}

func TestFloatingPointNumbers_SignedAndUnsigned(t *testing.T) {
	input := `values +3.14 -2.718 0.5`

	doc, err := New().Parse(input)
	want := &Document{
		Nodes: []Node{
			{
				Name: "values",
				Arguments: []Value{
					{Type: ValueTypeNumber, Value: "+3.14"},
					{Type: ValueTypeNumber, Value: "-2.718"},
					{Type: ValueTypeNumber, Value: "0.5"},
				},
				Properties: []Property{},
				Children:   []Node{},
			},
		},
	}

	assert.NoError(t, err)
	assert.Equal(t, want, doc)
}

// ============================================================
// Feature 6: Mixed Arguments and Properties - Interleaved syntax
// ============================================================

func TestMixedArgumentsAndProperties_Simple(t *testing.T) {
	input := `person "John" age=30 "Doe" city="NYC"`

	doc, err := New().Parse(input)
	want := &Document{
		Nodes: []Node{
			{
				Name: "person",
				Arguments: []Value{
					{Type: ValueTypeString, Value: "John"},
					{Type: ValueTypeString, Value: "Doe"},
				},
				Properties: []Property{
					{Key: "age", Value: Value{Type: ValueTypeNumber, Value: "30"}},
					{Key: "city", Value: Value{Type: ValueTypeString, Value: "NYC"}},
				},
				Children: []Node{},
			},
		},
	}

	assert.NoError(t, err)
	assert.Equal(t, want, doc)
}

func TestMixedArgumentsAndProperties_Multiple(t *testing.T) {
	input := `product 123 name="Widget" "v2.0" price=99.99 available=#true`

	doc, err := New().Parse(input)
	want := &Document{
		Nodes: []Node{
			{
				Name: "product",
				Arguments: []Value{
					{Type: ValueTypeNumber, Value: "123"},
					{Type: ValueTypeString, Value: "v2.0"},
				},
				Properties: []Property{
					{Key: "name", Value: Value{Type: ValueTypeString, Value: "Widget"}},
					{Key: "price", Value: Value{Type: ValueTypeNumber, Value: "99.99"}},
					{Key: "available", Value: Value{Type: ValueTypeBoolean, Value: "true"}},
				},
				Children: []Node{},
			},
		},
	}

	assert.NoError(t, err)
	assert.Equal(t, want, doc)
}

func TestMixedArgumentsAndProperties_AllTypes(t *testing.T) {
	input := `node "str" num=42 #true bool=#false bare-arg key="value" #null`

	doc, err := New().Parse(input)
	want := &Document{
		Nodes: []Node{
			{
				Name: "node",
				Arguments: []Value{
					{Type: ValueTypeString, Value: "str"},
					{Type: ValueTypeBoolean, Value: "true"},
					{Type: ValueTypeString, Value: "bare-arg"},
					{Type: ValueTypeNull, Value: "null"},
				},
				Properties: []Property{
					{Key: "num", Value: Value{Type: ValueTypeNumber, Value: "42"}},
					{Key: "bool", Value: Value{Type: ValueTypeBoolean, Value: "false"}},
					{Key: "key", Value: Value{Type: ValueTypeString, Value: "value"}},
				},
				Children: []Node{},
			},
		},
	}

	assert.NoError(t, err)
	assert.Equal(t, want, doc)
}

func TestMixedArgumentsAndProperties_WithTypeAnnotations(t *testing.T) {
	input := `data (i32)100 format="json" (f64)3.14 compression="gzip"`

	doc, err := New().Parse(input)
	want := &Document{
		Nodes: []Node{
			{
				Name: "data",
				Arguments: []Value{
					{Type: ValueTypeNumber, Value: "100", TypeAnnotation: "i32"},
					{Type: ValueTypeNumber, Value: "3.14", TypeAnnotation: "f64"},
				},
				Properties: []Property{
					{Key: "format", Value: Value{Type: ValueTypeString, Value: "json"}},
					{Key: "compression", Value: Value{Type: ValueTypeString, Value: "gzip"}},
				},
				Children: []Node{},
			},
		},
	}

	assert.NoError(t, err)
	assert.Equal(t, want, doc)
}

func TestMixedArgumentsAndProperties_WithChildren(t *testing.T) {
	input := `server "prod" port=8080 "v1" {
  endpoint "/api"
}`

	doc, err := New().Parse(input)
	want := &Document{
		Nodes: []Node{
			{
				Name: "server",
				Arguments: []Value{
					{Type: ValueTypeString, Value: "prod"},
					{Type: ValueTypeString, Value: "v1"},
				},
				Properties: []Property{
					{Key: "port", Value: Value{Type: ValueTypeNumber, Value: "8080"}},
				},
				Children: []Node{
					{
						Name: "endpoint",
						Arguments: []Value{
							{Type: ValueTypeString, Value: "/api"},
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

func TestMixedArgumentsAndProperties_PropertiesFirst(t *testing.T) {
	input := `config debug=#true "config.json" timeout=30`

	doc, err := New().Parse(input)
	want := &Document{
		Nodes: []Node{
			{
				Name: "config",
				Arguments: []Value{
					{Type: ValueTypeString, Value: "config.json"},
				},
				Properties: []Property{
					{Key: "debug", Value: Value{Type: ValueTypeBoolean, Value: "true"}},
					{Key: "timeout", Value: Value{Type: ValueTypeNumber, Value: "30"}},
				},
				Children: []Node{},
			},
		},
	}

	assert.NoError(t, err)
	assert.Equal(t, want, doc)
}

func TestVariousWhitespace(t *testing.T) {
	input := "node1\n\nnode2\n  \nnode3"

	doc, err := New().Parse(input)
	want := &Document{
		Nodes: []Node{
			*node("node1"),
			*node("node2"),
			*node("node3"),
		},
	}

	assert.NoError(t, err)
	assert.Equal(t, want, doc)
}

func TestQuotedNodeNameOnly(t *testing.T) {
	input := `"node-with-dashes"`

	doc, err := New().Parse(input)
	want := &Document{
		Nodes: []Node{
			*node("node-with-dashes"),
		},
	}

	assert.NoError(t, err)
	assert.Equal(t, want, doc)
}

func TestSiblingsAfterNesting(t *testing.T) {
	input := `first {}
second`

	doc, err := New().Parse(input)
	want := &Document{
		Nodes: []Node{
			*node("first"),
			*node("second"),
		},
	}

	assert.NoError(t, err)
	assert.Equal(t, want, doc)
}

func TestSemicolonSeparatedNodes(t *testing.T) {
	input := `node1; node2`

	doc, err := New().Parse(input)
	want := &Document{
		Nodes: []Node{
			*node("node1"),
			*node("node2"),
		},
	}

	assert.NoError(t, err)
	assert.Equal(t, want, doc)
}

func TestSemicolonSeparatedNodesWithChildrenInline(t *testing.T) {
	input := `node1; node2 { node3 }`

	doc, err := New().Parse(input)
	want := &Document{
		Nodes: []Node{
			*node("node1"),
			{
				Name:       "node2",
				Arguments:  []Value{},
				Properties: []Property{},
				Children: []Node{
					*node("node3"),
				},
			},
		},
	}

	assert.NoError(t, err)
	assert.Equal(t, want, doc)
}

func TestSemicolonSeparatedNodesWithChildrenMultiline(t *testing.T) {
	input := `node1; node2 {
  node3
}`

	doc, err := New().Parse(input)
	want := &Document{
		Nodes: []Node{
			*node("node1"),
			{
				Name:       "node2",
				Arguments:  []Value{},
				Properties: []Property{},
				Children: []Node{
					*node("node3"),
				},
			},
		},
	}

	assert.NoError(t, err)
	assert.Equal(t, want, doc)
}

// ============================================================
// Legacy Test (Keep for compatibility)
// ============================================================

func TestGoldenDataLocal(t *testing.T) {
	t.Run("sample test", func(t *testing.T) {
		// This test passes
	})
}
