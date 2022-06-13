package main

import (
	"bufio"
	"flag"
	"fmt"
	"io"
	"log"
	"os"
	"unicode"
	"unicode/utf8"

	"golang.org/x/exp/constraints"
	"golang.org/x/exp/slices"
	"golang.org/x/text/unicode/runenames"
)

func main() {
	byFlag := flag.String("by", "line", "token type: line, byte, or rune")
	flag.Parse()

	switch *byFlag {
	case "line":
		for _, e := range SortedEntries(LineFreq(os.Stdin)) {
			fmt.Printf("%d\t%s\n", e.Frequency, e.Token)
		}
	case "byte":
		for _, e := range SortedEntries(ByteFreq(os.Stdin)) {
			fmt.Printf("%d\t%02x\t%s\t%s\n", e.Frequency, e.Token, neatByte(e.Token), runenames.Name(rune(e.Token)))
		}
	case "rune":
		for _, e := range SortedEntries(RuneFreq(os.Stdin)) {
			fmt.Printf("%d\t%U\t%s\t%s\n", e.Frequency, e.Token, neatRune(e.Token), runenames.Name(e.Token))
		}
	}
}

func neatRune(r rune) string {
	x := utf8.RuneError
	if unicode.IsPrint(r) {
		x = r
	}
	return string(x)
}

func neatByte(b byte) string {
	x := utf8.RuneError
	if r := rune(b); unicode.IsPrint(r) {
		x = r
	}
	return string(x)
}

func LineFreq(r io.Reader) map[string]int {
	freq := make(map[string]int)
	s := bufio.NewScanner(r)
	for s.Scan() {
		freq[s.Text()]++ // max token length is 65535 bytes
	}
	if err := s.Err(); err != nil {
		log.Fatal(err) // handle this eventually
	}
	return freq
}

func ByteFreq(r io.Reader) map[byte]int {
	freq := make(map[byte]int)
	br := bufio.NewReader(r)
	for {
		line, err := br.ReadByte()
		if err != nil {
			break
		}
		freq[line]++
	}
	return freq
}

func RuneFreq(r io.Reader) map[rune]int {
	freq := make(map[rune]int)
	br := bufio.NewReader(r)
	for {
		rn, _, err := br.ReadRune()
		if err != nil {
			break
		}
		freq[rn]++
	}
	return freq
}

type Entry[T constraints.Ordered] struct {
	Token     T
	Frequency int
}

// more returns true if a > b.
func more[T constraints.Ordered](a, b Entry[T]) bool {
	return a.Frequency > b.Frequency || (a.Frequency == b.Frequency && a.Token < b.Token)
}

func SortedEntries[K constraints.Ordered](m map[K]int) []Entry[K] {
	entries := make([]Entry[K], 0)
	for k, v := range m {
		entries = append(entries, Entry[K]{k, v})
	}
	slices.SortFunc(entries, more[K])
	return entries
}
