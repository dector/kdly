# KDLy

A KDL v2 parser implementation written in Go.

## About

This project implements a parser for [KDL (KDL Document Language) v2](https://kdl.dev), a human-friendly document language with a focus on readability and ease of use.

## Implementation

The parser uses a manually written state machine for tokenization, providing fine-grained control over the lexical analysis phase.

## Technology

- Language: Go
- Tokenizer: Hand-written state machine
- Target Specification: KDL v2
- To run tests - use `task test`