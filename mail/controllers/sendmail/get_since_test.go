package sendmail

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

func TestSendMailController_GetSendMailsSince(t *testing.T) {
	gin.SetMode(gin.TestMode)

	tests := []struct {
		name              string
		sinceParam        string
		queryParams       string
		expectedStatus    int
		setupMock         func(*mocks.MockSendMailRepository, primitive.ObjectID, time.Time, int64, int64)
		setupAuth         func(*gin.Context, primitive.ObjectID)
		expectedSendMails int
	}{
		{
			name:              "Success with valid ISO8601 date",
			sinceParam:        "2024-01-01T00:00:00Z",
			queryParams:       "",
			expectedStatus:    http.StatusOK,
			expectedSendMails: 2,
			setupMock: func(mockRepo *mocks.MockSendMailRepository, userID primitive.ObjectID, sinceTime time.Time, page, limit int64) {
				mailID1 := primitive.NewObjectID()
				mailID2 := primitive.NewObjectID()

				// Create send mails for the authenticated user only (repository now filters by userID)
				sendMails := []*models.SendMail{
					{
						Mail: &models.Mail{
							ID:     &mailID1,
							UserID: userID,
							Headers: map[string]string{
								"Subject": "User Send Mail 1",
							},
						},
						SendStatus: models.SendStatusPending,
					},
					{
						Mail: &models.Mail{
							ID:     &mailID2,
							UserID: userID,
							Headers: map[string]string{
								"Subject": "User Send Mail 2",
							},
						},
						SendStatus: models.SendStatusPending,
					},
				}
				mockRepo.On("GetSince", mock.Anything, userID, sinceTime, page, limit).Return(sendMails, int64(2), nil)
			},
			setupAuth: func(c *gin.Context, userID primitive.ObjectID) {
				c.Set("authUser", &auth.UserAuthInfo{UserID: userID})
			},
		},
		{
			name:              "Success with pagination parameters",
			sinceParam:        "2024-01-01T00:00:00Z",
			queryParams:       "?page=2&size=15",
			expectedStatus:    http.StatusOK,
			expectedSendMails: 1,
			setupMock: func(mockRepo *mocks.MockSendMailRepository, userID primitive.ObjectID, sinceTime time.Time, page, limit int64) {
				mailID := primitive.NewObjectID()
				sendMails := []*models.SendMail{
					{
						Mail: &models.Mail{
							ID:     &mailID,
							UserID: userID,
							Headers: map[string]string{
								"Subject": "User Send Mail",
							},
						},
						SendStatus: models.SendStatusPending,
					},
				}
				mockRepo.On("GetSince", mock.Anything, userID, sinceTime, page, limit).Return(sendMails, int64(16), nil)
			},
			setupAuth: func(c *gin.Context, userID primitive.ObjectID) {
				c.Set("authUser", &auth.UserAuthInfo{UserID: userID})
			},
		},
		{
			name:              "Success with no send mails found",
			sinceParam:        "2024-01-01T00:00:00Z",
			queryParams:       "",
			expectedStatus:    http.StatusOK,
			expectedSendMails: 0,
			setupMock: func(mockRepo *mocks.MockSendMailRepository, userID primitive.ObjectID, sinceTime time.Time, page, limit int64) {
				sendMails := []*models.SendMail{}
				mockRepo.On("GetSince", mock.Anything, userID, sinceTime, page, limit).Return(sendMails, int64(0), nil)
			},
			setupAuth: func(c *gin.Context, userID primitive.ObjectID) {
				c.Set("authUser", &auth.UserAuthInfo{UserID: userID})
			},
		},
		{
			name:              "Error when repository fails",
			sinceParam:        "2024-01-01T00:00:00Z",
			queryParams:       "",
			expectedStatus:    http.StatusInternalServerError,
			expectedSendMails: 0,
			setupMock: func(mockRepo *mocks.MockSendMailRepository, userID primitive.ObjectID, sinceTime time.Time, page, limit int64) {
				mockRepo.On("GetSince", mock.Anything, userID, sinceTime, page, limit).Return(nil, int64(0), assert.AnError)
			},
			setupAuth: func(c *gin.Context, userID primitive.ObjectID) {
				c.Set("authUser", &auth.UserAuthInfo{UserID: userID})
			},
		},
		{
			name:              "Error when missing since parameter",
			sinceParam:        "",
			queryParams:       "",
			expectedStatus:    http.StatusBadRequest,
			expectedSendMails: 0,
			setupMock: func(mockRepo *mocks.MockSendMailRepository, userID primitive.ObjectID, sinceTime time.Time, page, limit int64) {
				// No mock setup needed as the request should fail before reaching the repository
			},
			setupAuth: func(c *gin.Context, userID primitive.ObjectID) {
				c.Set("authUser", &auth.UserAuthInfo{UserID: userID})
			},
		},
		{
			name:              "Error when invalid date format",
			sinceParam:        "invalid-date",
			queryParams:       "",
			expectedStatus:    http.StatusBadRequest,
			expectedSendMails: 0,
			setupMock: func(mockRepo *mocks.MockSendMailRepository, userID primitive.ObjectID, sinceTime time.Time, page, limit int64) {
				// No mock setup needed as the request should fail before reaching the repository
			},
			setupAuth: func(c *gin.Context, userID primitive.ObjectID) {
				c.Set("authUser", &auth.UserAuthInfo{UserID: userID})
			},
		},
		{
			name:              "Error when unauthorized",
			sinceParam:        "2024-01-01T00:00:00Z",
			queryParams:       "",
			expectedStatus:    http.StatusUnauthorized,
			expectedSendMails: 0,
			setupMock: func(mockRepo *mocks.MockSendMailRepository, userID primitive.ObjectID, sinceTime time.Time, page, limit int64) {
				// No mock setup needed as the request should fail before reaching the repository
			},
			setupAuth: func(c *gin.Context, userID primitive.ObjectID) {
				// No auth setup - should trigger unauthorized
			},
		},
		{
			name:              "Success with different date formats",
			sinceParam:        "2024-01-01T12:30:45+02:00",
			queryParams:       "",
			expectedStatus:    http.StatusOK,
			expectedSendMails: 1,
			setupMock: func(mockRepo *mocks.MockSendMailRepository, userID primitive.ObjectID, sinceTime time.Time, page, limit int64) {
				mailID := primitive.NewObjectID()
				sendMails := []*models.SendMail{
					{
						Mail: &models.Mail{
							ID:     &mailID,
							UserID: userID,
							Headers: map[string]string{
								"Subject": "User Send Mail",
							},
						},
						SendStatus: models.SendStatusPending,
					},
				}
				mockRepo.On("GetSince", mock.Anything, userID, sinceTime, page, limit).Return(sendMails, int64(1), nil)
			},
			setupAuth: func(c *gin.Context, userID primitive.ObjectID) {
				c.Set("authUser", &auth.UserAuthInfo{UserID: userID})
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			mockRepo := &mocks.MockSendMailRepository{}
			userID := primitive.NewObjectID()

			// Parse the since parameter to get the expected time
			var sinceTime time.Time
			if tt.sinceParam != "" {
				var err error
				sinceTime, err = time.Parse(time.RFC3339, tt.sinceParam)
				if err != nil {
					// Try parsing with timezone offset format
					sinceTime, err = time.Parse("2006-01-02T15:04:05-07:00", tt.sinceParam)
					if err != nil {
						// For invalid date tests, we don't need to parse
						sinceTime = time.Now()
					}
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

			controller := NewSendMailController(mockRepo, nil, nil, nil)

			router := gin.New()
			router.Use(func(c *gin.Context) {
				tt.setupAuth(c, userID)
				// Set pagination parameters in context
				c.Set("page", int(page))
				c.Set("size", int(limit))
				c.Next()
			})

			sendMailRoutes := router.Group("/mail/send")
			{
				sendMailRoutes.GET("/since", pagination.New(), controller.GetSendMailsSince)
			}

			// Create request with since parameter
			reqURL := "/mail/send/since"
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
				var response PaginatedSendMailResponse
				err := json.Unmarshal(w.Body.Bytes(), &response)
				assert.NoError(t, err)
				assert.Len(t, response.SendMails, tt.expectedSendMails)
				assert.NotNil(t, response.TotalCount)
				assert.NotNil(t, response.Page)
				assert.NotNil(t, response.Size)
				assert.NotNil(t, response.TotalPages)

				// Verify all returned send mails belong to the authenticated user
				for _, sendMail := range response.SendMails {
					assert.NotNil(t, sendMail.Mail)
					assert.Equal(t, userID, sendMail.Mail.UserID)
				}
			}

			mockRepo.AssertExpectations(t)
		})
	}
}
