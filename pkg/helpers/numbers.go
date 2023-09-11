package helpers

import (
	"strconv"
)

func AlwaysIntFromString(s string) int {
	i, err := strconv.Atoi(s)
	if err != nil {
		return 0
	}
	return i
}

func AlwaysInt64FromString(s string) int64 {
	i, err := strconv.ParseInt(s, 10, 64)
	if err != nil {
		return 0
	}
	return i
}

func AlwaysFloatFromString(s string) float32 {
	f, err := strconv.ParseFloat(s, 32)
	if err != nil {
		return 0
	}
	return float32(f)
}
