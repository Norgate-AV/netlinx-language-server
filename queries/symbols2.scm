(define_device_section) @device_section
(define_constant_section) @constant_section
(define_variable_section) @variable_section

[
  ; Pattern 1: [qualifier] [type] <identifier> = <value>
  (declaration
    (type_qualifier)? @qualifier
    type: _? @type
    declarator: (_
      declarator: (identifier) @identifier
      size: _? @size
      value: _? @value))

  ; Pattern 2: <identifier> = <value>
  (expression_statement
    (assignment_expression
      left: (identifier) @identifier
      right: _ @value))

  ; Pattern 3: <identifier>
  (expression_statement
    (identifier) @identifier)

  ; Pattern 4: <identifier>[] = <value>
  (expression_statement
    (assignment_expression
      left: (subscript_expression
        argument: (subscript_expression
          argument: (identifier) @identifier)
        index: (_)? @size)
      right: (_)? @value))

  ; (define_function
  ;   (function_definition
  ;     return_type: (_)? @return_type
  ;     name: @identifier))
      ; parameters: (parameter_list)
      ;   (parameter_declaration
      ;     type: @type
      ;     declarator: @parameter
      ;     size: (_)? @size)?)
]

