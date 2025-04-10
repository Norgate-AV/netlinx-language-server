package analysis

import (
	"github.com/Norgate-AV/netlinx-language-server/internal/lsp"
	tree_sitter "github.com/tree-sitter/go-tree-sitter"
)

func (s *State) ExtractSymbols(uri string) ([]lsp.DocumentSymbol, error) {
	// Get the syntax tree using your existing parser
	tree, ok := s.GetSyntaxTree(uri)
	if !ok {
		// s.Logger.Debug("No syntax tree found for "+uri, logrus.Fields{})
		return []lsp.DocumentSymbol{}, nil
	}

	root := tree.RootNode()
	if root == nil {
		// s.Logger.Debug("Root node is nil for "+uri, logrus.Fields{})
		return []lsp.DocumentSymbol{}, nil
	}

	s.Logger.Printf("DEBUG: We have a root node for %s", uri)

	// s.Logger.Printf("DEBUG: Root node kind: %s", root.Kind())

	// Get the document content for text extraction
	content, ok := s.GetDocument(uri)
	if !ok {
		// s.Logger.Printf("DEBUG: No content found for %s", uri)
		return []lsp.DocumentSymbol{}, nil
	}

	// Debug functions - add at top of file
	// dumpNode := func(node *tree_sitter.Node, prefix string) {
	// 	if node == nil {
	// 		return
	// 	}
	// 	log.Printf("%sNode kind: %s, text: %s",
	// 		prefix,
	// 		node.Kind(),
	// 		getNodeText(node, []byte(content))[0:min(30, len(getNodeText(node, []byte(content))))])

	// 	for i := uint(0); i < node.ChildCount(); i++ {
	// 		child := node.Child(i)
	// 		dumpNode(child, prefix+"  ")
	// 	}
	// }

	// // Dump the first few levels to see structure
	// log.Printf("DEBUG: Dumping AST structure")
	// dumpNode(root, "")

	// Walk the AST and extract symbols
	symbols := []lsp.DocumentSymbol{}

	// Extract program name
	programNameNode := findNode(root, "program_name_declaration")
	if programNameNode != nil {
		valueNode := findChildByType(programNameNode, "string")
		if valueNode != nil {
			name := getNodeText(valueNode, []byte(content))
			// Remove quotes from string
			if len(name) > 2 {
				name = name[1 : len(name)-1]
			}
			symbol := createSymbol("PROGRAM: "+name, lsp.SymbolKindFile, programNameNode)
			symbols = append(symbols, symbol)
		}
	}

	// Process different sections (DEFINE_DEVICE, DEFINE_CONSTANT, etc.)
	for _, sectionType := range []string{
		"define_device_section",
		"define_constant_section",
		"define_variable_section",
	} {
		sectionNodes := findNodes(root, sectionType)
		for _, sectionNode := range sectionNodes {
			// Create section symbol
			sectionName := getSectionName(sectionType)
			sectionSymbol := createSymbol(sectionName, lsp.SymbolKindNamespace, sectionNode)

			// Find declarations within this section
			childSymbols := extractDeclarations(sectionNode, []byte(content))
			sectionSymbol.Children = childSymbols

			symbols = append(symbols, sectionSymbol)
		}
	}

	// Extract functions
	functionNodes := findNodes(root, "function_declaration")
	for _, funcNode := range functionNodes {
		symbol := extractFunctionSymbol(funcNode, []byte(content))
		symbols = append(symbols, symbol)
	}

	return symbols, nil
}

// Helper functions to work with the AST
func findNode(node *tree_sitter.Node, nodeType string) *tree_sitter.Node {
	if node == nil {
		return nil
	}

	if node.Kind() == nodeType {
		return node
	}

	childCount := node.ChildCount()
	for i := uint(0); i < childCount; i++ {
		child := node.Child(i)
		if found := findNode(child, nodeType); found != nil {
			return found
		}
	}

	return nil
}

func findNodes(node *tree_sitter.Node, nodeType string) []*tree_sitter.Node {
	var results []*tree_sitter.Node

	if node == nil {
		return results
	}

	if node.Kind() == nodeType {
		results = append(results, node)
	}

	childCount := node.ChildCount()
	for i := uint(0); i < childCount; i++ {
		child := node.Child(i)
		results = append(results, findNodes(child, nodeType)...)
	}

	return results
}

func findChildByType(node *tree_sitter.Node, nodeType string) *tree_sitter.Node {
	if node == nil {
		return nil
	}

	childCount := node.ChildCount()
	for i := uint(0); i < childCount; i++ {
		child := node.Child(i)
		if child.Kind() == nodeType {
			return child
		}
	}

	return nil
}

func getNodeText(node *tree_sitter.Node, content []byte) string {
	if node == nil {
		return ""
	}

	start := node.StartByte()
	end := node.EndByte()

	if start > end || uint(len(content)) < uint(end) {
		return ""
	}

	return string(content[start:end])
}

func createSymbol(name string, kind lsp.SymbolKind, node *tree_sitter.Node) lsp.DocumentSymbol {
	startPos := node.StartPosition()
	endPos := node.EndPosition()

	return lsp.DocumentSymbol{
		Name: name,
		Kind: kind,
		Range: lsp.Range{
			Start: lsp.Position{Line: uint32(startPos.Row), Character: uint32(startPos.Column)},
			End:   lsp.Position{Line: uint32(endPos.Row), Character: uint32(endPos.Column)},
		},
		SelectionRange: lsp.Range{
			Start: lsp.Position{Line: uint32(startPos.Row), Character: uint32(startPos.Column)},
			End:   lsp.Position{Line: uint32(endPos.Row), Character: uint32(endPos.Column)},
		},
		Children: []lsp.DocumentSymbol{},
	}
}

func getSectionName(sectionType string) string {
	switch sectionType {
	case "define_device_section":
		return "DEFINE_DEVICE"
	case "define_constant_section":
		return "DEFINE_CONSTANT"
	case "define_variable_section":
		return "DEFINE_VARIABLE"
	default:
		return sectionType
	}
}

func extractDeclarations(node *tree_sitter.Node, content []byte) []lsp.DocumentSymbol {
	var symbols []lsp.DocumentSymbol

	// Look for variable declarations
	varNodes := findNodes(node, "variable_declaration")
	for _, varNode := range varNodes {
		// Find identifier node (variable name)
		idNode := findChildByType(varNode, "identifier")
		if idNode != nil {
			name := getNodeText(idNode, content)

			// Find type node
			typeNode := findChildByType(varNode, "type")
			var typeName string
			if typeNode != nil {
				typeName = getNodeText(typeNode, content)
			}

			symbol := createSymbol(name, lsp.SymbolKindVariable, varNode)
			symbol.Detail = typeName
			symbols = append(symbols, symbol)
		}
	}

	return symbols
}

func extractFunctionSymbol(node *tree_sitter.Node, content []byte) lsp.DocumentSymbol {
	// Get function name
	nameNode := findChildByType(node, "identifier")
	name := "function"
	if nameNode != nil {
		name = getNodeText(nameNode, content)
	}

	// Determine if it's a regular function or define_program
	kind := lsp.SymbolKindFunction
	typeNode := findChildByType(node, "define_program_keyword")
	if typeNode != nil {
		kind = lsp.SymbolKindModule
	}

	// Get return type if any
	returnTypeNode := findChildByType(node, "type")
	detail := ""
	if returnTypeNode != nil {
		detail = getNodeText(returnTypeNode, content)
	}

	symbol := createSymbol(name, kind, node)
	symbol.Detail = detail

	return symbol
}
