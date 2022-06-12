package main

import (
	"strings"
	"testing"

	"golang.org/x/exp/maps"
	"golang.org/x/exp/slices"
)

func TestLineFreq(t *testing.T) {
	cases := []struct {
		in   string
		want map[string]int
	}{
		{"", map[string]int{}},
		{"a", map[string]int{"a": 1}},
		{"a\nb", map[string]int{"a": 1, "b": 1}},
		{"a b\nc", map[string]int{"a b": 1, "c": 1}},
		{"a\n", map[string]int{"a": 1}},
	}
	for _, c := range cases {
		got := LineFreq(strings.NewReader(c.in))
		if !maps.Equal(got, c.want) {
			t.Errorf("LineFreq(%q) == %v, want %v", c.in, got, c.want)
		}
	}
}

func TestByteFreq(t *testing.T) {

}

func TestRuneFreq(t *testing.T) {

}

func TestSortedEntries_Line(t *testing.T) {
	cases := []struct {
		in   map[string]int
		want []Entry[string]
	}{
		{map[string]int{}, []Entry[string]{}},
		{map[string]int{"a": 1}, []Entry[string]{{"a", 1}}},
		{map[string]int{"a": 1, "b": 2}, []Entry[string]{{"b", 2}, {"a", 1}}},
	}
	for _, c := range cases {
		got := SortedEntries(c.in)
		if !slices.Equal(got, c.want) {
			t.Errorf("ToEntries(%v) == %v, want %v", c.in, got, c.want)
		}
	}
}

func TestSortedEntries_Byte(t *testing.T) {

}

func TestSortedEntries_Rune(t *testing.T) {

}

func FuzzLineFreq_TooManyKeys(f *testing.F) {
	f.Add("foo")
	f.Fuzz(func(t *testing.T, s string) {
		r := strings.NewReader(s)
		m := LineFreq(r)
		if len(m) > len(s) {
			t.Error("too many keys")
		}
	})
}
