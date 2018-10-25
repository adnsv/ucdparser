# Package ucdparser

Downloading or updating:

``` bash
go get -u github.com/adnsv/ucdparser
```

In your code:

``` go
import "github.com/adnsv/ucdparser"
```

## Overview

Package ucdparser is a simple generic parser for
unicode.org reference tables.

GoDoc reference:
[![GoDoc](<a href="https://godoc.org/github.com/adnsv/ucdparser?status.svg">https://godoc.org/github.com/adnsv/ucdparser?status.svg</a>)](<a href="https://godoc.org/github.com/adnsv/ucdparser">https://godoc.org/github.com/adnsv/ucdparser</a>)

## Index

- [Package Examples](#pkgeg)
- [func Parse](#001)
- [func SplitLine](#002)
- [type Line](#003)
  - [func (*Line) Int](#004)
  - [func (*Line) Rune](#005)
  - [func (*Line) RuneRange](#006)
  - [func (*Line) Runes](#007)
  - [func (*Line) String](#008)
  - [func (*Line) Uint](#009)

## <a name='eg'>Package Examples</a>

## Package Example

``` go
package ucdparser

import (
    "fmt"
    "log"
    "net/http"
)

const root = "https://www.unicode.org/Public/UCD/latest/ucd/"

func Example() {
    url := root + "NamesList.txt"

    fmt.Printf("Fetching '%s'...", url)
    resp, err := http.Get(url)
    if err != nil {
        log.Fatalf("HTTP GET: %v", err)
    }
    if resp.StatusCode != 200 {
        log.Fatalf("Bad GET status for %q: %q", url, resp.Status)
    }
    defer resp.Body.Close()

    Parse(resp.Body, func(ln *Line) {
        if len(ln.Fields) == 0 {
            fmt.Println("#", ln.Comment)
        } else {
            s := ""
            for i, f := range ln.Fields {
                if i == 0 {
                    l, h := ln.RuneRange(0)
                    s += fmt.Sprintf("|%#x..%#x", l, h)
                } else {
                    s += "|" + string(f)
                }
            }
            if ln.Comment != "" {
                s += "#" + ln.Comment
            }
            fmt.Println(s)
        }
    })
}
```

## func <a name='001'>Parse</a>

``` go
func Parse(r io.Reader, f func(ln *Line)) error
```

Parse calls f for each non-empty parsed line.

## func <a name='002'>SplitLine</a>

``` go
func SplitLine(s string, line *Line) bool
```

SplitLine decomposes string into the line structure.
Returns false if the string is empty.

## type <a name='003'>Line</a>

``` go
type Line struct {
    Index   int
    Fields  []string
    Part    string
    Comment string
    Err     error
}
```

Line is used as a callback parameter in the Parse func.

### method (*Line) <a name='004'>Int</a>

``` go
func (ln *Line) Int(i int) (v int)
```

Int parses and returns field i as an integer value.

### method (*Line) <a name='005'>Rune</a>

``` go
func (ln *Line) Rune(i int) (r rune)
```

Rune reads i-th field as a rune.

### method (*Line) <a name='006'>RuneRange</a>

``` go
func (ln *Line) RuneRange(i int) (first, last rune)
```

RuneRange interprets and returns field i as a range of runes.

### method (*Line) <a name='007'>Runes</a>

``` go
func (ln *Line) Runes(i int) (runes []rune)
```

Runes interprets and returns field i as a sequence of runes.

### method (*Line) <a name='008'>String</a>

``` go
func (ln *Line) String(i int) (v string)
```

String returns field i as a string value.

### method (*Line) <a name='009'>Uint</a>

``` go
func (ln *Line) Uint(i int) (v uint)
```

Uint parses and returns field i as an unsigned integer value.



