package main

import (
	"bytes"
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

func TestDistribution(t *testing.T) {
	tests := []struct {
		in   string
		by   string
		want map[string]int
	}{
		{"a a a", "word", map[string]int{"a": 3}},
		{"aa a", "word", map[string]int{"aa": 1, "a": 1}},
		{"a\nb\nc", "line", map[string]int{"a": 1, "b": 1, "c": 1}},
		{"abc", "byte", map[string]int{"a": 1, "b": 1, "c": 1}},
		{"aあい", "rune", map[string]int{"a": 1, "あ": 1, "い": 1}},
		{"", "line", map[string]int{}},
		{"hello, world!", "word", map[string]int{"hello,": 1, "world!": 1}},
		{"a\tb  c", "word", map[string]int{"a": 1, "b": 1, "c": 1}},
		{"a\n\nb\n", "line", map[string]int{"a": 1, "": 1, "b": 1}},
		{"aaab", "byte", map[string]int{"a": 3, "b": 1}},
		{"πππ", "rune", map[string]int{"π": 3}},
		{"", "word", map[string]int{}},
		{"well-being well", "word", map[string]int{"well-being": 1, "well": 1}},
		{"aa, aa.", "word", map[string]int{"aa,": 1, "aa.": 1}},
		{"", "rune", map[string]int{}},
		{"a\u00A0b", "word", map[string]int{"a": 1, "b": 1}},
		{"\x00a", "byte", map[string]int{"\x00": 1, "a": 1}},
		{"e\u0301", "rune", map[string]int{"e": 1, "\u0301": 1}},
		{"👩‍👩‍👧‍👦", "rune", map[string]int{"👩": 2, "\u200d": 3, "👧": 1, "👦": 1}},
		{"\xe3\x81\x82", "byte", map[string]int{"\xe3": 1, "\x81": 1, "\x82": 1}},
	}

	for _, tt := range tests {
		r := strings.NewReader(tt.in)
		got, err := distribution(r, splitFuncs[tt.by])
		if err != nil {
			t.Errorf("%q: -by %s: error: %v", tt.in, tt.by, err)
		}
		if diff := cmp.Diff(tt.want, got); diff != "" {
			t.Errorf("%q: -by %s: mismatch (-want +got):\n%s", tt.in, tt.by, diff)
		}
	}
}
