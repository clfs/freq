package main

import (
	"strings"
	"testing"

	"github.com/google/go-cmp/cmp"
)

func TestPretty(t *testing.T) {
	tests := []struct {
		by   string
		in   entry
		want string
	}{
		{
			"line",
			entry{key: "hello", count: 3},
			"3\thello",
		},
		{
			"word",
			entry{key: "world", count: 2},
			"2\tworld",
		},
		{
			"byte",
			entry{key: "A", count: 5},
			"5\t41\tA\tLATIN CAPITAL LETTER A",
		},
		{
			"byte",
			entry{key: "\x00", count: 1},
			"1\t00\t�\t<control>",
		},
		{
			"rune",
			entry{key: "B", count: 4},
			"4\tU+0042\tB\tLATIN CAPITAL LETTER B",
		},
		{
			"rune",
			entry{key: "\x00", count: 2},
			"2\tU+0000\t�\t<control>",
		},
		{
			"rune",
			entry{key: "π", count: 3},
			"3\tU+03C0\tπ\tGREEK SMALL LETTER PI",
		},
		{
			"rune",
			entry{key: "😀", count: 1},
			"1\tU+1F600\t😀\tGRINNING FACE",
		},
	}
	for _, tt := range tests {
		got := prettyFuncs[tt.by](tt.in)
		if diff := cmp.Diff(tt.want, got); diff != "" {
			t.Errorf(
				"key=%q, count=%d, by=%q: mismatch (-want +got):\n%s",
				tt.in.key, tt.in.count, tt.by, diff,
			)
		}
	}
}

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
			"",
			"word",
			map[string]int{},
		},
		{
			"",
			"byte",
			map[string]int{},
		},
		{
			"",
			"rune",
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
			"あ",
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
			t.Errorf("in=%q, by=%q: error: %v", tt.in, tt.by, err)
		}
		if diff := cmp.Diff(tt.want, got); diff != "" {
			t.Errorf("in=%q, by=%q: mismatch (-want +got):\n%s", tt.in, tt.by, diff)
		}
	}
}
