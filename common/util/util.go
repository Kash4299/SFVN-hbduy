package util

func IsValidDays(days int) bool {
	switch days {
	case 1, 7, 14, 30, 90, 180, 365:
		return true
	default:
		return false
	}
}
