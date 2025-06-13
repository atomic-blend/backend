package notes

import (
	"atomic_blend_api/auth"
	"atomic_blend_api/repositories"

	"github.com/gin-gonic/gin"
	"go.mongodb.org/mongo-driver/mongo"
)

// NoteController handles note related operations
type NoteController struct {
	noteRepo repositories.NoteRepositoryInterface
}

// NewNoteController creates a new note controller instance
func NewNoteController(noteRepo repositories.NoteRepositoryInterface) *NoteController {
	return &NoteController{
		noteRepo: noteRepo,
	}
}

// SetupRoutes sets up the note routes
func SetupRoutes(router *gin.Engine, database *mongo.Database) {
	noteRepo := repositories.NewNoteRepository(database)
	noteController := NewNoteController(noteRepo)
	setupNoteRoutes(router, noteController)
}

// SetupRoutesWithMock sets up the note routes with a mock repository for testing
func SetupRoutesWithMock(router *gin.Engine, noteRepo repositories.NoteRepositoryInterface) {
	noteController := NewNoteController(noteRepo)
	setupNoteRoutes(router, noteController)
}

// setupNoteRoutes sets up the routes for note controller
func setupNoteRoutes(router *gin.Engine, noteController *NoteController) {
	noteRoutes := router.Group("/notes")
	auth.RequireAuth(noteRoutes)
	{
		noteRoutes.GET("", noteController.GetAllNotes)
		noteRoutes.GET("/:id", noteController.GetNoteByID)
		noteRoutes.POST("", noteController.CreateNote)
		noteRoutes.PUT("/:id", noteController.UpdateNote)
		noteRoutes.DELETE("/:id", noteController.DeleteNote)
	}
}
