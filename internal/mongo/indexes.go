package mongo

import (
	"context"
	"fmt"

	"go.mongodb.org/mongo-driver/v2/bson"
	"go.mongodb.org/mongo-driver/v2/mongo"
	"go.mongodb.org/mongo-driver/v2/mongo/options"
)

type indexSpec struct {
	Collection string
	Keys       bson.D
	Unique     bool
	Name       string
}

func EnsureIndexes(ctx context.Context, db *mongo.Database) error {
	indexes := []indexSpec{
		// users
		{Collection: "users", Keys: bson.D{{Key: "id", Value: 1}}, Unique: true, Name: "uniq_id"},
		{Collection: "users", Keys: bson.D{{Key: "username", Value: 1}}, Unique: true, Name: "uniq_username"},
		{Collection: "users", Keys: bson.D{{Key: "email", Value: 1}}, Unique: true, Name: "uniq_email"},

		// skills
		{Collection: "skills", Keys: bson.D{{Key: "id", Value: 1}}, Unique: true, Name: "uniq_id"},
		{Collection: "skills", Keys: bson.D{{Key: "slug", Value: 1}}, Unique: true, Name: "uniq_slug"},

		// categories
		{Collection: "categories", Keys: bson.D{{Key: "id", Value: 1}}, Unique: true, Name: "uniq_id"},
		{Collection: "categories", Keys: bson.D{{Key: "slug", Value: 1}}, Unique: true, Name: "uniq_slug"},

		// videos
		{Collection: "videos", Keys: bson.D{{Key: "id", Value: 1}}, Unique: true, Name: "uniq_id"},
		{Collection: "videos", Keys: bson.D{{Key: "userId", Value: 1}, {Key: "created_at", Value: -1}}, Unique: false, Name: "by_user_created_at"},

		// user_skills UNIQUE(user_id, skill_id)
		{Collection: "user_skills", Keys: bson.D{{Key: "id", Value: 1}}, Unique: true, Name: "uniq_id"},
		{Collection: "user_skills", Keys: bson.D{{Key: "userId", Value: 1}, {Key: "skillId", Value: 1}}, Unique: true, Name: "uniq_user_skill"},

		// user_interests UNIQUE(user_id, category_id)
		{Collection: "user_interests", Keys: bson.D{{Key: "id", Value: 1}}, Unique: true, Name: "uniq_id"},
		{Collection: "user_interests", Keys: bson.D{{Key: "userId", Value: 1}, {Key: "categoryId", Value: 1}}, Unique: true, Name: "uniq_user_category"},

		// video_categories UNIQUE(video_id, category_id)
		{Collection: "video_categories", Keys: bson.D{{Key: "id", Value: 1}}, Unique: true, Name: "uniq_id"},
		{Collection: "video_categories", Keys: bson.D{{Key: "videoId", Value: 1}, {Key: "categoryId", Value: 1}}, Unique: true, Name: "uniq_video_category"},

		// video_tags UNIQUE(video_id, raw_tag)
		{Collection: "video_tags", Keys: bson.D{{Key: "id", Value: 1}}, Unique: true, Name: "uniq_id"},
		{Collection: "video_tags", Keys: bson.D{{Key: "videoId", Value: 1}, {Key: "rawTag", Value: 1}}, Unique: true, Name: "uniq_video_rawTag"},

		// follows UNIQUE(follower_id, following_id)
		{Collection: "follows", Keys: bson.D{{Key: "id", Value: 1}}, Unique: true, Name: "uniq_id"},
		{Collection: "follows", Keys: bson.D{{Key: "followerId", Value: 1}, {Key: "followingId", Value: 1}}, Unique: true, Name: "uniq_follower_following"},

		// likes UNIQUE(user_id, video_id)
		{Collection: "likes", Keys: bson.D{{Key: "id", Value: 1}}, Unique: true, Name: "uniq_id"},
		{Collection: "likes", Keys: bson.D{{Key: "userId", Value: 1}, {Key: "videoId", Value: 1}}, Unique: true, Name: "uniq_user_video"},

		// saves UNIQUE(user_id, video_id)
		{Collection: "saves", Keys: bson.D{{Key: "id", Value: 1}}, Unique: true, Name: "uniq_id"},
		{Collection: "saves", Keys: bson.D{{Key: "userId", Value: 1}, {Key: "videoId", Value: 1}}, Unique: true, Name: "uniq_user_video"},

		// comments
		{Collection: "comments", Keys: bson.D{{Key: "id", Value: 1}}, Unique: true, Name: "uniq_id"},
		{Collection: "comments", Keys: bson.D{{Key: "videoId", Value: 1}, {Key: "created_at", Value: -1}}, Unique: false, Name: "by_video_created_at"},
		{Collection: "comments", Keys: bson.D{{Key: "parentCommentId", Value: 1}}, Unique: false, Name: "by_parent"},

		// video_views
		{Collection: "video_views", Keys: bson.D{{Key: "id", Value: 1}}, Unique: true, Name: "uniq_id"},
		{Collection: "video_views", Keys: bson.D{{Key: "videoId", Value: 1}, {Key: "created_at", Value: -1}}, Unique: false, Name: "by_video_created_at"},

		// notifications
		{Collection: "notifications", Keys: bson.D{{Key: "id", Value: 1}}, Unique: true, Name: "uniq_id"},
		{Collection: "notifications", Keys: bson.D{{Key: "recipientId", Value: 1}, {Key: "isRead", Value: 1}, {Key: "created_at", Value: -1}}, Unique: false, Name: "inbox"},

		// reports
		{Collection: "reports", Keys: bson.D{{Key: "id", Value: 1}}, Unique: true, Name: "uniq_id"},
		{Collection: "reports", Keys: bson.D{{Key: "reporterId", Value: 1}, {Key: "created_at", Value: -1}}, Unique: false, Name: "by_reporter_created_at"},
	}

	for _, spec := range indexes {
		coll := db.Collection(spec.Collection)
		model := mongo.IndexModel{
			Keys: spec.Keys,
			Options: options.Index().
				SetName(spec.Name).
				SetUnique(spec.Unique),
		}
		if _, err := coll.Indexes().CreateOne(ctx, model); err != nil {
			return fmt.Errorf("create index %s on %s: %w", spec.Name, spec.Collection, err)
		}
	}

	return nil
}

