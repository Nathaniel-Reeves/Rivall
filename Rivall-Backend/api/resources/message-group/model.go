package messagegroup

import (
	"context"

	"github.com/rs/zerolog/log"
	"go.mongodb.org/mongo-driver/v2/bson"

	"Rivall-Backend/api/global"
	"Rivall-Backend/api/resources/users"
)

/*
Message Group struct represents a message group in the system
{
	"_id": "5f8a0b9b0f1b5b1b3c1b1b1b",
	"users": [<user>],
	"last_message": <message>,
	"messages": [<message>]
}
*/

const database string = "Rivall-DB"
const collectionName string = "MessageGroups"

type MessageGroup struct {
	ID          bson.ObjectID `json:"_id"           bson:"_id"`
	Users       []users.User  `json:"users"         bson:"users"`
	LastMessage Message       `json:"last_message"  bson:"last_message"`
	Messages    []Message     `json:"messages"      bson:"messages"`
	GroupName   string        `json:"group_name"    bson:"group_name"`
}

/*
Message struct represents a message in the system
{
	"_id": "5f8a0b9b0f1b5b1b3c1b1b1b",
	"user": <user>,
	"message_data": "Hello, world!",
	"timestamp": "2021-01-01T00:00:00Z"
	"message_type": Enum("TEXT", "IMAGE")
}
*/

type Message struct {
	ID          bson.ObjectID `json:"_id"           bson:"_id"`
	User        users.User    `json:"user"          bson:"user"`
	MessageData string        `json:"message_data"  bson:"message_data"`
	Timestamp   string        `json:"timestamp"     bson:"timestamp"`
	MessageType string        `json:"message_type"  bson:"message_type"`
}

func CreateMessageGroup(user users.User, messageGroup MessageGroup) error {
	collection := global.MongoClient.Database(database).Collection(collectionName)

	messageGroup.Users = append(messageGroup.Users, user)

	_, err := collection.InsertOne(context.TODO(), messageGroup)
	if err != nil {
		log.Error().Err(err).Msg("Failed to create message group")
		return err
	}

	log.Info().Msgf("Created message group with ID '%v'", messageGroup.ID)
	return nil
}

func AddMessageGroupToUsers(user users.User, messageGroup MessageGroup) error {
	collection := global.MongoClient.Database(database).Collection(collectionName)

	_, err := collection.UpdateMany(
		context.TODO(),
		bson.M{"_id": bson.M{"$in": messageGroup.Users}},
		bson.M{"$push": bson.M{"groups": messageGroup.ID}},
	)
	if err != nil {
		log.Error().Err(err).Msg("Failed to add message group to users")
		return err
	}

	log.Info().Msgf("Added message group with ID '%v' to users", messageGroup.ID)
	return nil
}
