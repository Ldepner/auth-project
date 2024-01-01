package models

type UserRecord struct {
	Email    string `bson:"email,omitempty"`
	Password string `bson:"password,omitempty"`
}
