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

run_import_basic = ["--no-ast", "test/import-basic.go"]

out_import_basic = """
Declarations: package
Function definitions:
Entry: 0
Exit: 1
Graph:
  0 => 2: (package_clause (package_identifier))
  2 => 3: (import_spec path: (interpreted_string_literal (interpreted_string_literal_content)))
  3 => 1: nil
"""

run_import_splat = ["--no-ast", "test/import-splat.go"]

out_import_splat = """
Declarations: package.*
Function definitions:
Entry: 0
Exit: 1
Graph:
  0 => 2: (package_clause (package_identifier))
  2 => 3: (import_spec name: (dot) path: (interpreted_string_literal (interpreted_string_literal_content)))
  3 => 1: nil
"""

run_import_named = ["--no-ast", "test/import-named.go"]

out_import_named = """
Declarations: myname
Function definitions:
Entry: 0
Exit: 1
Graph:
  0 => 2: (package_clause (package_identifier))
  2 => 3: (import_spec name: (package_identifier) path: (interpreted_string_literal (interpreted_string_literal_content)))
  3 => 1: nil
"""

run_import_block = ["--no-ast", "test/import-block.go"]

out_import_block = """
Declarations: myname pkg
Function definitions:
Entry: 0
Exit: 1
Graph:
  0 => 2: (package_clause (package_identifier))
  2 => 4: (import_spec name: (package_identifier) path: (interpreted_string_literal (interpreted_string_literal_content)))
  3 => 1: nil
  4 => 3: (import_spec path: (interpreted_string_literal (interpreted_string_literal_content)))
"""

run_import_blank = ["--no-ast", "test/import-blank.go"]

out_import_blank = """
Declarations:
Function definitions:
Entry: 0
Exit: 1
Graph:
  0 => 2: (package_clause (package_identifier))
  2 => 3: (import_spec name: (blank_identifier) path: (interpreted_string_literal (interpreted_string_literal_content)))
  3 => 1: nil
"""
