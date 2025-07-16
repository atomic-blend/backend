package tags

import (
	"atomic-blend/backend/productivity/models"

	"go.mongodb.org/mongo-driver/bson/primitive"
)

// removeTagFromSlice filters out a specific tag ID from a slice of Tag objects
func removeTagFromSlice(tags []*models.Tag, tagToRemove primitive.ObjectID) []*models.Tag {
	result := make([]*models.Tag, 0, len(tags))
	for _, tag := range tags {
		if tag.ID != nil && *tag.ID != tagToRemove {
			result = append(result, tag)
		}
	}
	return result
}
