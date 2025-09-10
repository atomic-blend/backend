package mail

import (
	"bytes"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/atomic-blend/backend/mail/models"
	"github.com/atomic-blend/backend/mail/tests/mocks"
	"github.com/atomic-blend/backend/shared/middlewares/auth"
	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

func TestMailController_PutMailActions(t *testing.T) {
	gin.SetMode(gin.TestMode)

	userID := primitive.NewObjectID()
	mailID := primitive.NewObjectID()

	tests := []struct {
		name           string
		payload        interface{}
		expectedStatus int
		setupMock      func(*mocks.MockMailRepository, primitive.ObjectID, primitive.ObjectID)
		setupAuth      func(*gin.Context, primitive.ObjectID)
	}{
		{
			name: "Success mark read",
			payload: PutActionsPayload{
				Read: []string{mailID.Hex()},
			},
			expectedStatus: http.StatusOK,
			setupMock: func(mockRepo *mocks.MockMailRepository, userID primitive.ObjectID, mailID primitive.ObjectID) {
				mail := &models.Mail{ID: &mailID, UserID: userID}
				mockRepo.On("GetByID", mock.Anything, mailID).Return(mail, nil)
				mockRepo.On("Update", mock.Anything, mock.MatchedBy(func(m *models.Mail) bool {
					return m.ID != nil && *m.ID == mailID && m.Read != nil && *m.Read == true
				})).Return(nil)
			},
			setupAuth: func(c *gin.Context, userID primitive.ObjectID) {
				c.Set("authUser", &auth.UserAuthInfo{UserID: userID})
			},
		},
		{
			name: "Success mark unread",
			payload: PutActionsPayload{
				Unread: []string{mailID.Hex()},
			},
			expectedStatus: http.StatusOK,
			setupMock: func(mockRepo *mocks.MockMailRepository, userID primitive.ObjectID, mailID primitive.ObjectID) {
				mail := &models.Mail{ID: &mailID, UserID: userID}
				mockRepo.On("GetByID", mock.Anything, mailID).Return(mail, nil)
				mockRepo.On("Update", mock.Anything, mock.MatchedBy(func(m *models.Mail) bool {
					return m.ID != nil && *m.ID == mailID && m.Read != nil && *m.Read == false
				})).Return(nil)
			},
			setupAuth: func(c *gin.Context, userID primitive.ObjectID) {
				c.Set("authUser", &auth.UserAuthInfo{UserID: userID})
			},
		},
		{
			name:           "Bad payload",
			payload:        "not a json",
			expectedStatus: http.StatusBadRequest,
			setupMock:      func(mockRepo *mocks.MockMailRepository, userID primitive.ObjectID, mailID primitive.ObjectID) {},
			setupAuth: func(c *gin.Context, userID primitive.ObjectID) {
				c.Set("authUser", &auth.UserAuthInfo{UserID: userID})
			},
		},
		{
			name:           "Unauthorized",
			payload:        PutActionsPayload{Read: []string{mailID.Hex()}},
			expectedStatus: http.StatusUnauthorized,
			setupMock:      func(mockRepo *mocks.MockMailRepository, userID primitive.ObjectID, mailID primitive.ObjectID) {},
			setupAuth:      func(c *gin.Context, userID primitive.ObjectID) {},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			mockRepo := &mocks.MockMailRepository{}
			tt.setupMock(mockRepo, userID, mailID)

			controller := NewMailController(mockRepo)

			router := gin.New()
			router.Use(func(c *gin.Context) {
				tt.setupAuth(c, userID)
				c.Next()
			})

			mailRoutes := router.Group("/mail")
			{
				mailRoutes.PUT("/actions", controller.PutMailActions)
			}

			var body bytes.Buffer
			if str, ok := tt.payload.(string); ok {
				body.WriteString(str)
			} else {
				json.NewEncoder(&body).Encode(tt.payload)
			}

			req, _ := http.NewRequest("PUT", "/mail/actions", &body)
			req.Header.Set("Content-Type", "application/json")
			w := httptest.NewRecorder()

			router.ServeHTTP(w, req)

			assert.Equal(t, tt.expectedStatus, w.Code)

			mockRepo.AssertExpectations(t)
		})
	}
}
