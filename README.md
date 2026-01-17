# freq

Compute frequency distributions over standard input.

Requirements:

- [Go](https://go.dev)

Install or update:

```plaintext
go install github.com/clfs/freq@latest
```

Uninstall:

```bash
rm -i $(which freq)
```

Usage:

```plaintext
$ freq -h
Usage of freq:
  -by string
        line, byte, rune, or word (default "line")
```

Examples:

```plaintext
$ ps -eo user | freq | head
660     calvin          
193     root            
18      _accessoryupdater
13      _rmd            
12      _cmiodalassistants
8       _locationd      
6       _coreaudiod     
5       _applepay       
5       _nsurlsessiond  
4       _driverkit
```

```plaintext
$ cat /bin/ls | freq -by byte | head
107087  0x00
1843    0xff
1628    0x01
1142    0x48
1125    0x03
925     0x5f
882     0x74
736     0x40
731     0x65
666     0x02
```

```plaintext
$ cat /usr/share/locale/zh_CN/LC_TIME | freq -by rune | head
58      "\n"
37      "月"
21      "%"
16      " "
15      "1"
7       "星"
7       "期"
6       "2"
3       "/"
3       "0"
```

```plaintext
$ man tar | freq -by word | head
326     the
129     and
121     to
119     is
92      or
86      of
76      tar
69      be
65      archive
62      a
```
