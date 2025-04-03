package habits

import (
	"atomic_blend_api/models"
	"atomic_blend_api/tests/mocks"
	"time"

	"github.com/gin-gonic/gin"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

// createTestHabit creates a test habit entity with default values
func createTestHabit() *models.Habit {
	name := "Test Habit"
	emoji := "🏃‍♂️"
	frequency := models.FrequencyDaily
	numberOfTimes := 3
	daysOfWeek := []int{1, 3, 5} // Monday, Wednesday, Friday
	startDate := time.Now().Format(time.RFC3339)
	endDate := time.Now().AddDate(0, 1, 0).Format(time.RFC3339) // One month later
	now := time.Now().Format(time.RFC3339)
	reminders := []string{"09:00", "18:00"}
	citation := "Test citation"

	return &models.Habit{
		ID:            primitive.NewObjectID(),
		UserID:        primitive.NewObjectID(),
		Name:          &name,
		Emoji:         &emoji,
		Frequency:     &frequency,
		NumberOfTimes: &numberOfTimes,
		DaysOfWeek:    &daysOfWeek,
		StartDate:     &startDate,
		EndDate:       &endDate,
		CreatedAt:     &now,
		UpdatedAt:     &now,
		Reminders:     reminders,
		Citation:      &citation,
	}
}

// setupTest creates a test router and mock repository for testing
func setupTest() (*gin.Engine, *mocks.MockHabitRepository) {
	gin.SetMode(gin.TestMode)
	router := gin.New()
	mockRepo := new(mocks.MockHabitRepository)
	habitController := NewHabitController(mockRepo)

	// Set up routes with middleware
	habitRoutes := router.Group("/habits")
	{
		habitRoutes.GET("", habitController.GetAllHabits)
		habitRoutes.GET("/:id", habitController.GetHabitByID)
		habitRoutes.POST("", habitController.CreateHabit)
		habitRoutes.PUT("/:id", habitController.UpdateHabit)
		habitRoutes.DELETE("/:id", habitController.DeleteHabit)
	}

	return router, mockRepo
}
