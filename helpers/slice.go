package helpers

func StringSliceContains(haystack []string, needle string ) bool {
	for _,x := range haystack {
		if x == needle {
			return true
		}
	}
	return false
}