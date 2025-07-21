"""
Test CFG generation for simple functions
"""


run_empty_file = ["--no-ast", "test/empty.go"]

out_empty_file = """
Declarations:
Function definitions:
Entry: 0
Exit: 1
Graph:
  0 => 1: nil
"""

run_package = ["--no-ast", "test/package.go"]

out_package = """
Declarations:
Function definitions:
Entry: 0
Exit: 1
Graph:
  0 => 2: (package_clause (package_identifier))
  2 => 1: nil
"""

