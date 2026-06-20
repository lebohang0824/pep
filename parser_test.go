package main

import (
	"encoding/json"
	"os"
	"testing"

	"pep/internal/lexer"
	"pep/internal/parser"
	"pep/internal/semantic"
	"pep/internal/generator"
)

func TestParseSample(t *testing.T) {
	input, err := os.ReadFile("testdata/valid.pep")
	if err != nil {
		t.Fatal(err)
	}

	l := lexer.NewLexer(string(input))
	p := parser.NewParser(l)
	app, errs := p.Parse()
	if len(errs) > 0 {
		for _, e := range errs {
			t.Logf("Parse error: %s", e)
		}
		t.Fatal("parsing failed")
	}

	semErrs := semantic.Analyze(app)
	if len(semErrs) > 0 {
		for _, e := range semErrs {
			t.Logf("Semantic error: %s", e)
		}
		t.Fatal("semantic analysis failed")
	}

	data, err := generator.Generate(app, true)
	if err != nil {
		t.Fatal(err)
	}

	expected, err := os.ReadFile("testdata/expected.json")
	if err != nil {
		t.Fatal(err)
	}

	var got, want interface{}
	if err := json.Unmarshal(data, &got); err != nil {
		t.Fatal(err)
	}
	if err := json.Unmarshal(expected, &want); err != nil {
		t.Fatal(err)
	}

	gotJSON, _ := json.MarshalIndent(got, "", "  ")
	wantJSON, _ := json.MarshalIndent(want, "", "  ")

	if string(gotJSON) != string(wantJSON) {
		t.Errorf("output mismatch:\ngot:\n%s\n\nwant:\n%s", string(gotJSON), string(wantJSON))
	}
}
