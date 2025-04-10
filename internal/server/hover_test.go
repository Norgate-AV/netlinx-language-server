package server

import (
	"testing"

	"github.com/Norgate-AV/netlinx-language-server/internal/analysis"
	"github.com/Norgate-AV/netlinx-language-server/internal/logger"
	"github.com/Norgate-AV/netlinx-language-server/internal/lsp"
)

func TestHover(t *testing.T) {
	log := logger.NewStdLogger()

	state := analysis.NewState()
	state.AddDocument("file:///test.axs", "PROGRAM_NAME='Test'\nDEFINE_VARIABLE\nINTEGER x")

	srv := NewServer(log, state)

	hover, err := srv.GetHoverInfo("file:///test.axs", lsp.Position{Line: 2, Character: 8})
	if err != nil {
		t.Fatalf("Unexpected error: %v", err)
	}

	if hover == nil || hover.Contents.Value == "" {
		t.Fatal("Expected non-empty hover content")
	}

	t.Logf("Hover content: %s", hover.Contents.Value)
}
