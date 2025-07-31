package middlewares

import (
	"html/template"
	"net/url"
	"strings"

	"github.com/gin-gonic/gin"
	strip "github.com/grokify/html-strip-tags-go"
)

func Secure() gin.HandlerFunc {
	return func(c *gin.Context) {
		if strings.Contains(c.Request.URL.Path, "/docs") {
			c.Next()
			return
		}

		c.Writer.Header().Set("X-XSS-Protection", "1; mode=block")
		c.Writer.Header().Set("X-Content-Type-Options", "nosniff")
		c.Writer.Header().Set("X-Frame-Options", "SAMEORIGIN")
		c.Next()
	}
}

func StripHTMLMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		query := c.Request.URL.Query()
		sanitized := url.Values{}

		for key, values := range query {
			for _, value := range values {
				safeValue := template.HTMLEscapeString(value)
				safeValue = strings.ReplaceAll(safeValue, "=", "")
				safeValue = strings.ReplaceAll(safeValue, "<", "")
				safeValue = strings.ReplaceAll(safeValue, ">", "")
				safeValue = strings.ReplaceAll(safeValue, "*", "")
				safeValue = strings.ReplaceAll(safeValue, " AND ", "")
				safeValue = strings.ReplaceAll(safeValue, " OR ", "")
				safeValue = strings.ReplaceAll(safeValue, " and ", "")
				safeValue = strings.ReplaceAll(safeValue, " or ", "")
				safeValue = strings.ReplaceAll(safeValue, " || ", "")
				safeValue = strings.ReplaceAll(safeValue, " && ", "")
				safeValue = strings.ReplaceAll(safeValue, "'", "")
				safeValue = strings.ReplaceAll(safeValue, "&#39;", "")
				safeValue = strip.StripTags(safeValue)
				sanitized.Add(key, safeValue)
			}
		}

		c.Request.URL.RawQuery = sanitized.Encode()
		c.Next()
	}
}
