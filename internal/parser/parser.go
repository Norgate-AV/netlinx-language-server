package parser

import (
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

func (p *Parser) Parse(code []byte) *tree_sitter.Tree {
	return p.parser.Parse(code, nil)
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
