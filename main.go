package main

import (
	"bufio"
	"cmp"
	"flag"
	"fmt"
	"io"
	"log"
	"os"
	"slices"
	"unicode"
	"unicode/utf8"

	"golang.org/x/text/unicode/runenames"
)

var splitFuncs = map[string]bufio.SplitFunc{
	"line": bufio.ScanLines,
	"byte": bufio.ScanBytes,
	"rune": bufio.ScanRunes,
	"word": bufio.ScanWords,
}

var prettyFuncs = map[string]func(entry) string{
	"line": prettyLine,
	"word": prettyWord,
	"byte": prettyByte,
	"rune": prettyRune,
}

func prettyLine(e entry) string {
	return fmt.Sprintf("%d\t%s", e.count, e.key)
}

func prettyWord(e entry) string {
	return fmt.Sprintf("%d\t%s", e.count, e.key)
}

func prettyByte(e entry) string {
	return fmt.Sprintf("%d\t%02x\t%s\t%s", e.count, e.key, neatByte(e.key[0]), runenames.Name(rune(e.key[0])))
}

func prettyRune(e entry) string {
	return fmt.Sprintf("%d\t%U\t%s\t%s", e.count, firstRune(e.key), neatRune(firstRune(e.key)), runenames.Name(firstRune(e.key)))
}

type unsupportedByError struct {
	by string
}

func (e unsupportedByError) Error() string {
	return fmt.Sprintf("unsupported -by value: %s", e.by)
}

type entry struct {
	key   string
	count int
}

type command struct {
	by string
	r  io.Reader
	w  io.Writer
}

func newCommand(by string, r io.Reader, w io.Writer) *command {
	return &command{by: by, r: r, w: w}
}

func (c *command) run() error {
	fn, ok := splitFuncs[c.by]
	if !ok {
		return unsupportedByError{by: c.by}
	}

	d, err := distribution(c.r, fn)
	if err != nil {
		return err
	}

	var entries []entry
	for k, v := range d {
		entries = append(entries, entry{key: k, count: v})
	}

	// 3, x
	// 2, a
	// 2, b
	// ...
	slices.SortFunc(entries, func(a, b entry) int {
		if a.count != b.count {
			return cmp.Compare(b.count, a.count)
		}
		return cmp.Compare(a.key, b.key)
	})

	pretty, ok := prettyFuncs[c.by]
	if !ok {
		return unsupportedByError{by: c.by}
	}

	for _, e := range entries {
		fmt.Fprintln(c.w, pretty(e))
	}

	return nil
}

func main() {
	log.SetFlags(0)

	byFlag := flag.String("by", "line", "line, byte, rune, or word")
	flag.Parse()

	cmd := newCommand(*byFlag, os.Stdin, os.Stdout)

	if err := cmd.run(); err != nil {
		log.Fatal(err)
	}
}

func distribution(r io.Reader, split bufio.SplitFunc) (map[string]int, error) {
	m := make(map[string]int)
	s := bufio.NewScanner(r)
	s.Split(split)
	for s.Scan() {
		m[s.Text()]++
	}
	return m, s.Err()
}

func firstRune(s string) rune {
	if s == "" {
		return utf8.RuneError
	}
	r, _ := utf8.DecodeRuneInString(s)
	return r
}

func neatRune(r rune) string {
	if !unicode.IsPrint(r) {
		return string(utf8.RuneError)
	}
	return string(r)
}

func neatByte(b byte) string {
	return neatRune(rune(b))
}
