# Golden Tests

This directory contains golden tests for the KDL parser.

## Structure

Each test consists of three files:

1. **Test Spec** (`testXXX.json`): Defines the test configuration
2. **Input File** (`testXXX-input.kdl`): The KDL input to parse
3. **Output File** (`testXXX-output.json`): The expected parse result

## Test Spec Format

```json
{
  "input": "test001-input.kdl",
  "output": "test001-output.json",
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
- `output`: Path to the expected output JSON file (relative to tests directory)
- `title`: Short description of what the test verifies
- `errors`: Array of expected errors (empty array or omitted for successful parse tests)

## Output Format

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
2. Create the input file: `testXXX-input.kdl`
3. Create the expected output: `testXXX-output.json`
4. Run tests to verify

The test runner will automatically discover and run all test spec files.
