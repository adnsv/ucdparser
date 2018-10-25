package ucdparser

import (
	"fmt"
	"strings"
	"testing"
)

func handle(input string) string {
	out := &strings.Builder{}
	Parse(strings.NewReader(input), func(ln *Line) {
		if ln.Part != "" {
			fmt.Fprintln(out, "@"+ln.Part)
		} else if len(ln.Fields) == 0 {
			// fmt.Fprintln(out, "#"+ln.Comment)
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
			fmt.Fprintln(out, s)
		}
	})
	return out.String()
}

func TestParser(t *testing.T) {
	type args struct {
		input string
		want  string
	}
	tests := []args{
		{in1, out1},
	}
	for i, tt := range tests {
		t.Run(fmt.Sprintf("test #%d", i), func(t *testing.T) {
			got := handle(tt.input)
			if got != tt.want {
				t.Errorf("\ngot:\n%v\nwant:\n%v", got, tt.want)
			}
		})
	}
	//	fmt.Print(out.String())
}

const in1 = `# Comments should be skipped
# rune;  bool;  uint; int; float; runes; # Y
0..0005; Y;     0;    2;      -5.25 ;  0 1 2 3 4 5;
6..0007; Yes  ; 6;    1;     -4.25  ;  0006 0007;
8;       T ;    8 ;   0 ;-3.25  ;;# T
9;       True  ;9  ;  -1;-2.25  ;  0009;

# more comments to be ignored
@Part0  

U+0A;       N;   10  ;   -2;  -1.25; ;# N
B;       No;   11 ;   -3;  -0.25; 
C;        False;12;   -4;   0.75;
D;        ;13;-5;1.75;

@Part1   # Another part. 
# We test part comments get removed by not commenting the the next line.
E..10FFFF; F;   14  ; -6;   2.75;
`

const out1 = `|0x0..0x5|Y|0|2|-5.25|0 1 2 3 4 5|
|0x6..0x7|Yes|6|1|-4.25|0006 0007|
|0x8..0x8|T|8|0|-3.25||#T
|0x9..0x9|True|9|-1|-2.25|0009|
@Part0
|0xa..0xa|N|10|-2|-1.25||#N
|0xb..0xb|No|11|-3|-0.25|
|0xc..0xc|False|12|-4|0.75|
|0xd..0xd||13|-5|1.75|
@Part1
|0xe..0x10ffff|F|14|-6|2.75|
`
