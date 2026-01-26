package readers

import (
	"strings"
	"testing"
)

func TestReadHTMLFile(t *testing.T) {
	reader := HTMLReader{}

	tests := []struct {
		name     string
		file     string
		contains []string
		excludes []string
	}{
		{
			name: "blog post",
			file: "../../test-data/html/blog-post.html",
			contains: []string{
				"Understanding Vector Embeddings for RAG",
				"What are Vector Embeddings?",
				"semantically similar",
				"Chunk size",
			},
			excludes: []string{
				"<script>",
				"<style>",
				"<nav>",
				"console.log",
			},
		},
		{
			name: "documentation",
			file: "../../test-data/html/documentation.html",
			contains: []string{
				"Search API Reference",
				"POST /api/search",
				"query",
				"limit",
			},
			excludes: []string{
				"<script>",
				"hljs.highlightAll",
			},
		},
		{
			name: "noisy page",
			file: "../../test-data/html/noisy-page.html",
			contains: []string{
				"Machine Learning Best Practices for Production",
				"Model Versioning",
				"Monitoring and Observability",
				"Data Quality Checks",
			},
			excludes: []string{
				"<script>",
				"displayAd",
				"gtag",
				"cookieConsent",
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			content := reader.Read(tt.file)

			if content == "" {
				t.Errorf("HTMLReader.Read(%s) returned an empty string", tt.file)
				return
			}

			// Check that expected content is present
			for _, expected := range tt.contains {
				if !strings.Contains(content, expected) {
					t.Errorf("HTMLReader.Read(%s) does not contain expected text: %s", tt.file, expected)
				}
			}

			// Check that unwanted content is removed
			for _, unwanted := range tt.excludes {
				if strings.Contains(content, unwanted) {
					t.Errorf("HTMLReader.Read(%s) still contains unwanted text: %s", tt.file, unwanted)
				}
			}
		})
	}
}
