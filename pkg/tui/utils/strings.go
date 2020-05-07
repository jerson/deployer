package utils

import (
	"github.com/spkg/bom"
)

// CleanString ...
func CleanString(s string) string {
	output := string(bom.Clean([]byte(s)))
	return NormalizeLinefeeds(output)
}
