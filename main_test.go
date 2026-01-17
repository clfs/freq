package main

import (
	"bufio"
	"bytes"
	"errors"
	"path/filepath"
	"strings"
	"testing"

	"github.com/google/go-cmp/cmp"
	"golang.org/x/tools/txtar"
)

type testCase struct {
	input []byte

	// A map from -by values to expected wants.
	wants map[string]string
}

func readTestCase(t *testing.T, filename string) testCase {
	t.Helper()
	a, err := txtar.ParseFile(filename)
	if err != nil {
		t.Fatal(err)
	}
	tc := testCase{
		input: a.Comment,
		wants: make(map[string]string),
	}
	for _, f := range a.Files {
		tc.wants[f.Name] = string(f.Data)
	}
	if len(tc.wants) == 0 {
		t.Fatalf("%s: no wanted outputs found", filename)
	}
	return tc
}

func TestCommand_Run(t *testing.T) {
	filenames, err := filepath.Glob("testdata/*")
	if err != nil {
		t.Fatal(err)
	}

	for _, filename := range filenames {
		tc := readTestCase(t, filename)
		t.Run(filename, func(t *testing.T) {
			for by, want := range tc.wants {
				t.Run(by, func(t *testing.T) {
					var (
						r = bytes.NewReader(tc.input)
						w = new(bytes.Buffer)
					)
					cmd := newCommand(by, r, w)
					if err := cmd.run(); err != nil {
						t.Fatalf("%s: -by %s: error: %v", filename, by, err)
					}
					got := w.String()
					if diff := cmp.Diff(want, got); diff != "" {
						t.Errorf("%s: -by %s: mismatch (-want +got):\n%s", filename, by, diff)
					}
				})
			}
		})
	}
}

func TestDistribution_TooLong(t *testing.T) {
	in := strings.Repeat("a", bufio.MaxScanTokenSize+1)
	_, err := distribution(strings.NewReader(in), "word")
	if !errors.Is(err, bufio.ErrTooLong) {
		t.Errorf("want bufio.ErrTooLong, got %v", err)
	}
}

func TestDistribution_UnsupportedBy(t *testing.T) {
	_, err := distribution(strings.NewReader("a"), "sentence")
	if _, ok := err.(unsupportedByError); !ok {
		t.Errorf("want unsupportedByError, got %v", err)
	}
}

func TestDistribution(t *testing.T) {
	tests := []struct {
		in   string
		by   string
		want map[string]int
	}{
		{"a a a", "word", map[string]int{"a": 3}},
		{"aa a", "word", map[string]int{"aa": 1, "a": 1}},
	}

	for _, tt := range tests {
		r := strings.NewReader(tt.in)
		got, err := distribution(r, tt.by)
		if err != nil {
			t.Errorf("%q: -by %s: error: %v", tt.in, tt.by, err)
		}
		if diff := cmp.Diff(tt.want, got); diff != "" {
			t.Errorf("%q: -by %s: mismatch (-want +got):\n%s", tt.in, tt.by, diff)
		}
	}
}
