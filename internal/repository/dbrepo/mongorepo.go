package dbrepo

import (
	"context"
	"errors"
	"github.com/Ldepner/auth-project/internal/helpers"
	"github.com/Ldepner/auth-project/internal/models"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"golang.org/x/crypto/bcrypt"
)

func (r *mongoDBRepo) GetUserRecordByEmail(email string) (*models.UserRecord, error) {
	var user *models.UserRecord
	users := r.DB.Collection("users")

	filter := bson.D{{"email", email}}
	err := users.FindOne(context.TODO(), filter).Decode(&user)
	if err != nil {
		if errors.Is(err, mongo.ErrNoDocuments) {
			return nil, &helpers.ErrRecordNotFound{Err: err}
		}
		return nil, err
	}

	return user, nil
}

func (r *mongoDBRepo) CreateUserRecord(user *models.UserRecord) error {
	users := r.DB.Collection("users")

	// Does record with email already exist?
	filter := bson.D{{"email", user.Email}}

	var blankUser *models.UserRecord
	err := users.FindOne(context.TODO(), filter).Decode(&blankUser)

	if len(blankUser.Email) > 0 {
		return &helpers.ErrDuplicateEmail{}
	}

	// hash password
	hashedBytes, err := bcrypt.GenerateFromPassword([]byte(user.Password), 14)
	if err != nil {
		return err
	}
	newUserRecord := models.UserRecord{
		Email:    user.Email,
		Password: string(hashedBytes),
	}
	_, err = users.InsertOne(context.TODO(), newUserRecord)
	if err != nil {
		return err
	}

	return nil
}
