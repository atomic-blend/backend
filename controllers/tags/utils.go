package tags

import "go.mongodb.org/mongo-driver/bson/primitive"

// removeTagFromSlice filters out a specific tag ID from a slice of tag IDs
func removeTagFromSlice(tags []primitive.ObjectID, tagToRemove primitive.ObjectID) []primitive.ObjectID {
	result := make([]primitive.ObjectID, 0, len(tags))
	for _, tag := range tags {
		if tag != tagToRemove {
			result = append(result, tag)
		}
	}
	return result
}
