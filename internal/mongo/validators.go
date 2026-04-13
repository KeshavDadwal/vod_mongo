package mongo

import (
	"context"
	"fmt"

	"go.mongodb.org/mongo-driver/v2/bson"
	"go.mongodb.org/mongo-driver/v2/mongo"
)

// EnsureValidators applies lightweight MongoDB JSON Schema validators where possible.
// Note: MongoDB cannot enforce foreign keys, and cannot easily enforce "followerId != followingId".
func EnsureValidators(ctx context.Context, db *mongo.Database) error {
	// reports CHECK: exactly one of videoId, reportedUserId
	// Implemented as: require reporterId + reason, and oneOf:
	// - videoId present and reportedUserId absent
	// - reportedUserId present and videoId absent
	reportsValidator := bson.D{
		{Key: "$jsonSchema", Value: bson.D{
			{Key: "bsonType", Value: "object"},
			{Key: "required", Value: bson.A{"id", "reporterId", "reason", "created_at"}},
			{Key: "properties", Value: bson.D{
				{Key: "id", Value: bson.D{{Key: "bsonType", Value: "string"}}},
				{Key: "reporterId", Value: bson.D{{Key: "bsonType", Value: "string"}}},
				{Key: "videoId", Value: bson.D{{Key: "bsonType", Value: bson.A{"string", "null"}}}},
				{Key: "reportedUserId", Value: bson.D{{Key: "bsonType", Value: bson.A{"string", "null"}}}},
				{Key: "reason", Value: bson.D{{Key: "bsonType", Value: "string"}}},
				{Key: "created_at", Value: bson.D{{Key: "bsonType", Value: "date"}}},
			}},
			{Key: "oneOf", Value: bson.A{
				bson.D{{Key: "required", Value: bson.A{"videoId"}}, {Key: "not", Value: bson.D{{Key: "required", Value: bson.A{"reportedUserId"}}}}},
				bson.D{{Key: "required", Value: bson.A{"reportedUserId"}}, {Key: "not", Value: bson.D{{Key: "required", Value: bson.A{"videoId"}}}}},
			}},
		}},
	}

	if err := collMod(ctx, db, "reports", reportsValidator); err != nil {
		return err
	}

	return nil
}

func collMod(ctx context.Context, db *mongo.Database, collection string, validator bson.D) error {
	cmd := bson.D{
		{Key: "collMod", Value: collection},
		{Key: "validator", Value: validator},
		{Key: "validationLevel", Value: "moderate"},
	}
	if err := db.RunCommand(ctx, cmd).Err(); err != nil {
		// If collection doesn't exist yet, create it with the validator.
		createCmd := bson.D{
			{Key: "create", Value: collection},
			{Key: "validator", Value: validator},
			{Key: "validationLevel", Value: "moderate"},
		}
		if err2 := db.RunCommand(ctx, createCmd).Err(); err2 != nil {
			return fmt.Errorf("apply validator to %s: %v (and create failed: %w)", collection, err, err2)
		}
	}
	return nil
}

