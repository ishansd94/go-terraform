package helpers

import (
	"regexp"
	"strings"
)

const AnsiCodes = "[\u001B\u009B][[\\]()#;?]*(?:(?:(?:[a-zA-Z\\d]*(?:;[a-zA-Z\\d]*)*)?\u0007)|(?:(?:\\d{1,4}(?:;\\d{0,4})*)?[\\dA-PRZcf-ntqry=><~]))"

func InStringSlice(haystack []string, needle string ) bool {
	for _,x := range haystack {
		if x == needle {
			return true
		}
	}
	return false
}

func TrimString(org string, trimmers map[string]string) string {
	str := org
	for oldstr, newstr := range trimmers{
		str = strings.Replace(str, oldstr, newstr ,-1)
	}
	return str
}

func RemoveAnsiCodes(str string) string {
	re := regexp.MustCompile(AnsiCodes)
	return re.ReplaceAllString(str, "")
}

func SanitizeHCL(out string) string {
	var str string

	splitFunc := func(c rune) bool {
		return c == '\n'
	}

	r := strings.FieldsFunc(out, splitFunc)
	r = r[1:]

	str = strings.Join(r, "\n")

	return RemoveAnsiCodes(str)
}