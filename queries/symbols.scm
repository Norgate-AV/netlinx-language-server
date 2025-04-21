; Sections
(section (define_device_section)) @section.define_device
(section (define_constant_section)) @section.define_constant
(section (define_variable_section)) @section.define_variable
(section (define_type_section)) @section.define_type
(section (define_start_section)) @section.define_start
(section (define_event_section)) @section.define_event
(section (define_program_section)) @section.define_program

; Pattern 1: [qualifier] [type] <identifier> = <value>
(declaration
  (storage_class_specifier)? @symbol.storage
  (type_qualifier)? @symbol.qualifier
  type: _? @symbol.type
  declarator: (_
    declarator: (identifier) @symbol.identifier
    size: _? @symbol.size
    value: _? @symbol.value)) @symbol.declaration

(declaration
  (storage_class_specifier)? @symbol.storage
  (type_qualifier)? @symbol.qualifier
  type: _? @symbol.type
  declarator: (identifier) @symbol.identifier
  value: _? @symbol.value) @symbol.declaration

; Pattern 2: <identifier> = <value>
(expression_statement
  (assignment_expression
    left: (identifier) @symbol.identifier
    right: _? @symbol.value)) @symbol.declaration

; Pattern 3: <identifier>
(expression_statement
  (identifier) @symbol.identifier) @symbol.declaration

; Pattern 4: <identifier>[] = <value>
(expression_statement
  (assignment_expression
    left: (subscript_expression
      argument: (subscript_expression
        argument: (identifier) @symbol.identifier)
      index: _? @symbol.size)
    right: _? @symbol.value)) @symbol.declaration

(type_definition
  (struct_specifier
    name: (type_identifier) @type.identifier
    body: _ @type.body)) @type.definition

(
  (define_function_keyword)
  (function_definition
    return_type: _? @function.return_type
    name: (identifier) @function.name
    parameters: (parameter_list
      ((parameter_declaration
        type: _? @function.parameter.type
        declarator: (identifier) @function.parameter.name))*))
) @function.definition
