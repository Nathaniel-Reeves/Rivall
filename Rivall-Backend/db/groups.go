package db

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

func ReadByGroupId(groupID string) Group {
	// Read a message group by its ID
	var group Group

	id, err := bson.ObjectIDFromHex(groupID)
	if err != nil {
		globals.Logger.Error().Err(err).Msg("Failed to convert group ID")
		return group
	}

	collection := globals.MongoClient.Database(Database).Collection("Groups")

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

	bsonAdminUserID, err := bson.ObjectIDFromHex(adminUserID)
	if err != nil {
		globals.Logger.Error().Err(err).Msg("Failed to convert admin user ID")
		return "", err
	}

	GroupID := bson.NewObjectID()
	group := Group{
		ID:           GroupID,
		AdminID:      bsonAdminUserID,
		GroupName:    groupName,
		GroupMembers: []bson.ObjectID{},
		LastMessage:  Message{},
		Messages:     []Message{},
		CreatedAt:    bson.Timestamp{},
	}

	collection := globals.MongoClient.Database(Database).Collection("Groups")

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

func GetGroupAdminID(groupID string) (string, error) {
	// Get the admin ID of a group
	collection := globals.MongoClient.Database(Database).Collection("Groups")

	bsonGroupID, err := bson.ObjectIDFromHex(groupID)
	if err != nil {
		globals.Logger.Error().Err(err).Msg("Failed to convert group ID")
		return "", err
	}

	var result Group
	err = collection.FindOne(context.Background(), bson.M{"_id": bsonGroupID}).Decode(&result)
	if err != nil {
		globals.Logger.Error().Err(err).Msg("Failed to get group admin ID")
		return "", err
	}

	return result.AdminID.Hex(), nil
}

func AddUserToGroup(groupID string, userID string) error {
	// Add a user to a message group

	bsonGroupID, err := bson.ObjectIDFromHex(groupID)
	if err != nil {
		globals.Logger.Error().Err(err).Msg("Failed to convert group ID")
		return err
	}

	collection := globals.MongoClient.Database(Database).Collection("Groups")

	filter := bson.D{{"_id", bsonGroupID}}
	update := bson.D{{"$push", bson.D{{"user_ids", userID}}}}

	_, err = collection.UpdateOne(context.Background(), filter, update)
	if err != nil {
		globals.Logger.Error().Err(err).Msg("Failed to add user to message group")
		return err
	}

	return nil
}

func InsertGroupMessage(GroupID string, message Message) error {
	// Insert a message into the direct messages group
	collection := globals.MongoClient.Database(Database).Collection("Groups")

	bsonGroupID, err := bson.ObjectIDFromHex(GroupID)
	if err != nil {
		globals.Logger.Error().Err(err).Msg("failed to convert direct message ID")
	}

	filter := bson.M{"_id": bsonGroupID}
	update := bson.M{"$push": bson.M{"messages": message}, "$set": bson.M{"last_message": message}}

	_, err = collection.UpdateOne(context.Background(), filter, update)
	if err != nil {
		return err
	}

	return nil
}

func GroupExists(groupID string) bool {
	// Check if a group exists in the database
	collection := globals.MongoClient.Database(Database).Collection("Groups")

	bsonGroupID, err := bson.ObjectIDFromHex(groupID)
	if err != nil {
		globals.Logger.Error().Err(err).Msg("failed to convert group ID")
		return false
	}

	var result Group
	err = collection.FindOne(context.Background(), bson.M{"_id": bsonGroupID}).Decode(&result)
	if err != nil {
		return false
	}

	return true
}

func UserInGroup(groupID string, userID string) bool {
	// Check if a user is in a group
	collection := globals.MongoClient.Database(Database).Collection("Groups")

	bsonGroupID, err := bson.ObjectIDFromHex(groupID)
	if err != nil {
		globals.Logger.Error().Err(err).Msg("failed to convert group ID")
		return false
	}

	var result Group
	err = collection.FindOne(context.Background(), bson.M{"_id": bsonGroupID}).Decode(&result)
	if err != nil {
		return false
	}

	// Check if the user ID is in the group members
	for _, member := range result.GroupMembers {
		if member.Hex() == userID {
			return true
		}
	}

	return false
}

func UserWasRequestedToJoinGroup(groupID string, userID string) bool {
	// Check if a user was requested to join a group
	collection := globals.MongoClient.Database(Database).Collection("Groups")

	bsonGroupID, err := bson.ObjectIDFromHex(groupID)
	if err != nil {
		globals.Logger.Error().Err(err).Msg("failed to convert group ID")
		return false
	}

	var result Group
	err = collection.FindOne(context.Background(), bson.M{"_id": bsonGroupID}).Decode(&result)
	if err != nil {
		return false
	}

	// Check if the user ID is in the group members
	for _, member := range result.GroupMembers {
		if member.Hex() == userID {
			return true
		}
	}

	return false
}

func GetGroupMembers(groupID string) ([]string, error) {
	// Get all group members
	collection := globals.MongoClient.Database(Database).Collection("Groups")

	bsonGroupID, err := bson.ObjectIDFromHex(groupID)
	if err != nil {
		globals.Logger.Error().Err(err).Msg("failed to convert group ID")
		return nil, err
	}

	var result Group
	err = collection.FindOne(context.Background(), bson.M{"_id": bsonGroupID}).Decode(&result)
	if err != nil {
		return nil, err
	}

	members := make([]string, len(result.GroupMembers))
	for i, member := range result.GroupMembers {
		members[i] = member.Hex()
	}

	return members, nil
}
