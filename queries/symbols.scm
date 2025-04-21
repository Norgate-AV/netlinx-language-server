(
  (section (define_device_section)) @device_section
  ([
    ; Pattern 1: [qualifier] [type] <identifier> = <value>
    (declaration
      (type_qualifier)? @qualifier
      type: _? @type
      declarator: (_
        declarator: (identifier) @identifier
        size: _? @size
        value: _? @value))

    (declaration
      (type_qualifier)? @qualifier
      type: _? @type
      declarator: (identifier) @identifier
      value: _? @value)

    ; Pattern 2: <identifier> = <value>
    (expression_statement
      (assignment_expression
        left: (identifier) @identifier
        right: _? @value))

    ; Pattern 3: <identifier>
    (expression_statement
      (identifier) @identifier)

    ; Pattern 4: <identifier>[] = <value>
    (expression_statement
      (assignment_expression
        left: (subscript_expression
          argument: (subscript_expression
            argument: (identifier) @identifier)
          index: _? @size)
        right: _? @value))
  ])+
)

(
  (section (define_constant_section)) @constant_section
  ([
    ; Pattern 1: [qualifier] [type] <identifier> = <value>
    (declaration
      (type_qualifier)? @qualifier
      type: _? @type
      declarator: (_
        declarator: (identifier) @identifier
        size: _? @size
        value: _? @value))

    (declaration
      (type_qualifier)? @qualifier
      type: _? @type
      declarator: (identifier) @identifier
      value: _? @value)

    ; Pattern 2: <identifier> = <value>
    (expression_statement
      (assignment_expression
        left: (identifier) @identifier
        right: _? @value))

    ; Pattern 3: <identifier>
    (expression_statement
      (identifier) @identifier)

    ; Pattern 4: <identifier>[] = <value>
    (expression_statement
      (assignment_expression
        left: (subscript_expression
          argument: (subscript_expression
            argument: (identifier) @identifier)
          index: _? @size)
        right: _? @value))
  ])+
)

(
  (section (define_variable_section)) @variable_section
  ([
    ; Pattern 1: [qualifier] [type] <identifier> = <value>
    (declaration
      (type_qualifier)? @qualifier
      type: _? @type
      declarator: (_
        declarator: (identifier) @identifier
        size: _? @size
        value: _? @value))

    (declaration
      (type_qualifier)? @qualifier
      type: _? @type
      declarator: (identifier) @identifier
      value: _? @value)

    ; Pattern 2: <identifier> = <value>
    (expression_statement
      (assignment_expression
        left: (identifier) @identifier
        right: _? @value))

    ; Pattern 3: <identifier>
    (expression_statement
      (identifier) @identifier)

    ; Pattern 4: <identifier>[] = <value>
    (expression_statement
      (assignment_expression
        left: (subscript_expression
          argument: (subscript_expression
            argument: (identifier) @identifier)
          index: _? @size)
        right: _? @value))
  ])+
)

(
  (section (define_type_section)) @type_section
  (type_definition
    (struct_specifier
      name: (type_identifier) @type_identifier
      body: _ @type_body))+
)

(
  (define_function_keyword) @function_definition
  (function_definition
    return_type: _? @return_type
    name: (identifier) @function_name
    parameters: (parameter_list
      ((parameter_declaration
        type: _? @param_type
        declarator: (identifier) @param_name))*))
)
