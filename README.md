# freq
Report frequency distributions from stdin.

## Setup

Install it:
```text
go install github.com/clfs/freq@latest
```

Update it:
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
$ echo "Hello...\nПривет...\nHello..." | freq
2	Hello...
1	Привет...
```

By bytes:
```text
$ head -c 10 /bin/ls | freq -by byte
4       00      �       <control>
1       01      �       <control>
1       02      �       <control>
1       ba      º       MASCULINE ORDINAL INDICATOR
1       be      ¾       VULGAR FRACTION THREE QUARTERS
1       ca      Ê       LATIN CAPITAL LETTER E WITH CIRCUMFLEX
1       fe      þ       LATIN SMALL LETTER THORN
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

The usual recommendation for showing frequency distributions is to:

1. Call `sort` to sort standard input.
2. Then, call `uniq -c` to deduplicate standard input while adding line counts.
3. Then, call `sort -nr` to reverse sort by the new line counts.

However, sorting standard input is expensive and unnecessary. You shouldn't
need to sort data to get a frequency distribution in the first place.

```text
$ time yes | head -n 123456789 | sort | uniq -c | sort -nr
123456789 y
yes  3.04s user 0.04s system 28% cpu 10.905 total
head -n 123456789  4.87s user 0.06s system 45% cpu 10.902 total
sort  31.87s user 0.70s system 87% cpu 37.152 total
uniq -c  16.17s user 0.05s system 43% cpu 37.152 total
sort -nr  0.00s user 0.00s system 0% cpu 37.151 total
```

In comparison, `freq` is faster and easier to remember.

```text
$ time yes | head -n 123456789 | freq
123456789	y
yes  3.10s user 0.04s system 63% cpu 4.977 total
head -n 123456789  4.93s user 0.04s system 99% cpu 4.971 total
freq  2.20s user 0.06s system 45% cpu 4.970 total
```