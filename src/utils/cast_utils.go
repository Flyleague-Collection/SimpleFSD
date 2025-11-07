package utils

import "strconv"

// StrToInt cast string to int
func StrToInt(str string, defaultValue int) int {
	result, err := strconv.Atoi(str)
	if err != nil {
		return defaultValue
	}
	return result
}

// StrToFloat cast string to float64
func StrToFloat(str string, defaultValue float64) float64 {
	result, err := strconv.ParseFloat(str, 64)
	if err != nil {
		return defaultValue
	}
	return result
}
