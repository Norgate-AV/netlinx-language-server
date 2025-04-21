package semantic

import (
	"github.com/Norgate-AV/netlinx-language-server/internal/logger"
	"github.com/Norgate-AV/netlinx-language-server/internal/parser"

	tree_sitter "github.com/tree-sitter/go-tree-sitter"
)

type Analyzer struct {
	parser   *parser.TreeSitter
	logger   logger.Logger
	document *Document
	content  []byte
}

func NewAnalyzer(parser *parser.TreeSitter, logger logger.Logger) *Analyzer {
	return &Analyzer{
		parser: parser,
		logger: logger,
	}
}

// Analyze performs semantic analysis on a NetLinx document
func (a *Analyzer) Analyze(uri string, content string, tree *tree_sitter.Tree) (*Document, error) {
	if tree == nil || tree.RootNode() == nil {
		return &Document{URI: uri, Symbols: make(map[string]*Symbol)}, nil
	}

	// Initialize document and state
	a.document = &Document{
		URI:           uri,
		Symbols:       make(map[string]*Symbol),
		SymbolsByKind: make(map[SymbolType][]*Symbol),
	}
	a.content = []byte(content)

	// Create global scope
	// a.document.GlobalScope = &Scope{
	// 	Symbols:   make(map[string]*Symbol),
	// 	Node:      tree.RootNode(),
	// 	Range:     nodeToRange(tree.RootNode()),
	// 	ScopeType: "global",
	// }

	a.document.Scopes = append(a.document.Scopes, a.document.GlobalScope)

	// First pass: Collect all declarations
	// declVisitor := &DeclarationVisitor{
	// 	analyzer:     a,
	// 	currentScope: a.document.GlobalScope,
	// }

	// declWalker := NewTreeWalker(declVisitor, a.content, a.logger)
	// declWalker.Walk(tree.RootNode())

	// // Second pass: Resolve references
	// refVisitor := &ReferenceVisitor{
	// 	analyzer:     a,
	// 	currentScope: a.document.GlobalScope,
	// }

	// refWalker := NewTreeWalker(refVisitor, a.content, a.logger)
	// refWalker.Walk(tree.RootNode())

	return a.document, nil
}

// findAllSections locates all section nodes in the AST
// func (a *Analyzer) findAllSections(root *tree_sitter.Node, content []byte) []*Section {
// 	var sections []*Section

// 	query, err := parser.CreateQuery(`
//         (define_device_section) @device_section
//         (define_constant_section) @constant_section
//         (define_variable_section) @variable_section
//     `)
// 	if err != nil {
// 		a.logger.Printf("Failed to create section query: %v", err)
// 		return sections
// 	}
// 	defer query.Close()

// 	cursor := parser.CreateQueryCursor()
// 	defer cursor.Close()

// 	captures := cursor.Captures(query, root, content)

// 	for match, index := captures.Next(); match != nil; match, index = captures.Next() {
// 		node := match.Captures[index].Node
// 		var sectionKind string

// 		// switch query.CaptureNameForId(uint(index)) {
// 		// case "device_section":
// 		// 	sectionKind = "DEFINE_DEVICE"
// 		// case "constant_section":
// 		// 	sectionKind = "DEFINE_CONSTANT"
// 		// case "variable_section":
// 		// 	sectionKind = "DEFINE_VARIABLE"
// 		// }

// 		sections = append(sections, &Section{
// 			Kind:  sectionKind,
// 			Node:  &node,
// 			Range: nodeToRange(&node),
// 		})
// 	}

// 	return sections
// }

// processDeviceSection analyzes device declarations within a DEFINE_DEVICE section
// func (a *Analyzer) processDeviceSection(doc *Document, section *Section, content []byte) {
// 	// Query for device declarations within this section
// 	query, err := parser.CreateQuery(`
//         (section
//             name: (section_name) @section_name
//             body: (section_body
//                 (assignment_statement
//                     left: (identifier) @device_name
//                     right: (device_number) @device_value
//                 )
//             )
//         ) @device_section
//     `)
// 	if err != nil {
// 		a.logger.Printf("Failed to create device section query: %v", err)
// 		return
// 	}
// 	defer query.Close()

// 	cursor := parser.CreateQueryCursor()
// 	defer cursor.Close()

// 	// cursor.Exec(query, section.Node, content)

// 	// for match := cursor.NextMatch(); match != nil; match = cursor.NextMatch() {
// 	// 	for i := 0; i < int(match.CaptureCount); i++ {
// 	// 		capture := match.Captures[i]
// 	// 		name := query.CaptureNameForId(uint(i))

// 	// 		if name == "device_name" {
// 	// 			deviceName := getNodeText(capture.Node, content)
// 	// 			deviceValue := ""

// 	// 			// Find the corresponding device value
// 	// 			for j := 0; j < int(match.CaptureCount); j++ {
// 	// 				if query.CaptureNameForId(uint(j)) == "device_value" {
// 	// 					deviceValue = getNodeText(match.Captures[j].Node, content)
// 	// 					break
// 	// 				}
// 	// 			}

// 	// 			symbol := &Symbol{
// 	// 				Name:    deviceName,
// 	// 				Type:    DeviceSymbol,
// 	// 				Value:   deviceValue,
// 	// 				Range:   nodeToRange(capture.Node),
// 	// 				Section: section,
// 	// 				Node:    capture.Node,
// 	// 			}

// 	// 			doc.Symbols[deviceName] = symbol
// 	// 		}
// 	// 	}
// 	// }
// }
