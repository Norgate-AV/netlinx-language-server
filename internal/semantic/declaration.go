package semantic

// DeclarationVisitor collects declarations during tree traversal
// type DeclarationVisitor struct {
// 	analyzer     *Analyzer
// 	currentScope *Scope
// }

// Visit handles a single node in the declaration phase
// func (v *DeclarationVisitor) Visit(node *tree_sitter.Node, content []byte) bool {
// 	nodeKind := node.Kind()

// 	// Handle different node types for declarations
// 	switch nodeKind {
// 	case "program_name_declaration":
// 		// v.handleProgramNameDecl(node, content)

// 	case "define_device_section", "define_constant_section", "define_variable_section":
// 		// v.handleSectionDecl(node, content)

// 	case "function_declaration":
// 		// v.handleFunctionDecl(node, content)

// 	case "variable_declaration":
// 		// v.handleVariableDecl(node, content)

// 	case "{", "if_statement", "while_statement", "for_statement", "switch_statement":
// 		// Create a new nested scope
// 		scope := &Scope{
// 			Parent:  v.currentScope,
// 			Symbols: make(map[string]*Symbol),
// 			Node:    node,
// 			Range:   nodeToRange(node),
// 			// ScopeType: getScopeTypeForNode(node),
// 		}
// 		v.currentScope.Children = append(v.currentScope.Children, scope)
// 		v.analyzer.document.Scopes = append(v.analyzer.document.Scopes, scope)

// 		// Save current scope, enter new scope, process children, restore scope
// 		// oldScope := v.currentScope
// 		v.currentScope = scope
// 		// Continue traversal inside the scope
// 		return true

// 	case "}", "end_if", "end_while", "end_for", "end_switch":
// 		// Exit current scope
// 		if v.currentScope.Parent != nil {
// 			v.currentScope = v.currentScope.Parent
// 		}
// 	}

// 	// Continue traversal for other nodes
// 	return true
// }

// func (v *DeclarationVisitor) handleProgramNameDecl(node *tree_sitter.Node, content []byte) {
// 	// Extract program name value
// 	valueNode := findChildByType(node, "string")
// 	if valueNode == nil {
// 		return
// 	}

// 	name := getNodeText(valueNode, content)
// 	// Trim quotes
// 	name = name[1 : len(name)-1]
// 	v.analyzer.document.ProgramName = name
// }

// Additional handler methods for various declarations...
