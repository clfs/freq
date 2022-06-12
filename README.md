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
4	0x00	�	<control>
1	0x01	�	<control>
1	0x02	�	<control>
1	0xBA	º	MASCULINE ORDINAL INDICATOR
1	0xBE	¾	VULGAR FRACTION THREE QUARTERS
1	0xCA	Ê	LATIN CAPITAL LETTER E WITH CIRCUMFLEX
1	0xFE	þ	LATIN SMALL LETTER THORN
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