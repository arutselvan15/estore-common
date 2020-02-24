// Package helper helper
package helper

// ContainsString check for a string in a slice
func ContainsString(slice []string, s string) bool {
	// check items
	for _, item := range slice {
		if item == s {
			return true
		}
	}

	return false
}

// RemoveString removes string from slice
func RemoveString(slice []string, s string) []string {
	result := []string{}

	for _, item := range slice {
		if item == s {
			continue
		}

		result = append(result, item)
	}

	return result
}
