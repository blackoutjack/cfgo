"""
Test AST generation for basic statements and expressions
"""


run_empty_file = ["--no-cfg", "test/empty.go"]

out_empty_file = """
(source_file)
"""

run_package = ["--no-cfg", "test/package.go"]

out_package = """
(source_file (package_clause (package_identifier)))
"""

