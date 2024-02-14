package cache

import "strings"

// ParseCacheControl parses the Cache-Control header and returns a map of directives.
func ParseCacheControl(header string) map[string]string {
	directives := make(map[string]string)
	for _, directive := range strings.Split(header, ",") {
		parts := strings.SplitN(strings.TrimSpace(directive), "=", 2)
		if len(parts) == 2 {
			directives[parts[0]] = parts[1]
		} else {
			directives[parts[0]] = ""
		}
	}
	return directives
}
