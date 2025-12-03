# Test Failure Analysis and Implementation Plan

This document analyzes all failing tests and determines whether the test needs fixing or the parser implementation is missing functionality.

## Summary Statistics
- **Total Failing Tests**: 89
- **Tests Needing Implementation**: Most tests (multiline strings, escapes, numbers, etc.)
- **Tests That May Need Fixing**: TBD after implementation

---

## Multiline Raw String Tests (9 tests)

### Test_TODO_MultilineRawString
**Status**: ❌ FAILING
**Location**: pkg/parser/parser_todo_test.go:2025
**Issue**: Parser includes `""` delimiters in output instead of just content
**Expected**: `"hey\neveryone\nhow goes?\n"`
**Actual**: `""\nhey\neveryone\nhow goes?\n""`
**Action**: **IMPLEMENTATION MISSING** - Need to implement multiline raw string parsing (`#"""..."""#` syntax)

### Test_TODO_MultilineRawStringContainingQuotes
**Status**: ❌ FAILING
**Location**: pkg/parser/parser_todo_test.go:2051
**Issue**: Parser fails with error on line 4, doesn't handle `##"""` prefix correctly
**Expected**: String containing quotes with `##` prefix for disambiguation
**Actual**: Parse error
**Action**: **IMPLEMENTATION MISSING** - Need to support multiple `#` prefixes for raw strings

### Test_TODO_MultilineRawStringEmpty
**Status**: ❌ FAILING
**Location**: pkg/parser/parser_todo_test.go:2077
**Issue**: Empty multiline raw string includes delimiters
**Expected**: Empty string `""`
**Actual**: `""\n""`
**Action**: **IMPLEMENTATION MISSING** - Same as above

### Test_TODO_MultilineRawStringEmptyIndented
**Status**: ❌ FAILING
**Location**: pkg/parser/parser_todo_test.go:2099
**Issue**: Empty indented multiline raw string includes delimiters and tab
**Expected**: Empty string `""`
**Actual**: `""\n\t""`
**Action**: **IMPLEMENTATION MISSING** - Same as above

### Test_TODO_MultilineRawStringIndented
**Status**: ❌ FAILING
**Location**: pkg/parser/parser_todo_test.go:2121
**Issue**: Doesn't strip common indentation properly
**Expected**: `"  hey\n everyone\n   how goes?\n"` (2 spaces dedented)
**Actual**: `""\n    hey\n   everyone\n     how goes?\n  ""`
**Action**: **IMPLEMENTATION MISSING** - Need to implement indentation stripping based on closing delimiter position

### Test_TODO_MultilineRawStringNonMatchingPrefixCharacterErrorFail
**Status**: ❌ FAILING
**Location**: pkg/parser/parser_todo_test.go:2147
**Issue**: Should error on mixed tabs/spaces but doesn't
**Expected**: Error
**Actual**: No error
**Action**: **IMPLEMENTATION MISSING** - Need to validate consistent indentation character type

### Test_TODO_MultilineRawStringNonMatchingPrefixCountErrorFail
**Status**: ❌ FAILING
**Location**: pkg/parser/parser_todo_test.go:2161
**Issue**: Should error on inconsistent indentation levels but doesn't
**Expected**: Error
**Actual**: No error
**Action**: **IMPLEMENTATION MISSING** - Need to validate consistent indentation levels

### Test_TODO_MultilineRawStringSingleLineErrFail
**Status**: ❌ FAILING
**Location**: pkg/parser/parser_todo_test.go:2175
**Issue**: Should error when multiline raw string is on single line
**Expected**: Error
**Actual**: No error
**Action**: **IMPLEMENTATION MISSING** - Need to enforce multiline requirement for `"""` syntax

### Test_TODO_MultilineRawStringSingleQuoteErrFail
**Status**: ❌ FAILING
**Location**: pkg/parser/parser_todo_test.go:2184
**Issue**: Should error when using single quote instead of triple quotes
**Expected**: Error
**Actual**: No error
**Action**: **IMPLEMENTATION MISSING** - Need to reject `#"...\n..."#` syntax (must be triple quotes)

---

## Multiline String Tests (7 tests)

### Test_TODO_MultilineString
**Status**: ❌ FAILING (likely)
**Location**: pkg/parser/parser_todo_test.go:2147
**Action**: **IMPLEMENTATION MISSING** - Need to implement multiline string parsing (`"""..."""` syntax without `#`)

### Test_TODO_MultilineStringContainingQuotes
**Action**: **IMPLEMENTATION MISSING** - Same as above

### Test_TODO_MultilineStringEmpty
**Action**: **IMPLEMENTATION MISSING** - Same as above

### Test_TODO_MultilineStringEmptyIndented
**Action**: **IMPLEMENTATION MISSING** - Same as above

### Test_TODO_MultilineStringIndented
**Action**: **IMPLEMENTATION MISSING** - Same as above, with indentation handling

### Test_TODO_MultilineStringSingleLineErrFail
**Action**: **IMPLEMENTATION MISSING** - Need to enforce multiline requirement

### Test_TODO_MultilineStringSingleQuoteErrFail
**Action**: **IMPLEMENTATION MISSING** - Need to reject single quote syntax

---

## Escape Sequence Tests (10+ tests)

### Test_TODO_AllEscapes
**Status**: ❌ FAILING (likely)
**Location**: pkg/parser/parser_todo_test.go:8
**Action**: **IMPLEMENTATION MISSING** - Need to implement all escape sequences: `\"`, `\\`, `\b`, `\f`, `\n`, `\r`, `\t`, `\s`

