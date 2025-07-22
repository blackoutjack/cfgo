# Golze

Golze is a dataflow engine for Go (currently in a CFG-only proof-of-concept
phase), written in Go with the
[go-tree-sitter](https://github.com/tree-sitter/go-tree-sitter)
interface to the parse tree provided by the
[tree-sitter-go](https://github.com/tree-sitter/tree-sitter-go)
grammar.

## Getting started

```
# Clone/build the repo
git clone https://github.com/blackoutjack/golze.git
cd golze && go build

# Install the test framework
python -m venv ./.venv/golze
. .venv/golze/bin/activate
pip install -r requirements.txt

# Run the tests
python -m test
```

## Current support

This is a work in progress. Only the following node kinds are supported.
- package declarations
- import declarations (all formats)
- var declarations
- const declarations
- function declarations

Stay tuned for more.
