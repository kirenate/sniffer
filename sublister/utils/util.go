package utils

import (
	"slices"
	"strconv"
	"strings"
)

func Format(s string, query interface{}, format string) string {
	switch i := query.(type) {
	case string:
		s = strings.ReplaceAll(s, format, i)
	case int:
		s = strings.ReplaceAll(s, format, strconv.Itoa(i))
	}
	return s
}

func DeleteRepetitions(s []string) []string {
	slices.Sort(s)
	var newSlice []string
	for _, v := range s {
		if slices.Contains(newSlice, v) {
			continue
		}
		newSlice = append(newSlice, v)
	}
	return newSlice
}
