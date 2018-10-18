package utils

// Here lies everything that is considered helper function

// IsLetter checks if it's a letter
// this is also checked on the server side too, but well behaving client should not burden the server
func IsLetter(s string) bool {

	//should be just one
	if len(s) > 1 {
		return false
	}

	//should not be a strange character
	for _, r := range s {
		if (r < 'a' || r > 'z') && (r < 'A' || r > 'Z') {
			return false
		}
	}
	return true
}
