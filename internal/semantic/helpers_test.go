package semantic_test

import (
	"testing"

	"github.com/Norgate-AV/netlinx-language-server/internal/semantic"
)

func TestSizeOfDataType(t *testing.T) {
	tests := []struct {
		dataType string
		expected uint
	}{
		{"char", 1},
		{"widechar", 2},
		{"integer", 2},
		{"sinteger", 2},
		{"long", 4},
		{"slong", 4},
		{"float", 8},
		{"double", 8},
		{"dev", 6},
		{"devchan", 8},
		{"devlev", 8},
	}

	for _, test := range tests {
		t.Run(test.dataType, func(t *testing.T) {
			result := semantic.SizeOfDataType(test.dataType)
			if result != test.expected {
				t.Errorf("Expected %d, got %d", test.expected, result)
			}
		})
	}
}
