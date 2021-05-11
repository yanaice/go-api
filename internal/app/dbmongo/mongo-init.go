package dbmongo

import (
	"context"
	"go-starter-project/internal/app/config"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"log"
)

var mdb *mongo.Database

func Init() {
	clientOptions := options.Client()
	clientOptions.ApplyURI(config.Conf.MongoDB.URI)
	if config.Conf.MongoDB.Username != "" {
		clientOptions.SetAuth(
			options.Credential{
				Username:   config.Conf.MongoDB.Username,
				Password:   config.Conf.MongoDB.Password,
				AuthSource: config.Conf.MongoDB.Schema,
			},
		)
	}

	client, err := mongo.Connect(context.TODO(), clientOptions)

	if err != nil {
		log.Fatal(err)
	}

	err = client.Ping(context.TODO(), nil)

	if err != nil {
		log.Fatal(err)
	}

	mdb = client.Database(config.Conf.MongoDB.Schema)
}
