// internal/analysis/semantic/analyzer.go
package semantic

import (
	"github.com/Norgate-AV/netlinx-language-server/internal/lsp"
	"github.com/Norgate-AV/netlinx-language-server/parser"
	"github.com/sirupsen/logrus"
	tree_sitter "github.com/tree-sitter/go-tree-sitter"
)

type Analyzer struct {
	parser *parser.TreeSitter
	logger logrus.FieldLogger
}

func NewAnalyzer(parser *parser.TreeSitter, logger logrus.FieldLogger) *Analyzer {
	return &Analyzer{
		parser: parser,
		logger: logger,
	}
}

// Analyze performs semantic analysis on a NetLinx document
func (a *Analyzer) Analyze(uri string, content string, tree *tree_sitter.Tree) (*Document, error) {
	doc := &Document{
		URI:     uri,
		Symbols: make(map[string]*Symbol),
	}

	if tree == nil || tree.RootNode() == nil {
		return doc, nil
	}

	root := tree.RootNode()
	contentBytes := []byte(content)

	// First pass: find all sections
	doc.Sections = a.findAllSections(root, contentBytes)

	// Second pass: process each section
	for _, section := range doc.Sections {
		switch section.Kind {
		case "DEFINE_DEVICE":
			a.processDeviceSection(doc, section, contentBytes)
		case "DEFINE_CONSTANT":
			a.processConstantSection(doc, section, contentBytes)
		case "DEFINE_VARIABLE":
			a.processVariableSection(doc, section, contentBytes)
		}
	}

	// Process functions and other top-level declarations
	a.processFunctions(doc, root, contentBytes)

	return doc, nil
}

// findAllSections locates all section nodes in the AST
func (a *Analyzer) findAllSections(root *tree_sitter.Node, content []byte) []*Section {
	var sections []*Section

	query, err := parser.CreateQuery(`
        (define_device_section) @device_section
        (define_constant_section) @constant_section
        (define_variable_section) @variable_section
    `)
	if err != nil {
		a.logger.Errorf("Failed to create section query: %v", err)
		return sections
	}
	defer query.Close()

	cursor := parser.CreateQueryCursor()
	defer cursor.Close()

	captures := cursor.Captures(query, root, content)

	for match, index := captures.Next(); match != nil; match, index = captures.Next() {
		node := match.Captures[index].Node
		var sectionKind string

		switch query.CaptureNameForId(uint(index)) {
		case "device_section":
			sectionKind = "DEFINE_DEVICE"
		case "constant_section":
			sectionKind = "DEFINE_CONSTANT"
		case "variable_section":
			sectionKind = "DEFINE_VARIABLE"
		}

		sections = append(sections, &Section{
			Kind:  sectionKind,
			Node:  node,
			Range: nodeToRange(node),
		})
	}

	return sections
}

// processDeviceSection analyzes device declarations within a DEFINE_DEVICE section
func (a *Analyzer) processDeviceSection(doc *Document, section *Section, content []byte) {
	// Query for device declarations within this section
	query, err := parser.CreateQuery(`
        (section
            name: (section_name) @section_name
            body: (section_body
                (assignment_statement
                    left: (identifier) @device_name
                    right: (device_number) @device_value
                )
            )
        ) @device_section
    `)
	if err != nil {
		a.logger.Errorf("Failed to create device section query: %v", err)
		return
	}
	defer query.Close()

	cursor := parser.CreateQueryCursor()
	defer cursor.Close()

	cursor.Exec(query, section.Node, content)

	for match := cursor.NextMatch(); match != nil; match = cursor.NextMatch() {
		for i := 0; i < int(match.CaptureCount); i++ {
			capture := match.Captures[i]
			name := query.CaptureNameForId(uint(i))

			if name == "device_name" {
				deviceName := getNodeText(capture.Node, content)
				deviceValue := ""

				// Find the corresponding device value
				for j := 0; j < int(match.CaptureCount); j++ {
					if query.CaptureNameForId(uint(j)) == "device_value" {
						deviceValue = getNodeText(match.Captures[j].Node, content)
						break
					}
				}

				symbol := &Symbol{
					Name:    deviceName,
					Type:    DeviceSymbol,
					Value:   deviceValue,
					Range:   nodeToRange(capture.Node),
					Section: section,
					Node:    capture.Node,
				}

				doc.Symbols[deviceName] = symbol
			}
		}
	}
}

// Additional helper methods (processConstantSection, processVariableSection, etc.)
// ...

// Helper to convert node to LSP Range
func nodeToRange(node *tree_sitter.Node) lsp.Range {
	start := node.StartPoint()
	end := node.EndPoint()

	return lsp.Range{
		Start: lsp.Position{
			Line:      uint32(start.Row),
			Character: uint32(start.Column),
		},
		End: lsp.Position{
			Line:      uint32(end.Row),
			Character: uint32(end.Column),
		},
	}
}

func getNodeText(node *tree_sitter.Node, content []byte) string {
	start := node.StartByte()
	end := node.EndByte()

	if start >= end || uint32(len(content)) < end {
		return ""
	}

	return string(content[start:end])
}
