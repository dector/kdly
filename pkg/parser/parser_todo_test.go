package parser

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func Test_TODO_AllEscapes(t *testing.T) {
	input := `node "\"\\\b\f\n\r\t\s"
`

	doc, err := New().Parse(input)

	_ = doc
	_ = err
	t.Skip("TODO")
}

func Test_AllNodeFields(t *testing.T) {
	input := `node arg prop=val {
    inner_node
}
`

	doc, err := New().Parse(input)
	want := &Document{
		Nodes: []Node{
			{
				Name: "node",
				Arguments: []Value{
					{Type: ValueTypeString, Value: "arg"},
				},
				Properties: []Property{
					{Key: "prop", Value: Value{Type: ValueTypeString, Value: "val"}},
				},
				Children: []Node{
					{
						Name:       "inner_node",
						Arguments:  []Value{},
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

func Test_ArgAndPropSameName(t *testing.T) {
	input := `node arg arg=val
`

	doc, err := New().Parse(input)
	want := &Document{
		Nodes: []Node{
			{
				Name: "node",
				Arguments: []Value{
					{Type: ValueTypeString, Value: "arg"},
				},
				Properties: []Property{
					{Key: "arg", Value: Value{Type: ValueTypeString, Value: "val"}},
				},
				Children: []Node{},
			},
		},
	}

	assert.NoError(t, err)
	assert.Equal(t, want, doc)
}

func Test_ArgBare(t *testing.T) {
	input := `node a
`

	doc, err := New().Parse(input)
	want := &Document{
		Nodes: []Node{
			{
				Name: "node",
				Arguments: []Value{
					{Type: ValueTypeString, Value: "a"},
				},
				Properties: []Property{},
				Children:   []Node{},
			},
		},
	}

	assert.NoError(t, err)
	assert.Equal(t, want, doc)
}

func Test_ArgFalseType(t *testing.T) {
	input := `node (type)#false
`

	doc, err := New().Parse(input)
	want := &Document{
		Nodes: []Node{
			{
				Name: "node",
				Arguments: []Value{
					{Type: ValueTypeBoolean, Value: "false", TypeAnnotation: "type"},
				},
				Properties: []Property{},
				Children:   []Node{},
			},
		},
	}

	assert.NoError(t, err)
	assert.Equal(t, want, doc)
}

func Test_ArgFloatType(t *testing.T) {
	input := `node (type)2.5`

	doc, err := New().Parse(input)
	want := &Document{
		Nodes: []Node{
			{
				Name: "node",
				Arguments: []Value{
					{Type: ValueTypeNumber, Value: "2.5", TypeAnnotation: "type"},
				},
				Properties: []Property{},
				Children:   []Node{},
			},
		},
	}

	assert.NoError(t, err)
	assert.Equal(t, want, doc)
}

func Test_ArgHexType(t *testing.T) {
	input := `node (type)0x10
`

	doc, err := New().Parse(input)
	want := &Document{
		Nodes: []Node{
			{
				Name: "node",
				Arguments: []Value{
					{Type: ValueTypeNumber, Value: "0x10", TypeAnnotation: "type"},
				},
				Properties: []Property{},
				Children:   []Node{},
			},
		},
	}

	assert.NoError(t, err)
	assert.Equal(t, want, doc)
}

func Test_ArgNullType(t *testing.T) {
	input := `node (type)#null
`

	doc, err := New().Parse(input)
	want := &Document{
		Nodes: []Node{
			{
				Name: "node",
				Arguments: []Value{
					{Type: ValueTypeNull, Value: "null", TypeAnnotation: "type"},
				},
				Properties: []Property{},
				Children:   []Node{},
			},
		},
	}

	assert.NoError(t, err)
	assert.Equal(t, want, doc)
}

func Test_ArgRawStringType(t *testing.T) {
	input := `node (type)#"str"#
`

	doc, err := New().Parse(input)
	want := &Document{
		Nodes: []Node{
			{
				Name: "node",
				Arguments: []Value{
					{Type: ValueTypeString, Value: "str", TypeAnnotation: "type"},
				},
				Properties: []Property{},
				Children:   []Node{},
			},
		},
	}

	assert.NoError(t, err)
	assert.Equal(t, want, doc)
}

func Test_ArgStringType(t *testing.T) {
	input := `node (type)"str"
`

	doc, err := New().Parse(input)
	want := &Document{
		Nodes: []Node{
			{
				Name: "node",
				Arguments: []Value{
					{Type: ValueTypeString, Value: "str", TypeAnnotation: "type"},
				},
				Properties: []Property{},
				Children:   []Node{},
			},
		},
	}

	assert.NoError(t, err)
	assert.Equal(t, want, doc)
}

func Test_ArgTrueType(t *testing.T) {
	input := `node (type)#true
`

	doc, err := New().Parse(input)
	want := &Document{
		Nodes: []Node{
			{
				Name: "node",
				Arguments: []Value{
					{Type: ValueTypeBoolean, Value: "true", TypeAnnotation: "type"},
				},
				Properties: []Property{},
				Children:   []Node{},
			},
		},
	}

	assert.NoError(t, err)
	assert.Equal(t, want, doc)
}

func Test_ArgType(t *testing.T) {
	input := `node (type)arg
`

	doc, err := New().Parse(input)
	want := &Document{
		Nodes: []Node{
			{
				Name: "node",
				Arguments: []Value{
					{Type: ValueTypeString, Value: "arg", TypeAnnotation: "type"},
				},
				Properties: []Property{},
				Children:   []Node{},
			},
		},
	}

	assert.NoError(t, err)
	assert.Equal(t, want, doc)
}

func Test_ArgZeroType(t *testing.T) {
	input := `node (type)0
`

	doc, err := New().Parse(input)
	want := &Document{
		Nodes: []Node{
			{
				Name: "node",
				Arguments: []Value{
					{Type: ValueTypeNumber, Value: "0", TypeAnnotation: "type"},
				},
				Properties: []Property{},
				Children:   []Node{},
			},
		},
	}

	assert.NoError(t, err)
	assert.Equal(t, want, doc)
}

func Test_AsteriskInBlockComment(t *testing.T) {
	input := `node /* * */`

	doc, err := New().Parse(input)
	want := &Document{
		Nodes: []Node{
			{
				Name:       "node",
				Arguments:  []Value{},
				Properties: []Property{},
				Children:   []Node{},
			},
		},
	}

	assert.NoError(t, err)
	assert.Equal(t, want, doc)
}

func Test_BareEmoji(t *testing.T) {
	input := `😁 happy!
`

	doc, err := New().Parse(input)
	want := &Document{
		Nodes: []Node{
			{
				Name: "😁",
				Arguments: []Value{
					{Type: ValueTypeString, Value: "happy!"},
				},
				Properties: []Property{},
				Children:   []Node{},
			},
		},
	}

	assert.NoError(t, err)
	assert.Equal(t, want, doc)
}

func Test_BareIdentDot(t *testing.T) {
	input := `node .`

	doc, err := New().Parse(input)
	want := &Document{
		Nodes: []Node{
			{
				Name: "node",
				Arguments: []Value{
					{Type: ValueTypeString, Value: "."},
				},
				Properties: []Property{},
				Children:   []Node{},
			},
		},
	}

	assert.NoError(t, err)
	assert.Equal(t, want, doc)
}

func Test_BareIdentNumericDotFail(t *testing.T) {
	input := `node .0n`

	doc, err := New().Parse(input)

	_ = doc
	assert.Error(t, err)
}

func Test_BareIdentNumericFail(t *testing.T) {
	input := `node 0n`

	doc, err := New().Parse(input)

	_ = doc
	assert.Error(t, err)
}

func Test_BareIdentNumericSignFail(t *testing.T) {
	input := `node +0n`

	doc, err := New().Parse(input)

	_ = doc
	assert.Error(t, err)
}

func Test_BareIdentSign(t *testing.T) {
	input := `node +`

	doc, err := New().Parse(input)
	want := &Document{
		Nodes: []Node{
			{
				Name: "node",
				Arguments: []Value{
					{Type: ValueTypeString, Value: "+"},
				},
				Properties: []Property{},
				Children:   []Node{},
			},
		},
	}

	assert.NoError(t, err)
	assert.Equal(t, want, doc)
}

func Test_BareIdentSignDot(t *testing.T) {
	input := `node +.`

	doc, err := New().Parse(input)
	want := &Document{
		Nodes: []Node{
			{
				Name: "node",
				Arguments: []Value{
					{Type: ValueTypeString, Value: "+."},
				},
				Properties: []Property{},
				Children:   []Node{},
			},
		},
	}

	assert.NoError(t, err)
	assert.Equal(t, want, doc)
}

func Test_Binary(t *testing.T) {
	input := `node 0b10`

	doc, err := New().Parse(input)
	want := &Document{
		Nodes: []Node{
			{
				Name: "node",
				Arguments: []Value{
					{Type: ValueTypeNumber, Value: "2"},
				},
				Properties: []Property{},
				Children:   []Node{},
			},
		},
	}

	assert.NoError(t, err)
	assert.Equal(t, want, doc)
}

func Test_BinaryTrailingUnderscore(t *testing.T) {
	input := `node 0b10_`

	doc, err := New().Parse(input)
	want := &Document{
		Nodes: []Node{
			{
				Name: "node",
				Arguments: []Value{
					{Type: ValueTypeNumber, Value: "2"},
				},
				Properties: []Property{},
				Children:   []Node{},
			},
		},
	}

	assert.NoError(t, err)
	assert.Equal(t, want, doc)
}

func Test_BinaryUnderscore(t *testing.T) {
	input := `node 0b1_0
`

	doc, err := New().Parse(input)
	want := &Document{
		Nodes: []Node{
			{
				Name: "node",
				Arguments: []Value{
					{Type: ValueTypeNumber, Value: "2"},
				},
				Properties: []Property{},
				Children:   []Node{},
			},
		},
	}

	assert.NoError(t, err)
	assert.Equal(t, want, doc)
}

func Test_BlankArgType(t *testing.T) {
	input := `node ("")10`

	doc, err := New().Parse(input)
	want := &Document{
		Nodes: []Node{
			{
				Name: "node",
				Arguments: []Value{
					{Type: ValueTypeNumber, Value: "10", TypeAnnotation: ""},
				},
				Properties: []Property{},
				Children:   []Node{},
			},
		},
	}

	assert.NoError(t, err)
	assert.Equal(t, want, doc)
}

func Test_BlankNodeType(t *testing.T) {
	input := `("")node
`

	doc, err := New().Parse(input)
	want := &Document{
		Nodes: []Node{
			{
				Name:           "node",
				TypeAnnotation: "",
				Arguments:      []Value{},
				Properties:     []Property{},
				Children:       []Node{},
			},
		},
	}

	assert.NoError(t, err)
	assert.Equal(t, want, doc)
}

func Test_BlankPropType(t *testing.T) {
	input := `node key=("")#true
`

	doc, err := New().Parse(input)
	want := &Document{
		Nodes: []Node{
			{
				Name: "node",
				Arguments: []Value{},
				Properties: []Property{
					{Key: "key", Value: Value{Type: ValueTypeBoolean, Value: "true", TypeAnnotation: ""}},
				},
				Children: []Node{},
			},
		},
	}

	assert.NoError(t, err)
	assert.Equal(t, want, doc)
}

func Test_BlockComment(t *testing.T) {
	input := `node /* comment */ arg
`

	doc, err := New().Parse(input)
	want := &Document{
		Nodes: []Node{
			{
				Name: "node",
				Arguments: []Value{
					{Type: ValueTypeString, Value: "arg"},
				},
				Properties: []Property{},
				Children:   []Node{},
			},
		},
	}

	assert.NoError(t, err)
	assert.Equal(t, want, doc)
}

func Test_BlockCommentAfterNode(t *testing.T) {
	input := `node /* hey */ arg
`

	doc, err := New().Parse(input)
	want := &Document{
		Nodes: []Node{
			{
				Name: "node",
				Arguments: []Value{
					{Type: ValueTypeString, Value: "arg"},
				},
				Properties: []Property{},
				Children:   []Node{},
			},
		},
	}

	assert.NoError(t, err)
	assert.Equal(t, want, doc)
}

func Test_BlockCommentBeforeNode(t *testing.T) {
	input := `/* hey */ node`

	doc, err := New().Parse(input)
	want := &Document{
		Nodes: []Node{
			{
				Name:       "node",
				Arguments:  []Value{},
				Properties: []Property{},
				Children:   []Node{},
			},
		},
	}

	assert.NoError(t, err)
	assert.Equal(t, want, doc)
}

func Test_BlockCommentBeforeNodeNoSpace(t *testing.T) {
	input := `/* hey*/node
`

	doc, err := New().Parse(input)
	want := &Document{
		Nodes: []Node{
			{
				Name:       "node",
				Arguments:  []Value{},
				Properties: []Property{},
				Children:   []Node{},
			},
		},
	}

	assert.NoError(t, err)
	assert.Equal(t, want, doc)
}

func Test_BlockCommentNewline(t *testing.T) {
	input := `/* hey */
`

	doc, err := New().Parse(input)
	want := &Document{
		Nodes: []Node{},
	}

	assert.NoError(t, err)
	assert.Equal(t, want, doc)
}

func Test_TODO_BomInitial(t *testing.T) {
	input := string([]byte{0xef, 0xbb, 0xbf, 0x6e, 0x6f, 0x64, 0x65, 0x20, 0x61, 0x72, 0x67, 0x0a})

	doc, err := New().Parse(input)

	_ = doc
	_ = err
	t.Skip("TODO")
}

func Test_TODO_BomLaterFail(t *testing.T) {
	input := string([]byte{0x6e, 0x6f, 0x64, 0x65, 0x20, 0xef, 0xbb, 0xbf, 0x61, 0x72, 0x67, 0x0a})

	doc, err := New().Parse(input)

	_ = doc
	_ = err
	t.Skip("TODO")
}

func Test_BooleanArg(t *testing.T) {
	input := `node #false #true
`

	doc, err := New().Parse(input)
	want := &Document{
		Nodes: []Node{
			{
				Name: "node",
				Arguments: []Value{
					{Type: ValueTypeBoolean, Value: "false"},
					{Type: ValueTypeBoolean, Value: "true"},
				},
				Properties: []Property{},
				Children:   []Node{},
			},
		},
	}

	assert.NoError(t, err)
	assert.Equal(t, want, doc)
}

func Test_BooleanProp(t *testing.T) {
	input := `node prop1=#true prop2=#false
`

	doc, err := New().Parse(input)
	want := &Document{
		Nodes: []Node{
			{
				Name:      "node",
				Arguments: []Value{},
				Properties: []Property{
					{Key: "prop1", Value: Value{Type: ValueTypeBoolean, Value: "true"}},
					{Key: "prop2", Value: Value{Type: ValueTypeBoolean, Value: "false"}},
				},
				Children: []Node{},
			},
		},
	}

	assert.NoError(t, err)
	assert.Equal(t, want, doc)
}

func Test_BracesInBareId(t *testing.T) {
	input := `foo123{bar}
`

	doc, err := New().Parse(input)
	want := &Document{
		Nodes: []Node{
			{
				Name:       "foo123",
				Arguments:  []Value{},
				Properties: []Property{},
				Children: []Node{
					{
						Name:       "bar",
						Arguments:  []Value{},
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

func Test_ChevronsInBareId(t *testing.T) {
	input := `foo123<bar>foo weeee
`

	doc, err := New().Parse(input)
	want := &Document{
		Nodes: []Node{
			{
				Name: "foo123<bar>foo",
				Arguments: []Value{
					{Type: ValueTypeString, Value: "weeee"},
				},
				Properties: []Property{},
				Children:   []Node{},
			},
		},
	}

	assert.NoError(t, err)
	assert.Equal(t, want, doc)
}

func Test_CommaInBareId(t *testing.T) {
	input := `foo123,bar weeee
`

	doc, err := New().Parse(input)
	want := &Document{
		Nodes: []Node{
			{
				Name: "foo123,bar",
				Arguments: []Value{
					{Type: ValueTypeString, Value: "weeee"},
				},
				Properties: []Property{},
				Children:   []Node{},
			},
		},
	}

	assert.NoError(t, err)
	assert.Equal(t, want, doc)
}

func Test_CommentAfterArgType(t *testing.T) {
	input := `node (type)/*hey*/10
`

	doc, err := New().Parse(input)
	want := &Document{
		Nodes: []Node{
			{
				Name: "node",
				Arguments: []Value{
					{Type: ValueTypeNumber, Value: "10", TypeAnnotation: "type"},
				},
				Properties: []Property{},
				Children:   []Node{},
			},
		},
	}

	assert.NoError(t, err)
	assert.Equal(t, want, doc)
}

func Test_CommentAfterNodeType(t *testing.T) {
	input := `(type)/*hey*/node
`

	doc, err := New().Parse(input)
	want := &Document{
		Nodes: []Node{
			{
				Name:           "node",
				TypeAnnotation: "type",
				Arguments:      []Value{},
				Properties:     []Property{},
				Children:       []Node{},
			},
		},
	}

	assert.NoError(t, err)
	assert.Equal(t, want, doc)
}

func Test_CommentAfterPropType(t *testing.T) {
	input := `node key=(type)/*hey*/10
`

	doc, err := New().Parse(input)
	want := &Document{
		Nodes: []Node{
			{
				Name:      "node",
				Arguments: []Value{},
				Properties: []Property{
					{Key: "key", Value: Value{Type: ValueTypeNumber, Value: "10", TypeAnnotation: "type"}},
				},
				Children: []Node{},
			},
		},
	}

	assert.NoError(t, err)
	assert.Equal(t, want, doc)
}

func Test_CommentAndNewline(t *testing.T) {
	input := `node1 //
node2
`

	doc, err := New().Parse(input)
	want := &Document{
		Nodes: []Node{
			{
				Name:       "node1",
				Arguments:  []Value{},
				Properties: []Property{},
				Children:   []Node{},
			},
			{
				Name:       "node2",
				Arguments:  []Value{},
				Properties: []Property{},
				Children:   []Node{},
			},
		},
	}

	assert.NoError(t, err)
	assert.Equal(t, want, doc)
}

func Test_CommentInArgType(t *testing.T) {
	input := `node (type/*hey*/)10
`

	doc, err := New().Parse(input)
	want := &Document{
		Nodes: []Node{
			{
				Name: "node",
				Arguments: []Value{
					{Type: ValueTypeNumber, Value: "10", TypeAnnotation: "type"},
				},
				Properties: []Property{},
				Children:   []Node{},
			},
		},
	}

	assert.NoError(t, err)
	assert.Equal(t, want, doc)
}

func Test_CommentInNodeType(t *testing.T) {
	input := `(type/*hey*/)node
`

	doc, err := New().Parse(input)
	want := &Document{
		Nodes: []Node{
			{
				Name:           "node",
				TypeAnnotation: "type",
				Arguments:      []Value{},
				Properties:     []Property{},
				Children:       []Node{},
			},
		},
	}

	assert.NoError(t, err)
	assert.Equal(t, want, doc)
}

func Test_CommentInPropType(t *testing.T) {
	input := `node key=(type/*hey*/)10
`

	doc, err := New().Parse(input)
	want := &Document{
		Nodes: []Node{
			{
				Name:      "node",
				Arguments: []Value{},
				Properties: []Property{
					{Key: "key", Value: Value{Type: ValueTypeNumber, Value: "10", TypeAnnotation: "type"}},
				},
				Children: []Node{},
			},
		},
	}

	assert.NoError(t, err)
	assert.Equal(t, want, doc)
}

func Test_CommentedArg(t *testing.T) {
	input := `node /- arg1 arg2
`

	doc, err := New().Parse(input)
	want := &Document{
		Nodes: []Node{
			{
				Name: "node",
				Arguments: []Value{
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

func Test_CommentedChild(t *testing.T) {
	input := `node arg /- {
     inner_node
}
`

	doc, err := New().Parse(input)
	want := &Document{
		Nodes: []Node{
			{
				Name: "node",
				Arguments: []Value{
					{Type: ValueTypeString, Value: "arg"},
				},
				Properties: []Property{},
				Children:   []Node{},
			},
		},
	}

	assert.NoError(t, err)
	assert.Equal(t, want, doc)
}

func Test_CommentedLine(t *testing.T) {
	input := `// node_1
node_2`

	doc, err := New().Parse(input)
	want := &Document{
		Nodes: []Node{
			{
				Name:       "node_2",
				Arguments:  []Value{},
				Properties: []Property{},
				Children:   []Node{},
			},
		},
	}

	assert.NoError(t, err)
	assert.Equal(t, want, doc)
}

func Test_CommentedNode(t *testing.T) {
	input := `/- node_1
node_2
/- node_3
`

	doc, err := New().Parse(input)
	want := &Document{
		Nodes: []Node{
			{
				Name:       "node_2",
				Arguments:  []Value{},
				Properties: []Property{},
				Children:   []Node{},
			},
		},
	}

	assert.NoError(t, err)
	assert.Equal(t, want, doc)
}

func Test_CommentedProp(t *testing.T) {
	input := `node /- prop=val arg
`

	doc, err := New().Parse(input)
	want := &Document{
		Nodes: []Node{
			{
				Name: "node",
				Arguments: []Value{
					{Type: ValueTypeString, Value: "arg"},
				},
				Properties: []Property{},
				Children:   []Node{},
			},
		},
	}

	assert.NoError(t, err)
	assert.Equal(t, want, doc)
}

func Test_CrlfBetweenNodes(t *testing.T) {
	input := "node1\r\nnode2\r\n"

	doc, err := New().Parse(input)
	want := &Document{
		Nodes: []Node{
			{
				Name:       "node1",
				Arguments:  []Value{},
				Properties: []Property{},
				Children:   []Node{},
			},
			{
				Name:       "node2",
				Arguments:  []Value{},
				Properties: []Property{},
				Children:   []Node{},
			},
		},
	}

	assert.NoError(t, err)
	assert.Equal(t, want, doc)
}

func Test_DashDash(t *testing.T) {
	input := `node --
`

	doc, err := New().Parse(input)
	want := &Document{
		Nodes: []Node{
			{
				Name: "node",
				Arguments: []Value{
					{Type: ValueTypeString, Value: "--"},
				},
				Properties: []Property{},
				Children:   []Node{},
			},
		},
	}

	assert.NoError(t, err)
	assert.Equal(t, want, doc)
}

func Test_TODO_DotButNoFractionBeforeExponentFail(t *testing.T) {
	input := `node 1.e7`

	doc, err := New().Parse(input)

	_ = doc
	_ = err
	t.Skip("TODO")
}

func Test_TODO_DotButNoFractionFail(t *testing.T) {
	input := `node 1.`

	doc, err := New().Parse(input)

	_ = doc
	_ = err
	t.Skip("TODO")
}

func Test_TODO_DotInExponentFail(t *testing.T) {
	input := `node 1.0.0`

	doc, err := New().Parse(input)

	_ = doc
	_ = err
	t.Skip("TODO")
}

func Test_TODO_DotZeroFail(t *testing.T) {
	input := `node .0`

	doc, err := New().Parse(input)

	_ = doc
	_ = err
	t.Skip("TODO")
}

func Test_Emoji(t *testing.T) {
	input := `node 😀
`

	doc, err := New().Parse(input)
	want := &Document{
		Nodes: []Node{
			{
				Name: "node",
				Arguments: []Value{
					{Type: ValueTypeString, Value: "😀"},
				},
				Properties: []Property{},
				Children:   []Node{},
			},
		},
	}

	assert.NoError(t, err)
	assert.Equal(t, want, doc)
}

func Test_Empty(t *testing.T) {
	input := ``

	doc, err := New().Parse(input)
	want := &Document{
		Nodes: []Node{},
	}

	assert.NoError(t, err)
	assert.Equal(t, want, doc)
}

func Test_TODO_EmptyArgTypeFail(t *testing.T) {
	input := `node ()10
`

	doc, err := New().Parse(input)

	_ = doc
	_ = err
	t.Skip("TODO")
}

func Test_EmptyChild(t *testing.T) {
	input := `node {
}`

	doc, err := New().Parse(input)
	want := &Document{
		Nodes: []Node{
			{
				Name:       "node",
				Arguments:  []Value{},
				Properties: []Property{},
				Children:   []Node{},
			},
		},
	}

	assert.NoError(t, err)
	assert.Equal(t, want, doc)
}

func Test_EmptyChildDifferentLines(t *testing.T) {
	input := `node {
}`

	doc, err := New().Parse(input)
	want := &Document{
		Nodes: []Node{
			{
				Name:       "node",
				Arguments:  []Value{},
				Properties: []Property{},
				Children:   []Node{},
			},
		},
	}

	assert.NoError(t, err)
	assert.Equal(t, want, doc)
}

func Test_EmptyChildSameLine(t *testing.T) {
	input := `node {}`

	doc, err := New().Parse(input)
	want := &Document{
		Nodes: []Node{
			{
				Name:       "node",
				Arguments:  []Value{},
				Properties: []Property{},
				Children:   []Node{},
			},
		},
	}

	assert.NoError(t, err)
	assert.Equal(t, want, doc)
}

func Test_EmptyChildWhitespace(t *testing.T) {
	input := `node {

     }`

	doc, err := New().Parse(input)
	want := &Document{
		Nodes: []Node{
			{
				Name:       "node",
				Arguments:  []Value{},
				Properties: []Property{},
				Children:   []Node{},
			},
		},
	}

	assert.NoError(t, err)
	assert.Equal(t, want, doc)
}

func Test_EmptyLineComment(t *testing.T) {
	input := `//
node`

	doc, err := New().Parse(input)
	want := &Document{
		Nodes: []Node{
			{
				Name:       "node",
				Arguments:  []Value{},
				Properties: []Property{},
				Children:   []Node{},
			},
		},
	}

	assert.NoError(t, err)
	assert.Equal(t, want, doc)
}

func Test_TODO_EmptyNodeTypeFail(t *testing.T) {
	input := `()node
`

	doc, err := New().Parse(input)

	_ = doc
	_ = err
	t.Skip("TODO")
}

func Test_TODO_EmptyPropTypeFail(t *testing.T) {
	input := `node key=()#false
`

	doc, err := New().Parse(input)

	_ = doc
	_ = err
	t.Skip("TODO")
}

func Test_TODO_EmptyQuotedNodeId(t *testing.T) {
	input := `"" arg
`

	doc, err := New().Parse(input)

	_ = doc
	_ = err
	t.Skip("TODO")
}

func Test_TODO_EmptyQuotedPropKey(t *testing.T) {
	input := `node ""=empty
`

	doc, err := New().Parse(input)

	_ = doc
	_ = err
	t.Skip("TODO")
}

func Test_TODO_EmptyStringArg(t *testing.T) {
	input := `node ""
`

	doc, err := New().Parse(input)

	_ = doc
	_ = err
	t.Skip("TODO")
}

func Test_TODO_EofAfterEscape(t *testing.T) {
	input := `node \`

	doc, err := New().Parse(input)

	_ = doc
	_ = err
	t.Skip("TODO")
}

func Test_TODO_ErrBackslashInBareIdFail(t *testing.T) {
	input := `foo123\bar weeee
`

	doc, err := New().Parse(input)

	_ = doc
	_ = err
	t.Skip("TODO")
}

func Test_TODO_EscMultipleNewlines(t *testing.T) {
	input := `node "1\


2"
`

	doc, err := New().Parse(input)

	_ = doc
	_ = err
	t.Skip("TODO")
}

func Test_TODO_EscNewlineInString(t *testing.T) {
	input := `node "hello\nworld"`

	doc, err := New().Parse(input)

	_ = doc
	_ = err
	t.Skip("TODO")
}

func Test_TODO_EscUnicodeInString(t *testing.T) {
	input := `node "hello\u{0a}world"
`

	doc, err := New().Parse(input)

	_ = doc
	_ = err
	t.Skip("TODO")
}

func Test_TODO_EscapedWhitespace(t *testing.T) {
	input := `// All of these strings are the same
node \
	"Hello\n\tWorld" \
	"""
	Hello
		World
	""" \
	"Hello\n\      \tWorld" \
	"Hello\n\
    \tWorld" \
	"Hello\n\t\
        World"

// Note that this file deliberately mixes space and newline indentation for
// test purposes
`

	doc, err := New().Parse(input)

	_ = doc
	_ = err
	t.Skip("TODO")
}

func Test_TODO_Escline(t *testing.T) {
	input := `node \
    arg
`

	doc, err := New().Parse(input)

	_ = doc
	_ = err
	t.Skip("TODO")
}

func Test_TODO_EsclineAfterSemicolon(t *testing.T) {
	input := `node; \
node
`

	doc, err := New().Parse(input)

	_ = doc
	_ = err
	t.Skip("TODO")
}

func Test_TODO_EsclineAlone(t *testing.T) {
	input := `\
`

	doc, err := New().Parse(input)

	_ = doc
	_ = err
	t.Skip("TODO")
}

func Test_TODO_EsclineEmptyLine(t *testing.T) {
	input := `\

node
`

	doc, err := New().Parse(input)

	_ = doc
	_ = err
	t.Skip("TODO")
}

func Test_TODO_EsclineEndOfNode(t *testing.T) {
	input := `a \

b
`

	doc, err := New().Parse(input)

	_ = doc
	_ = err
	t.Skip("TODO")
}

func Test_TODO_EsclineInChildBlock(t *testing.T) {
	input := `parent {
    child
    \ // comment
    child
}
`

	doc, err := New().Parse(input)

	_ = doc
	_ = err
	t.Skip("TODO")
}

func Test_TODO_EsclineLineComment(t *testing.T) {
	input := `node \   // comment
    arg \// comment
    arg2
`

	doc, err := New().Parse(input)

	_ = doc
	_ = err
	t.Skip("TODO")
}

func Test_TODO_EsclineNode(t *testing.T) {
	input := `node1
\
node2
`

	doc, err := New().Parse(input)

	_ = doc
	_ = err
	t.Skip("TODO")
}

func Test_TODO_EsclineNodeType(t *testing.T) {
	input := `\
(type)node
`

	doc, err := New().Parse(input)

	_ = doc
	_ = err
	t.Skip("TODO")
}

func Test_TODO_EsclineSlashdash(t *testing.T) {
	input := `node
\
/-
node
`

	doc, err := New().Parse(input)

	_ = doc
	_ = err
	t.Skip("TODO")
}

func Test_FalsePrefixInBareId(t *testing.T) {
	input := `false_id
`

	doc, err := New().Parse(input)
	want := &Document{
		Nodes: []Node{
			{
				Name:       "false_id",
				Arguments:  []Value{},
				Properties: []Property{},
				Children:   []Node{},
			},
		},
	}

	assert.NoError(t, err)
	assert.Equal(t, want, doc)
}

func Test_FalsePrefixInPropKey(t *testing.T) {
	input := `node false_id=1
`

	doc, err := New().Parse(input)
	want := &Document{
		Nodes: []Node{
			{
				Name:      "node",
				Arguments: []Value{},
				Properties: []Property{
					{Key: "false_id", Value: Value{Type: ValueTypeNumber, Value: "1"}},
				},
				Children: []Node{},
			},
		},
	}

	assert.NoError(t, err)
	assert.Equal(t, want, doc)
}

func Test_TODO_FalsePropKeyFail(t *testing.T) {
	input := `node false=1
`

	doc, err := New().Parse(input)

	_ = doc
	_ = err
	t.Skip("TODO")
}

func Test_TODO_FloatingPointKeywordIdentifierStringsFail(t *testing.T) {
	input := `floats inf -inf nan
`

	doc, err := New().Parse(input)

	_ = doc
	_ = err
	t.Skip("TODO")
}

func Test_TODO_FloatingPointKeywords(t *testing.T) {
	input := `floats #inf #-inf #nan
`

	doc, err := New().Parse(input)

	_ = doc
	_ = err
	t.Skip("TODO")
}

func Test_TODO_HashInIdFail(t *testing.T) {
	input := `foo#bar weee
`

	doc, err := New().Parse(input)

	_ = doc
	_ = err
	t.Skip("TODO")
}

func Test_Hex(t *testing.T) {
	input := `node 0xabcdef1234567890`

	doc, err := New().Parse(input)
	want := &Document{
		Nodes: []Node{
			{
				Name: "node",
				Arguments: []Value{
					{Type: ValueTypeNumber, Value: "0xabcdef1234567890"},
				},
				Properties: []Property{},
				Children:   []Node{},
			},
		},
	}

	assert.NoError(t, err)
	assert.Equal(t, want, doc)
}

func Test_HexInt(t *testing.T) {
	input := `node 0xABCDEF0123456789abcdef
`

	doc, err := New().Parse(input)
	want := &Document{
		Nodes: []Node{
			{
				Name: "node",
				Arguments: []Value{
					{Type: ValueTypeNumber, Value: "0xABCDEF0123456789abcdef"},
				},
				Properties: []Property{},
				Children:   []Node{},
			},
		},
	}

	assert.NoError(t, err)
	assert.Equal(t, want, doc)
}

func Test_TODO_HexIntUnderscores(t *testing.T) {
	input := `node 0xABC_def_0123`

	doc, err := New().Parse(input)

	_ = doc
	_ = err
	t.Skip("TODO")
}

func Test_TODO_HexLeadingZero(t *testing.T) {
	input := `node 0x01`

	doc, err := New().Parse(input)

	_ = doc
	_ = err
	t.Skip("TODO")
}

func Test_TODO_IllegalCharInBinaryFail(t *testing.T) {
	input := `node 0bx01
`

	doc, err := New().Parse(input)

	_ = doc
	_ = err
	t.Skip("TODO")
}

func Test_TODO_IllegalCharInHexFail(t *testing.T) {
	input := `node 0x10g10`

	doc, err := New().Parse(input)

	_ = doc
	_ = err
	t.Skip("TODO")
}

func Test_TODO_IllegalCharInOctalFail(t *testing.T) {
	input := `node 0o45678`

	doc, err := New().Parse(input)

	_ = doc
	_ = err
	t.Skip("TODO")
}

func Test_InitialSlashdash(t *testing.T) {
	input := `/-node here
another-node
`

	doc, err := New().Parse(input)
	want := &Document{
		Nodes: []Node{
			{
				Name:       "another-node",
				Arguments:  []Value{},
				Properties: []Property{},
				Children:   []Node{},
			},
		},
	}

	assert.NoError(t, err)
	assert.Equal(t, want, doc)
}

func Test_TODO_IntMultipleUnderscore(t *testing.T) {
	input := `node 1_2_3_4`

	doc, err := New().Parse(input)

	_ = doc
	_ = err
	t.Skip("TODO")
}

func Test_JustBlockComment(t *testing.T) {
	input := `/* hey */`

	doc, err := New().Parse(input)
	want := &Document{
		Nodes: []Node{},
	}

	assert.NoError(t, err)
	assert.Equal(t, want, doc)
}

func Test_JustChild(t *testing.T) {
	input := `node {
    inner_node
}`

	doc, err := New().Parse(input)
	want := &Document{
		Nodes: []Node{
			{
				Name:       "node",
				Arguments:  []Value{},
				Properties: []Property{},
				Children: []Node{
					{
						Name:       "inner_node",
						Arguments:  []Value{},
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

func Test_JustNewline(t *testing.T) {
	input := `
`

	doc, err := New().Parse(input)
	want := &Document{
		Nodes: []Node{},
	}

	assert.NoError(t, err)
	assert.Equal(t, want, doc)
}

func Test_JustNodeId(t *testing.T) {
	input := `node`

	doc, err := New().Parse(input)
	want := &Document{
		Nodes: []Node{
			{
				Name:       "node",
				Arguments:  []Value{},
				Properties: []Property{},
				Children:   []Node{},
			},
		},
	}

	assert.NoError(t, err)
	assert.Equal(t, want, doc)
}

func Test_JustSpace(t *testing.T) {
	input := ` `

	doc, err := New().Parse(input)
	want := &Document{
		Nodes: []Node{},
	}

	assert.NoError(t, err)
	assert.Equal(t, want, doc)
}

func Test_TODO_JustSpaceInArgTypeFail(t *testing.T) {
	input := `node ( )false
`

	doc, err := New().Parse(input)

	_ = doc
	_ = err
	t.Skip("TODO")
}

func Test_TODO_JustSpaceInNodeTypeFail(t *testing.T) {
	input := `( )node
`

	doc, err := New().Parse(input)

	_ = doc
	_ = err
	t.Skip("TODO")
}

func Test_TODO_JustSpaceInPropTypeFail(t *testing.T) {
	input := `node key=( )0x10
`

	doc, err := New().Parse(input)

	_ = doc
	_ = err
	t.Skip("TODO")
}

func Test_TODO_JustTypeNoArgFail(t *testing.T) {
	input := `node (type)
`

	doc, err := New().Parse(input)

	_ = doc
	_ = err
	t.Skip("TODO")
}

func Test_TODO_JustTypeNoNodeIdFail(t *testing.T) {
	input := `(type)
`

	doc, err := New().Parse(input)

	_ = doc
	_ = err
	t.Skip("TODO")
}

func Test_TODO_JustTypeNoPropFail(t *testing.T) {
	input := `node key=(type)
`

	doc, err := New().Parse(input)

	_ = doc
	_ = err
	t.Skip("TODO")
}

func Test_LeadingNewline(t *testing.T) {
	input := `
node`

	doc, err := New().Parse(input)
	want := &Document{
		Nodes: []Node{
			{
				Name:       "node",
				Arguments:  []Value{},
				Properties: []Property{},
				Children:   []Node{},
			},
		},
	}

	assert.NoError(t, err)
	assert.Equal(t, want, doc)
}

func Test_TODO_LeadingZeroBinary(t *testing.T) {
	input := `node 0b01
`

	doc, err := New().Parse(input)

	_ = doc
	_ = err
	t.Skip("TODO")
}

func Test_TODO_LeadingZeroInt(t *testing.T) {
	input := `node 011
`

	doc, err := New().Parse(input)

	_ = doc
	_ = err
	t.Skip("TODO")
}

func Test_TODO_LeadingZeroOct(t *testing.T) {
	input := `node 0o01
`

	doc, err := New().Parse(input)

	_ = doc
	_ = err
	t.Skip("TODO")
}

func Test_TODO_LegacyRawStringFail(t *testing.T) {
	input := `node r"foo"
`

	doc, err := New().Parse(input)

	_ = doc
	_ = err
	t.Skip("TODO")
}

func Test_TODO_LegacyRawStringHashFail(t *testing.T) {
	input := `node r#"foo"#
`

	doc, err := New().Parse(input)

	_ = doc
	_ = err
	t.Skip("TODO")
}

func Test_MultilineComment(t *testing.T) {
	input := `node /*
some
comments
*/ arg
`

	doc, err := New().Parse(input)
	want := &Document{
		Nodes: []Node{
			{
				Name: "node",
				Arguments: []Value{
					{Type: ValueTypeString, Value: "arg"},
				},
				Properties: []Property{},
				Children:   []Node{},
			},
		},
	}

	assert.NoError(t, err)
	assert.Equal(t, want, doc)
}

func Test_TODO_MultilineNodes(t *testing.T) {
	input := `node \
    arg1 \// comment
    arg2
`

	doc, err := New().Parse(input)

	_ = doc
	_ = err
	t.Skip("TODO")
}

func Test_TODO_MultilineRawString(t *testing.T) {
	input := `node #"""
hey
everyone
how goes?
"""#
`

	doc, err := New().Parse(input)

	_ = doc
	_ = err
	t.Skip("TODO")
}

func Test_TODO_MultilineRawStringContainingQuotes(t *testing.T) {
	input := `node ##"""
"""triple-quote"""
##"too few quotes"##
#"""too few #"""#
"""##
`

	doc, err := New().Parse(input)

	_ = doc
	_ = err
	t.Skip("TODO")
}

func Test_TODO_MultilineRawStringEmpty(t *testing.T) {
	input := `node #"""
"""#`

	doc, err := New().Parse(input)

	_ = doc
	_ = err
	t.Skip("TODO")
}

func Test_TODO_MultilineRawStringEmptyIndented(t *testing.T) {
	input := `node #"""
	"""#`

	doc, err := New().Parse(input)

	_ = doc
	_ = err
	t.Skip("TODO")
}

func Test_TODO_MultilineRawStringIndented(t *testing.T) {
	input := `node #"""
    hey
   everyone
     how goes?
  """#
`

	doc, err := New().Parse(input)

	_ = doc
	_ = err
	t.Skip("TODO")
}

func Test_TODO_MultilineRawStringNonMatchingPrefixCharacterErrorFail(t *testing.T) {
	input := `node #"""
    hey
   everyone
	   how goes?
  """#
`

	doc, err := New().Parse(input)

	_ = doc
	_ = err
	t.Skip("TODO")
}

func Test_TODO_MultilineRawStringNonMatchingPrefixCountErrorFail(t *testing.T) {
	input := `node #"""
    hey
 everyone
     how goes?
  """#
`

	doc, err := New().Parse(input)

	_ = doc
	_ = err
	t.Skip("TODO")
}

func Test_TODO_MultilineRawStringSingleLineErrFail(t *testing.T) {
	input := `node #"""one line"""#`

	doc, err := New().Parse(input)

	_ = doc
	_ = err
	t.Skip("TODO")
}

func Test_TODO_MultilineRawStringSingleQuoteErrFail(t *testing.T) {
	input := `node #"
hey
everyone
how goes?
"#
`

	doc, err := New().Parse(input)

	_ = doc
	_ = err
	t.Skip("TODO")
}

func Test_TODO_MultilineString(t *testing.T) {
	input := `node """
hey
everyone
how goes?
"""
`

	doc, err := New().Parse(input)

	_ = doc
	_ = err
	t.Skip("TODO")
}

func Test_TODO_MultilineStringContainingQuotes(t *testing.T) {
	input := `node """
this string contains "quotes", twice""
"""
`

	doc, err := New().Parse(input)

	_ = doc
	_ = err
	t.Skip("TODO")
}

func Test_TODO_MultilineStringDoubleBackslash(t *testing.T) {
	input := `node """
a\\ b
a\\\ b
"""
`

	doc, err := New().Parse(input)

	_ = doc
	_ = err
	t.Skip("TODO")
}

func Test_TODO_MultilineStringEmpty(t *testing.T) {
	input := `node """
"""`

	doc, err := New().Parse(input)

	_ = doc
	_ = err
	t.Skip("TODO")
}

func Test_TODO_MultilineStringEmptyIndented(t *testing.T) {
	input := `node """
	"""`

	doc, err := New().Parse(input)

	_ = doc
	_ = err
	t.Skip("TODO")
}

func Test_TODO_MultilineStringEscapeDelimiter(t *testing.T) {
	input := `node """
\"""
"""
`

	doc, err := New().Parse(input)

	_ = doc
	_ = err
	t.Skip("TODO")
}

func Test_TODO_MultilineStringEscapeInClosingLine(t *testing.T) {
	input := `node """
  foo \
bar
  baz
  \   """
`

	doc, err := New().Parse(input)

	_ = doc
	_ = err
	t.Skip("TODO")
}

func Test_TODO_MultilineStringEscapeInClosingLineShallow(t *testing.T) {
	input := `node """
  foo \
bar
  baz
\   """
`

	doc, err := New().Parse(input)

	_ = doc
	_ = err
	t.Skip("TODO")
}

func Test_TODO_MultilineStringEscapeNewlineAtEnd(t *testing.T) {
	input := `node """
    a
   \
"""
`

	doc, err := New().Parse(input)

	_ = doc
	_ = err
	t.Skip("TODO")
}

func Test_TODO_MultilineStringEscapeNewlineAtEndFail(t *testing.T) {
	input := `node """
a
   \
"""
`

	doc, err := New().Parse(input)

	_ = doc
	_ = err
	t.Skip("TODO")
}

func Test_TODO_MultilineStringFinalWhitespaceEscapeFail(t *testing.T) {
	input := `node """
  foo
  bar\
  """`

	doc, err := New().Parse(input)

	_ = doc
	_ = err
	t.Skip("TODO")
}

func Test_TODO_MultilineStringIndented(t *testing.T) {
	input := `node """
    hey
   everyone
     how goes?
  """
`

	doc, err := New().Parse(input)

	_ = doc
	_ = err
	t.Skip("TODO")
}

func Test_TODO_MultilineStringNonLiteralPrefixFail(t *testing.T) {
	input := `node """
\s escaped prefix
  literal prefix
  """
`

	doc, err := New().Parse(input)

	_ = doc
	_ = err
	t.Skip("TODO")
}

func Test_TODO_MultilineStringNonMatchingPrefixCharacterErrorFail(t *testing.T) {
	input := `node """
    hey
   everyone
	   how goes?
  """
`

	doc, err := New().Parse(input)

	_ = doc
	_ = err
	t.Skip("TODO")
}

func Test_TODO_MultilineStringNonMatchingPrefixCountErrorFail(t *testing.T) {
	input := `node """
    hey
 everyone
     how goes?
  """
`

	doc, err := New().Parse(input)

	_ = doc
	_ = err
	t.Skip("TODO")
}

func Test_TODO_MultilineStringSingleLineErrFail(t *testing.T) {
	input := `node """one line"""`

	doc, err := New().Parse(input)

	_ = doc
	_ = err
	t.Skip("TODO")
}

func Test_TODO_MultilineStringSingleQuoteErrFail(t *testing.T) {
	input := `node "
hey
everyone
how goes?
"
`

	doc, err := New().Parse(input)

	_ = doc
	_ = err
	t.Skip("TODO")
}

func Test_TODO_MultilineStringWhitespaceOnly(t *testing.T) {
	input := `// This file deliberately contains unusual whitespace
// The first two strings are empty
node """
  	""" """
 	 \
     
 	 """ """
       
 """\
    \ // The next two strings contains only whitespace
    """
   

      \s
    """ #"""
    

  """#
`

	doc, err := New().Parse(input)

	_ = doc
	_ = err
	t.Skip("TODO")
}

func Test_TODO_MultilineStringWrappedBinary(t *testing.T) {
	input := `node """
    dead\
    beef
    """
`

	doc, err := New().Parse(input)

	_ = doc
	_ = err
	t.Skip("TODO")
}

func Test_TODO_MultipleDotsInFloatBeforeExponentFail(t *testing.T) {
	input := `node 1.0.0e7`

	doc, err := New().Parse(input)

	_ = doc
	_ = err
	t.Skip("TODO")
}

func Test_TODO_MultipleDotsInFloatFail(t *testing.T) {
	input := `node 1.0.0`

	doc, err := New().Parse(input)

	_ = doc
	_ = err
	t.Skip("TODO")
}

func Test_TODO_MultipleEsInFloatFail(t *testing.T) {
	input := `node 1.0E10e10
`

	doc, err := New().Parse(input)

	_ = doc
	_ = err
	t.Skip("TODO")
}

func Test_TODO_MultipleXInHexFail(t *testing.T) {
	input := `node 0xx10`

	doc, err := New().Parse(input)

	_ = doc
	_ = err
	t.Skip("TODO")
}

func Test_NegativeExponent(t *testing.T) {
	input := `node 1.0e-10`

	doc, err := New().Parse(input)
	want := &Document{
		Nodes: []Node{
			{
				Name: "node",
				Arguments: []Value{
					{Type: ValueTypeNumber, Value: "1.0e-10"},
				},
				Properties: []Property{},
				Children:   []Node{},
			},
		},
	}

	assert.NoError(t, err)
	assert.Equal(t, want, doc)
}

func Test_NegativeFloat(t *testing.T) {
	input := `node -1.0 key=-10.0`

	doc, err := New().Parse(input)
	want := &Document{
		Nodes: []Node{
			{
				Name: "node",
				Arguments: []Value{
					{Type: ValueTypeNumber, Value: "-1.0"},
				},
				Properties: []Property{
					{Key: "key", Value: Value{Type: ValueTypeNumber, Value: "-10.0"}},
				},
				Children: []Node{},
			},
		},
	}

	assert.NoError(t, err)
	assert.Equal(t, want, doc)
}

func Test_NegativeInt(t *testing.T) {
	input := `node -10 prop=-15`

	doc, err := New().Parse(input)
	want := &Document{
		Nodes: []Node{
			{
				Name: "node",
				Arguments: []Value{
					{Type: ValueTypeNumber, Value: "-10"},
				},
				Properties: []Property{
					{Key: "prop", Value: Value{Type: ValueTypeNumber, Value: "-15"}},
				},
				Children: []Node{},
			},
		},
	}

	assert.NoError(t, err)
	assert.Equal(t, want, doc)
}

func Test_NestedBlockComment(t *testing.T) {
	input := `node /* hi /* there */ everyone */ arg
`

	doc, err := New().Parse(input)
	want := &Document{
		Nodes: []Node{
			{
				Name: "node",
				Arguments: []Value{
					{Type: ValueTypeString, Value: "arg"},
				},
				Properties: []Property{},
				Children:   []Node{},
			},
		},
	}

	assert.NoError(t, err)
	assert.Equal(t, want, doc)
}

func Test_NestedChildren(t *testing.T) {
	input := `node1 {
    node2 {
        node
    }
}`

	doc, err := New().Parse(input)
	want := &Document{
		Nodes: []Node{
			{
				Name:       "node1",
				Arguments:  []Value{},
				Properties: []Property{},
				Children: []Node{
					{
						Name:       "node2",
						Arguments:  []Value{},
						Properties: []Property{},
						Children: []Node{
							{
								Name:       "node",
								Arguments:  []Value{},
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

func Test_TODO_NestedComments(t *testing.T) {
	input := `node /*/* nested */*/ arg
`

	doc, err := New().Parse(input)

	_ = doc
	_ = err
	t.Skip("TODO")
}

func Test_TODO_NestedMultilineBlockComment(t *testing.T) {
	input := `node /*
hey /*
how's
*/
    it going
    */ arg
`

	doc, err := New().Parse(input)

	_ = doc
	_ = err
	t.Skip("TODO")
}

func Test_NewlineBetweenNodes(t *testing.T) {
	input := `node1
node2
`

	doc, err := New().Parse(input)
	want := &Document{
		Nodes: []Node{
			{
				Name:       "node1",
				Arguments:  []Value{},
				Properties: []Property{},
				Children:   []Node{},
			},
			{
				Name:       "node2",
				Arguments:  []Value{},
				Properties: []Property{},
				Children:   []Node{},
			},
		},
	}

	assert.NoError(t, err)
	assert.Equal(t, want, doc)
}

func Test_TODO_NewlinesInBlockComment(t *testing.T) {
	input := `node /* hey so
I was thinking
about newts */ arg
`

	doc, err := New().Parse(input)

	_ = doc
	_ = err
	t.Skip("TODO")
}

func Test_TODO_NoDecimalExponent(t *testing.T) {
	input := `node 1e10`

	doc, err := New().Parse(input)

	_ = doc
	_ = err
	t.Skip("TODO")
}

func Test_TODO_NoDigitsInHexFail(t *testing.T) {
	input := `node 0x`

	doc, err := New().Parse(input)

	_ = doc
	_ = err
	t.Skip("TODO")
}

func Test_TODO_NoIntegerDigitFail(t *testing.T) {
	input := `node .1`

	doc, err := New().Parse(input)

	_ = doc
	_ = err
	t.Skip("TODO")
}

func Test_TODO_NoSolidusEscapeFail(t *testing.T) {
	input := `node "\/"
`

	doc, err := New().Parse(input)

	_ = doc
	_ = err
	t.Skip("TODO")
}

func Test_NodeFalse(t *testing.T) {
	input := `node #false
`

	doc, err := New().Parse(input)
	want := &Document{
		Nodes: []Node{
			{
				Name: "node",
				Arguments: []Value{
					{Type: ValueTypeBoolean, Value: "false"},
				},
				Properties: []Property{},
				Children:   []Node{},
			},
		},
	}

	assert.NoError(t, err)
	assert.Equal(t, want, doc)
}

func Test_NodeTrue(t *testing.T) {
	input := `node #true
`

	doc, err := New().Parse(input)
	want := &Document{
		Nodes: []Node{
			{
				Name: "node",
				Arguments: []Value{
					{Type: ValueTypeBoolean, Value: "true"},
				},
				Properties: []Property{},
				Children:   []Node{},
			},
		},
	}

	assert.NoError(t, err)
	assert.Equal(t, want, doc)
}

func Test_NodeType(t *testing.T) {
	input := `(type)node`

	doc, err := New().Parse(input)
	want := &Document{
		Nodes: []Node{
			{
				Name:           "node",
				TypeAnnotation: "type",
				Arguments:      []Value{},
				Properties:     []Property{},
				Children:       []Node{},
			},
		},
	}

	assert.NoError(t, err)
	assert.Equal(t, want, doc)
}

func Test_NullArg(t *testing.T) {
	input := `node #null
`

	doc, err := New().Parse(input)
	want := &Document{
		Nodes: []Node{
			{
				Name: "node",
				Arguments: []Value{
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

func Test_NullPrefixInBareId(t *testing.T) {
	input := `null_id
`

	doc, err := New().Parse(input)
	want := &Document{
		Nodes: []Node{
			{
				Name:       "null_id",
				Arguments:  []Value{},
				Properties: []Property{},
				Children:   []Node{},
			},
		},
	}

	assert.NoError(t, err)
	assert.Equal(t, want, doc)
}

func Test_NullPrefixInPropKey(t *testing.T) {
	input := `node null_id=1
`

	doc, err := New().Parse(input)
	want := &Document{
		Nodes: []Node{
			{
				Name:      "node",
				Arguments: []Value{},
				Properties: []Property{
					{Key: "null_id", Value: Value{Type: ValueTypeNumber, Value: "1"}},
				},
				Children: []Node{},
			},
		},
	}

	assert.NoError(t, err)
	assert.Equal(t, want, doc)
}

func Test_NullProp(t *testing.T) {
	input := `node prop=#null
`

	doc, err := New().Parse(input)
	want := &Document{
		Nodes: []Node{
			{
				Name:      "node",
				Arguments: []Value{},
				Properties: []Property{
					{Key: "prop", Value: Value{Type: ValueTypeNull, Value: "null"}},
				},
				Children: []Node{},
			},
		},
	}

	assert.NoError(t, err)
	assert.Equal(t, want, doc)
}

func Test_TODO_NullPropKeyFail(t *testing.T) {
	input := `node null=1
`

	doc, err := New().Parse(input)

	_ = doc
	_ = err
	t.Skip("TODO")
}

func Test_NumericArg(t *testing.T) {
	input := `node 15.7`

	doc, err := New().Parse(input)
	want := &Document{
		Nodes: []Node{
			{
				Name: "node",
				Arguments: []Value{
					{Type: ValueTypeNumber, Value: "15.7"},
				},
				Properties: []Property{},
				Children:   []Node{},
			},
		},
	}

	assert.NoError(t, err)
	assert.Equal(t, want, doc)
}

func Test_NumericProp(t *testing.T) {
	input := `node prop=10.0`

	doc, err := New().Parse(input)
	want := &Document{
		Nodes: []Node{
			{
				Name:      "node",
				Arguments: []Value{},
				Properties: []Property{
					{Key: "prop", Value: Value{Type: ValueTypeNumber, Value: "10.0"}},
				},
				Children: []Node{},
			},
		},
	}

	assert.NoError(t, err)
	assert.Equal(t, want, doc)
}

func Test_Octal(t *testing.T) {
	input := `node 0o76543210`

	doc, err := New().Parse(input)
	want := &Document{
		Nodes: []Node{
			{
				Name: "node",
				Arguments: []Value{
					{Type: ValueTypeNumber, Value: "0o76543210"},
				},
				Properties: []Property{},
				Children:   []Node{},
			},
		},
	}

	assert.NoError(t, err)
	assert.Equal(t, want, doc)
}

func Test_OnlyCr(t *testing.T) {
	input := "\r"

	doc, err := New().Parse(input)
	want := &Document{
		Nodes: []Node{},
	}

	assert.NoError(t, err)
	assert.Equal(t, want, doc)
}

func Test_OnlyLineComment(t *testing.T) {
	input := `// hi`

	doc, err := New().Parse(input)
	want := &Document{
		Nodes: []Node{},
	}

	assert.NoError(t, err)
	assert.Equal(t, want, doc)
}

func Test_OnlyLineCommentCrlf(t *testing.T) {
	input := "// comment\r\n"

	doc, err := New().Parse(input)
	want := &Document{
		Nodes: []Node{},
	}

	assert.NoError(t, err)
	assert.Equal(t, want, doc)
}

func Test_OnlyLineCommentNewline(t *testing.T) {
	input := `// hiiii
`

	doc, err := New().Parse(input)
	want := &Document{
		Nodes: []Node{},
	}

	assert.NoError(t, err)
	assert.Equal(t, want, doc)
}

func Test_TODO_OptionalChildSemicolon(t *testing.T) {
	input := `node {foo;bar;baz}
`

	doc, err := New().Parse(input)

	_ = doc
	_ = err
	t.Skip("TODO")
}

func Test_TODO_ParensInBareIdFail(t *testing.T) {
	input := `foo123(bar)foo weeee
`

	doc, err := New().Parse(input)

	_ = doc
	_ = err
	t.Skip("TODO")
}

func Test_TODO_ParseAllArgTypes(t *testing.T) {
	input := `node 1 1.0 1.0e10 1.0e-10 0x01 0o07 0b10 arg "arg" #"arg\"# #true #false #null
`

	doc, err := New().Parse(input)

	_ = doc
	_ = err
	t.Skip("TODO")
}

func Test_PositiveExponent(t *testing.T) {
	input := `node 1.0e+10`

	doc, err := New().Parse(input)
	want := &Document{
		Nodes: []Node{
			{
				Name: "node",
				Arguments: []Value{
					{Type: ValueTypeNumber, Value: "1.0e+10"},
				},
				Properties: []Property{},
				Children:   []Node{},
			},
		},
	}

	assert.NoError(t, err)
	assert.Equal(t, want, doc)
}

func Test_PositiveInt(t *testing.T) {
	input := `node +10`

	doc, err := New().Parse(input)
	want := &Document{
		Nodes: []Node{
			{
				Name: "node",
				Arguments: []Value{
					{Type: ValueTypeNumber, Value: "+10"},
				},
				Properties: []Property{},
				Children:   []Node{},
			},
		},
	}

	assert.NoError(t, err)
	assert.Equal(t, want, doc)
}

func Test_PreserveDuplicateNodes(t *testing.T) {
	input := `node
node
`

	doc, err := New().Parse(input)
	want := &Document{
		Nodes: []Node{
			{
				Name:       "node",
				Arguments:  []Value{},
				Properties: []Property{},
				Children:   []Node{},
			},
			{
				Name:       "node",
				Arguments:  []Value{},
				Properties: []Property{},
				Children:   []Node{},
			},
		},
	}

	assert.NoError(t, err)
	assert.Equal(t, want, doc)
}

func Test_PreserveNodeOrder(t *testing.T) {
	input := `node2
node5
node1`

	doc, err := New().Parse(input)
	want := &Document{
		Nodes: []Node{
			{
				Name:       "node2",
				Arguments:  []Value{},
				Properties: []Property{},
				Children:   []Node{},
			},
			{
				Name:       "node5",
				Arguments:  []Value{},
				Properties: []Property{},
				Children:   []Node{},
			},
			{
				Name:       "node1",
				Arguments:  []Value{},
				Properties: []Property{},
				Children:   []Node{},
			},
		},
	}

	assert.NoError(t, err)
	assert.Equal(t, want, doc)
}

func Test_PropFalseType(t *testing.T) {
	input := `node key=(type)#false
`

	doc, err := New().Parse(input)
	want := &Document{
		Nodes: []Node{
			{
				Name:      "node",
				Arguments: []Value{},
				Properties: []Property{
					{Key: "key", Value: Value{Type: ValueTypeBoolean, Value: "false", TypeAnnotation: "type"}},
				},
				Children: []Node{},
			},
		},
	}

	assert.NoError(t, err)
	assert.Equal(t, want, doc)
}

func Test_PropFloatType(t *testing.T) {
	input := `node key=(type)2.5E10
`

	doc, err := New().Parse(input)
	want := &Document{
		Nodes: []Node{
			{
				Name:      "node",
				Arguments: []Value{},
				Properties: []Property{
					{Key: "key", Value: Value{Type: ValueTypeNumber, Value: "2.5E10", TypeAnnotation: "type"}},
				},
				Children: []Node{},
			},
		},
	}

	assert.NoError(t, err)
	assert.Equal(t, want, doc)
}

func Test_PropHexType(t *testing.T) {
	input := `node key=(type)0x10
`

	doc, err := New().Parse(input)
	want := &Document{
		Nodes: []Node{
			{
				Name:      "node",
				Arguments: []Value{},
				Properties: []Property{
					{Key: "key", Value: Value{Type: ValueTypeNumber, Value: "0x10", TypeAnnotation: "type"}},
				},
				Children: []Node{},
			},
		},
	}

	assert.NoError(t, err)
	assert.Equal(t, want, doc)
}

func Test_PropIdentifierType(t *testing.T) {
	input := `node key=(type)str
`

	doc, err := New().Parse(input)
	want := &Document{
		Nodes: []Node{
			{
				Name:      "node",
				Arguments: []Value{},
				Properties: []Property{
					{Key: "key", Value: Value{Type: ValueTypeString, Value: "str", TypeAnnotation: "type"}},
				},
				Children: []Node{},
			},
		},
	}

	assert.NoError(t, err)
	assert.Equal(t, want, doc)
}

func Test_PropNullType(t *testing.T) {
	input := `node key=(type)#null
`

	doc, err := New().Parse(input)
	want := &Document{
		Nodes: []Node{
			{
				Name:      "node",
				Arguments: []Value{},
				Properties: []Property{
					{Key: "key", Value: Value{Type: ValueTypeNull, Value: "null", TypeAnnotation: "type"}},
				},
				Children: []Node{},
			},
		},
	}

	assert.NoError(t, err)
	assert.Equal(t, want, doc)
}

func Test_TODO_PropRawStringType(t *testing.T) {
	input := `node key=(type)#"str"#
`

	doc, err := New().Parse(input)

	_ = doc
	_ = err
	t.Skip("TODO")
}

func Test_PropStringType(t *testing.T) {
	input := `node key=(type)"str"
`

	doc, err := New().Parse(input)
	want := &Document{
		Nodes: []Node{
			{
				Name:      "node",
				Arguments: []Value{},
				Properties: []Property{
					{Key: "key", Value: Value{Type: ValueTypeString, Value: "str", TypeAnnotation: "type"}},
				},
				Children: []Node{},
			},
		},
	}

	assert.NoError(t, err)
	assert.Equal(t, want, doc)
}

func Test_PropTrueType(t *testing.T) {
	input := `node key=(type)#true
`

	doc, err := New().Parse(input)
	want := &Document{
		Nodes: []Node{
			{
				Name:      "node",
				Arguments: []Value{},
				Properties: []Property{
					{Key: "key", Value: Value{Type: ValueTypeBoolean, Value: "true", TypeAnnotation: "type"}},
				},
				Children: []Node{},
			},
		},
	}

	assert.NoError(t, err)
	assert.Equal(t, want, doc)
}

func Test_TODO_PropType(t *testing.T) {
	input := `node key=(type)#true
`

	doc, err := New().Parse(input)

	_ = doc
	_ = err
	t.Skip("TODO")
}

func Test_PropZeroType(t *testing.T) {
	input := `node key=(type)0
`

	doc, err := New().Parse(input)
	want := &Document{
		Nodes: []Node{
			{
				Name:      "node",
				Arguments: []Value{},
				Properties: []Property{
					{Key: "key", Value: Value{Type: ValueTypeNumber, Value: "0", TypeAnnotation: "type"}},
				},
				Children: []Node{},
			},
		},
	}

	assert.NoError(t, err)
	assert.Equal(t, want, doc)
}

func Test_QuestionMarkBeforeNumber(t *testing.T) {
	input := `node ?15
`

	doc, err := New().Parse(input)
	want := &Document{
		Nodes: []Node{
			{
				Name: "node",
				Arguments: []Value{
					{Type: ValueTypeString, Value: "?15"},
				},
				Properties: []Property{},
				Children:   []Node{},
			},
		},
	}

	assert.NoError(t, err)
	assert.Equal(t, want, doc)
}

func Test_TODO_QuoteInBareIdFail(t *testing.T) {
	input := `foo123"bar weeee
`

	doc, err := New().Parse(input)

	_ = doc
	_ = err
	t.Skip("TODO")
}

func Test_TODO_QuotedArgType(t *testing.T) {
	input := `node ("type/")10`

	doc, err := New().Parse(input)

	_ = doc
	_ = err
	t.Skip("TODO - parser does not support quoted type annotations yet")
}

func Test_QuotedNodeName(t *testing.T) {
	input := `"0node"`

	doc, err := New().Parse(input)
	want := &Document{
		Nodes: []Node{
			{
				Name:       "0node",
				Arguments:  []Value{},
				Properties: []Property{},
				Children:   []Node{},
			},
		},
	}

	assert.NoError(t, err)
	assert.Equal(t, want, doc)
}

func Test_TODO_QuotedNodeType(t *testing.T) {
	input := `("type/")node
`

	doc, err := New().Parse(input)

	_ = doc
	_ = err
	t.Skip("TODO - parser does not support quoted type annotations yet")
}

func Test_QuotedNumeric(t *testing.T) {
	input := `node prop="10.0"`

	doc, err := New().Parse(input)
	want := &Document{
		Nodes: []Node{
			{
				Name:      "node",
				Arguments: []Value{},
				Properties: []Property{
					{Key: "prop", Value: Value{Type: ValueTypeString, Value: "10.0"}},
				},
				Children: []Node{},
			},
		},
	}

	assert.NoError(t, err)
	assert.Equal(t, want, doc)
}

func Test_TODO_QuotedPropName(t *testing.T) {
	input := `node "0prop"=val
`

	doc, err := New().Parse(input)

	_ = doc
	_ = err
	t.Skip("TODO - parser does not support quoted property keys yet")
}

func Test_TODO_QuotedPropType(t *testing.T) {
	input := `node key=("type/")#true
`

	doc, err := New().Parse(input)

	_ = doc
	_ = err
	t.Skip("TODO")
}

func Test_RNode(t *testing.T) {
	input := `r "arg"
`

	doc, err := New().Parse(input)
	want := &Document{
		Nodes: []Node{
			{
				Name: "r",
				Arguments: []Value{
					{Type: ValueTypeString, Value: "arg"},
				},
				Properties: []Property{},
				Children:   []Node{},
			},
		},
	}

	assert.NoError(t, err)
	assert.Equal(t, want, doc)
}

func Test_TODO_RawArgType(t *testing.T) {
	input := `node (type)#true
`

	doc, err := New().Parse(input)

	_ = doc
	_ = err
	t.Skip("TODO")
}

func Test_TODO_RawNodeName(t *testing.T) {
	input := `#"\node"#
`

	doc, err := New().Parse(input)

	_ = doc
	_ = err
	t.Skip("TODO")
}

func Test_TODO_RawNodeType(t *testing.T) {
	input := `(type)node`

	doc, err := New().Parse(input)

	_ = doc
	_ = err
	t.Skip("TODO")
}

func Test_TODO_RawPropType(t *testing.T) {
	input := `node key=(type)#true
`

	doc, err := New().Parse(input)

	_ = doc
	_ = err
	t.Skip("TODO")
}

func Test_TODO_RawStringArg(t *testing.T) {
	input := `node_1 #""arg\n"and #stuff"#
node_2 ##"#"arg\n"#and #stuff"##
`

	doc, err := New().Parse(input)

	_ = doc
	_ = err
	t.Skip("TODO")
}

func Test_TODO_RawStringBackslash(t *testing.T) {
	input := `node #"\n"#
`

	doc, err := New().Parse(input)

	_ = doc
	_ = err
	t.Skip("TODO")
}

func Test_TODO_RawStringHashNoEsc(t *testing.T) {
	input := `node #"#"#
`

	doc, err := New().Parse(input)

	_ = doc
	_ = err
	t.Skip("TODO")
}

func Test_TODO_RawStringJustBackslash(t *testing.T) {
	input := `node #"\"#
`

	doc, err := New().Parse(input)

	_ = doc
	_ = err
	t.Skip("TODO")
}

func Test_TODO_RawStringJustQuoteFail(t *testing.T) {
	input := `// This fails because ` + "`" + `"""` + "`" + ` MUST be followed by a newline.
node #"""#
`

	doc, err := New().Parse(input)

	_ = doc
	_ = err
	t.Skip("TODO")
}

func Test_TODO_RawStringMultipleHash(t *testing.T) {
	input := `node ###""#"##"###
`

	doc, err := New().Parse(input)

	_ = doc
	_ = err
	t.Skip("TODO")
}

func Test_TODO_RawStringNewline(t *testing.T) {
	input := `node #"""
hello
world
"""#
`

	doc, err := New().Parse(input)

	_ = doc
	_ = err
	t.Skip("TODO")
}

func Test_TODO_RawStringProp(t *testing.T) {
	input := `node_1 prop=#""arg#"\n"#
node_2 prop=##"#"arg#"#\n"##
`

	doc, err := New().Parse(input)

	_ = doc
	_ = err
	t.Skip("TODO")
}

func Test_TODO_RawStringQuote(t *testing.T) {
	input := `node #"a"b"#
`

	doc, err := New().Parse(input)

	_ = doc
	_ = err
	t.Skip("TODO")
}

func Test_RepeatedArg(t *testing.T) {
	input := `node arg arg
`

	doc, err := New().Parse(input)
	want := &Document{
		Nodes: []Node{
			{
				Name: "node",
				Arguments: []Value{
					{Type: ValueTypeString, Value: "arg"},
					{Type: ValueTypeString, Value: "arg"},
				},
				Properties: []Property{},
				Children:   []Node{},
			},
		},
	}

	assert.NoError(t, err)
	assert.Equal(t, want, doc)
}

func Test_RepeatedProp(t *testing.T) {
	input := `node prop=10 prop=11`

	doc, err := New().Parse(input)
	want := &Document{
		Nodes: []Node{
			{
				Name:      "node",
				Arguments: []Value{},
				Properties: []Property{
					{Key: "prop", Value: Value{Type: ValueTypeNumber, Value: "11"}},
				},
				Children: []Node{},
			},
		},
	}

	assert.NoError(t, err)
	assert.Equal(t, want, doc)
}

func Test_SameNameNodes(t *testing.T) {
	input := `node
node
`

	doc, err := New().Parse(input)
	want := &Document{
		Nodes: []Node{
			{
				Name:       "node",
				Arguments:  []Value{},
				Properties: []Property{},
				Children:   []Node{},
			},
			{
				Name:       "node",
				Arguments:  []Value{},
				Properties: []Property{},
				Children:   []Node{},
			},
		},
	}

	assert.NoError(t, err)
	assert.Equal(t, want, doc)
}

func Test_SciNotationLarge(t *testing.T) {
	input := `node prop=1.23E+1000`

	doc, err := New().Parse(input)
	want := &Document{
		Nodes: []Node{
			{
				Name:      "node",
				Arguments: []Value{},
				Properties: []Property{
					{Key: "prop", Value: Value{Type: ValueTypeNumber, Value: "1.23E+1000"}},
				},
				Children: []Node{},
			},
		},
	}

	assert.NoError(t, err)
	assert.Equal(t, want, doc)
}

func Test_SciNotationSmall(t *testing.T) {
	input := `node prop=1.23E-1000`

	doc, err := New().Parse(input)
	want := &Document{
		Nodes: []Node{
			{
				Name:      "node",
				Arguments: []Value{},
				Properties: []Property{
					{Key: "prop", Value: Value{Type: ValueTypeNumber, Value: "1.23E-1000"}},
				},
				Children: []Node{},
			},
		},
	}

	assert.NoError(t, err)
	assert.Equal(t, want, doc)
}

func Test_TODO_SemicolonAfterChild(t *testing.T) {
	input := `node {
     childnode
};
`

	doc, err := New().Parse(input)

	_ = doc
	_ = err
	t.Skip("TODO")
}

func Test_TODO_SemicolonInChild(t *testing.T) {
	input := `node1 {
      node2;
}`

	doc, err := New().Parse(input)

	_ = doc
	_ = err
	t.Skip("TODO")
}

func Test_TODO_SemicolonMissingAfterChildrenFail(t *testing.T) {
	input := `foo123{bar}foo weeee
`

	doc, err := New().Parse(input)

	_ = doc
	_ = err
	t.Skip("TODO")
}

func Test_SemicolonSeparated(t *testing.T) {
	input := `node1;node2`

	doc, err := New().Parse(input)
	want := &Document{
		Nodes: []Node{
			{
				Name:       "node1",
				Arguments:  []Value{},
				Properties: []Property{},
				Children:   []Node{},
			},
			{
				Name:       "node2",
				Arguments:  []Value{},
				Properties: []Property{},
				Children:   []Node{},
			},
		},
	}

	assert.NoError(t, err)
	assert.Equal(t, want, doc)
}

func Test_SemicolonSeparatedNodes(t *testing.T) {
	input := `node1; node2; `

	doc, err := New().Parse(input)
	want := &Document{
		Nodes: []Node{
			{
				Name:       "node1",
				Arguments:  []Value{},
				Properties: []Property{},
				Children:   []Node{},
			},
			{
				Name:       "node2",
				Arguments:  []Value{},
				Properties: []Property{},
				Children:   []Node{},
			},
		},
	}

	assert.NoError(t, err)
	assert.Equal(t, want, doc)
}

func Test_SemicolonTerminated(t *testing.T) {
	input := `node1;`

	doc, err := New().Parse(input)
	want := &Document{
		Nodes: []Node{
			{
				Name:       "node1",
				Arguments:  []Value{},
				Properties: []Property{},
				Children:   []Node{},
			},
		},
	}

	assert.NoError(t, err)
	assert.Equal(t, want, doc)
}

func Test_SingleArg(t *testing.T) {
	input := `node arg
`

	doc, err := New().Parse(input)
	want := &Document{
		Nodes: []Node{
			{
				Name: "node",
				Arguments: []Value{
					{Type: ValueTypeString, Value: "arg"},
				},
				Properties: []Property{},
				Children:   []Node{},
			},
		},
	}

	assert.NoError(t, err)
	assert.Equal(t, want, doc)
}

func Test_SingleProp(t *testing.T) {
	input := `node prop=val
`

	doc, err := New().Parse(input)
	want := &Document{
		Nodes: []Node{
			{
				Name:      "node",
				Arguments: []Value{},
				Properties: []Property{
					{Key: "prop", Value: Value{Type: ValueTypeString, Value: "val"}},
				},
				Children: []Node{},
			},
		},
	}

	assert.NoError(t, err)
	assert.Equal(t, want, doc)
}

func Test_TODO_SlashInBareIdFail(t *testing.T) {
	input := `foo123/bar weeee
`

	doc, err := New().Parse(input)

	_ = doc
	_ = err
	t.Skip("TODO")
}

func Test_TODO_SlashdashAfterArgTypeFail(t *testing.T) {
	input := `node (ty)/-arg1 arg2
`

	doc, err := New().Parse(input)

	_ = doc
	_ = err
	t.Skip("TODO")
}

func Test_TODO_SlashdashAfterNodeTypeFail(t *testing.T) {
	input := `(ty)/-node
other-node
`

	doc, err := New().Parse(input)

	_ = doc
	_ = err
	t.Skip("TODO")
}

func Test_TODO_SlashdashAfterPropKeyFail(t *testing.T) {
	input := `node key /- = value
`

	doc, err := New().Parse(input)

	_ = doc
	_ = err
	t.Skip("TODO")
}

func Test_TODO_SlashdashAfterPropValTypeFail(t *testing.T) {
	input := `node key=(ty)/-val other-arg
`

	doc, err := New().Parse(input)

	_ = doc
	_ = err
	t.Skip("TODO")
}

func Test_TODO_SlashdashAfterTypeFail(t *testing.T) {
	input := `node (type) /- arg1 arg2
`

	doc, err := New().Parse(input)

	_ = doc
	_ = err
	t.Skip("TODO")
}

func Test_TODO_SlashdashArgAfterNewlineEsc(t *testing.T) {
	input := `node \
    /- arg arg2
`

	doc, err := New().Parse(input)

	_ = doc
	_ = err
	t.Skip("TODO")
}

func Test_TODO_SlashdashArgBeforeNewlineEsc(t *testing.T) {
	input := `node /-    \
    arg
`

	doc, err := New().Parse(input)

	_ = doc
	_ = err
	t.Skip("TODO")
}

func Test_TODO_SlashdashBeforeChildrenEndFail(t *testing.T) {
	input := `node {
    child1
    /-
}
`

	doc, err := New().Parse(input)

	_ = doc
	_ = err
	t.Skip("TODO")
}

func Test_TODO_SlashdashBeforeEofFail(t *testing.T) {
	input := `node foo /-
`

	doc, err := New().Parse(input)

	_ = doc
	_ = err
	t.Skip("TODO")
}

func Test_TODO_SlashdashBeforePropValueFail(t *testing.T) {
	input := `node key = /-val etc
`

	doc, err := New().Parse(input)

	_ = doc
	_ = err
	t.Skip("TODO")
}

func Test_TODO_SlashdashBeforeSemicolonFail(t *testing.T) {
	input := `node foo /-;
`

	doc, err := New().Parse(input)

	_ = doc
	_ = err
	t.Skip("TODO")
}

func Test_TODO_SlashdashBetweenChildBlocksFail(t *testing.T) {
	input := `node { one } /- { two } { three }
`

	doc, err := New().Parse(input)

	_ = doc
	_ = err
	t.Skip("TODO")
}

func Test_TODO_SlashdashChild(t *testing.T) {
	input := `node /- {
    node2
}
`

	doc, err := New().Parse(input)

	_ = doc
	_ = err
	t.Skip("TODO")
}

func Test_TODO_SlashdashChildBlockBeforeEntryErrFail(t *testing.T) {
	input := `node /-{
    child
} foo {
    bar
}
`

	doc, err := New().Parse(input)

	_ = doc
	_ = err
	t.Skip("TODO")
}

func Test_TODO_SlashdashEmptyChild(t *testing.T) {
	input := `node /- {
}
`

	doc, err := New().Parse(input)

	_ = doc
	_ = err
	t.Skip("TODO")
}

func Test_TODO_SlashdashEsclineBeforeArgType(t *testing.T) {
	input := `node /-\
(ty)arg1 arg2
`

	doc, err := New().Parse(input)

	_ = doc
	_ = err
	t.Skip("TODO")
}

func Test_TODO_SlashdashEsclineBeforeChildren(t *testing.T) {
	input := `node arg1 /-\
{
}
`

	doc, err := New().Parse(input)

	_ = doc
	_ = err
	t.Skip("TODO")
}

func Test_TODO_SlashdashEsclineBeforeNode(t *testing.T) {
	input := `/-\
node1
node2
`

	doc, err := New().Parse(input)

	_ = doc
	_ = err
	t.Skip("TODO")
}

func Test_TODO_SlashdashFalseNode(t *testing.T) {
	input := `node foo /-
not-a-node bar
`

	doc, err := New().Parse(input)

	_ = doc
	_ = err
	t.Skip("TODO")
}

func Test_TODO_SlashdashFullNode(t *testing.T) {
	input := `/- node 1.0 "a" b="""
b
"""
`

	doc, err := New().Parse(input)

	_ = doc
	_ = err
	t.Skip("TODO")
}

func Test_TODO_SlashdashInSlashdash(t *testing.T) {
	input := `/- node1 /- 1.0
node2`

	doc, err := New().Parse(input)

	_ = doc
	_ = err
	t.Skip("TODO")
}

func Test_TODO_SlashdashInsideArgTypeFail(t *testing.T) {
	input := `node (/-bad)nope
`

	doc, err := New().Parse(input)

	_ = doc
	_ = err
	t.Skip("TODO")
}

func Test_TODO_SlashdashInsideNodeTypeFail(t *testing.T) {
	input := `(/-ty)node
`

	doc, err := New().Parse(input)

	_ = doc
	_ = err
	t.Skip("TODO")
}

func Test_TODO_SlashdashMultiLineCommentEntry(t *testing.T) {
	input := `node 1 /- /*
multi
line
comment
here
*/ 2 3
`

	doc, err := New().Parse(input)

	_ = doc
	_ = err
	t.Skip("TODO")
}

func Test_TODO_SlashdashMultiLineCommentInline(t *testing.T) {
	input := `node 1 /-/*two*/2 3
`

	doc, err := New().Parse(input)

	_ = doc
	_ = err
	t.Skip("TODO")
}

func Test_TODO_SlashdashMultipleChildBlocks(t *testing.T) {
	input := `node foo /-{
    one
} \
/-{
    two
} {
    three
} /-{
    four
}
`

	doc, err := New().Parse(input)

	_ = doc
	_ = err
	t.Skip("TODO")
}

func Test_TODO_SlashdashNegativeNumber(t *testing.T) {
	input := `node /--1.0 2.0`

	doc, err := New().Parse(input)

	_ = doc
	_ = err
	t.Skip("TODO")
}

func Test_TODO_SlashdashNewlineBeforeChildren(t *testing.T) {
	input := `node 1 2 /-
{
    child
}
`

	doc, err := New().Parse(input)

	_ = doc
	_ = err
	t.Skip("TODO")
}

func Test_TODO_SlashdashNewlineBeforeEntry(t *testing.T) {
	input := `node 1 /-
2 3
`

	doc, err := New().Parse(input)

	_ = doc
	_ = err
	t.Skip("TODO")
}

func Test_TODO_SlashdashNewlineBeforeNode(t *testing.T) {
	input := `/-
node 1 2 3
`

	doc, err := New().Parse(input)

	_ = doc
	_ = err
	t.Skip("TODO")
}

func Test_TODO_SlashdashNodeInChild(t *testing.T) {
	input := `node1 {
    /- node2
}`

	doc, err := New().Parse(input)

	_ = doc
	_ = err
	t.Skip("TODO")
}

func Test_TODO_SlashdashNodeWithChild(t *testing.T) {
	input := `/- node {
   node2
}`

	doc, err := New().Parse(input)

	_ = doc
	_ = err
	t.Skip("TODO")
}

func Test_TODO_SlashdashOnlyNode(t *testing.T) {
	input := `/-node
`

	doc, err := New().Parse(input)

	_ = doc
	_ = err
	t.Skip("TODO")
}

func Test_TODO_SlashdashOnlyNodeWithSpace(t *testing.T) {
	input := `/- node`

	doc, err := New().Parse(input)

	_ = doc
	_ = err
	t.Skip("TODO")
}

func Test_TODO_SlashdashProp(t *testing.T) {
	input := `node /- key=value arg
`

	doc, err := New().Parse(input)

	_ = doc
	_ = err
	t.Skip("TODO")
}

func Test_TODO_SlashdashRawPropKey(t *testing.T) {
	input := `node /- key=value
`

	doc, err := New().Parse(input)

	_ = doc
	_ = err
	t.Skip("TODO")
}

func Test_TODO_SlashdashRepeatedProp(t *testing.T) {
	input := `node arg=correct /- arg=wrong
`

	doc, err := New().Parse(input)

	_ = doc
	_ = err
	t.Skip("TODO")
}

func Test_TODO_SlashdashSingleLineCommentEntry(t *testing.T) {
	input := `node 1 /- // stuff
2 3
`

	doc, err := New().Parse(input)

	_ = doc
	_ = err
	t.Skip("TODO")
}

func Test_TODO_SlashdashSingleLineCommentNode(t *testing.T) {
	input := `/- // this is a comment
node1
node2
`

	doc, err := New().Parse(input)

	_ = doc
	_ = err
	t.Skip("TODO")
}

func Test_TODO_SpaceAfterArgType(t *testing.T) {
	input := `node (type) 10
`

	doc, err := New().Parse(input)

	_ = doc
	_ = err
	t.Skip("TODO")
}

func Test_TODO_SpaceAfterNodeType(t *testing.T) {
	input := `(type) node
`

	doc, err := New().Parse(input)

	_ = doc
	_ = err
	t.Skip("TODO")
}

func Test_TODO_SpaceAfterPropType(t *testing.T) {
	input := `node key=(type) #false
`

	doc, err := New().Parse(input)

	_ = doc
	_ = err
	t.Skip("TODO")
}

func Test_TODO_SpaceAroundPropMarker(t *testing.T) {
	input := `node foo = bar
`

	doc, err := New().Parse(input)

	_ = doc
	_ = err
	t.Skip("TODO")
}

func Test_TODO_SpaceInArgType(t *testing.T) {
	input := `node (type )#false
`

	doc, err := New().Parse(input)

	_ = doc
	_ = err
	t.Skip("TODO")
}

func Test_TODO_SpaceInNodeType(t *testing.T) {
	input := `( type)node
`

	doc, err := New().Parse(input)

	_ = doc
	_ = err
	t.Skip("TODO")
}

func Test_TODO_SpaceInPropType(t *testing.T) {
	input := `node key=(type )#false
`

	doc, err := New().Parse(input)

	_ = doc
	_ = err
	t.Skip("TODO")
}

func Test_TODO_SquareBracketInBareIdFail(t *testing.T) {
	input := `foo123[bar]foo weeee
`

	doc, err := New().Parse(input)

	_ = doc
	_ = err
	t.Skip("TODO")
}

func Test_StringArg(t *testing.T) {
	input := `node "arg"`

	doc, err := New().Parse(input)
	want := &Document{
		Nodes: []Node{
			{
				Name: "node",
				Arguments: []Value{
					{Type: ValueTypeString, Value: "arg"},
				},
				Properties: []Property{},
				Children:   []Node{},
			},
		},
	}

	assert.NoError(t, err)
	assert.Equal(t, want, doc)
}

func Test_TODO_StringEscapedLiteralWhitespace(t *testing.T) {
	input := `node "Hello \
World \          Stuff"
`

	doc, err := New().Parse(input)

	_ = doc
	_ = err
	t.Skip("TODO")
}

func Test_StringProp(t *testing.T) {
	input := `node prop="val"`

	doc, err := New().Parse(input)
	want := &Document{
		Nodes: []Node{
			{
				Name:      "node",
				Arguments: []Value{},
				Properties: []Property{
					{Key: "prop", Value: Value{Type: ValueTypeString, Value: "val"}},
				},
				Children: []Node{},
			},
		},
	}

	assert.NoError(t, err)
	assert.Equal(t, want, doc)
}

func Test_TODO_TabSpace(t *testing.T) {
	input := `node
`

	doc, err := New().Parse(input)

	_ = doc
	_ = err
	t.Skip("TODO")
}

func Test_TODO_TrailingCrlf(t *testing.T) {
	input := `node
`

	doc, err := New().Parse(input)

	_ = doc
	_ = err
	t.Skip("TODO")
}

func Test_TODO_TrailingUnderscoreHex(t *testing.T) {
	input := `node 0x123abc_`

	doc, err := New().Parse(input)

	_ = doc
	_ = err
	t.Skip("TODO")
}

func Test_TODO_TrailingUnderscoreOctal(t *testing.T) {
	input := `node 0o123_
`

	doc, err := New().Parse(input)

	_ = doc
	_ = err
	t.Skip("TODO")
}

func Test_TruePrefixInBareId(t *testing.T) {
	input := `true_id
`

	doc, err := New().Parse(input)
	want := &Document{
		Nodes: []Node{
			{
				Name:       "true_id",
				Arguments:  []Value{},
				Properties: []Property{},
				Children:   []Node{},
			},
		},
	}

	assert.NoError(t, err)
	assert.Equal(t, want, doc)
}

func Test_TruePrefixInPropKey(t *testing.T) {
	input := `node true_id=1
`

	doc, err := New().Parse(input)
	want := &Document{
		Nodes: []Node{
			{
				Name:      "node",
				Arguments: []Value{},
				Properties: []Property{
					{Key: "true_id", Value: Value{Type: ValueTypeNumber, Value: "1"}},
				},
				Children: []Node{},
			},
		},
	}

	assert.NoError(t, err)
	assert.Equal(t, want, doc)
}

func Test_TODO_TruePropKeyFail(t *testing.T) {
	input := `node true=1
`

	doc, err := New().Parse(input)

	_ = doc
	_ = err
	t.Skip("TODO")
}

func Test_TwoNodes(t *testing.T) {
	input := `node1
node2
`

	doc, err := New().Parse(input)
	want := &Document{
		Nodes: []Node{
			{
				Name:       "node1",
				Arguments:  []Value{},
				Properties: []Property{},
				Children:   []Node{},
			},
			{
				Name:       "node2",
				Arguments:  []Value{},
				Properties: []Property{},
				Children:   []Node{},
			},
		},
	}

	assert.NoError(t, err)
	assert.Equal(t, want, doc)
}

func Test_TODO_TypeBeforePropKeyFail(t *testing.T) {
	input := `node (type)key=10
`

	doc, err := New().Parse(input)

	_ = doc
	_ = err
	t.Skip("TODO")
}

func Test_TODO_UnbalancedRawHashesFail(t *testing.T) {
	input := `node ##"foo"#
`

	doc, err := New().Parse(input)

	_ = doc
	_ = err
	t.Skip("TODO")
}

func Test_TODO_UnderscoreAtStartOfFractionFail(t *testing.T) {
	input := `node 1._7`

	doc, err := New().Parse(input)

	_ = doc
	_ = err
	t.Skip("TODO")
}

func Test_TODO_UnderscoreAtStartOfHexFail(t *testing.T) {
	input := `node 0x_10`

	doc, err := New().Parse(input)

	_ = doc
	_ = err
	t.Skip("TODO")
}

func Test_TODO_UnderscoreBeforeNumber(t *testing.T) {
	input := `node _15
`

	doc, err := New().Parse(input)

	_ = doc
	_ = err
	t.Skip("TODO")
}

func Test_TODO_UnderscoreInExponent(t *testing.T) {
	input := `node 1.0e-10_0
`

	doc, err := New().Parse(input)

	_ = doc
	_ = err
	t.Skip("TODO")
}

func Test_TODO_UnderscoreInFloat(t *testing.T) {
	input := `node 1_1.0
`

	doc, err := New().Parse(input)

	_ = doc
	_ = err
	t.Skip("TODO")
}

func Test_TODO_UnderscoreInFraction(t *testing.T) {
	input := `node 1.0_2`

	doc, err := New().Parse(input)

	_ = doc
	_ = err
	t.Skip("TODO")
}

func Test_TODO_UnderscoreInInt(t *testing.T) {
	input := `node 1_0
`

	doc, err := New().Parse(input)

	_ = doc
	_ = err
	t.Skip("TODO")
}

func Test_TODO_UnderscoreInOctal(t *testing.T) {
	input := `node 0o012_3456_7`

	doc, err := New().Parse(input)

	_ = doc
	_ = err
	t.Skip("TODO")
}

func Test_TODO_UnicodeDeleteFail(t *testing.T) {
	input := `// 0x007F (Delete)
node1 arg
`

	doc, err := New().Parse(input)

	_ = doc
	_ = err
	t.Skip("TODO")
}

func Test_TODO_UnicodeEscapedAboveMaxFail(t *testing.T) {
	input := `no "Higher than max Unicode Scalar Value \u{10FFFF} \u{11FFFF}"
`

	doc, err := New().Parse(input)

	_ = doc
	_ = err
	t.Skip("TODO")
}

func Test_TODO_UnicodeEscapedH1Fail(t *testing.T) {
	input := `no "Surrogates high\u{D800}"
`

	doc, err := New().Parse(input)

	_ = doc
	_ = err
	t.Skip("TODO")
}

func Test_TODO_UnicodeEscapedH2Fail(t *testing.T) {
	input := `no "Surrogates high\u{D911}"
`

	doc, err := New().Parse(input)

	_ = doc
	_ = err
	t.Skip("TODO")
}

func Test_TODO_UnicodeEscapedH3Fail(t *testing.T) {
	input := `no "Surrogates high\u{DABB}"
`

	doc, err := New().Parse(input)

	_ = doc
	_ = err
	t.Skip("TODO")
}

func Test_TODO_UnicodeEscapedH4Fail(t *testing.T) {
	input := `no "Surrogates high\u{DBFF}"
`

	doc, err := New().Parse(input)

	_ = doc
	_ = err
	t.Skip("TODO")
}

func Test_TODO_UnicodeEscapedL1Fail(t *testing.T) {
	input := `no "Surrogates low\u{DC00}"
`

	doc, err := New().Parse(input)

	_ = doc
	_ = err
	t.Skip("TODO")
}

func Test_TODO_UnicodeEscapedL2Fail(t *testing.T) {
	input := `no "Surrogates low\u{DEAD}"
`

	doc, err := New().Parse(input)

	_ = doc
	_ = err
	t.Skip("TODO")
}

func Test_TODO_UnicodeEscapedL3Fail(t *testing.T) {
	input := `eno "Surrogates low\u{DFFF}"
`

	doc, err := New().Parse(input)

	_ = doc
	_ = err
	t.Skip("TODO")
}

func Test_TODO_UnicodeEscapedTooLongLead0Fail(t *testing.T) {
	input := `no "Even with leading 0s Unicode Scalar Value escapes must ≤6: \u{0012345}"
`

	doc, err := New().Parse(input)

	_ = doc
	_ = err
	t.Skip("TODO")
}

func Test_TODO_UnicodeFsiFail(t *testing.T) {
	input := `// 0x2068
node1 ⁨arg
`

	doc, err := New().Parse(input)

	_ = doc
	_ = err
	t.Skip("TODO")
}

func Test_TODO_UnicodeLreFail(t *testing.T) {
	input := `// 0x202A
node1 ‪arg
`

	doc, err := New().Parse(input)

	_ = doc
	_ = err
	t.Skip("TODO")
}

func Test_TODO_UnicodeLriFail(t *testing.T) {
	input := `// 0x2066
node1⁦arg
`

	doc, err := New().Parse(input)

	_ = doc
	_ = err
	t.Skip("TODO")
}

func Test_TODO_UnicodeLrmFail(t *testing.T) {
	input := `// 0x200E
node ‎arg
`

	doc, err := New().Parse(input)

	_ = doc
	_ = err
	t.Skip("TODO")
}

func Test_TODO_UnicodeLroFail(t *testing.T) {
	input := `// 0x202D
node ‭arg
`

	doc, err := New().Parse(input)

	_ = doc
	_ = err
	t.Skip("TODO")
}

func Test_TODO_UnicodePdfFail(t *testing.T) {
	input := `// 0x202C
node ‬arg
`

	doc, err := New().Parse(input)

	_ = doc
	_ = err
	t.Skip("TODO")
}

func Test_TODO_UnicodePdiFail(t *testing.T) {
	input := `// 0x2069
node ⁩arg
`

	doc, err := New().Parse(input)

	_ = doc
	_ = err
	t.Skip("TODO")
}

func Test_TODO_UnicodeRleFail(t *testing.T) {
	input := `// 0x202B
node1 ‫arg
`

	doc, err := New().Parse(input)

	_ = doc
	_ = err
	t.Skip("TODO")
}

func Test_TODO_UnicodeRliFail(t *testing.T) {
	input := `// 0x2067
node1 ⁧arg
`

	doc, err := New().Parse(input)

	_ = doc
	_ = err
	t.Skip("TODO")
}

func Test_TODO_UnicodeRlmFail(t *testing.T) {
	input := `// 0x200F
node ‏arg
`

	doc, err := New().Parse(input)

	_ = doc
	_ = err
	t.Skip("TODO")
}

func Test_TODO_UnicodeRloFail(t *testing.T) {
	input := `// 0x202E
node ‮arg
`

	doc, err := New().Parse(input)

	_ = doc
	_ = err
	t.Skip("TODO")
}

func Test_TODO_UnicodeSilly(t *testing.T) {
	input := `ノード　お名前=ฅ^•ﻌ•^ฅ
`

	doc, err := New().Parse(input)

	_ = doc
	_ = err
	t.Skip("TODO")
}

func Test_TODO_UnicodeUnder0x20Fail(t *testing.T) {
	input := `// 0x0019
node1 arg
`

	doc, err := New().Parse(input)

	_ = doc
	_ = err
	t.Skip("TODO")
}

func Test_TODO_UnterminatedEmptyNodeFail(t *testing.T) {
	input := `node {
`

	doc, err := New().Parse(input)

	_ = doc
	_ = err
	t.Skip("TODO")
}

func Test_TODO_UnusualBareIdCharsInQuotedId(t *testing.T) {
	input := `"foo123~!@$%^&*.:'|?+<>,` + "`" + `-_" weeee
`

	doc, err := New().Parse(input)

	_ = doc
	_ = err
	t.Skip("TODO")
}

func Test_TODO_UnusualCharsInBareId(t *testing.T) {
	input := `foo123~!@$%^&*.:'|?+<>,` + "`" + `-_ weeee
`

	doc, err := New().Parse(input)

	_ = doc
	_ = err
	t.Skip("TODO")
}

func Test_TODO_VerticalTabWhitespace(t *testing.T) {
	input := `node argnode2 arg2
`

	doc, err := New().Parse(input)

	_ = doc
	_ = err
	t.Skip("TODO")
}

func Test_TODO_ZeroFloat(t *testing.T) {
	input := `node 0.0
`

	doc, err := New().Parse(input)

	_ = doc
	_ = err
	t.Skip("TODO")
}

func Test_ZeroInt(t *testing.T) {
	input := `node 0
`

	doc, err := New().Parse(input)
	want := &Document{
		Nodes: []Node{
			{
				Name: "node",
				Arguments: []Value{
					{Type: ValueTypeNumber, Value: "0"},
				},
				Properties: []Property{},
				Children:   []Node{},
			},
		},
	}

	assert.NoError(t, err)
	assert.Equal(t, want, doc)
}

func Test_TODO_ZeroSpaceBeforeFirstArgFail(t *testing.T) {
	input := `node"string"
`

	doc, err := New().Parse(input)

	_ = doc
	_ = err
	t.Skip("TODO")
}

func Test_TODO_ZeroSpaceBeforePropFail(t *testing.T) {
	input := `node foo="value"bar=5
`

	doc, err := New().Parse(input)

	_ = doc
	_ = err
	t.Skip("TODO")
}

func Test_TODO_ZeroSpaceBeforeSecondArgFail(t *testing.T) {
	input := `node "string"1
`

	doc, err := New().Parse(input)

	_ = doc
	_ = err
	t.Skip("TODO")
}

func Test_TODO_ZeroSpaceBeforeSlashdashArg(t *testing.T) {
	input := `node "string"/-1
`

	doc, err := New().Parse(input)

	_ = doc
	_ = err
	t.Skip("TODO")
}

func Test_TODO_ZeroSpaceBeforeSlashdashChildren(t *testing.T) {
	input := `node "string"/-{}
node "string" {}/-{}
`

	doc, err := New().Parse(input)

	_ = doc
	_ = err
	t.Skip("TODO")
}

func Test_TODO_ZeroSpaceBeforeSlashdashProp(t *testing.T) {
	input := `node "string"/-foo=1
`

	doc, err := New().Parse(input)

	_ = doc
	_ = err
	t.Skip("TODO")
}
