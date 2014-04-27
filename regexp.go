package main

import (
	"fmt"
	"regexp"
)

// The same regular expression as above but using a named capturing group
var call = regexp.MustCompile("(?:^leiser define ([a-zA-Z0-9]+): .*)")

func main() {
	message := "leiser define string: anything"

	matches := call.FindStringSubmatch(message)

	if len(matches) == 2 {
		fmt.Println(matches)
	}
}
