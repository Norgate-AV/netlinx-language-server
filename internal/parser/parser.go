package parser

import (
	"errors"
	"fmt"

	tree_sitter_netlinx "github.com/norgate-av/tree-sitter-netlinx/bindings/go"
	tree_sitter "github.com/tree-sitter/go-tree-sitter"
)

type Parser struct {
	parser   *tree_sitter.Parser
	language *tree_sitter.Language
}

type Tree = tree_sitter.Tree

func NewParser() *Parser {
	language := tree_sitter.NewLanguage(tree_sitter_netlinx.Language())

	parser := tree_sitter.NewParser()
	parser.SetLanguage(language)

	return &Parser{
		parser:   parser,
		language: language,
	}
}

func (p *Parser) Parse(code []byte) (*tree_sitter.Tree, error) {
	tree := p.parser.Parse(code, nil)
	if tree == nil {
		return nil, errors.New("parse error: tree is nil")
	}

	root := tree.RootNode()
	if root == nil {
		return nil, errors.New("parse error: root node is nil")
	}

	if root.HasError() {
		node := findFirstErrorNode(root)
		if node != nil {
			start := node.StartPosition()
			return tree, fmt.Errorf("syntax error at line %d, column %d", start.Row+1, start.Column+1)
		}

		return tree, errors.New("syntax error detected")
	}

	return tree, nil
}

func (p *Parser) Close() {
	if p.parser == nil {
		return
	}

	p.parser.Close()
}

func (p *Parser) GetLanguage() *tree_sitter.Language {
	return p.language
}

func findFirstErrorNode(node *tree_sitter.Node) *tree_sitter.Node {
	if !node.HasError() {
		return nil
	}

	count := node.ChildCount()

	for i := uint(0); i < count; i++ {
		child := node.Child(i)

		if child.HasError() {
			return findFirstErrorNode(child)
		}
	}

	return node
}
