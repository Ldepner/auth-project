package dbrepo

import (
	"github.com/Ldepner/auth-project/internal/config"
	"github.com/Ldepner/auth-project/internal/repository"
	"go.mongodb.org/mongo-driver/mongo"
)

type mongoDBRepo struct {
	App *config.AppConfig
	DB  *mongo.Database
}

func NewMongoRepo(conn *mongo.Database, a *config.AppConfig) repository.DBRepo {
	return &mongoDBRepo{
		App: a,
		DB:  conn,
	}
}
