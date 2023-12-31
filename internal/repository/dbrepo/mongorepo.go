package dbrepo

import (
	"context"
	"github.com/Ldepner/auth-project/internal/models"
	"go.mongodb.org/mongo-driver/bson"
)

func (r *mongoDBRepo) GetUserRecordByEmail(email string) (*models.UserRecord, error) {
	var user *models.UserRecord
	users := r.DB.Collection("users")

	filter := bson.D{{"email", email}}
	err := users.FindOne(context.TODO(), filter).Decode(&user)
	if err != nil {
		return nil, err
	}

	return user, nil
}

func (r *mongoDBRepo) CreateUserRecord(user *models.UserRecord) error {

	return nil
}
