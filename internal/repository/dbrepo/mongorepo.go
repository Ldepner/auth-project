package dbrepo

import (
	"context"
	"errors"
	"github.com/Ldepner/auth-project/internal/helpers"
	"github.com/Ldepner/auth-project/internal/models"
	"github.com/go-webauthn/webauthn/webauthn"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"golang.org/x/crypto/bcrypt"
	"time"
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

func (r *mongoDBRepo) GetUserRecordByID(id string) (*models.UserRecord, error) {
	var user *models.UserRecord
	users := r.DB.Collection("users")
	userID, err := primitive.ObjectIDFromHex(id)

	filter := bson.D{{"_id", userID}}
	err = users.FindOne(context.TODO(), filter).Decode(&user)
	if err != nil {
		if errors.Is(err, mongo.ErrNoDocuments) {
			return nil, &helpers.ErrRecordNotFound{Err: err}
		}
		return nil, err
	}

	return user, nil
}

func (r *mongoDBRepo) CreateUserRecord(user *models.UserRecord) (string, error) {
	users := r.DB.Collection("users")

	// Does record with email already exist?
	filter := bson.D{{"email", user.Email}}

	var blankUser *models.UserRecord
	err := users.FindOne(context.TODO(), filter).Decode(&blankUser)

	if len(blankUser.Email) > 0 {
		return "", &helpers.ErrDuplicateEmail{}
	}

	// hash password
	hashedBytes, err := bcrypt.GenerateFromPassword([]byte(user.Password), 14)
	if err != nil {
		return "", err
	}
	newUserRecord := models.UserRecord{
		Email:    user.Email,
		Password: string(hashedBytes),
	}
	result, err := users.InsertOne(context.TODO(), newUserRecord)
	if err != nil {
		return "", err
	}

	id := result.InsertedID.(primitive.ObjectID).Hex()

	return id, nil
}

func (r *mongoDBRepo) UpdateUserRecord(userID string, field string, updatedValue any) error {
	users := r.DB.Collection("users")
	id, _ := primitive.ObjectIDFromHex(userID)
	filter := bson.D{{"_id", id}}
	// Creates instructions to add the "avg_rating" field to documents
	update := bson.D{{"$set", bson.D{{field, updatedValue}}}}
	// Updates the first document that has the specified "_id" value
	_, err := users.UpdateOne(context.TODO(), filter, update)
	if err != nil {
		panic(err)
	}
	return nil
}

func (r *mongoDBRepo) CreateSession(userID string, lastActiveAt, expiresAt time.Time, authenticated bool) (string, error) {
	sessions := r.DB.Collection("sessions")

	newSession := models.Session{
		UserID:        userID,
		Authenticated: authenticated,
		LastActiveAt:  lastActiveAt,
		ExpiresAt:     expiresAt,
	}
	result, err := sessions.InsertOne(context.TODO(), newSession)
	if err != nil {
		return "", err
	}
	token := result.InsertedID.(primitive.ObjectID).Hex()

	return token, nil
}

func (r *mongoDBRepo) DeleteSession(token string) error {
	id, err := primitive.ObjectIDFromHex(token)
	if err != nil {
		return err
	}

	sessions := r.DB.Collection("sessions")
	filter := bson.D{{"_id", id}}
	_, err = sessions.DeleteOne(context.TODO(), filter, nil)
	if err != nil {
		return err
	}

	return nil
}

func (r *mongoDBRepo) UpdateSessionAuthenticated(token string, authenticated bool) error {
	sessions := r.DB.Collection("sessions")
	id, _ := primitive.ObjectIDFromHex(token)
	filter := bson.D{{"_id", id}}
	// Creates instructions to add the "avg_rating" field to documents
	update := bson.D{{"$set", bson.D{{"authenticated", authenticated}}}}
	// Updates the first document that has the specified "_id" value
	_, err := sessions.UpdateOne(context.TODO(), filter, update)
	if err != nil {
		panic(err)
	}
	return nil
}

func (r *mongoDBRepo) UpdateWebAuthNSession(token string, sessionType string, webauthnSessionData *webauthn.SessionData) error {
	sessions := r.DB.Collection("sessions")
	id, _ := primitive.ObjectIDFromHex(token)
	filter := bson.D{{"_id", id}}
	// Creates instructions to add the "avg_rating" field to documents
	update := bson.D{{"$set", bson.D{{sessionType, webauthnSessionData}}}}
	// Updates the first document that has the specified "_id" value
	_, err := sessions.UpdateOne(context.TODO(), filter, update)
	if err != nil {
		panic(err)
	}
	return nil
}

func (r *mongoDBRepo) GetSessionByID(token string) (*models.Session, error) {
	var session *models.Session

	id, err := primitive.ObjectIDFromHex(token)
	if err != nil {
		return nil, err
	}

	sessions := r.DB.Collection("sessions")
	filter := bson.D{{"_id", id}}
	err = sessions.FindOne(context.TODO(), filter, nil).Decode(&session)
	if err != nil {
		if errors.Is(err, mongo.ErrNoDocuments) {
			return nil, &helpers.ErrRecordNotFound{Err: err}
		}
		return nil, err
	}

	return session, nil
}
