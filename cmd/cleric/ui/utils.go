package ui

func isValidServerName(name string) bool {
	// Check for empty string
	if len(name) == 0 {
		return false
	}

	// Check if name starts with a number
	if name[0] >= '0' && name[0] <= '9' {
		return false
	}

	// Check if name contains only alphanumeric characters and underscore
	for _, char := range name {
		if !((char >= 'a' && char <= 'z') ||
			(char >= 'A' && char <= 'Z') ||
			(char >= '0' && char <= '9') ||
			char == '_') {
			return false
		}
	}
	return true
}
