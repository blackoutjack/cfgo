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
Declarations: package:module
Function definitions:
Entry: 0
Exit: 1
Graph:
  0 => 3: (package_clause (package_identifier))
  2 => 1: nil
  3 => 2: (import_spec path: (interpreted_string_literal (interpreted_string_literal_content)))
"""

run_import_splat = ["--no-ast", "test/import-splat.go"]

out_import_splat = """
Declarations: package.*:any
Function definitions:
Entry: 0
Exit: 1
Graph:
  0 => 3: (package_clause (package_identifier))
  2 => 1: nil
  3 => 2: (import_spec name: (dot) path: (interpreted_string_literal (interpreted_string_literal_content)))
"""

run_import_named = ["--no-ast", "test/import-named.go"]

out_import_named = """
Declarations: myname:module
Function definitions:
Entry: 0
Exit: 1
Graph:
  0 => 3: (package_clause (package_identifier))
  2 => 1: nil
  3 => 2: (import_spec name: (package_identifier) path: (interpreted_string_literal (interpreted_string_literal_content)))
"""

run_import_block = ["--no-ast", "test/import-block.go"]

out_import_block = """
Declarations: myname:module pkg:module
Function definitions:
Entry: 0
Exit: 1
Graph:
  0 => 3: (package_clause (package_identifier))
  2 => 1: nil
  3 => 4: (import_spec name: (package_identifier) path: (interpreted_string_literal (interpreted_string_literal_content)))
  4 => 2: (import_spec path: (interpreted_string_literal (interpreted_string_literal_content)))
"""

run_import_blank = ["--no-ast", "test/import-blank.go"]

out_import_blank = """
Declarations:
Function definitions:
Entry: 0
Exit: 1
Graph:
  0 => 3: (package_clause (package_identifier))
  2 => 1: nil
  3 => 2: (import_spec name: (blank_identifier) path: (interpreted_string_literal (interpreted_string_literal_content)))
"""

run_var_declaration = ["--no-ast", "test/var-declaration.go"]

out_var_declaration = """
Declarations: myvar:string one:int two:int inited:string implType:int a:string b:int c:string d:string
Function definitions:
Entry: 0
Exit: 1
Graph:
  0 => 3: (package_clause (package_identifier))
  2 => 1: nil
  3 => 4: (var_declaration (var_spec name: (identifier) type: (type_identifier)))
  4 => 5: (var_declaration (var_spec name: (identifier) name: (identifier) type: (type_identifier)))
  5 => 6: (var_declaration (var_spec name: (identifier) type: (type_identifier) value: (expression_list (interpreted_string_literal (interpreted_string_literal_content)))))
  6 => 7: (var_declaration (var_spec name: (identifier) value: (expression_list (int_literal))))
  7 => 8: (var_declaration (var_spec name: (identifier) name: (identifier) value: (expression_list (interpreted_string_literal (interpreted_string_literal_content)) (int_literal))))
  8 => 2: (var_declaration (var_spec name: (identifier) name: (identifier) type: (type_identifier) value: (expression_list (interpreted_string_literal (interpreted_string_literal_content)) (interpreted_string_literal (interpreted_string_literal_content)))))
"""

run_const_declaration = ["--no-ast", "test/const-declaration.go"]

out_const_declaration = """
Declarations: myconst:int implicitType:string constA:string constB:int constC:string constD:string
Function definitions:
Entry: 0
Exit: 1
Graph:
  0 => 3: (package_clause (package_identifier))
  2 => 1: nil
  3 => 4: (const_declaration (const_spec name: (identifier) type: (type_identifier) value: (expression_list (int_literal))))
  4 => 5: (const_declaration (const_spec name: (identifier) value: (expression_list (interpreted_string_literal (interpreted_string_literal_content)))))
  5 => 6: (const_declaration (const_spec name: (identifier) name: (identifier) value: (expression_list (interpreted_string_literal (interpreted_string_literal_content)) (int_literal))))
  6 => 2: (const_declaration (const_spec name: (identifier) name: (identifier) type: (type_identifier) value: (expression_list (interpreted_string_literal (interpreted_string_literal_content)) (interpreted_string_literal (interpreted_string_literal_content)))))
