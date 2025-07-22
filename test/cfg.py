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

run_var_declaration = ["--no-ast", "test/var-declaration.go"]

out_var_declaration = """
Declarations: myvar one two inited implType
Function definitions:
Entry: 0
Exit: 1
Graph:
  0 => 2: (package_clause (package_identifier))
  2 => 3: (var_declaration (var_spec name: (identifier) type: (type_identifier)))
  3 => 4: (var_declaration (var_spec name: (identifier) name: (identifier) type: (type_identifier)))
  4 => 5: (var_declaration (var_spec name: (identifier) type: (type_identifier) value: (expression_list (interpreted_string_literal (interpreted_string_literal_content)))))
  5 => 6: (var_declaration (var_spec name: (identifier) value: (expression_list (int_literal))))
  6 => 1: nil
"""

run_const_declaration = ["--no-ast", "test/const-declaration.go"]

out_const_declaration = """
Declarations: myconst implicitType
Function definitions:
Entry: 0
Exit: 1
Graph:
  0 => 2: (package_clause (package_identifier))
  2 => 3: (const_declaration (const_spec name: (identifier) type: (type_identifier) value: (expression_list (int_literal))))
  3 => 4: (const_declaration (const_spec name: (identifier) value: (expression_list (interpreted_string_literal (interpreted_string_literal_content)))))
  4 => 1: nil
"""
