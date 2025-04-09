package models

import (
	"Rivall-Backend/globals"
	"context"

	"go.mongodb.org/mongo-driver/v2/bson"
)

const collectionName string = "Groups"

type Group struct {
	ID           bson.ObjectID   `json:"_id"           bson:"_id"`
	AdminID      bson.ObjectID   `json:"admin_id"      bson:"admin_id"`
	GroupMembers []bson.ObjectID `json:"users"         bson:"users"`
	LastMessage  Message         `json:"last_message"  bson:"last_message"`
	Messages     []Message       `json:"messages"      bson:"messages"`
	GroupName    string          `json:"group_name"    bson:"group_name"`
	CreatedAt    bson.Timestamp  `json:"created_at"    bson:"created_at"`
}

type Message struct {
	ID          bson.ObjectID   `json:"_id"           bson:"_id"`
	User        bson.ObjectID   `json:"user"          bson:"user"`
	MessageData string          `json:"message_data"  bson:"message_data"`
	Timestamp   string          `json:"timestamp"     bson:"timestamp"`
	MessageType string          `json:"message_type"  bson:"message_type"`
	SeenBy      []bson.ObjectID `json:"seen_by"     bson:"seen_by"`
}

func ReadByGroupId(groupID string) Group {
	// Read a message group by its ID
	var group Group

	id, err := bson.ObjectIDFromHex(groupID)
	if err != nil {
		globals.Logger.Error().Err(err).Msg("Failed to convert group ID")
		return group
	}

	collection := globals.MongoClient.Database(database).Collection(collectionName)

	filter := bson.D{{"_id", id}}
	err = collection.FindOne(context.Background(), filter).Decode(&group)
	if err != nil {
		globals.Logger.Error().Err(err).Msg("Failed to read message group")
		return group
	}

	return group
}

func CreateGroup(groupName string, adminUserID string) (string, error) {
	// Create new message group in database
	// Only add the Admin user to the group

	convertedAdminUserID, err := bson.ObjectIDFromHex(adminUserID)
	if err != nil {
		globals.Logger.Error().Err(err).Msg("Failed to convert admin user ID")
		return "", err
	}

	GroupID := bson.NewObjectID()
	group := Group{
		ID:           GroupID,
		AdminID:      convertedAdminUserID,
		GroupName:    groupName,
		GroupMembers: []bson.ObjectID{},
		LastMessage:  Message{},
		Messages:     []Message{},
		CreatedAt:    bson.Timestamp{},
	}

	collection := globals.MongoClient.Database(database).Collection(collectionName)

	result, err := collection.InsertOne(context.Background(), group)
	if err != nil {
		globals.Logger.Error().Err(err).Msg("Failed to create new message group")
		return "", err
	}

	// Ensure result.InsertedID is a string
	if oid, ok := result.InsertedID.(bson.ObjectID); ok {
		globals.Logger.Info().Msgf("Created group with ID: %v", oid.Hex())
		return oid.Hex(), nil
	}
	globals.Logger.Error().Msg("Failed to convert inserted ID to string")
	return "", err
}

func AddUserToGroup(groupID string, userID string) error {
	// Add a user to a message group

	id, err := bson.ObjectIDFromHex(groupID)
	if err != nil {
		globals.Logger.Error().Err(err).Msg("Failed to convert group ID")
		return err
	}

	collection := globals.MongoClient.Database(database).Collection(collectionName)

	filter := bson.D{{"_id", id}}
	update := bson.D{{"$push", bson.D{{"user_ids", userID}}}}

	_, err = collection.UpdateOne(context.Background(), filter, update)
	if err != nil {
		globals.Logger.Error().Err(err).Msg("Failed to add user to message group")
		return err
	}

	return nil
}
