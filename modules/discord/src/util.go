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

func reduce(arr []interface{}, chunkSize int) [][]interface{} {
	chunks := make([][]interface{}, 0)
	chunk := make([]interface{}, 0)
	for i, item := range arr {
		chunkIndex := i / chunkSize
		if chunkIndex >= len(chunks) {
			chunks = append(chunks, chunk)
			chunk = make([]interface{}, 0)
		}
		chunk = append(chunk, item)
		chunks[chunkIndex] = chunk
	}
	if len(chunk) > 0 {
		chunks = append(chunks, chunk)
	}
	return chunks
}
