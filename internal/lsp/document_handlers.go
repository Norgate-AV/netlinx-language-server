package lsp

import (
	"context"
	"encoding/json"

	"github.com/Norgate-AV/netlinx-language-server/internal/jsonrpc"
)

func (h *Handler) handleTextDocumentDidOpen(ctx context.Context, conn *jsonrpc.Conn, req *jsonrpc.Request) {
	var params DidOpenTextDocumentParams
	if err := json.Unmarshal(*req.Params, &params); err != nil {
		h.logger.Error("Failed to parse didOpen params", err)
		return
	}

	doc := &Document{
		URI:        params.TextDocument.URI,
		LanguageID: params.TextDocument.LanguageID,
		Version:    params.TextDocument.Version,
		Content:    params.TextDocument.Text,
	}

	h.server.docs.Add(doc)
	h.logger.Debug("Document opened", "uri", doc.URI)

	// Run diagnostics on the newly opened document
	// This would be where you'd implement NetLinx-specific syntax checking
}

func (h *Handler) handleTextDocumentDidChange(ctx context.Context, conn *jsonrpc.Conn, req *jsonrpc.Request) {
	var params DidChangeTextDocumentParams
	if err := json.Unmarshal(*req.Params, &params); err != nil {
		h.logger.Error("Failed to parse didChange params", err)
		return
	}

	if len(params.ContentChanges) == 0 {
		return
	}

	// For full sync, just replace the entire document content
	newContent := params.ContentChanges[0].Text
	uri := params.TextDocument.URI
	version := params.TextDocument.Version

	if ok := h.server.docs.Update(uri, version, newContent); !ok {
		h.logger.Warn("Failed to update document - not found", "uri", uri)
		return
	}

	h.logger.Debug("Document updated", "uri", uri, "version", version)

	// Run diagnostics on the updated document
}

func (h *Handler) handleTextDocumentDidClose(ctx context.Context, conn *jsonrpc.Conn, req *jsonrpc.Request) {
	var params DidCloseTextDocumentParams
	if err := json.Unmarshal(*req.Params, &params); err != nil {
		h.logger.Error("Failed to parse didClose params", err)
		return
	}

	uri := params.TextDocument.URI
	h.server.docs.Remove(uri)
	h.logger.Debug("Document closed", "uri", uri)
}
