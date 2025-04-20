; Device section queries
; File: queries/device.scm

; Track the current section
(define_device_section) @device_section

; Pattern 1: [qualifier] [type] <identifier> = <value>
((declaration
  (type_qualifier)? @qualifier
  type: _? @type
  declarator: (init_declarator
    declarator: (identifier) @identifier
    value: _ @value)))

; Pattern 2: <identifier> = <value>
((expression_statement
  (assignment_expression
    left: (identifier) @identifier
    right: _ @value)))
