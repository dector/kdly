package parser

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

// ============================================================
// Node Type Annotation Tests
// Tests for type annotations on node names: (type)nodename
// ============================================================

func TestNodeTypeAnnotation_Basic(t *testing.T) {
	input := `(contributor)person name="Foo McBar"`

	doc, err := New().Parse(input)

	want := &Document{
		Nodes: []Node{
			{
				Name:           "person",
				TypeAnnotation: "contributor",
				Arguments:      []Value{},
				Properties: []Property{
					{
						Key: "name",
						Value: Value{
							Type:  ValueTypeString,
							Value: "Foo McBar",
						},
					},
				},
				Children: []Node{},
			},
		},
	}

	assert.NoError(t, err)
	assert.Equal(t, want, doc)
}

func TestNodeTypeAnnotation_WithoutTypeAnnotation(t *testing.T) {
	input := `person name="Foo McBar"`

	doc, err := New().Parse(input)

	want := &Document{
		Nodes: []Node{
			{
				Name:           "person",
				TypeAnnotation: "",
				Arguments:      []Value{},
				Properties: []Property{
					{
						Key: "name",
						Value: Value{
							Type:  ValueTypeString,
							Value: "Foo McBar",
						},
					},
				},
				Children: []Node{},
			},
		},
	}

	assert.NoError(t, err)
	assert.Equal(t, want, doc)
}

func TestNodeTypeAnnotation_QuotedNodeName(t *testing.T) {
	input := `(author)"my node" key="value"`

	doc, err := New().Parse(input)

	want := &Document{
		Nodes: []Node{
			{
				Name:           "my node",
				TypeAnnotation: "author",
				Arguments:      []Value{},
				Properties: []Property{
					{
						Key: "key",
						Value: Value{
							Type:  ValueTypeString,
							Value: "value",
						},
					},
				},
				Children: []Node{},
			},
		},
	}

	assert.NoError(t, err)
	assert.Equal(t, want, doc)
}

func TestNodeTypeAnnotation_WithArguments(t *testing.T) {
	input := `(person)author "John Doe" age=30`

	doc, err := New().Parse(input)

	want := &Document{
		Nodes: []Node{
			{
				Name:           "author",
				TypeAnnotation: "person",
				Arguments: []Value{
					{
						Type:  ValueTypeString,
						Value: "John Doe",
					},
				},
				Properties: []Property{
					{
						Key: "age",
						Value: Value{
							Type:  ValueTypeNumber,
							Value: "30",
						},
					},
				},
				Children: []Node{},
			},
		},
	}

	assert.NoError(t, err)
	assert.Equal(t, want, doc)
}

