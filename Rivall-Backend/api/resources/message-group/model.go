package messagegroup

import (
	"context"

	"go.mongodb.org/mongo-driver/v2/bson"

	"Rivall-Backend/api/global"
	"Rivall-Backend/api/resources/users"
)

const database string = "Rivall-DB"
const collectionName string = "MessageGroups"

type MessageGroup struct {
	ID             bson.ObjectID   `json:"_id"           bson:"_id"`
	GroupMemberIDs []bson.ObjectID `json:"user_ids"      bson:"user_ids"`
	GroupMembers   []users.User    `json:"users"         bson:"users"`
	LastMessage    Message         `json:"last_message"  bson:"last_message"`
	Messages       []Message       `json:"messages"      bson:"messages"`
	GroupName      string          `json:"group_name"    bson:"group_name"`
}

type Message struct {
	ID          bson.ObjectID   `json:"_id"           bson:"_id"`
	User        users.User      `json:"user"          bson:"user"`
	MessageData string          `json:"message_data"  bson:"message_data"`
	Timestamp   string          `json:"timestamp"     bson:"timestamp"`
	MessageType string          `json:"message_type"  bson:"message_type"`
	SeenBy      []bson.ObjectID `json:"seen_by"     bson:"seen_by"`
}

func ReadByGroupId(groupID string) MessageGroup {
	// Read a message group by its ID
	var messageGroup MessageGroup

	id, err := bson.ObjectIDFromHex(groupID)
	if err != nil {
		global.Logger.Error().Err(err).Msg("Failed to convert group ID")
		return messageGroup
	}

	collection := global.MongoClient.Database(database).Collection(collectionName)

	filter := bson.D{{"_id", id}}
	err = collection.FindOne(context.Background(), filter).Decode(&messageGroup)
	if err != nil {
		global.Logger.Error().Err(err).Msg("Failed to read message group")
		return messageGroup
	}

	return messageGroup
}

func CreateMessageGroup(groupName string, adminUserID string) (string, error) {
	// Create new message group in database
	// Only add the Admin user to the group

	messageGroup := MessageGroup{
		GroupName:      groupName,
		GroupMemberIDs: []bson.ObjectID{},
		GroupMembers:   []users.User{},
		LastMessage:    Message{},
		Messages:       []Message{},
	}

	// Add the admin user to the group
	adminId, err := bson.ObjectIDFromHex(adminUserID)
	if err != nil {
		global.Logger.Error().Err(err).Msg("Failed to convert admin user ID")
		return "", err
	}

	messageGroup.GroupMemberIDs = append(messageGroup.GroupMemberIDs, adminId)

	collection := global.MongoClient.Database(database).Collection(collectionName)

	result, err := collection.InsertOne(context.Background(), messageGroup)
	if err != nil {
		global.Logger.Error().Err(err).Msg("Failed to create new message group")
		return "", err
	}

	// Add the admin user to the group
	usersCollection := global.MongoClient.Database(database).Collection(users.CollectionName)

	filter := bson.D{{"_id", adminId}}
	update := bson.D{{"$push", bson.D{{"message_group_ids", result.InsertedID}}}}

	_, err = usersCollection.UpdateOne(context.TODO(), filter, update)
	if err != nil {
		global.Logger.Error().Err(err).Msg("Failed to add admin user to message group")
		return "", err
	}

	// Ensure result.InsertedID is a string
	global.Logger.Info().Msgf("Inserted message group with ID %v", result.InsertedID)
	return result.InsertedID.(string), nil
}

func AddUserToMessageGroup(groupID string, userID string) error {
	// Add a user to a message group

	id, err := bson.ObjectIDFromHex(groupID)
	if err != nil {
		global.Logger.Error().Err(err).Msg("Failed to convert group ID")
		return err
	}

	collection := global.MongoClient.Database(database).Collection(collectionName)

	filter := bson.D{{"_id", id}}
	update := bson.D{{"$push", bson.D{{"user_ids", userID}}}}

	_, err = collection.UpdateOne(context.Background(), filter, update)
	if err != nil {
		global.Logger.Error().Err(err).Msg("Failed to add user to message group")
		return err
	}

	return nil
}
