package cache

import (
	"fmt"
	"strings"
)

// ReconstructCacheControl reconstructs the Cache-Control header from a map of directives.
func ReconstructCacheControl(cacheControl map[string]string) string {
	var parts []string
	for k, v := range cacheControl {
		if v != "" {
			parts = append(parts, fmt.Sprintf("%s=%s", k, v))
		} else {
			parts = append(parts, k)
		}
	}
	return strings.Join(parts, ", ")
}
