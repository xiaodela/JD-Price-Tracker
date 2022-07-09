package util

func Min(s []float64) float64 {
	minVal := s[0]
	if len(s) == 1 {
		return minVal
	}
	for _, val := range s[1:] {
		if val < minVal {
			minVal = val
		}
	}
	return minVal
}
