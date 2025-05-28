package time_entry

import (
	"atomic_blend_api/auth"
	"context"
	"net/http"

	"github.com/gin-gonic/gin"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

func setupTestRouter() *gin.Engine {
	gin.SetMode(gin.TestMode)
	return gin.New()
}

func setupAuthenticatedRequest(req *http.Request, userID primitive.ObjectID) {
	// Simulate an authenticated user context
	ctx := req.Context()
	ctx = context.WithValue(ctx, "authUser", &auth.UserAuthInfo{
		UserID: userID,
	})
	*req = *req.WithContext(ctx)
}








