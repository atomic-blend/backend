package draftmail

import (
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"net/url"
	"strings"
	"testing"
	"time"

	"github.com/atomic-blend/backend/mail/models"
	"github.com/atomic-blend/backend/mail/tests/mocks"
	"github.com/atomic-blend/backend/shared/middlewares/auth"
	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"github.com/webstradev/gin-pagination/v2/pkg/pagination"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

func TestDraftMailController_GetDraftMailsSince(t *testing.T) {
	gin.SetMode(gin.TestMode)

	tests := []struct {
		name               string
		sinceParam         string
		queryParams        string
		expectedStatus     int
		setupMock          func(*mocks.MockDraftMailRepository, primitive.ObjectID, time.Time, int64, int64)
		setupAuth          func(*gin.Context, primitive.ObjectID)
		expectedDraftMails int
	}{
		{
			name:               "Success with valid ISO8601 date",
			sinceParam:         "2024-01-01T00:00:00Z",
			queryParams:        "",
			expectedStatus:     http.StatusOK,
			expectedDraftMails: 2,
			setupMock: func(mockRepo *mocks.MockDraftMailRepository, userID primitive.ObjectID, sinceTime time.Time, page, limit int64) {
				mailID1 := primitive.NewObjectID()
				mailID2 := primitive.NewObjectID()

				// Create draft mails for the authenticated user only (repository now filters by userID)
				draftMails := []*models.SendMail{
					{
						Mail: &models.Mail{
							ID:     &mailID1,
							UserID: userID,
							Headers: map[string]string{
								"Subject": "User Draft Mail 1",
							},
						},
						SendStatus: models.SendStatusPending,
					},
					{
						Mail: &models.Mail{
							ID:     &mailID2,
							UserID: userID,
							Headers: map[string]string{
								"Subject": "User Draft Mail 2",
							},
						},
						SendStatus: models.SendStatusPending,
					},
				}
				mockRepo.On("GetSince", mock.Anything, userID, sinceTime, page, limit).Return(draftMails, int64(2), nil)
			},
			setupAuth: func(c *gin.Context, userID primitive.ObjectID) {
				c.Set("authUser", &auth.UserAuthInfo{UserID: userID})
			},
		},
		{
			name:               "Success with pagination parameters",
			sinceParam:         "2024-01-01T00:00:00Z",
			queryParams:        "?page=2&size=15",
			expectedStatus:     http.StatusOK,
			expectedDraftMails: 1,
			setupMock: func(mockRepo *mocks.MockDraftMailRepository, userID primitive.ObjectID, sinceTime time.Time, page, limit int64) {
				mailID := primitive.NewObjectID()
				draftMails := []*models.SendMail{
					{
						Mail: &models.Mail{
							ID:     &mailID,
							UserID: userID,
							Headers: map[string]string{
								"Subject": "User Draft Mail",
							},
						},
						SendStatus: models.SendStatusPending,
					},
				}
				mockRepo.On("GetSince", mock.Anything, userID, sinceTime, page, limit).Return(draftMails, int64(16), nil)
			},
			setupAuth: func(c *gin.Context, userID primitive.ObjectID) {
				c.Set("authUser", &auth.UserAuthInfo{UserID: userID})
			},
		},
		{
			name:               "Success with no draft mails found",
			sinceParam:         "2024-01-01T00:00:00Z",
			queryParams:        "",
			expectedStatus:     http.StatusOK,
			expectedDraftMails: 0,
			setupMock: func(mockRepo *mocks.MockDraftMailRepository, userID primitive.ObjectID, sinceTime time.Time, page, limit int64) {
				draftMails := []*models.SendMail{}
				mockRepo.On("GetSince", mock.Anything, userID, sinceTime, page, limit).Return(draftMails, int64(0), nil)
			},
			setupAuth: func(c *gin.Context, userID primitive.ObjectID) {
				c.Set("authUser", &auth.UserAuthInfo{UserID: userID})
			},
		},
		{
			name:               "Error when repository fails",
			sinceParam:         "2024-01-01T00:00:00Z",
			queryParams:        "",
			expectedStatus:     http.StatusInternalServerError,
			expectedDraftMails: 0,
			setupMock: func(mockRepo *mocks.MockDraftMailRepository, userID primitive.ObjectID, sinceTime time.Time, page, limit int64) {
				mockRepo.On("GetSince", mock.Anything, userID, sinceTime, page, limit).Return(nil, int64(0), assert.AnError)
			},
			setupAuth: func(c *gin.Context, userID primitive.ObjectID) {
				c.Set("authUser", &auth.UserAuthInfo{UserID: userID})
			},
		},
		{
			name:               "Error when missing since parameter",
			sinceParam:         "",
			queryParams:        "",
			expectedStatus:     http.StatusBadRequest,
			expectedDraftMails: 0,
			setupMock: func(mockRepo *mocks.MockDraftMailRepository, userID primitive.ObjectID, sinceTime time.Time, page, limit int64) {
				// No mock setup needed as the request should fail before reaching the repository
			},
			setupAuth: func(c *gin.Context, userID primitive.ObjectID) {
				c.Set("authUser", &auth.UserAuthInfo{UserID: userID})
			},
		},
		{
			name:               "Error when invalid date format",
			sinceParam:         "invalid-date",
			queryParams:        "",
			expectedStatus:     http.StatusBadRequest,
			expectedDraftMails: 0,
			setupMock: func(mockRepo *mocks.MockDraftMailRepository, userID primitive.ObjectID, sinceTime time.Time, page, limit int64) {
				// No mock setup needed as the request should fail before reaching the repository
			},
			setupAuth: func(c *gin.Context, userID primitive.ObjectID) {
				c.Set("authUser", &auth.UserAuthInfo{UserID: userID})
			},
		},
		{
			name:               "Error when unauthorized",
			sinceParam:         "2024-01-01T00:00:00Z",
			queryParams:        "",
			expectedStatus:     http.StatusUnauthorized,
			expectedDraftMails: 0,
			setupMock: func(mockRepo *mocks.MockDraftMailRepository, userID primitive.ObjectID, sinceTime time.Time, page, limit int64) {
				// No mock setup needed as the request should fail before reaching the repository
			},
			setupAuth: func(c *gin.Context, userID primitive.ObjectID) {
				// No auth setup - should trigger unauthorized
			},
		},
		{
			name:               "Success with different date formats",
			sinceParam:         "2024-01-01T12:30:45+02:00",
			queryParams:        "",
			expectedStatus:     http.StatusOK,
			expectedDraftMails: 1,
			setupMock: func(mockRepo *mocks.MockDraftMailRepository, userID primitive.ObjectID, sinceTime time.Time, page, limit int64) {
				mailID := primitive.NewObjectID()
				draftMails := []*models.SendMail{
					{
						Mail: &models.Mail{
							ID:     &mailID,
							UserID: userID,
							Headers: map[string]string{
								"Subject": "User Draft Mail",
							},
						},
						SendStatus: models.SendStatusPending,
					},
				}
				mockRepo.On("GetSince", mock.Anything, userID, sinceTime, page, limit).Return(draftMails, int64(1), nil)
			},
			setupAuth: func(c *gin.Context, userID primitive.ObjectID) {
				c.Set("authUser", &auth.UserAuthInfo{UserID: userID})
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			mockRepo := &mocks.MockDraftMailRepository{}
			userID := primitive.NewObjectID()

			// Parse the since parameter to get the expected time
			var sinceTime time.Time
			if tt.sinceParam != "" {
				var err error
				sinceTime, err = time.Parse(time.RFC3339, tt.sinceParam)
				if err != nil {
					// For invalid date tests, we don't need to parse
					sinceTime = time.Now()
				}
			}

			// Default pagination parameters
			page := int64(1)
			limit := int64(10)

			// Parse pagination parameters from queryParams if provided
			if tt.queryParams != "" {
				if strings.Contains(tt.queryParams, "page=2") {
					page = 2
				}
				if strings.Contains(tt.queryParams, "size=15") {
					limit = 15
				}
			}

			tt.setupMock(mockRepo, userID, sinceTime, page, limit)

			controller := NewDraftMailController(mockRepo, nil, nil, nil)

			router := gin.New()
			router.Use(func(c *gin.Context) {
				tt.setupAuth(c, userID)
				// Set pagination parameters in context
				c.Set("page", int(page))
				c.Set("size", int(limit))
				c.Next()
			})

			draftMailRoutes := router.Group("/mail/draft")
			{
				draftMailRoutes.GET("/since", pagination.New(), controller.GetDraftMailsSince)
			}

			// Create request with since parameter
			reqURL := "/mail/draft/since"
			if tt.sinceParam != "" {
				reqURL += "?since=" + url.QueryEscape(tt.sinceParam)
			}
			if tt.queryParams != "" {
				if tt.sinceParam != "" {
					reqURL += "&" + strings.TrimPrefix(tt.queryParams, "?")
				} else {
					reqURL += tt.queryParams
				}
			}

			req, _ := http.NewRequest("GET", reqURL, nil)
			w := httptest.NewRecorder()

			router.ServeHTTP(w, req)

			if w.Code != tt.expectedStatus {
				t.Logf("Test: %s, Expected status %d, got %d. Response body: %s", tt.name, tt.expectedStatus, w.Code, w.Body.String())
			}
			assert.Equal(t, tt.expectedStatus, w.Code)

			if tt.expectedStatus == http.StatusOK {
				var response PaginatedDraftMailResponse
				err := json.Unmarshal(w.Body.Bytes(), &response)
				assert.NoError(t, err)
				assert.Len(t, response.DraftMails, tt.expectedDraftMails)
				assert.NotNil(t, response.TotalCount)
				assert.NotNil(t, response.Page)
				assert.NotNil(t, response.Size)
				assert.NotNil(t, response.TotalPages)

				// Verify all returned draft mails belong to the authenticated user
				for _, draftMail := range response.DraftMails {
					assert.NotNil(t, draftMail.Mail)
					assert.Equal(t, userID, draftMail.Mail.UserID)
				}
			}

			mockRepo.AssertExpectations(t)
		})
	}
}
