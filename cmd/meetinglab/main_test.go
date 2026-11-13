package main

import (
	"bytes"
	"encoding/json"
	"os"
	"strings"
	"testing"

	"example.com/rfc-go-labs/lab"
)

func TestExample(t *testing.T) {
	input, err := os.ReadFile("../../examples/rfc21.json")
	if err != nil {
		t.Fatal(err)
	}
	var output bytes.Buffer
	if err := run(nil, bytes.NewReader(input), &output); err != nil {
		t.Fatal(err)
	}
	var got lab.Report
	if err := json.Unmarshal(output.Bytes(), &got); err != nil {
		t.Fatal(err)
	}
	if len(got.Decisions) != 0 || len(got.Questions) != 1 || len(got.Notes) != 3 ||
		len(got.Actions) != 1 || got.Actions[0].Owner != "George Gregg" || len(got.Actions[0].Findings) != 1 ||
		got.Actions[0].Findings[0] != "Missing deadline" {
		t.Fatalf("unexpected RFC 21 report: %+v", got)
	}
}

func TestCLIRejectsUnknownAndTrailingInput(t *testing.T) {
	for _, input := range []string{
		`{"extra":true}`,
		`{} {}`,
		strings.Repeat(" ", (64<<10)+1),
	} {
		var output bytes.Buffer
		if err := run(nil, strings.NewReader(input), &output); err == nil || output.Len() != 0 {
			t.Fatalf("invalid input of %d bytes produced %d bytes and error %v", len(input), output.Len(), err)
		}
	}
}

func TestHTMLUsesEscapedText(t *testing.T) {
	const input = `{"title":"<script>","as_of":"2026-11-13","items":[{"id":"x","kind":"note","summary":"A & B","source":"Minutes"}]}`
	var output bytes.Buffer
	if err := run([]string{"-html"}, strings.NewReader(input), &output); err != nil {
		t.Fatal(err)
	}
	if strings.Contains(output.String(), "<script>") || !strings.Contains(output.String(), "&lt;script&gt;") ||
		!strings.Contains(output.String(), "A &amp; B") {
		t.Fatalf("expected escaped text in report: %s", output.String())
	}
}
