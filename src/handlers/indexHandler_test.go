package handlers

import (
	"reflect"
	"testing"
)

func TestChunkText(t *testing.T) {
	tests := []struct {
		name         string
		text         string
		maxChunkSize int
		expected     []string
	}{
		{
			name:         "empty text",
			text:         "",
			maxChunkSize: 100,
			expected:     []string{},
		},
		{
			name:         "single paragraph smaller than chunk size",
			text:         "This is a small paragraph.",
			maxChunkSize: 100,
			expected:     []string{"This is a small paragraph."},
		},
		{
			name:         "multiple paragraphs within chunk size",
			text:         "First paragraph.\n\nSecond paragraph.",
			maxChunkSize: 100,
			expected:     []string{"First paragraph.\n\nSecond paragraph."},
		},
		{
			name:         "multiple paragraphs exceeding chunk size",
			text:         "This is the first paragraph with some content.\n\nThis is the second paragraph that should be in a separate chunk because combined they exceed the maximum chunk size.",
			maxChunkSize: 50,
			expected: []string{
				"This is the first paragraph with some content.",
				"This is the second paragraph that should be in a separate chunk because combined they exceed the maximum chunk size.",
			},
		},
		{
			name:         "three paragraphs with mixed sizes",
			text:         "Short.\n\nMedium length paragraph here.\n\nThis is a longer paragraph that contains more text and should be split into its own chunk.",
			maxChunkSize: 60,
			expected: []string{
				"Short.\n\nMedium length paragraph here.",
				"This is a longer paragraph that contains more text and should be split into its own chunk.",
			},
		},
		{
			name:         "text with empty paragraphs",
			text:         "First.\n\n\n\nSecond.",
			maxChunkSize: 100,
			expected:     []string{"First.\n\nSecond."},
		},
		{
			name:         "text with whitespace-only paragraphs",
			text:         "First.\n\n   \n\nSecond.",
			maxChunkSize: 100,
			expected:     []string{"First.\n\nSecond."},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := chunkText(tt.text, tt.maxChunkSize)

			if len(result) == 0 && len(tt.expected) == 0 {
				return
			}

			if !reflect.DeepEqual(result, tt.expected) {
				t.Errorf("chunkText() = %v, want %v", result, tt.expected)
			}
		})
	}
}
