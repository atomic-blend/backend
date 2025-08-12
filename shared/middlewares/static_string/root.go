package staticstringmiddleware

import "github.com/gin-gonic/gin"

// RequireStaticStringMiddleware is a middleware that checks if a static string is present in the request
func RequireStaticStringMiddleware(group *gin.RouterGroup, staticString string) *gin.RouterGroup {
	group.Use(New(staticString))
	return group
}