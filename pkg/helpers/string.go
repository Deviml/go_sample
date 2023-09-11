package helpers

import "strings"

func SplitStringQuery(s string) []string {
	if s == "" {
		return []string{}
	}
	return strings.Split(s, "|")
}
