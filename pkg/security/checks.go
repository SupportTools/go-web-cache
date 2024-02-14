package security

import (
	"net/http"
	"strings"
)

// HasWordPressLoginCookie checks if the request contains a WordPress login cookie.
func HasWordPressLoginCookie(req *http.Request) bool {
	for _, cookie := range req.Cookies() {
		if strings.HasPrefix(cookie.Name, "wordpress_logged_in_") {
			return true
		}
	}
	return false
}
