package users

import (
	"github.com/gin-gonic/gin"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

// deletePersonalData handles the deletion of all user personal data
// This includes tasks and any other personal data associated with the user
func (c *UserController) DeletePersonalData(ctx *gin.Context, userID primitive.ObjectID) error {
	// Use the task repository factory to create a task repository
	taskRepo := c.getTaskRepository()

	// Get all tasks for the user
	tasks, err := taskRepo.GetAll(ctx, &userID)
	if err != nil {
		return err
	}

	// Delete each task
	for _, task := range tasks {
		if err := taskRepo.Delete(ctx, task.ID); err != nil {
			return err
		}
	}

	// TODO: delete habits, tags and other personal data

	return nil
}