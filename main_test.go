package main

import (
	"bytes"
	"path/filepath"
	"testing"

	"github.com/google/go-cmp/cmp"
	"golang.org/x/tools/txtar"
)

type testCase struct {
	input []byte

	// A map from -by values to expected outputs.
	outputs map[string]string
}

func readTestCase(t *testing.T, filename string) testCase {
	t.Helper()
	a, err := txtar.ParseFile(filename)
	if err != nil {
		t.Fatal(err)
	}
	tc := testCase{
		input:   a.Comment,
		outputs: make(map[string]string),
	}
	for _, f := range a.Files {
		tc.outputs[f.Name] = string(f.Data)
	}
	if len(tc.outputs) == 0 {
		t.Fatalf("%s: no outputs found", filename)
	}
	return tc
}

func TestCommandRun(t *testing.T) {
	filenames, err := filepath.Glob("testdata/*")
	if err != nil {
		t.Fatal(err)
	}
	for _, filename := range filenames {
		tc := readTestCase(t, filename)
		t.Run(filename, func(t *testing.T) {
			for by, want := range tc.outputs {
				t.Run(by, func(t *testing.T) {
					var buf bytes.Buffer
					cmd := command{
						by: by,
						r:  bytes.NewReader(tc.input),
						w:  &buf,
					}
					if err := cmd.run(); err != nil {
						t.Fatalf("%s: -by %s: error: %v", filename, by, err)
					}
					got := buf.String()
					if diff := cmp.Diff(want, got); diff != "" {
						t.Errorf("%s: -by %s: mismatch (-want +got):\n%s", filename, by, diff)
					}
				})
			}
		})
	}
}
