package semantic

// ReferenceVisitor resolves references during tree traversal
// type ReferenceVisitor struct {
// 	analyzer     *Analyzer
// 	currentScope *Scope
// }

// Visit handles a single node in the reference resolution phase
// func (v *ReferenceVisitor) Visit(node *tree_sitter.Node, content []byte) bool {
// 	nodeKind := node.Kind()

// 	// Track scope changes just like in declaration visitor
// 	switch nodeKind {
// 	case "{", "if_statement", "while_statement", "for_statement", "switch_statement":
// 		// Find corresponding scope and enter it
// 		for _, scope := range v.analyzer.document.Scopes {
// 			if scope.Node == node {
// 				v.currentScope = scope
// 				break
// 			}
// 		}

// 	case "}", "end_if", "end_while", "end_for", "end_switch":
// 		// Exit current scope
// 		if v.currentScope.Parent != nil {
// 			v.currentScope = v.currentScope.Parent
// 		}

// 	case "identifier":
// 		// This might be a reference to a declared symbol
// 		v.resolveReference(node, content)
// 	}

// 	return true
// }

// func (v *ReferenceVisitor) resolveReference(node *tree_sitter.Node, content []byte) {
// 	// Skip identifiers that are part of declarations
// 	// if isPartOfDeclaration(node) {
// 	// 	return
// 	// }

// 	name := getNodeText(node, content)

// 	// Resolve symbol from current scope upward
// 	symbol := v.lookupSymbolInScope(v.currentScope, name)
// 	if symbol != nil {
// 		// Create and add reference
// 		// ref := &Reference{
// 		// 	Symbol:  symbol,
// 		// 	Node:    node,
// 		// 	Range:   nodeToRange(node),
// 		// 	// IsWrite: isWriteContext(node),
// 		// }

// 		// symbol.References = append(symbol.References, ref)
// 	}
// }

// func (v *ReferenceVisitor) lookupSymbolInScope(scope *Scope, name string) *Symbol {
// 	// Look in current scope
// 	if symbol, exists := scope.Symbols[name]; exists {
// 		return symbol
// 	}

// 	// Look in parent scopes
// 	if scope.Parent != nil {
// 		return v.lookupSymbolInScope(scope.Parent, name)
// 	}

// 	// Try global symbol table as last resort
// 	return v.analyzer.document.Symbols[name]
// }

// Other helper functions for reference resolution...
