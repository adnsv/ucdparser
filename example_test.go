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
