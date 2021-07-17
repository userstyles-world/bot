package utils

func ContainsInSlice(slide []string, value string) bool {
	for _, v := range slide {
		if v == value {
			return true
		}
	}
	return false
}
