package main

import (
	"bufio"
	"cmp"
	"flag"
	"fmt"
	"log"
	"os"
	"slices"
)

type mode struct {
	splitFunc bufio.SplitFunc
	format    string
}

var supportedModes = map[string]mode{
	"line": {bufio.ScanLines, "%d\t%s\n"},
	"byte": {bufio.ScanBytes, "%d\t0x%02x\n"},
	"rune": {bufio.ScanRunes, "%d\t%q\n"},
	"word": {bufio.ScanWords, "%d\t%s\n"},
}

func main() {
	log.SetFlags(0)

	byFlag := flag.String("by", "line", "line, byte, rune, or word")
	flag.Parse()

	mode, ok := supportedModes[*byFlag]
	if !ok {
		log.Fatalf("error: unsupported mode %q", *byFlag)
	}

	distribution := make(map[string]int)

	s := bufio.NewScanner(os.Stdin)
	s.Split(mode.splitFunc)
	for s.Scan() {
		distribution[s.Text()]++
	}
	if err := s.Err(); err != nil {
		log.Fatalf("error: %v", err)
	}

	type entry struct {
		key   string
		count int
	}

	entries := make([]entry, 0, len(distribution))
	for k, v := range distribution {
		entries = append(entries, entry{k, v})
	}

	// 3 x
	// 2 a
	// 2 b
	// ...
	slices.SortFunc(entries, func(a, b entry) int {
		c := cmp.Compare(b.count, a.count)
		if c == 0 {
			return cmp.Compare(a.key, b.key)
		}
		return c
	})

	for _, e := range entries {
		fmt.Printf(mode.format, e.count, e.key)
	}
}
