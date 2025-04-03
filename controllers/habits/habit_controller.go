package habits

import (
	"atomic_blend_api/auth"
	"atomic_blend_api/repositories"

	"github.com/gin-gonic/gin"
	"go.mongodb.org/mongo-driver/mongo"
)

// HabitController handles habit related operations
type HabitController struct {
	habitRepo repositories.HabitRepositoryInterface
}

// NewHabitController creates a new habit controller instance
func NewHabitController(habitRepo repositories.HabitRepositoryInterface) *HabitController {
	return &HabitController{
		habitRepo: habitRepo,
	}
}

// SetupRoutes sets up the routes for the habit controller
func SetupRoutes(router *gin.Engine, database *mongo.Database) {
	habitRepo := repositories.NewHabitRepository(database)
	habitController := NewHabitController(habitRepo)
	setupHabitRoutes(router, habitController)
}

// SetupRoutesWithMock sets up the habit routes with a mock repository for testing
func SetupRoutesWithMock(router *gin.Engine, habitRepo repositories.HabitRepositoryInterface) {
	habitController := NewHabitController(habitRepo)
	setupHabitRoutes(router, habitController)
}

// setupHabitRoutes sets up the routes for habit controller
func setupHabitRoutes(router *gin.Engine, habitController *HabitController) {
	habitRoutes := router.Group("/habits")
	auth.RequireAuth(habitRoutes)
	{
		habitRoutes.POST("", habitController.CreateHabit)
		habitRoutes.GET("", habitController.GetAllHabits)
		habitRoutes.GET("/:id", habitController.GetHabitByID)
		habitRoutes.PUT("/:id", habitController.UpdateHabit)
		habitRoutes.DELETE("/:id", habitController.DeleteHabit)
	}
}
