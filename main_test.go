package main

import (
	"strings"
	"testing"

	"github.com/google/go-cmp/cmp"
)

func TestDistribution(t *testing.T) {
	tests := []struct {
		in   string
		by   string
		want map[string]int
	}{
		{
			"a a a",
			"word",
			map[string]int{
				"a": 3,
			},
		},
		{
			"aa a",
			"word",
			map[string]int{
				"aa": 1,
				"a":  1,
			},
		},
		{
			"a\nb\nc",
			"line",
			map[string]int{
				"a": 1,
				"b": 1,
				"c": 1,
			},
		},
		{
			"abc",
			"byte",
			map[string]int{
				"a": 1,
				"b": 1,
				"c": 1,
			},
		},
		{
			"aあい",
			"rune",
			map[string]int{
				"a": 1,
				"あ": 1,
				"い": 1,
			},
		},
		{
			"",
			"line",
			map[string]int{},
		},
		{
			"hello, world!",
			"word",
			map[string]int{
				"hello,": 1,
				"world!": 1,
			},
		},
		{
			"a\tb  c",
			"word",
			map[string]int{
				"a": 1,
				"b": 1,
				"c": 1,
			},
		},
		{
			"a\n\nb\n",
			"line",
			map[string]int{
				"a": 1,
				"":  1,
				"b": 1,
			},
		},
		{
			"aaab",
			"byte",
			map[string]int{
				"a": 3,
				"b": 1,
			},
		},
		{
			"πππ",
			"rune",
			map[string]int{
				"π": 3,
			},
		},
		{
			"",
			"word",
			map[string]int{},
		},
		{
			"well-being well",
			"word",
			map[string]int{
				"well-being": 1,
				"well":       1,
			},
		},
		{
			"aa, aa.",
			"word",
			map[string]int{
				"aa,": 1,
				"aa.": 1,
			},
		},
		{
			"",
			"rune",
			map[string]int{},
		},
		{
			"a\u00A0b",
			"word",
			map[string]int{
				"a": 1,
				"b": 1,
			},
		},
		{
			"\x00a",
			"byte",
			map[string]int{
				"\x00": 1,
				"a":    1,
			},
		},
		{
			"e\u0301",
			"rune",
			map[string]int{
				"e":      1,
				"\u0301": 1,
			},
		},
		{
			"👩‍👩‍👧‍👦",
			"rune",
			map[string]int{
				"👩":      2,
				"\u200d": 3,
				"👧":      1,
				"👦":      1,
			},
		},
		{
			"\xe3\x81\x82",
			"byte",
			map[string]int{
				"\xe3": 1,
				"\x81": 1,
				"\x82": 1,
			},
		},
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
