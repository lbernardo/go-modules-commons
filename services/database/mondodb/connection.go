package mondodb

import (
	"context"
	"github.com/spf13/viper"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"go.uber.org/zap"
)

func NewConnection(config *viper.Viper, logger *zap.Logger) *mongo.Client {
	client, err := mongo.Connect(context.Background(), options.Client().ApplyURI(config.GetString("app.database.mongodb.uri")))
	if err != nil {
		logger.Fatal("failed to connect to database", zap.Error(err))
		return nil
	}
	return client
}

func NewDatabase(client *mongo.Client, config *viper.Viper) *mongo.Database {
	return client.Database(config.GetString("app.database.mongodb.name"))
}
