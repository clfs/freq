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
368     calvin          
159     root            
13      _accessoryupdater
8       _cmiodalassistants
7       _locationd      
7       _rmd            
7       _softwareupdate 
5       _coreaudiod     
5       _spotlight      
4       _nsurlsessiond
```

```plaintext
$ cat /bin/ls | freq -by byte | head
138812  0x00
1846    0xff
1625    0x01
1407    0x03
1157    0x48
929     0x5f
899     0x74
883     0x20
757     0x40
720     0x65
```

```plaintext
$ cat /usr/share/locale/zh_CN/LC_TIME | freq -by rune | head
58      "\n"
25      "月"
21      "%"
16      " "
7       "星"
7       "期"
6       "一"
6       "二"
6       "十"
5       "1"
```

```plaintext
$ man tar | freq -by word | head
299     the
128     and
113     is
113     to
81      or
76      of
76      tar
65      archive
60      be
60      in
```
