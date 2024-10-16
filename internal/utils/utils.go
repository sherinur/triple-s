package utils

func ValidateBucketName(name string) bool {
	// TODO: Bucket names must be unique across the system.
	// TODO: Names should be between 3 and 63 characters long.
	// TODO: Only lowercase letters, numbers, hyphens (-), and dots (.) are allowed.
	// TODO: Must not be formatted as an IP address (e.g., 192.168.0.1).
	// TODO: Must not begin or end with a hyphen and must not contain two consecutive periods or dashes.

	return true
}

func IsUniqueBucketName(name string) bool {
	return true
}
