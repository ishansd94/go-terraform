package helpers

func InStringSlice(haystack []string, needle string ) bool {
	for _,x := range haystack {
		if x == needle {
			return true
		}
	}
	return false
}