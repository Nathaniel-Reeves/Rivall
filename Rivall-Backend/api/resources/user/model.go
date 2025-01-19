package user

import (
	"context"
	"time"

	"Rivall-Backend/api/utils"

	"go.mongodb.org/mongo-driver/bson/primitive"
)

type User struct {
	ID        primitive.ObjectID `bson:"_id,omitempty" json:"_id,omitempty"`
	Username  string
	FirstName string
	LastName  string
	Email     string
	Password  string
	Image     string
	CreatedAt time.Time
	UpdatedAt time.Time
}

func Insert(user User) error {
	collection := utils.MongoClient.Database("rivall").Collection("users")
	inserted, err := collection.InsertOne(context.TODO(), user)
	if err != nil {
		utils.Logger.Error().Err(err).Msg("Failed to insert user")
	}
	utils.Logger.Info().Msgf("Inserted user with ID %v", inserted.InsertedID)
	return err
}
