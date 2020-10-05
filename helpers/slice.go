package helpers

import (
	"strings"
)

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