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
			name: "Success archive mail",
			payload: PutActionsPayload{
				Archived: []string{mailID.Hex()},
			},
			expectedStatus: http.StatusOK,
			setupMock: func(mockRepo *mocks.MockMailRepository, userID primitive.ObjectID, mailID primitive.ObjectID) {
				mail := &models.Mail{ID: &mailID, UserID: userID}
				mockRepo.On("GetByID", mock.Anything, mailID).Return(mail, nil)
				mockRepo.On("Update", mock.Anything, mock.MatchedBy(func(m *models.Mail) bool {
					return m.ID != nil && *m.ID == mailID && m.Archived != nil && *m.Archived == true && m.Trashed != nil && *m.Trashed == false
				})).Return(nil)
			},
			setupAuth: func(c *gin.Context, userID primitive.ObjectID) {
				c.Set("authUser", &auth.UserAuthInfo{UserID: userID})
			},
		},
		{
			name: "Success unarchive mail",
			payload: PutActionsPayload{
				Unarchived: []string{mailID.Hex()},
			},
			expectedStatus: http.StatusOK,
			setupMock: func(mockRepo *mocks.MockMailRepository, userID primitive.ObjectID, mailID primitive.ObjectID) {
				mail := &models.Mail{ID: &mailID, UserID: userID}
				mockRepo.On("GetByID", mock.Anything, mailID).Return(mail, nil)
				mockRepo.On("Update", mock.Anything, mock.MatchedBy(func(m *models.Mail) bool {
					return m.ID != nil && *m.ID == mailID && m.Archived != nil && *m.Archived == false
				})).Return(nil)
			},
			setupAuth: func(c *gin.Context, userID primitive.ObjectID) {
				c.Set("authUser", &auth.UserAuthInfo{UserID: userID})
			},
		},
		{
			name: "Success multiple archive operations",
			payload: PutActionsPayload{
				Archived:   []string{mailID.Hex()},
				Unarchived: []string{mailID.Hex()},
			},
			expectedStatus: http.StatusOK,
			setupMock: func(mockRepo *mocks.MockMailRepository, userID primitive.ObjectID, mailID primitive.ObjectID) {
				mail := &models.Mail{ID: &mailID, UserID: userID}
				// Expect two calls - one for archive, one for unarchive
				mockRepo.On("GetByID", mock.Anything, mailID).Return(mail, nil).Twice()
				mockRepo.On("Update", mock.Anything, mock.MatchedBy(func(m *models.Mail) bool {
					return m.ID != nil && *m.ID == mailID && m.Archived != nil && *m.Archived == true && m.Trashed != nil && *m.Trashed == false
				})).Return(nil).Once()
				mockRepo.On("Update", mock.Anything, mock.MatchedBy(func(m *models.Mail) bool {
					return m.ID != nil && *m.ID == mailID && m.Archived != nil && *m.Archived == false
				})).Return(nil).Once()
			},
			setupAuth: func(c *gin.Context, userID primitive.ObjectID) {
				c.Set("authUser", &auth.UserAuthInfo{UserID: userID})
			},
		},
		{
			name: "Success combined read and archive operations",
			payload: PutActionsPayload{
				Read:     []string{mailID.Hex()},
				Archived: []string{mailID.Hex()},
			},
			expectedStatus: http.StatusOK,
			setupMock: func(mockRepo *mocks.MockMailRepository, userID primitive.ObjectID, mailID primitive.ObjectID) {
				mail := &models.Mail{ID: &mailID, UserID: userID}
				// Expect two calls - one for read, one for archive
				mockRepo.On("GetByID", mock.Anything, mailID).Return(mail, nil).Twice()
				mockRepo.On("Update", mock.Anything, mock.MatchedBy(func(m *models.Mail) bool {
					return m.ID != nil && *m.ID == mailID && m.Read != nil && *m.Read == true
				})).Return(nil).Once()
				mockRepo.On("Update", mock.Anything, mock.MatchedBy(func(m *models.Mail) bool {
					return m.ID != nil && *m.ID == mailID && m.Archived != nil && *m.Archived == true && m.Trashed != nil && *m.Trashed == false
				})).Return(nil).Once()
			},
			setupAuth: func(c *gin.Context, userID primitive.ObjectID) {
				c.Set("authUser", &auth.UserAuthInfo{UserID: userID})
			},
		},
		{
			name: "Success combined unread and unarchive operations",
			payload: PutActionsPayload{
				Unread:     []string{mailID.Hex()},
				Unarchived: []string{mailID.Hex()},
			},
			expectedStatus: http.StatusOK,
			setupMock: func(mockRepo *mocks.MockMailRepository, userID primitive.ObjectID, mailID primitive.ObjectID) {
				mail := &models.Mail{ID: &mailID, UserID: userID}
				// Expect two calls - one for unread, one for unarchive
				mockRepo.On("GetByID", mock.Anything, mailID).Return(mail, nil).Twice()
				mockRepo.On("Update", mock.Anything, mock.MatchedBy(func(m *models.Mail) bool {
					return m.ID != nil && *m.ID == mailID && m.Read != nil && *m.Read == false
				})).Return(nil).Once()
				mockRepo.On("Update", mock.Anything, mock.MatchedBy(func(m *models.Mail) bool {
					return m.ID != nil && *m.ID == mailID && m.Archived != nil && *m.Archived == false
				})).Return(nil).Once()
			},
			setupAuth: func(c *gin.Context, userID primitive.ObjectID) {
				c.Set("authUser", &auth.UserAuthInfo{UserID: userID})
			},
		},
		{
			name: "Success all operations combined",
			payload: PutActionsPayload{
				Read:       []string{mailID.Hex()},
				Unread:     []string{mailID.Hex()},
				Archived:   []string{mailID.Hex()},
				Unarchived: []string{mailID.Hex()},
			},
			expectedStatus: http.StatusOK,
			setupMock: func(mockRepo *mocks.MockMailRepository, userID primitive.ObjectID, mailID primitive.ObjectID) {
				mail := &models.Mail{ID: &mailID, UserID: userID}
				// Expect four calls - one for each operation
				mockRepo.On("GetByID", mock.Anything, mailID).Return(mail, nil).Times(4)
				mockRepo.On("Update", mock.Anything, mock.MatchedBy(func(m *models.Mail) bool {
					return m.ID != nil && *m.ID == mailID && m.Read != nil && *m.Read == true
				})).Return(nil).Once()
				mockRepo.On("Update", mock.Anything, mock.MatchedBy(func(m *models.Mail) bool {
					return m.ID != nil && *m.ID == mailID && m.Read != nil && *m.Read == false
				})).Return(nil).Once()
				mockRepo.On("Update", mock.Anything, mock.MatchedBy(func(m *models.Mail) bool {
					return m.ID != nil && *m.ID == mailID && m.Archived != nil && *m.Archived == true && m.Trashed != nil && *m.Trashed == false
				})).Return(nil).Once()
				mockRepo.On("Update", mock.Anything, mock.MatchedBy(func(m *models.Mail) bool {
					return m.ID != nil && *m.ID == mailID && m.Archived != nil && *m.Archived == false
				})).Return(nil).Once()
			},
			setupAuth: func(c *gin.Context, userID primitive.ObjectID) {
				c.Set("authUser", &auth.UserAuthInfo{UserID: userID})
			},
		},
		{
			name: "Archive operation with invalid mail ID",
			payload: PutActionsPayload{
				Archived: []string{"invalid-id"},
			},
			expectedStatus: http.StatusOK,
			setupMock: func(mockRepo *mocks.MockMailRepository, userID primitive.ObjectID, mailID primitive.ObjectID) {
				// No mock expectations since invalid ID should be skipped
			},
			setupAuth: func(c *gin.Context, userID primitive.ObjectID) {
				c.Set("authUser", &auth.UserAuthInfo{UserID: userID})
			},
		},
		{
			name: "Archive operation with mail not found",
			payload: PutActionsPayload{
				Archived: []string{mailID.Hex()},
			},
			expectedStatus: http.StatusOK,
			setupMock: func(mockRepo *mocks.MockMailRepository, userID primitive.ObjectID, mailID primitive.ObjectID) {
				mockRepo.On("GetByID", mock.Anything, mailID).Return(nil, assert.AnError)
			},
			setupAuth: func(c *gin.Context, userID primitive.ObjectID) {
				c.Set("authUser", &auth.UserAuthInfo{UserID: userID})
			},
		},
		{
			name: "Archive operation with mail belonging to different user",
			payload: PutActionsPayload{
				Archived: []string{mailID.Hex()},
			},
			expectedStatus: http.StatusOK,
			setupMock: func(mockRepo *mocks.MockMailRepository, userID primitive.ObjectID, mailID primitive.ObjectID) {
				differentUserID := primitive.NewObjectID()
				mail := &models.Mail{ID: &mailID, UserID: differentUserID}
				mockRepo.On("GetByID", mock.Anything, mailID).Return(mail, nil)
			},
			setupAuth: func(c *gin.Context, userID primitive.ObjectID) {
				c.Set("authUser", &auth.UserAuthInfo{UserID: userID})
			},
		},
		{
			name: "Archive operation with repository update error",
			payload: PutActionsPayload{
				Archived: []string{mailID.Hex()},
			},
			expectedStatus: http.StatusInternalServerError,
			setupMock: func(mockRepo *mocks.MockMailRepository, userID primitive.ObjectID, mailID primitive.ObjectID) {
				mail := &models.Mail{ID: &mailID, UserID: userID}
				mockRepo.On("GetByID", mock.Anything, mailID).Return(mail, nil)
				mockRepo.On("Update", mock.Anything, mock.MatchedBy(func(m *models.Mail) bool {
					return m.ID != nil && *m.ID == mailID && m.Archived != nil && *m.Archived == true && m.Trashed != nil && *m.Trashed == false
				})).Return(assert.AnError)
			},
			setupAuth: func(c *gin.Context, userID primitive.ObjectID) {
				c.Set("authUser", &auth.UserAuthInfo{UserID: userID})
			},
		},
		{
			name: "Unarchive operation with invalid mail ID",
			payload: PutActionsPayload{
				Unarchived: []string{"invalid-id"},
			},
			expectedStatus: http.StatusOK,
			setupMock: func(mockRepo *mocks.MockMailRepository, userID primitive.ObjectID, mailID primitive.ObjectID) {
				// No mock expectations since invalid ID should be skipped
			},
			setupAuth: func(c *gin.Context, userID primitive.ObjectID) {
				c.Set("authUser", &auth.UserAuthInfo{UserID: userID})
			},
		},
		{
			name: "Unarchive operation with mail not found",
			payload: PutActionsPayload{
				Unarchived: []string{mailID.Hex()},
			},
			expectedStatus: http.StatusOK,
			setupMock: func(mockRepo *mocks.MockMailRepository, userID primitive.ObjectID, mailID primitive.ObjectID) {
				mockRepo.On("GetByID", mock.Anything, mailID).Return(nil, assert.AnError)
			},
			setupAuth: func(c *gin.Context, userID primitive.ObjectID) {
				c.Set("authUser", &auth.UserAuthInfo{UserID: userID})
			},
		},
		{
			name: "Unarchive operation with mail belonging to different user",
			payload: PutActionsPayload{
				Unarchived: []string{mailID.Hex()},
			},
			expectedStatus: http.StatusOK,
			setupMock: func(mockRepo *mocks.MockMailRepository, userID primitive.ObjectID, mailID primitive.ObjectID) {
				differentUserID := primitive.NewObjectID()
				mail := &models.Mail{ID: &mailID, UserID: differentUserID}
				mockRepo.On("GetByID", mock.Anything, mailID).Return(mail, nil)
			},
			setupAuth: func(c *gin.Context, userID primitive.ObjectID) {
				c.Set("authUser", &auth.UserAuthInfo{UserID: userID})
			},
		},
		{
			name: "Unarchive operation with repository update error",
			payload: PutActionsPayload{
				Unarchived: []string{mailID.Hex()},
			},
			expectedStatus: http.StatusInternalServerError,
			setupMock: func(mockRepo *mocks.MockMailRepository, userID primitive.ObjectID, mailID primitive.ObjectID) {
				mail := &models.Mail{ID: &mailID, UserID: userID}
				mockRepo.On("GetByID", mock.Anything, mailID).Return(mail, nil)
				mockRepo.On("Update", mock.Anything, mock.MatchedBy(func(m *models.Mail) bool {
					return m.ID != nil && *m.ID == mailID && m.Archived != nil && *m.Archived == false
				})).Return(assert.AnError)
			},
			setupAuth: func(c *gin.Context, userID primitive.ObjectID) {
				c.Set("authUser", &auth.UserAuthInfo{UserID: userID})
			},
		},
		{
			name: "Empty payload with all fields",
			payload: PutActionsPayload{
				Read:       []string{},
				Unread:     []string{},
				Archived:   []string{},
				Unarchived: []string{},
			},
			expectedStatus: http.StatusOK,
			setupMock: func(mockRepo *mocks.MockMailRepository, userID primitive.ObjectID, mailID primitive.ObjectID) {
				// No mock expectations since all arrays are empty
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
