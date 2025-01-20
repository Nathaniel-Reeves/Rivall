package user

import (
	"Rivall-Backend/api/global"
	"context"
	"time"

	"github.com/rs/zerolog/log"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
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
const colName string = "Users"

type User struct {
	ID        primitive.ObjectID `bson:"_id,omitempty" json:"_id,omitempty"`
	Username  string             `bson:"username,omitempty" json:"username,omitempty"`
	FirstName string             `bson:"first_name,omitempty" json:"first_name,omitempty"`
	LastName  string             `bson:"last_name,omitempty" json:"last_name,omitempty"`
	Email     string             `bson:"email,omitempty" json:"email,omitempty"`
	Password  string             `bson:"password,omitempty" json:"password,omitempty"`
	Image     string             `bson:"image,omitempty" json:"image,omitempty"`
	CreatedAt time.Time          `bson:"createdat,omitempty" json:"createdat,omitempty"`
	UpdatedAt time.Time          `bson:"updatedat,omitempty" json:"updatedat,omitempty"`
}

func ReadId(id string) (User, error) {
	var result User

	filter := bson.D{{"_id", id}}

	collection := global.MongoClient.Database(database).Collection(colName)
	err := collection.FindOne(context.TODO(), filter).Decode(&result)
	if err != nil {
		log.Error().Err(err).Msg("Failed to read user")
	}
	return result, err
}

func ReadUsername(username string) (User, error) {
	var result User

	filter := bson.D{{"username", username}}

	collection := global.MongoClient.Database(database).Collection(colName)
	err := collection.FindOne(context.TODO(), filter).Decode(&result)
	if err != nil {
		log.Error().Err(err).Msg("Failed to read user")
	}
	return result, err
}

func CreateUser(user User) error {
	collection := global.MongoClient.Database(database).Collection(colName)
	inserted, err := collection.InsertOne(context.TODO(), user)
	if err != nil {
		log.Error().Err(err).Msg("Failed to insert user")
	}
	log.Info().Msgf("Inserted user with ID %v", inserted.InsertedID)
	return err
}
