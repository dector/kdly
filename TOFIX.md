# Failing Tests

## Critical Issue
- **PANIC**: `raw_string_just_quote_fail` - slice bounds out of range [4:1] in parser.go:273

## Failed Tests (115 total)

### Identifier/Bareword Issues
- [ ] bare_ident_sign
- [ ] bare_ident_sign/kdl
- [ ] bare_ident_sign_dot
- [ ] bare_ident_sign_dot/kdl
- [ ] chevrons_in_bare_id
- [ ] chevrons_in_bare_id/kdl
- [ ] comma_in_bare_id
- [ ] comma_in_bare_id/kdl
- [ ] braces_in_bare_id
- [ ] hash_in_id_fail

### Type Annotation Issues
- [ ] blank_arg_type
- [ ] blank_node_type
- [ ] blank_prop_type
- [ ] comment_after_arg_type
- [ ] comment_after_node_type
- [ ] comment_after_prop_type
- [ ] comment_in_arg_type
- [ ] comment_in_node_type
- [ ] comment_in_prop_type
- [ ] empty_arg_type_fail
- [ ] empty_node_type_fail
- [ ] empty_prop_type_fail
- [ ] just_space_in_arg_type_fail
- [ ] just_space_in_node_type_fail
- [ ] just_space_in_prop_type_fail
- [ ] just_type_no_arg_fail
- [ ] just_type_no_node_id_fail
- [ ] just_type_no_prop_fail
- [ ] quoted_arg_type
- [ ] quoted_node_type
- [ ] quoted_prop_type
- [ ] prop_float_type
- [ ] prop_float_type/kdl

### String/Escaping Issues
- [ ] bom_initial
- [ ] bom_initial/kdl
- [ ] escaped_whitespace
- [ ] escaped_whitespace/kdl
- [ ] esc_unicode_in_string
- [ ] esc_unicode_in_string/kdl
- [ ] esc_multiple_newlines
- [ ] esc_multiple_newlines/kdl
- [ ] multiline_raw_string_containing_quotes
- [ ] multiline_raw_string_indented
- [ ] multiline_raw_string_indented/kdl
- [ ] multiline_string_double_backslash
- [ ] multiline_string_double_backslash/kdl
- [ ] multiline_string_escape_delimiter
- [ ] multiline_string_escape_in_closing_line
- [ ] multiline_string_escape_in_closing_line/kdl
- [ ] multiline_string_escape_in_closing_line_shallow
- [ ] multiline_string_escape_in_closing_line_shallow/kdl
- [ ] multiline_string_escape_newline_at_end
- [ ] multiline_string_escape_newline_at_end/kdl
- [ ] multiline_string_indented
- [ ] multiline_string_indented/kdl
- [ ] multiline_string_whitespace_only
- [ ] multiline_string_wrapped_binary
- [ ] multiline_string_wrapped_binary/kdl
- [ ] raw_node_name
- [ ] raw_string_arg
- [ ] raw_string_just_quote_fail

### Escape Line (\\) Issues
- [ ] escline_after_semicolon
- [ ] escline_alone
- [ ] escline_empty_line
- [ ] escline_end_of_node
- [ ] escline_end_of_node/kdl
- [ ] escline_in_child_block
- [ ] escline_line_comment
- [ ] escline_line_comment/kdl
- [ ] escline_node
- [ ] escline_node_type
- [ ] escline_slashdash

### Number Format Issues
- [ ] hex
- [ ] hex_int
- [ ] illegal_char_in_binary_fail
- [ ] multiple_x_in_hex_fail
- [ ] no_digits_in_hex_fail
- [ ] negative_exponent
- [ ] negative_exponent/kdl
- [ ] negative_float
- [ ] negative_float/kdl
- [ ] no_decimal_exponent
- [ ] no_decimal_exponent/kdl
- [ ] positive_exponent
- [ ] positive_exponent/kdl
- [ ] floating_point_keywords

### Comment/Formatting Issues
- [ ] comment_and_newline
- [ ] comment_and_newline/kdl
- [ ] commented_child
- [ ] commented_child/kdl
- [ ] crlf_between_nodes
- [ ] crlf_between_nodes/kdl
- [ ] newline_between_nodes
- [ ] newline_between_nodes/kdl
- [ ] multiline_nodes
- [ ] multiline_nodes/kdl

### Property/Value Issues
- [ ] empty_quoted_node_id
- [ ] empty_quoted_prop_key
- [ ] numeric_prop
- [ ] numeric_prop/kdl
- [ ] quoted_prop_name
- [ ] optional_child_semicolon

### Data Preservation Issues
- [ ] preserve_duplicate_nodes
- [ ] preserve_duplicate_nodes/kdl
- [ ] preserve_node_order
- [ ] preserve_node_order/kdl
- [ ] parse_all_arg_types
- [ ] parse_all_arg_types/kdl
