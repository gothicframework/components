package gothicComponents

import (
	"bytes"
	"context"
	"strings"
	"testing"
)

// TestOptimizedImage_AltURLEncoding asserts that the hx-get URL query parameter
// `alt` is properly URL-encoded so that special characters in alt text do not
// break the query string — the image server receives the decoded original value.
// The hx-get path is reached when IsFirstLoad=true and Priority=false (the
// blur-up placeholder → HTMX full-resolve flow).
func TestOptimizedImage_AltURLEncoding(t *testing.T) {
	props := OptimizedImageProps{
		IsFirstLoad:  true,
		Priority:     false,
		ImgName:      "photo",
		ImgExtension: "jpeg",
		Alt:          `a & b <q>`,
	}
	var buf bytes.Buffer
	if err := OptimizedImage(props).Render(context.Background(), &buf); err != nil {
		t.Fatalf("OptimizedImage.Render: %v", err)
	}
	html := buf.String()

	// URL-encoding rules: space → +, & → %26, < → %3C, > → %3E
	if !strings.Contains(html, "alt=a+%26+b+%3Cq%3E") {
		t.Errorf("hx-get should contain URL-encoded alt text `a+%%26+b+%%3Cq%%3E`, got:\n%s", html)
	}
}

// TestOptimizedImage_AltRenderedPriority asserts the <img> alt attribute in the
// priority (above-the-fold) path contains the literal (non-encoded) alt text —
// only the hx-get URL parameter is encoded. The priority path is reached when
// IsFirstLoad=true and Priority=true.
func TestOptimizedImage_AltRenderedPriority(t *testing.T) {
	props := OptimizedImageProps{
		IsFirstLoad:  true,
		Priority:     true,
		ImgName:      "photo",
		ImgExtension: "jpeg",
		Alt:          `a & b <q>`,
	}
	var buf bytes.Buffer
	if err := OptimizedImage(props).Render(context.Background(), &buf); err != nil {
		t.Fatalf("OptimizedImage.Render: %v", err)
	}
	html := buf.String()

	// The HTML alt attribute should contain the literal text (templ auto-escapes
	// the & to &amp; etc. for HTML safety, but does NOT URL-encode it).
	if !strings.Contains(html, "alt=\"a &amp; b &lt;q&gt;\"") {
		t.Errorf("expected alt attribute with HTML-escaped value, got:\n%s", html)
	}
}

// TestOptimizedImage_AltRenderedNonFirstLoad asserts the non-first-load path
// (plain <img>, no HTMX) also renders the alt text correctly.
func TestOptimizedImage_AltRenderedNonFirstLoad(t *testing.T) {
	props := OptimizedImageProps{
		IsFirstLoad:  false,
		ImgName:      "photo",
		ImgExtension: "jpeg",
		Alt:          `hello world`,
	}
	var buf bytes.Buffer
	if err := OptimizedImage(props).Render(context.Background(), &buf); err != nil {
		t.Fatalf("OptimizedImage.Render: %v", err)
	}
	html := buf.String()

	if !strings.Contains(html, `alt="hello world"`) {
		t.Errorf("expected alt text `hello world`, got:\n%s", html)
	}
}
