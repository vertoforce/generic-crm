package sqlcrm

import (
	"regexp"
	"strings"
	"unicode/utf8"
)

var (
	space = regexp.MustCompile(`\s+`)
)

// CleanText from spaces, newlines, excessive whitespace, etc
func CleanText(s string) string {
	s = toValidUTF8(s)
	s = space.ReplaceAllString(s, " ")
	s = strings.ReplaceAll(s, "\ua020", "")
	s = strings.Trim(s, " \n\ua020")
	return s
}

// CleanNumber runs CleanText then removes dollar signs, commas
func CleanNumber(s string) string {
	s = strings.ReplaceAll(s, "$", "")
	s = strings.ReplaceAll(s, "*", "")
	s = strings.ReplaceAll(s, "USD", "")
	s = strings.ReplaceAll(s, ",", "")
	s = strings.ReplaceAll(s, "%", "")
	s = CleanText(s)
	return s
}

// toValidUTF8 Converts a string ot valid utf8 if it is not, or just returns the original
func toValidUTF8(s string) string {
	if utf8.ValidString(s) {
		return s
	}

	v := make([]rune, 0, len(s))
	for i, r := range s {
		if r == utf8.RuneError {
			_, size := utf8.DecodeRuneInString(s[i:])
			if size == 1 {
				continue
			}
		}
		v = append(v, r)
	}
	s = string(v)

	return s
}

// StringToASCIIBytes Fake converting UTF-8 internal string representation to standard
// ASCII bytes for serial connections.
func StringToASCIIBytes(s string) []byte {
	t := make([]byte, utf8.RuneCountInString(s))
	i := 0
	for _, r := range s {
		t[i] = byte(r)
		i++
	}
	return t
}
