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

type unsupportedByError string

func (e unsupportedByError) Error() string {
	return fmt.Sprintf("unsupported -by value: %s", string(e))
}

type command struct {
	by string
	r  io.Reader
	w  io.Writer
}

func newCommand(by string, r io.Reader, w io.Writer) *command {
	return &command{
		by: by,
		r:  r,
		w:  w,
	}
}

func (c *command) run() error {
	d, err := distribution(c.r, c.by)
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

	switch c.by {
	case "line", "word":
		for _, e := range entries {
			fmt.Fprintf(c.w, "%d\t%s\n", e.count, e.key)
		}
	case "byte":
		for _, e := range entries {
			fmt.Fprintf(c.w, "%d\t%02x\t%s\t%s\n", e.count, e.key, neatByte(e.key[0]), runenames.Name(rune(e.key[0])))
		}
	case "rune":
		for _, e := range entries {
			fmt.Fprintf(c.w, "%d\t%U\t%s\t%s\n", e.count, firstRune(e.key), neatRune(firstRune(e.key)), runenames.Name(firstRune(e.key)))
		}
	default:
		return unsupportedByError(c.by)
	}

	return nil
}

type entry struct {
	key   string
	count int
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

func distribution(r io.Reader, by string) (map[string]int, error) {
	var splitFunc bufio.SplitFunc

	switch by {
	case "line":
		splitFunc = bufio.ScanLines
	case "byte":
		splitFunc = bufio.ScanBytes
	case "rune":
		splitFunc = bufio.ScanRunes
	case "word":
		splitFunc = bufio.ScanWords
	default:
		return nil, unsupportedByError(by)
	}

	m := make(map[string]int)

	s := bufio.NewScanner(r)
	s.Split(splitFunc)

	for s.Scan() {
		m[s.Text()]++
	}

	if err := s.Err(); err != nil {
		return nil, err
	}

	return m, nil
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
