package main

import "unicode/utf8"

func limitString(s string, maxLength int) string {
	if utf8.RuneCountInString(s) <= maxLength {
		return s
	}
	runes := []rune(s)
	return string(runes[:maxLength-3]) + "..."
}

func GroupBy(xs []map[string]interface{}, key string) map[interface{}][]map[string]interface{} {
	rv := make(map[interface{}][]map[string]interface{})
	for _, x := range xs {
		k := x[key]
		rv[k] = append(rv[k], x)
	}
	return rv
}

func ShortText(s string, i int) string {
	if len(s) < i {
		return s
	}
	if utf8.ValidString(s[:i]) {
		return s[:i]
	}
	return s[:i+1]
}
