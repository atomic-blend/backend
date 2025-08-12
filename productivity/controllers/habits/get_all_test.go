package habits

import (
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"github.com/atomic-blend/backend/shared/middlewares/auth"
	"github.com/atomic-blend/backend/productivity/models"
	"testing"

	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

func TestGetAllHabits(t *testing.T) {
	router, mockRepo := setupTest()

	t.Run("successful get all habits - with habits", func(t *testing.T) {
		// Create authenticated user
		userID := primitive.NewObjectID()

		// Create test habits
		habit1 := createTestHabit()
		habit1.UserID = userID
		name1 := "Habit 1"
		habit1.Name = &name1

		habit2 := createTestHabit()
		habit2.UserID = userID
		name2 := "Habit 2"
		habit2.Name = &name2

		habits := []*models.Habit{habit1, habit2}

		mockRepo.On("GetAll", mock.Anything, mock.MatchedBy(func(userID *primitive.ObjectID) bool {
			return true
		})).Return(habits, nil).Once()

		// Mock the GetEntriesByHabitID calls for each habit
		mockRepo.On("GetEntriesByHabitID", mock.Anything, habit1.ID).Return([]models.HabitEntry{}, nil).Once()
		mockRepo.On("GetEntriesByHabitID", mock.Anything, habit2.ID).Return([]models.HabitEntry{}, nil).Once()

		w := httptest.NewRecorder()
		req, _ := http.NewRequest("GET", "/habits", nil)

		// Create a new context with the request and set auth user
		ctx, _ := gin.CreateTestContext(w)
		ctx.Request = req
		ctx.Set("authUser", &auth.UserAuthInfo{UserID: userID})

		// Call the controller directly with our context that has auth
		controller := NewHabitController(mockRepo)
		controller.GetAllHabits(ctx)

		assert.Equal(t, http.StatusOK, w.Code)
		var response []*models.Habit
		err := json.Unmarshal(w.Body.Bytes(), &response)
		assert.NoError(t, err)
		assert.Equal(t, 2, len(response))
		assert.Equal(t, *habit1.Name, *response[0].Name)
		assert.Equal(t, *habit2.Name, *response[1].Name)
	})

	t.Run("successful get all habits - empty list", func(t *testing.T) {
		// Create authenticated user
		userID := primitive.NewObjectID()

		// Return nil to simulate no habits
		mockRepo.On("GetAll", mock.Anything, mock.MatchedBy(func(userID *primitive.ObjectID) bool {
			return true
		})).Return(nil, nil).Once()

		// No need to mock GetEntriesByHabitID here since there are no habits

		w := httptest.NewRecorder()
		req, _ := http.NewRequest("GET", "/habits", nil)

		// Create a new context with the request and set auth user
		ctx, _ := gin.CreateTestContext(w)
		ctx.Request = req
		ctx.Set("authUser", &auth.UserAuthInfo{UserID: userID})

		// Call the controller directly with our context that has auth
		controller := NewHabitController(mockRepo)
		controller.GetAllHabits(ctx)

		assert.Equal(t, http.StatusOK, w.Code)
		// Check that the response is an empty array, not null
		assert.Equal(t, "[]", w.Body.String())
	})

	t.Run("unauthorized access - no auth user", func(t *testing.T) {
		w := httptest.NewRecorder()
		req, _ := http.NewRequest("GET", "/habits", nil)

		// Call the endpoint without authentication
		router.ServeHTTP(w, req)

		assert.Equal(t, http.StatusUnauthorized, w.Code)
		var response map[string]interface{}
		err := json.Unmarshal(w.Body.Bytes(), &response)
		assert.NoError(t, err)
		assert.Equal(t, "Authentication required", response["error"])
	})
}