func TestNodeTypeAnnotation_WithTypedArguments(t *testing.T) {
	input := `(config)server (url)"https://example.com" (port)8080`

	doc, err := New().Parse(input)

	want := &Document{
		Nodes: []Node{
			{
				Name:           "server",
				TypeAnnotation: "config",
				Arguments: []Value{
					{
						Type:           ValueTypeString,
						Value:          "https://example.com",
						TypeAnnotation: "url",
					},
					{
						Type:           ValueTypeNumber,
						Value:          "8080",
						TypeAnnotation: "port",
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

func TestNodeTypeAnnotation_WithChildren(t *testing.T) {
	input := `(root)parent {
  (child-type)child1 x=1
  (child-type)child2 x=2
}`

	doc, err := New().Parse(input)

	want := &Document{
		Nodes: []Node{
			{
				Name:           "parent",
				TypeAnnotation: "root",
				Arguments:      []Value{},
				Properties:     []Property{},
				Children: []Node{
					{
						Name:           "child1",
						TypeAnnotation: "child-type",
						Arguments:      []Value{},
						Properties: []Property{
							{
								Key: "x",
								Value: Value{
									Type:  ValueTypeNumber,
									Value: "1",
								},
							},
						},
						Children: []Node{},
					},
					{
						Name:           "child2",
						TypeAnnotation: "child-type",
						Arguments:      []Value{},
						Properties: []Property{
							{
								Key: "x",
								Value: Value{
									Type:  ValueTypeNumber,
									Value: "2",
								},
							},
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

func TestNodeTypeAnnotation_MultipleNodes(t *testing.T) {
	input := `(author)person name="Alice"
(contributor)person name="Bob"
(reviewer)person name="Charlie"`

	doc, err := New().Parse(input)

	want := &Document{
		Nodes: []Node{
			{
				Name:           "person",
				TypeAnnotation: "author",
				Arguments:      []Value{},
				Properties: []Property{
					{
						Key: "name",
						Value: Value{
							Type:  ValueTypeString,
							Value: "Alice",
						},
					},
				},
				Children: []Node{},
			},
			{
				Name:           "person",
				TypeAnnotation: "contributor",
				Arguments:      []Value{},
				Properties: []Property{
					{
						Key: "name",
						Value: Value{
							Type:  ValueTypeString,
							Value: "Bob",
						},
					},
				},
				Children: []Node{},
			},
			{
				Name:           "person",
				TypeAnnotation: "reviewer",
				Arguments:      []Value{},
				Properties: []Property{
					{
						Key: "name",
						Value: Value{
							Type:  ValueTypeString,
							Value: "Charlie",
						},
					},
				},
				Children: []Node{},
			},
		},
	}

	assert.NoError(t, err)
	assert.Equal(t, want, doc)
}

func TestNodeTypeAnnotation_WithSlashdash(t *testing.T) {
	input := `/- (commented)node1 x=1
(author)node2 x=2`

	doc, err := New().Parse(input)

	want := &Document{
		Nodes: []Node{
			{
				Name:           "node2",
				TypeAnnotation: "author",
				Arguments:      []Value{},
				Properties: []Property{
					{
						Key: "x",
						Value: Value{
							Type:  ValueTypeNumber,
							Value: "2",
						},
					},
				},
				Children: []Node{},
			},
		},
	}

	assert.NoError(t, err)
	assert.Equal(t, want, doc)
}

func TestNodeTypeAnnotation_MixedAnnotatedAndNot(t *testing.T) {
	input := `(typed)node1 x=1
node2 x=2
(typed)node3 x=3`

	doc, err := New().Parse(input)

	want := &Document{
		Nodes: []Node{
			{
				Name:           "node1",
				TypeAnnotation: "typed",
				Arguments:      []Value{},
				Properties: []Property{
					{
						Key: "x",
						Value: Value{
							Type:  ValueTypeNumber,
							Value: "1",
						},
					},
				},
				Children: []Node{},
			},
			{
				Name:           "node2",
				TypeAnnotation: "",
				Arguments:      []Value{},
				Properties: []Property{
					{
						Key: "x",
						Value: Value{
							Type:  ValueTypeNumber,
							Value: "2",
						},
					},
				},
				Children: []Node{},
			},
			{
				Name:           "node3",
				TypeAnnotation: "typed",
				Arguments:      []Value{},
				Properties: []Property{
					{
						Key: "x",
						Value: Value{
							Type:  ValueTypeNumber,
							Value: "3",
						},
					},
				},
				Children: []Node{},
			},
		},
	}

	assert.NoError(t, err)
	assert.Equal(t, want, doc)
}

func TestNodeTypeAnnotation_NestedChildren(t *testing.T) {
	input := `(database)config {
  (connection)pool {
    (setting)max-size 100
  }
}`

	doc, err := New().Parse(input)

	want := &Document{
		Nodes: []Node{
			{
				Name:           "config",
				TypeAnnotation: "database",
				Arguments:      []Value{},
				Properties:     []Property{},
				Children: []Node{
					{
						Name:           "pool",
						TypeAnnotation: "connection",
						Arguments:      []Value{},
						Properties:     []Property{},
						Children: []Node{
							{
								Name:           "max-size",
								TypeAnnotation: "setting",
								Arguments: []Value{
									{
										Type:  ValueTypeNumber,
										Value: "100",
									},
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

// ============================================================
// Error Cases
// ============================================================

func TestNodeTypeAnnotation_UnterminatedTypeAnnotation(t *testing.T) {
	input := `(author node name="test"`

	_, err := New().Parse(input)
	assert.Error(t, err)
	assert.Contains(t, err.Error(), "unterminated type annotation")
}

func TestNodeTypeAnnotation_EmptyTypeAnnotation(t *testing.T) {
	input := `()node name="test"`

	_, err := New().Parse(input)
	assert.Error(t, err)
	// Empty identifier should cause an error
}

// ============================================================
// Real-World Examples
// ============================================================

func TestNodeTypeAnnotation_RealWorld_PackageManifest(t *testing.T) {
	input := `(package)manifest {
  (metadata)name "my-app"
  (metadata)version "1.0.0"
  (dependency)requires {
    (library)lodash version="^4.17.0"
    (library)react version="^18.0.0"
  }
}`

	doc, err := New().Parse(input)

	want := &Document{
		Nodes: []Node{
			{
				Name:           "manifest",
				TypeAnnotation: "package",
				Arguments:      []Value{},
				Properties:     []Property{},
				Children: []Node{
					{
						Name:           "name",
						TypeAnnotation: "metadata",
						Arguments: []Value{
							{
								Type:  ValueTypeString,
								Value: "my-app",
							},
						},
						Properties: []Property{},
						Children:   []Node{},
					},
					{
						Name:           "version",
						TypeAnnotation: "metadata",
						Arguments: []Value{
							{
								Type:  ValueTypeString,
								Value: "1.0.0",
							},
						},
						Properties: []Property{},
						Children:   []Node{},
					},
					{
						Name:           "requires",
						TypeAnnotation: "dependency",
						Arguments:      []Value{},
						Properties:     []Property{},
						Children: []Node{
							{
								Name:           "lodash",
								TypeAnnotation: "library",
								Arguments:      []Value{},
								Properties: []Property{
									{
										Key: "version",
										Value: Value{
											Type:  ValueTypeString,
											Value: "^4.17.0",
										},
									},
								},
								Children: []Node{},
							},
							{
								Name:           "react",
								TypeAnnotation: "library",
								Arguments:      []Value{},
								Properties: []Property{
									{
										Key: "version",
										Value: Value{
											Type:  ValueTypeString,
											Value: "^18.0.0",
										},
									},
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

func TestNodeTypeAnnotation_RealWorld_DataSchema(t *testing.T) {
	input := `(entity)User {
  (field)id type="uuid" required=#true
  (field)name type="string" max-length=255
  (field)email type="email" unique=#true
  (relation)posts type="one-to-many" target="Post"
}`

	doc, err := New().Parse(input)

	want := &Document{
		Nodes: []Node{
			{
				Name:           "User",
				TypeAnnotation: "entity",
				Arguments:      []Value{},
				Properties:     []Property{},
				Children: []Node{
					{
						Name:           "id",
						TypeAnnotation: "field",
						Arguments:      []Value{},
						Properties: []Property{
							{
								Key: "type",
								Value: Value{
									Type:  ValueTypeString,
									Value: "uuid",
								},
							},
							{
								Key: "required",
								Value: Value{
									Type:  ValueTypeBoolean,
									Value: "true",
								},
							},
						},
						Children: []Node{},
					},
					{
						Name:           "name",
						TypeAnnotation: "field",
						Arguments:      []Value{},
						Properties: []Property{
							{
								Key: "type",
								Value: Value{
									Type:  ValueTypeString,
									Value: "string",
								},
							},
							{
								Key: "max-length",
								Value: Value{
									Type:  ValueTypeNumber,
									Value: "255",
								},
							},
						},
						Children: []Node{},
					},
					{
						Name:           "email",
						TypeAnnotation: "field",
						Arguments:      []Value{},
						Properties: []Property{
							{
								Key: "type",
								Value: Value{
									Type:  ValueTypeString,
									Value: "email",
								},
							},
							{
								Key: "unique",
								Value: Value{
									Type:  ValueTypeBoolean,
									Value: "true",
								},
							},
						},
						Children: []Node{},
					},
					{
						Name:           "posts",
						TypeAnnotation: "relation",
						Arguments:      []Value{},
						Properties: []Property{
							{
								Key: "type",
								Value: Value{
									Type:  ValueTypeString,
									Value: "one-to-many",
								},
							},
							{
								Key: "target",
								Value: Value{
									Type:  ValueTypeString,
									Value: "Post",
								},
							},
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
