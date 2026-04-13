package graph

import (
	"context"
	"errors"
	"fmt"
	"strings"
	"time"

	"github.com/google/uuid"
	"github.com/vod/graph/model"
	"go.mongodb.org/mongo-driver/v2/mongo"
	"golang.org/x/crypto/bcrypt"
)

type userDocument struct {
	ID               string    `bson:"id"`
	Username         string    `bson:"username"`
	Email            string    `bson:"email"`
	PasswordHash     string    `bson:"passwordHash"`
	DisplayName      *string   `bson:"displayName,omitempty"`
	AccountType      string    `bson:"accountType"`
	Role             string    `bson:"role"`
	FollowersCount   int       `bson:"followersCount"`
	FollowingCount   int       `bson:"followingCount"`
	VideosCount      int       `bson:"videosCount"`
	LikesReceived    int       `bson:"likesReceivedCount"`
	IsVerified       bool      `bson:"isVerified"`
	IsActive         bool      `bson:"isActive"`
	OpenToWork       bool      `bson:"openToWork"`
	CreatedAt        time.Time `bson:"created_at"`
	UpdatedAt        time.Time `bson:"updated_at"`
}

func signup(ctx context.Context, db *mongo.Database, input model.SignupInput) (*model.AuthPayload, error) {
	if db == nil {
		return nil, errors.New("database is not configured")
	}

	username := strings.TrimSpace(input.Username)
	email := strings.ToLower(strings.TrimSpace(input.Email))
	password := strings.TrimSpace(input.Password)
	if username == "" || email == "" || password == "" {
		return nil, errors.New("username, email, and password are required")
	}
	if len(password) < 8 {
		return nil, errors.New("password must be at least 8 characters")
	}

	hash, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	if err != nil {
		return nil, fmt.Errorf("hash password: %w", err)
	}

	now := time.Now().UTC()
	userID := uuid.NewString()
	doc := userDocument{
		ID:             userID,
		Username:       username,
		Email:          email,
		PasswordHash:   string(hash),
		DisplayName:    input.DisplayName,
		AccountType:    "candidate",
		Role:           "user",
		FollowersCount: 0,
		FollowingCount: 0,
		VideosCount:    0,
		LikesReceived:  0,
		IsVerified:     false,
		IsActive:       true,
		OpenToWork:     false,
		CreatedAt:      now,
		UpdatedAt:      now,
	}

	if _, err := db.Collection("users").InsertOne(ctx, doc); err != nil {
		if mongo.IsDuplicateKeyError(err) {
			return nil, errors.New("username or email already exists")
		}
		return nil, fmt.Errorf("create user: %w", err)
	}

	token := uuid.NewString()
	return &model.AuthPayload{
		Token: token,
		User: &model.PublicUser{
			ID:          userID,
			Username:    username,
			Email:       email,
			DisplayName: input.DisplayName,
			CreatedAt:   now.Format(time.RFC3339),
		},
	}, nil
}

