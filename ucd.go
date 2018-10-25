package ucdparser

import (
	"bufio"
	"fmt"
	"io"
	"strconv"
	"strings"
)

// Line is used as a callback parameter in the Parse func.
type Line struct {
	Index   int
	Fields  []string
	Part    string
	Comment string
	Err     error
}

// Rune reads i-th field as a rune.
func (ln *Line) Rune(i int) (r rune) {
	if ln.Err != nil {
		return
	}
	if i >= len(ln.Fields) {
		ln.Err = fmt.Errorf("invalid field index %d", i)
	} else {
		r, ln.Err = parseRune(ln.Fields[i])
	}
	return
}

// Runes interprets and returns field i as a sequence of runes.
func (ln *Line) Runes(i int) (runes []rune) {
	if ln.Err != nil {
		return
	}
	if i >= len(ln.Fields) {
		ln.Err = fmt.Errorf("invalid field index %d", i)
	} else {
		for _, s := range strings.Split(ln.Fields[i], " ") {
			if s == "" {
				continue
			}
			r, err := parseRune(s)
			if err != nil {
				ln.Err = err
				return
			}
			runes = append(runes, r)
		}
	}
	return
}

// RuneRange interprets and returns field i as a range of runes.
func (ln *Line) RuneRange(i int) (first, last rune) {
	if ln.Err != nil {
		return
	}
	if i >= len(ln.Fields) {
		ln.Err = fmt.Errorf("invalid field index %d", i)
	} else {
		s := ln.Fields[i]
		if p := strings.Index(s, ".."); p >= 0 {
			first, ln.Err = parseRune(s[:p])
			if ln.Err == nil {
				last, ln.Err = parseRune(s[p+2:])
			}
		} else {
			first, ln.Err = parseRune(s)
			last = first
		}
	}
	return
}

// String returns field i as a string value.
func (ln *Line) String(i int) (v string) {
	if ln.Err != nil {
		return
	}
	if i >= len(ln.Fields) {
		ln.Err = fmt.Errorf("invalid field index %d", i)
	} else {
		v = ln.Fields[i]
	}
	return
}

// Int parses and returns field i as an integer value.
func (ln *Line) Int(i int) (v int) {
	if ln.Err != nil {
		return
	}
	if i >= len(ln.Fields) {
		ln.Err = fmt.Errorf("invalid field index %d", i)
	} else {
		var x int64
		x, ln.Err = strconv.ParseInt(ln.Fields[i], 10, 64)
		return int(x)
	}
	return
}

// Uint parses and returns field i as an unsigned integer value.
func (ln *Line) Uint(i int) (v uint) {
	if ln.Err != nil {
		return
	}
	if i >= len(ln.Fields) {
		ln.Err = fmt.Errorf("invalid field index %d", i)
	} else {
		var x uint64
		x, ln.Err = strconv.ParseUint(ln.Fields[i], 10, 64)
		return uint(x)
	}
	return
}

// Parse calls f for each non-empty parsed line.
func Parse(r io.Reader, f func(ln *Line)) error {
	scanner := bufio.NewScanner(r)
	line := &Line{Index: 0}
	for scanner.Scan() {
		if SplitLine(scanner.Text(), line) {
			f(line)
			if line.Err != nil {
				return fmt.Errorf("error in line %d: %s", line.Index, line.Err)
			}
		}
	}
	return nil
}

// SplitLine decomposes string into the line structure.
// Returns false if the string is empty.
func SplitLine(s string, line *Line) bool {
	line.Index++
	line.Comment = line.Comment[:0]
	line.Part = line.Part[:0]
	line.Fields = line.Fields[:0]
	s = strings.TrimSpace(s)
	if p := strings.IndexByte(s, '#'); p >= 0 {
		line.Comment = strings.TrimSpace(s[p+1:])
		s = strings.TrimSpace(s[:p])
	}
	if len(s) > 0 {
		if s[0] == '@' {
			line.Part = strings.TrimSpace(s[1:])
		} else {
			line.Fields = strings.Split(s, ";")
			for i := range line.Fields {
				line.Fields[i] = strings.TrimSpace(line.Fields[i])
			}
		}
	}
	return true
}

// parseRune parses a [U+]HEX codepoint string into a rune.
func parseRune(s string) (rune, error) {
	if len(s) > 2 && s[0] == 'U' && s[1] == '+' {
		s = s[2:]
	}
	v, err := strconv.ParseUint(s, 16, 32)
	return rune(v), err
}
