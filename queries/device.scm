; Device section queries
; File: queries/device.scm

; Track the current section
(define_device_section) @device_section
(define_constant_section) @constant_section
(define_variable_section) @variable_section

; Pattern 1: [qualifier] [type] <identifier> = <value>
((declaration
  (type_qualifier)? @qualifier
  type: _? @type
  (type_identifier)? @type
  declarator: (_
    declarator: (identifier) @identifier
    size: _? @size
    value: _? @value)))

; Pattern 2: <identifier> = <value>
((expression_statement
  (assignment_expression
    left: (identifier) @identifier
    right: _ @value)))

; Pattern 3: <identifier>
(expression_statement
  (identifier) @identifier)

; Pattern 4: [qualifier] [custom_type] <identifier> = <value>
; ((declaration
;   (type_qualifier)? @qualifier
;   (type_identifier)? @type
;   declarator: (_
;     declarator: (identifier) @identifier
;     size: _? @size
;     value: _ @value)))
