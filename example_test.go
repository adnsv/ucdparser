package ucdparser

import (
	"fmt"
	"log"
)

const root = "https://www.unicode.org/Public/UCD/latest/ucd/"

func Example() {

	// fetch from remote

	url := root + "NamesList.txt"
	fmt.Printf("Fetching '%s'...", url)
	data, err := Fetch(url)
	if err != nil {
		log.Fatal(err)
	}
	defer data.Close()

	// parse

	Parse(data, func(ln *Line) {
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