"""

run_func_declaration_basic = ["--no-ast", "test/func-declaration-basic.go"]

out_func_declaration_basic = """
Declarations: main:func()
Function definitions: main
Entry: 0
Exit: 1
Graph:
  0 => 3: (package_clause (package_identifier))
  2 => 1: nil
  3 => 2: (function_declaration name: (identifier) parameters: (parameter_list) body: (block))

Function: main
Declarations:
Function definitions:
Entry: 4
Exit: 5
Graph:
  4 => 5: nil
"""

run_func_declaration_generic = ["--no-ast", "test/func-declaration-generic.go"]

out_func_declaration_generic = """
Declarations: generic:func[E<any,F<~string](E,F)(F,E)
Function definitions: generic
Entry: 0
Exit: 1
Graph:
  0 => 3: (package_clause (package_identifier))
  2 => 1: nil
  3 => 2: (function_declaration name: (identifier) type_parameters: (type_parameter_list (type_parameter_declaration name: (identifier) type: (type_constraint (type_identifier))) (type_parameter_declaration name: (identifier) type: (type_constraint (negated_type (type_identifier))))) parameters: (parameter_list (parameter_declaration name: (identifier) type: (type_identifier)) (parameter_declaration name: (identifier) type: (type_identifier))) result: (parameter_list (parameter_declaration type: (type_identifier)) (parameter_declaration type: (type_identifier))) body: (block (return_statement (expression_list (identifier) (identifier)))))

Function: generic
Declarations:
Function definitions:
Entry: 4
Exit: 5
Graph:
  4 => 6: (return_statement (expression_list (identifier) (identifier)))
  6 => 5: nil
"""

run_func_declaration_params = ["--no-ast", "test/func-declaration-params.go"]

out_func_declaration_params = """
Declarations: withParams:func(string,int,int)
Function definitions: withParams
Entry: 0
Exit: 1
Graph:
  0 => 3: (package_clause (package_identifier))
  2 => 1: nil
  3 => 2: (function_declaration name: (identifier) parameters: (parameter_list (parameter_declaration name: (identifier) type: (type_identifier)) (parameter_declaration name: (identifier) name: (identifier) type: (type_identifier))) body: (block))

Function: withParams
Declarations:
Function definitions:
Entry: 4
Exit: 5
Graph:
  4 => 5: nil
"""

run_func_declaration_result = ["--no-ast", "test/func-declaration-result.go"]

out_func_declaration_result = """
Declarations: numstring:func()(int,string)
Function definitions: numstring
Entry: 0
Exit: 1
Graph:
  0 => 3: (package_clause (package_identifier))
  2 => 1: nil
  3 => 2: (function_declaration name: (identifier) parameters: (parameter_list) result: (parameter_list (parameter_declaration type: (type_identifier)) (parameter_declaration type: (type_identifier))) body: (block (return_statement (expression_list (int_literal) (interpreted_string_literal (interpreted_string_literal_content))))))

Function: numstring
Declarations:
Function definitions:
Entry: 4
Exit: 5
Graph:
  4 => 6: (return_statement (expression_list (int_literal) (interpreted_string_literal (interpreted_string_literal_content))))
  6 => 5: nil
"""

run_short_var_declaration = ["--no-ast", "test/short-var-declaration.go"]

out_short_var_declaration = """
Declarations: contain:func()
Function definitions: contain
Entry: 0
Exit: 1
Graph:
  0 => 3: (package_clause (package_identifier))
  2 => 1: nil
  3 => 2: (function_declaration name: (identifier) parameters: (parameter_list) body: (block (short_var_declaration left: (expression_list (identifier)) right: (expression_list (interpreted_string_literal))) (short_var_declaration left: (expression_list (identifier) (identifier)) right: (expression_list (int_literal) (int_literal)))))

Function: contain
Declarations: myvar:string one:int two:int
Function definitions:
Entry: 4
Exit: 5
Graph:
  4 => 7: (short_var_declaration left: (expression_list (identifier)) right: (expression_list (interpreted_string_literal)))
  6 => 5: nil
  7 => 6: (short_var_declaration left: (expression_list (identifier) (identifier)) right: (expression_list (int_literal) (int_literal)))
"""
