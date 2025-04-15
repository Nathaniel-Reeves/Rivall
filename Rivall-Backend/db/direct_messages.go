package db

import (
	"Rivall-Backend/globals"
	"context"
	"time"

	"go.mongodb.org/mongo-driver/v2/bson"
)

type DirectMessages struct {
	ID                 bson.ObjectID `json:"_id" bson:"_id"`
	UserAID            bson.ObjectID `json:"user_a_id" bson:"user_a_id"`
	UserALastSeenIndex int           `json:"user_a_last_seen_index" bson:"user_a_last_seen_index"`
	UserBID            bson.ObjectID `json:"user_b_id" bson:"user_b_id"`
	UserBLastSeenIndex int           `json:"user_b_last_seen_index" bson:"user_b_last_seen_index"`
	CreatedAt          time.Time     `json:"created_at" bson:"created_at"`
	Messages           []Message     `json:"messages" bson:"messages"`
	LastMessage        Message       `json:"last_message" bson:"last_message"`
}

func CreateDirectMessages(userAID string, userBID string) (string, error) {
	// Create new direct message group in database
	// Only add the Admin user to the group
	collection := globals.MongoClient.Database(Database).Collection("DirectMessages")

	// Convert user IDs to bson.ObjectID
	bsonUserAID, err := bson.ObjectIDFromHex(userAID)
	if err != nil {
		globals.Logger.Error().Err(err).Msg("failed to convert user A ID")
		return "", err
	}
	bsonUserBID, err := bson.ObjectIDFromHex(userBID)
	if err != nil {
		globals.Logger.Error().Err(err).Msg("failed to convert user B ID")
		return "", err
	}

	directMessage := DirectMessages{
		ID:                 bson.NewObjectID(),
		UserAID:            bsonUserAID,
		UserALastSeenIndex: 0,
		UserBID:            bsonUserBID,
		UserBLastSeenIndex: 0,
		CreatedAt:          time.Now(),
		Messages:           []Message{},
		LastMessage: Message{
			ID:          bson.NewObjectID(),
			UserID:      bsonUserAID,
			MessageData: "",
			Timestamp:   time.Now().Format(time.RFC3339),
			MessageType: "text",
			SeenBy:      []bson.ObjectID{},
		},
	}

	result, err := collection.InsertOne(context.Background(), directMessage)
	if err != nil {
		return "", err
	}

	if oid, ok := result.InsertedID.(bson.ObjectID); ok {
		globals.Logger.Info().Msgf("Created direct message group with ID: %v", oid.Hex())
		return oid.Hex(), nil
	}

	return "", nil
}

func ReadDirectMessages(directMessageID string) (DirectMessages, error) {
	// Read direct messages group from database
	collection := globals.MongoClient.Database(Database).Collection("DirectMessages")

	bsonDirectMessageID, err := bson.ObjectIDFromHex(directMessageID)
	if err != nil {
		globals.Logger.Error().Err(err).Msg("failed to convert direct message ID")
		return DirectMessages{}, err
	}

	var directMessages DirectMessages
	err = collection.FindOne(context.Background(), bson.D{{"_id", bsonDirectMessageID}}).Decode(&directMessages)
	if err != nil {
		return DirectMessages{}, err
	}

	return directMessages, nil
}

func DirectMessageExists(directMessageID string) bool {
	// Check if a direct message group exists in the database
	collection := globals.MongoClient.Database(Database).Collection("DirectMessages")

	bsonDirectMessageID, err := bson.ObjectIDFromHex(directMessageID)
	if err != nil {
		globals.Logger.Error().Err(err).Msg("failed to convert direct message ID")
		return false
	}

	var result DirectMessages
	err = collection.FindOne(context.Background(), bson.M{"_id": bsonDirectMessageID}).Decode(&result)
	if err != nil {
		return false
	}

	return true
}

func UserInDirectMessage(directMessageID string, userID string) bool {
	collection := globals.MongoClient.Database(Database).Collection("DirectMessages")

	bsonDirectMessageID, err := bson.ObjectIDFromHex(directMessageID)
	if err != nil {
		globals.Logger.Error().Err(err).Msg("failed to convert direct message ID")
		return false
	}

	var result DirectMessages

	filter := bson.D{{"_id", bsonDirectMessageID}}
	err = collection.FindOne(context.TODO(), filter).Decode(&result)
	globals.Logger.Info().Msgf("Result: %v", result)
	if err != nil {
		globals.Logger.Error().Err(err).Msg("failed to find direct message group")
		return false
	}

	UserA := result.UserAID.Hex()
	globals.Logger.Info().Msgf("UserAID: %s", UserA)
	UserB := result.UserBID.Hex()
	globals.Logger.Info().Msgf("UserBID: %s", UserB)

	// Check if the user ID matches either UserAID or UserBID
	return UserA == userID || UserB == userID
}

func InsertMessage(directMessageID string, message Message) error {
	// Insert a message into the direct messages group
	collection := globals.MongoClient.Database(Database).Collection("DirectMessages")

	bsonDirectMessageID, err := bson.ObjectIDFromHex(directMessageID)
	if err != nil {
		globals.Logger.Error().Err(err).Msg("failed to convert direct message ID")
	}

	filter := bson.M{"_id": bsonDirectMessageID}
	update := bson.M{"$push": bson.M{"messages": message}, "$set": bson.M{"last_message": message}}

	_, err = collection.UpdateOne(context.Background(), filter, update)
	if err != nil {
		return err
	}

	return nil
}
