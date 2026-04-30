package graph

import (
	"context"
	"errors"
	"fmt"
	"strings"
	"time"

	"github.com/vod/graph/model"
	"github.com/vod/internal/auth"
	"go.mongodb.org/mongo-driver/v2/bson"
	"go.mongodb.org/mongo-driver/v2/mongo"
	"golang.org/x/crypto/bcrypt"
)

func Signin(ctx context.Context, db *mongo.Database, input model.SigninInput) (*model.AuthPayload, error) {
	if db == nil {
		return nil, errors.New("database is not configured")
	}

	email := strings.ToLower(strings.TrimSpace(input.Email))
	if email == "" {
		return nil, errors.New("email is required")
	}

	password := strings.TrimSpace(input.Password)
	if password == "" {
		return nil, errors.New("password is required")
	}

	var user userDocument
	err := db.Collection("users").FindOne(ctx, bson.M{"email": email}).Decode(&user)
	if err != nil {
		if errors.Is(err, mongo.ErrNoDocuments) {
			return nil, errors.New("user not found")
		}
		return nil, err
	}

	if err := bcrypt.CompareHashAndPassword([]byte(user.PasswordHash), []byte(password)); err != nil {
		return nil, errors.New("invalid password")
	}

	token, err := auth.GenerateSignupToken(user.ID, user.Username, user.Email)
	if err != nil {
		return nil, fmt.Errorf("generate paseto token: %w", err)
	}

	return &model.AuthPayload{
		Token: token,
		User: &model.PublicUser{
			ID:          user.ID,
			Username:    user.Username,
			Email:       user.Email,
			DisplayName: user.DisplayName,
			CreatedAt:   user.CreatedAt.Format(time.RFC3339),
		},
	}, nil
}
