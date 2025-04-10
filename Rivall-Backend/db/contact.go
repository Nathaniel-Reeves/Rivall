package db

import (
	"Rivall-Backend/globals"
	"context"

	"github.com/rs/zerolog/log"
	"go.mongodb.org/mongo-driver/v2/bson"
)

type Contact struct {
	ID              bson.ObjectID `json:"_id"          bson:"_id"`
	DirectMessageID bson.ObjectID `json:"direct_message_id" bson:"direct_message_id"`
	ContactID       bson.ObjectID `json:"contact_id"   bson:"contact_id"`
}

func CreateContact(userID string, contactID string) error {
	collection := globals.MongoClient.Database(Database).Collection("Users")

	// Setup Direct Message Line for contact
	directMessageID, err := CreateDirectMessages(userID, contactID)
	if err != nil {
		log.Error().Err(err).Msg("Failed to create direct message line")
	}

	bsonDirectMessageID, err := bson.ObjectIDFromHex(directMessageID)
	if err != nil {
		log.Error().Err(err).Msg("Failed to convert direct message ID")
	}

	bsonContactID, err := bson.ObjectIDFromHex(contactID)
	if err != nil {
		log.Error().Err(err).Msg("Failed to convert contact ID")
		return err
	}

	bsonUserID, err := bson.ObjectIDFromHex(userID)
	if err != nil {
		log.Error().Err(err).Msg("Failed to convert user ID")
	}

	contact := Contact{
		ID:              bson.NewObjectID(),
		DirectMessageID: bsonDirectMessageID,
		ContactID:       bsonContactID,
	}

	filter := bson.D{{"_id", bsonUserID}}
	update := bson.D{{"$push", bson.D{{"contact_ids", contact.ID}}}}
	_, err = collection.UpdateOne(context.TODO(), filter, update)
	if err != nil {
		log.Error().Err(err).Msg("Failed to add contact to user")
	}

	filter = bson.D{{"_id", bsonContactID}}
	update = bson.D{{"$push", bson.D{{"contact_ids", contact.ID}}}}
	_, err = collection.UpdateOne(context.TODO(), filter, update)
	if err != nil {
		log.Error().Err(err).Msg("Failed to add user to contact")
	}

	return err
}
