# Golden Tests

This directory contains golden tests for the KDL parser.

## Structure

Each test consists of:

1. **Test Spec** (`testXXX.json`): Defines the test configuration
2. **Input File** (`testXXX-in.kdl`): The KDL input to parse
3. **Output Files** (multiple formats):
   - JSON: `testXXX-out.json` - The expected parse result in JSON format
   - XML: `testXXX-out.xml` - The expected parse result in XML format

## Test Spec Format

```json
{
  "input": "test001-in.kdl",
  "output": {
    "json": "test001-out.json",
    "xml": "test001-out.xml"
  },
  "title": "short test description",
  "errors": [
    {
      "type": "error-type",
      "line": 1,
      "pos": 5
    }
  ]
}
```

Fields:
- `input`: Path to the input KDL file (relative to tests directory)
- `output`: Map of format to output file path (e.g., `{"json": "test001-out.json", "xml": "test001-out.xml"}`)
- `title`: Short description of what the test verifies
- `errors`: Array of expected errors (empty array or omitted for successful parse tests)

## Output Formats

### JSON Output Format

The output JSON uses a minimalist format to represent the parsed AST:

```json
{
  "nodes": [
    {
      "name": "node-name",
      "args": [
        {"type": "string", "value": "hello"},
        {"type": "number", "value": 42},
        {"type": "boolean", "value": true},
        {"type": "null", "value": null}
      ],
      "props": {
        "key": {"type": "string", "value": "value"}
      },
      "children": [
        {"name": "child-node"}
      ],
      "typeAnnotation": "type-name"
    }
  ]
}
```

Fields are omitted if empty (e.g., no `args` field if the node has no arguments).

### XML Output Format

The output XML represents the parsed AST in an XML structure:

```xml
<?xml version="1.0" encoding="UTF-8"?>
<document>
  <node name="node-name" typeAnnotation="type-name">
    <arg type="string">hello</arg>
    <arg type="number">42</arg>
    <arg type="boolean">true</arg>
    <arg type="null"></arg>
    <prop name="key" type="string">value</prop>
    <node name="child-node"></node>
  </node>
</document>
```

Since KDL is structurally similar to XML, the XML format provides a natural representation of the parsed tree structure.

## Error Format

For tests that expect parse errors:

```json
{
  "type": "error-type-identifier",
  "line": 1,
  "pos": 5
}
```

Fields:
- `type`: Error type identifier (kebab-case)
- `line`: 1-based line number where error occurred
- `pos`: 1-based column number where error occurred

## Running Tests

Run all golden tests:
```bash
go test -v -run TestGoldenTests
```

Run a specific test:
```bash
go test -v -run "TestGoldenTests/simple_node"
```

## Adding New Tests

1. Create a new test spec file: `testXXX.json` (increment number)
2. Create the input file: `testXXX-in.kdl`
3. Create the expected output files:
   - `testXXX-out.json` (JSON format)
   - `testXXX-out.xml` (XML format)
4. Update the test spec to reference both output files
5. Run tests to verify

The test runner will automatically discover and run all test spec files, testing against all specified output formats.
