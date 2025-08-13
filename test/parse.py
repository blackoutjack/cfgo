"""
Test AST generation for basic statements and expressions
"""


run_empty_file = ["--no-cfg", "test/empty.go"]

out_empty_file = """
source_file
"""

run_package = ["--no-cfg", "test/package.go"]

out_package = """
source_file
  package_clause
    package
    package_identifier
"""

run_import_basic = ["--no-cfg", "test/import-basic.go"]

out_import_basic = """
source_file
  package_clause
    package
    package_identifier
  import_declaration
    import
    import_spec
      path:interpreted_string_literal
        "
        interpreted_string_literal_content
        "
"""

run_import_splat = ["--no-cfg", "test/import-splat.go"]

out_import_splat = """
source_file
  package_clause
    package
    package_identifier
  import_declaration
    import
    import_spec
      name:dot
        .
      path:interpreted_string_literal
        "
        interpreted_string_literal_content
        "
"""

run_import_named = ["--no-cfg", "test/import-named.go"]

out_import_named = """
source_file
  package_clause
    package
    package_identifier
  import_declaration
    import
    import_spec
      name:package_identifier
      path:interpreted_string_literal
        "
        interpreted_string_literal_content
        "
"""

run_import_block = ["--no-cfg", "test/import-block.go"]

out_import_block = """
source_file
  package_clause
    package
    package_identifier
  import_declaration
    import
    import_spec_list
      (
      import_spec
        name:package_identifier
        path:interpreted_string_literal
          "
          interpreted_string_literal_content
          "
      import_spec
        path:interpreted_string_literal
          "
          interpreted_string_literal_content
          "
      )
"""

run_import_blank = ["--no-cfg", "test/import-blank.go"]

out_import_blank = """
source_file
  package_clause
    package
    package_identifier
  import_declaration
    import
    import_spec
      name:blank_identifier
      path:interpreted_string_literal
        "
        interpreted_string_literal_content
        "
"""

run_var_declaration = ["--no-cfg", "test/var-declaration.go"]

out_var_declaration = """
source_file
  package_clause
    package
    package_identifier
  var_declaration
    var
    var_spec
      name:identifier
      type:type_identifier
  var_declaration
    var
    var_spec
      name:identifier
      ,
      name:identifier
      type:type_identifier
  var_declaration
    var
    var_spec
      name:identifier
      type:type_identifier
      =
      value:expression_list
        interpreted_string_literal
          "
          interpreted_string_literal_content
          "
  var_declaration
    var
    var_spec
      name:identifier
      =
      value:expression_list
        int_literal
  var_declaration
    var
    var_spec
      name:identifier
      ,
      name:identifier
      =
      value:expression_list
        interpreted_string_literal
          "
          interpreted_string_literal_content
          "
        ,
        int_literal
  var_declaration
    var
    var_spec
      name:identifier
      ,
      name:identifier
      type:type_identifier
      =
      value:expression_list
        interpreted_string_literal
          "
          interpreted_string_literal_content
          "
        ,
        interpreted_string_literal
          "
          interpreted_string_literal_content
          "
"""

run_const_declaration = ["--no-cfg", "test/const-declaration.go"]

out_const_declaration = """
source_file
  package_clause
    package
    package_identifier
  const_declaration
    const
    const_spec
      name:identifier
      type:type_identifier
      =
      value:expression_list
        int_literal
  const_declaration
    const
    const_spec
      name:identifier
      =
      value:expression_list
        interpreted_string_literal
          "
          interpreted_string_literal_content
          "
  const_declaration
    const
    const_spec
      name:identifier
      name:,
      name:identifier
      =
      value:expression_list
        interpreted_string_literal
          "
          interpreted_string_literal_content
          "
        ,
        int_literal
  const_declaration
    const
    const_spec
      name:identifier
      name:,
      name:identifier
      type:type_identifier
      =
      value:expression_list
        interpreted_string_literal
          "
          interpreted_string_literal_content
          "
        ,
        interpreted_string_literal
          "
          interpreted_string_literal_content
          "
"""

run_func_declaration_basic = ["--no-cfg", "test/func-declaration-basic.go"]

out_func_declaration_basic = """
source_file
  package_clause
    package
    package_identifier
  function_declaration
    func
    name:identifier
    parameters:parameter_list
      (
      )
    body:block
      {
      }
"""

run_func_declaration_generic = ["--no-cfg", "test/func-declaration-generic.go"]

out_func_declaration_generic = """
source_file
  package_clause
    package
    package_identifier
  function_declaration
    func
    name:identifier
    type_parameters:type_parameter_list
      [
      type_parameter_declaration
        name:identifier
        type:type_constraint
          type_identifier
      ,
      type_parameter_declaration
        name:identifier
        type:type_constraint
          negated_type
            ~
            type_identifier
      ]
    parameters:parameter_list
      (
      parameter_declaration
        name:identifier
        type:type_identifier
      ,
      parameter_declaration
        name:identifier
        type:type_identifier
      )
    result:parameter_list
      (
      parameter_declaration
        type:type_identifier
      ,
      parameter_declaration
        type:type_identifier
      )
    body:block
      {
      return_statement
        return
        expression_list
          identifier
          ,
          identifier
      }
"""


run_func_declaration_params = ["--no-cfg", "test/func-declaration-params.go"]

out_func_declaration_params = """
source_file
  package_clause
    package
    package_identifier
  function_declaration
    func
    name:identifier
    parameters:parameter_list
      (
      parameter_declaration
        name:identifier
        type:type_identifier
      ,
      parameter_declaration
        name:identifier
        ,
        name:identifier
        type:type_identifier
      )
    body:block
      {
      }
"""

run_func_declaration_result = ["--no-cfg", "test/func-declaration-result.go"]

out_func_declaration_result = """
source_file
  package_clause
    package
    package_identifier
  function_declaration
    func
    name:identifier
    parameters:parameter_list
      (
      )
    result:parameter_list
      (
      parameter_declaration
        type:type_identifier
      ,
      parameter_declaration
        type:type_identifier
      )
    body:block
      {
      return_statement
        return
        expression_list
          int_literal
          ,
          interpreted_string_literal
            "
            interpreted_string_literal_content
            "
      }
"""

run_short_var_declaration = ["--no-cfg", "test/short-var-declaration.go"]

out_short_var_declaration = """
source_file
  package_clause
    package
    package_identifier
  function_declaration
    func
    name:identifier
    parameters:parameter_list
      (
      )
    body:block
      {
      short_var_declaration
        left:expression_list
          identifier
        :=
        right:expression_list
          interpreted_string_literal
            "
            "
      short_var_declaration
        left:expression_list
          identifier
          ,
          identifier
        :=
        right:expression_list
          int_literal
          ,
          int_literal
      }
"""

