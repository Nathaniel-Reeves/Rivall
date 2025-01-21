package user

import (
	"Rivall-Backend/api/global"
	"context"

	"github.com/davecgh/go-spew/spew"
	"github.com/rs/zerolog/log"
	"go.mongodb.org/mongo-driver/v2/bson"
)

/*
User struct represents a user in the system
{
	"_id": "5f8a0b9b0f1b5b1b3c1b1b1b",
	"username": "nreeves",
	"firstName": "Nathaniel",
	"lastName": "Reeves",
	"email": "
	"password": "password",
	"image": "https://www.google.com",
	"createdAt": "2020-10-17T00:00:00Z",
	"updatedAt": "2020-10-17T00:00:00Z"
}
*/

const database string = "Rivall-DB"
const collectionName string = "Users"

type User struct {
	ID        bson.ObjectID `json:"_id"           bson:"_id"`
	Username  string        `json:"username"      bson:"username"`
	FirstName string        `json:"first_name"    bson:"first_name"`
	LastName  string        `json:"last_name"     bson:"last_name"`
	Email     string        `json:"email"         bson:"email"`
	Password  string        `json:"password"      bson:"password"`
	Image     string        `json:"image"         bson:"image"`
}

func ReadId(id string) User {
	var result User

	log.Debug().Msgf("Reading user with ID '%v'", id)
	i, _ := bson.ObjectIDFromHex(id)
	filter := bson.D{{"_id", i}}

	collection := global.MongoClient.Database(database).Collection(collectionName)
	err := collection.FindOne(context.TODO(), filter).Decode(&result)
	if err != nil {
		log.Error().Err(err).Msg("Failed to read user")

	}
	return result
}

func ReadUsername(username string) User {
	var result User

	log.Debug().Msgf("Reading user with username '%v'", username)
	filter := bson.D{}

	collection := global.MongoClient.Database(database).Collection(collectionName)
	err := collection.FindOne(context.TODO(), filter).Decode(&result)
	if err != nil {
		log.Error().Err(err).Msg("Failed to read user")
	}
	return result
}

func CreateUser(user User) error {
	collection := global.MongoClient.Database(database).Collection(collectionName)

	if log.Debug().Enabled() {
		log.Debug().Msgf("Creating user...")
		spew.Dump(user)
	}

	inserted, err := collection.InsertOne(context.TODO(), user)
	if err != nil {
		log.Error().Err(err).Msg("Failed to insert user")
	}
	log.Info().Msgf("Inserted user with ID %v", inserted.InsertedID)
	return err
}
