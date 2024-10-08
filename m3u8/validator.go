package m3u8

import "regexp"

var tsRegex = regexp.MustCompile(`\.ts`)

type ParseValidationFunc = func(string) bool

func DefaultParseValidationFunc(text string) bool {
	return tsRegex.MatchString(text)
}