### Test_TODO_EscapeBackslash through Test_TODO_EscapeUnicode6
**Action**: **IMPLEMENTATION MISSING** - Individual escape sequence implementations

---

## Number Parsing Tests (30+ tests)

### Test_Binary
**Status**: ❌ FAILING
**Location**: pkg/parser/parser_todo_test.go:445
**Issue**: Binary number not being converted
**Expected**: `"2"` (decimal value)
**Actual**: `"0b10"` (raw string)
**Action**: **IMPLEMENTATION MISSING** - Need to parse and convert binary numbers (`0b` prefix)

### Test_BinaryTrailingUnderscore
**Action**: **IMPLEMENTATION MISSING** - Need to handle underscores in binary numbers

### Test_Hexadecimal, Test_HexadecimalTrailingUnderscore
**Action**: **IMPLEMENTATION MISSING** - Need to parse and convert hex numbers (`0x` prefix)

### Test_Octal, Test_OctalTrailingUnderscore
**Action**: **IMPLEMENTATION MISSING** - Need to parse and convert octal numbers (`0o` prefix)

### Test_SignedFloat, Test_SignedFloatExponent, etc.
**Action**: **IMPLEMENTATION MISSING** - Need to handle signed floats with exponents

---

## Bare Identifier Tests (6+ tests)

### Test_BareIdentDot
**Status**: ❌ FAILING
**Location**: pkg/parser/parser_todo_test.go:355
**Issue**: Bare identifier `.` is not being parsed as an argument
**Expected**: Argument with value `"."`
**Actual**: No arguments
**Action**: **IMPLEMENTATION MISSING** - Need to parse `.` as a valid bare identifier

### Test_BareIdentNumericDotFail
**Status**: ❌ FAILING
**Issue**: Should error on `.0n` but doesn't
**Action**: **IMPLEMENTATION MISSING** - Need to validate bare identifiers cannot start with numeric patterns

### Test_BareIdentNumericFail
**Status**: ❌ FAILING
**Issue**: Should error on `0n` but doesn't
**Action**: **IMPLEMENTATION MISSING** - Same validation needed

### Test_BareIdentNumericSignFail
**Status**: ❌ FAILING
**Issue**: Should error on `+0n` but doesn't
**Action**: **IMPLEMENTATION MISSING** - Same validation needed

### Test_BareIdentSignDot
**Status**: ❌ FAILING
**Location**: pkg/parser/parser_todo_test.go:424
**Issue**: `+.` being parsed as number instead of string
**Expected**: String `"+."`
**Actual**: Number `"+"`
**Action**: **IMPLEMENTATION MISSING** - Need to parse `+.` as bare identifier, not number

---

## Raw String Tests (5+ tests)

### Test_TODO_RawStringBackslash
**Action**: **IMPLEMENTATION MISSING** - Raw strings should not process escape sequences

### Test_TODO_RawStringMultiHash
**Action**: **IMPLEMENTATION MISSING** - Need to support multiple `#` prefixes (`##"..."##`)

### Test_TODO_RawStringNewLine, Test_TODO_RawStringQuote, etc.
**Action**: **IMPLEMENTATION MISSING** - Raw string edge cases

---

## Type Annotation Tests (10+ tests)

### Test_TODO_ArgBareNumericWsDashType
**Action**: **IMPLEMENTATION MISSING** - Type annotations with whitespace and dashes

### Test_TODO_NodeBoolDashType through Test_TODO_PropRawStringType
**Action**: **IMPLEMENTATION MISSING** - Various type annotation combinations

---

## Implementation Priority

### Phase 1: Critical Foundation (Required for most tests)
1. **Multiline Raw Strings** (`#"""..."""#`) - 9 tests
   - Parse triple-quote delimiters
   - Handle multiple `#` prefixes
   - Implement indentation stripping
   - Validate consistent indentation
   - Enforce multiline requirement

2. **Multiline Strings** (`"""..."""`) - 7 tests
   - Same as raw strings but with escape processing

3. **Escape Sequences** - 10+ tests
   - `\"`, `\\`, `\b`, `\f`, `\n`, `\r`, `\t`, `\s`
   - Unicode escapes: `\u{...}`

### Phase 2: Number System (30+ tests)
1. **Number Base Conversion**
   - Binary: `0b` prefix
   - Octal: `0o` prefix
   - Hexadecimal: `0x` prefix
   - Handle underscores in numbers

2. **Float Handling**
   - Signed floats
   - Exponent notation
   - Edge cases (infinity, NaN)

### Phase 3: Identifier Validation (6+ tests)
1. **Bare Identifier Rules**
   - Allow `.` as valid identifier
   - Allow `+.`, `-.` as identifiers
   - Reject numeric-looking patterns (`0n`, `.0n`, `+0n`)

### Phase 4: Raw Strings (5+ tests)
1. **Single-line Raw Strings** (`#"..."#`)
   - No escape processing
   - Multiple `#` support
   - Special character handling

### Phase 5: Type Annotations (10+ tests)
1. **Type Annotation Edge Cases**
   - Whitespace handling
   - Special characters in types

---

## Test Validation Notes

After implementing each phase, we should:
1. Run the tests for that phase
2. Check if any test expectations are incorrect
3. Verify against KDL v2 specification
4. Update tests if specification differs from expectations

---

## Next Steps

1. **Immediate**: Review KDL v2 specification for multiline string syntax
2. **Start**: Implement multiline raw string tokenization in the lexer
3. **Then**: Implement multiline string parsing in the parser
4. **Finally**: Work through each phase systematically
