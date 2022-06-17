# freq
Compute frequency distributions from stdin.

## Setup

Install or update it:
```text
go install github.com/clfs/freq@latest
```

Uninstall it:
```text
rm $(which freq)
```

## Examples

Usage:
```text
$ freq -h
Usage of freq:
  -by string
        token type: line, byte, or rune (default "line")
```

By lines:
```text
$ awk '{print $1}' go.sum | freq
2	golang.org/x/exp
2	golang.org/x/text
```

By bytes:
```text
$ head -c 10 /bin/ls | freq -by byte
4	00	�	<control>
1	01	�	<control>
1	02	�	<control>
1	ba	º	MASCULINE ORDINAL INDICATOR
1	be	¾	VULGAR FRACTION THREE QUARTERS
1	ca	Ê	LATIN CAPITAL LETTER E WITH CIRCUMFLEX
1	fe	þ	LATIN SMALL LETTER THORN
```

By runes:
```text
$ echo -n "Hello? Привет? Hello?" | freq -by rune
4	U+006C	l	LATIN SMALL LETTER L
3	U+003F	?	QUESTION MARK
2	U+0020	 	SPACE
2	U+0048	H	LATIN CAPITAL LETTER H
2	U+0065	e	LATIN SMALL LETTER E
2	U+006F	o	LATIN SMALL LETTER O
1	U+041F	П	CYRILLIC CAPITAL LETTER PE
1	U+0432	в	CYRILLIC SMALL LETTER VE
1	U+0435	е	CYRILLIC SMALL LETTER IE
1	U+0438	и	CYRILLIC SMALL LETTER I
1	U+0440	р	CYRILLIC SMALL LETTER ER
1	U+0442	т	CYRILLIC SMALL LETTER TE
```

## Why write this?

This is the usual recommendation for computing frequency distributions:

```bash
cat example.txt | sort | uniq -c | sort -nr
```

However, it's not necessary to sort the input. Sorting is slow and uses a lot of
memory. In comparison, `freq` is faster and has more features.

Some people also recommend using `awk`, like this:

```bash
awk ' { tot[$0]++ } END { for (i in tot) print tot[i],i } ' example.txt | sort
```

It doesn't use any sorting, but it seems hard to remember.