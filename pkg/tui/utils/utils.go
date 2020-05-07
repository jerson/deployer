package utils

import (
	"strings"
	"time"
)

// RunEvery ...
func RunEvery(duration time.Duration, function func() error) {
	for {
		_ = function()
		time.Sleep(duration)
	}
}

// Max returns the maximum of two integers
func Max(x, y int) int {
	if x > y {
		return x
	}
	return y
}

// NormalizeLinefeeds - Removes all Windows and Mac style line feeds
func NormalizeLinefeeds(str string) string {
	str = strings.Replace(str, "\r\n", "\n", -1)
	str = strings.Replace(str, "\r", "", -1)
	return str
}
