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

run_import_basic = ["--no-cfg", "test/import-basic.go"]

out_import_basic = """
(source_file (package_clause (package_identifier)) (import_declaration (import_spec path: (interpreted_string_literal (interpreted_string_literal_content)))))
"""

run_import_splat = ["--no-cfg", "test/import-splat.go"]

out_import_splat = """
(source_file (package_clause (package_identifier)) (import_declaration (import_spec name: (dot) path: (interpreted_string_literal (interpreted_string_literal_content)))))
"""

run_import_named = ["--no-cfg", "test/import-named.go"]

out_import_named = """
(source_file (package_clause (package_identifier)) (import_declaration (import_spec name: (package_identifier) path: (interpreted_string_literal (interpreted_string_literal_content)))))
"""

run_import_block = ["--no-cfg", "test/import-block.go"]

out_import_block = """
(source_file (package_clause (package_identifier)) (import_declaration (import_spec_list (import_spec name: (package_identifier) path: (interpreted_string_literal (interpreted_string_literal_content))) (import_spec path: (interpreted_string_literal (interpreted_string_literal_content))))))
"""

run_import_blank = ["--no-cfg", "test/import-blank.go"]

out_import_blank = """
(source_file (package_clause (package_identifier)) (import_declaration (import_spec name: (blank_identifier) path: (interpreted_string_literal (interpreted_string_literal_content)))))
"""
