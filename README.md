# CFGo

CFGo is a Control Flow Graph for Go files, currently in the proof-of-concept
phase.  The CFG is generated directly from the parse tree provided by
[tree-sitter-go](https://github.com/tree-sitter/tree-sitter-go) rather than
creating an intermediate AST first.

## Getting started

```
# Clone/build the repo
git clone https://github.com/blackoutjack/cfgo.git
cd cfgo && go build

# Install the test framework
python -m venv ./.venv/cfgo
. .venv/cfgo/bin/activate
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
