package utils

func Difference(a []string, b []string) []string {
	m := make(map[string]bool)
	for _, item := range b {
		m[item] = true
	}

	var diff []string
	for _, item := range a {
		if _, found := m[item]; !found {
			diff = append(diff, item)
		}
	}

	return diff
}
