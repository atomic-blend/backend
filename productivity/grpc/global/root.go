package user

import (
	"github.com/atomic-blend/backend/productivity/repositories"
)

// grpcServer is the gRPC server for the productivity service
type grpcServer struct {
	taskRepo      repositories.TaskRepositoryInterface
	habitRepo     repositories.HabitRepositoryInterface
	noteRepo      repositories.NoteRepositoryInterface
	tagRepo       repositories.TagRepositoryInterface
	folderRepo    repositories.FolderRepositoryInterface
	timeEntryRepo repositories.TimeEntryRepositoryInterface
}

// NewGrpcServer create a new instance of GrpcServer
func NewGrpcServer(taskRepo repositories.TaskRepositoryInterface, habitRepo repositories.HabitRepositoryInterface, noteRepo repositories.NoteRepositoryInterface, tagRepo repositories.TagRepositoryInterface, folderRepo repositories.FolderRepositoryInterface, timeEntryRepo repositories.TimeEntryRepositoryInterface) *grpcServer {
	return &grpcServer{
		taskRepo:      taskRepo,
		habitRepo:     habitRepo,
		noteRepo:      noteRepo,
		tagRepo:       tagRepo,
		folderRepo:    folderRepo,
		timeEntryRepo: timeEntryRepo,
	}
}
