package rendering

import (
	"fmt"
	"regexp"
)

const (
	// language=regexp
	notesRegex = `(?i)notes?:`
	// language=regexp
	splitFormat = `\r?\n%s\r?\n`
)

var (
	htmlElementAttributesRegexp = regexp.MustCompile(`(?P<key>[a-z]+(-[a-z]+)*)="(?P<value>.+)"`)
	notesRegexp                 = regexp.MustCompile(fmt.Sprintf(`^%s`, notesRegex))
	notesLineRegexp             = regexp.MustCompile(fmt.Sprintf(`\r?\n%s\r?\n`, notesRegex))
)
