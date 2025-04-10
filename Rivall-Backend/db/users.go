package db

import (
	"Rivall-Backend/globals"
	"context"
	"errors"

	"github.com/rs/zerolog/log"
	"go.mongodb.org/mongo-driver/v2/bson"
	"golang.org/x/crypto/bcrypt"
)

type User struct {
	ID           bson.ObjectID   `json:"_id"           bson:"_id"`
	FirstName    string          `json:"first_name"    bson:"first_name"`
	LastName     string          `json:"last_name"     bson:"last_name"`
	Email        string          `json:"email"         bson:"email"`
	Password     string          `json:"password" bson:"password"`
	AvatarImage  string          `json:"avatar_image"  bson:"avatar_image"`
	GroupIDs     []bson.ObjectID `bson:"group_ids"`
	ContactIDs   []bson.ObjectID `bson:"contact_ids"`
	OTP          string          `json:"otp"`
	RefreshToken string          `json:"refresh_token" bson:"refresh_token"`
	// Contacts are not stored on the Database, they are fetched from the contact_ids
	Contacts      []Contact      `json:"contacts"      bson:"omitempty"`
	GroupRequests []GroupRequest `json:"group_requests" bson:"group_requests"`
}

type GroupRequest struct {
	ID            bson.ObjectID  `json:"_id"        bson:"_id"`
	SendUserID    bson.ObjectID  `json:"send_user_id" bson:"send_user_id"`
	RecieveUserID bson.ObjectID  `json:"receive_user_id" bson:"receive_user_id"`
	GroupID       bson.ObjectID  `json:"group_id"  bson:"group_id"`
	GroupName     string         `json:"group_name" bson:"group_name"`
	Message       string         `json:"message" bson:"message"`
	Timestamp     bson.Timestamp `json:"timestamp" bson:"timestamp"`
	Status        int8           `json:"status" bson:"status"`
}

func ReadByUserId(id string) User {
	var result User

	log.Debug().Msgf("Reading user with ID '%v'", id)
	i, _ := bson.ObjectIDFromHex(id)
	filter := bson.D{{"_id", i}}

	collection := globals.MongoClient.Database(Database).Collection("Users")
	err := collection.FindOne(context.TODO(), filter).Decode(&result)
	if err != nil {
		log.Error().Err(err).Msg("Failed to read user")
		return result
	}

	// fetch contacts
	for _, contactID := range result.ContactIDs {
		contact := ReadContactById(contactID.Hex())
		if contact.ID == bson.NilObjectID {
			log.Error().Msg("Failed to read contact")
			continue
		}
		result.Contacts = append(result.Contacts, contact)
	}

	// fetch message group requests
	for _, requests := range result.GroupRequests {
		request := ReadGroupRequestById(requests.ID.Hex())
		if request.ID == bson.NilObjectID {
			log.Error().Msg("Failed to read message group request")
			continue
		}
		result.GroupRequests = append(result.GroupRequests, request)
	}

	return result
}

func ReadGroupRequestById(id string) GroupRequest {
	var result GroupRequest

	log.Debug().Msgf("Reading message group request with ID '%v'", id)
	i, _ := bson.ObjectIDFromHex(id)
	filter := bson.D{{"_id", i}}

	collection := globals.MongoClient.Database(Database).Collection("Users")
	err := collection.FindOne(context.TODO(), filter).Decode(&result)
	if err != nil {
		log.Error().Err(err).Msg("Failed to read message group request")
	}

	return result
}

func ReadContactById(id string) Contact {
	var result Contact

	log.Debug().Msgf("Reading contact with ID '%v'", id)
	i, _ := bson.ObjectIDFromHex(id)
	filter := bson.D{{"_id", i}}

	collection := globals.MongoClient.Database(Database).Collection("Users")
	err := collection.FindOne(context.TODO(), filter).Decode(&result)
	if err != nil {
		log.Error().Err(err).Msg("Failed to read contact")
	}
	return result
}

func ReadByUserEmail(email string) User {
	var result User

	log.Debug().Msgf("Reading user with email '%v'", email)
	filter := bson.D{{"email", email}}

	collection := globals.MongoClient.Database(Database).Collection("Users")
	err := collection.FindOne(context.TODO(), filter).Decode(&result)
	if err != nil {
		log.Error().Err(err).Msg("Failed to read user")
	}
	return result
}

func HashUserPassword(user User) User {
	// Convert the plain pwd string to a byte slice
	pwd := []byte(user.Password)

	// Use GenerateFromPassword to hash & salt pwd.
	// MinCost is just an integer constant provided by the bcrypt
	// package along with DefaultCost & MaxCost.
	// The cost can be any value you want provided it isn't lower
	// than the MinCost (4)
	hash, err := bcrypt.GenerateFromPassword(pwd, 12)
	if err != nil {
		log.Error().Err(err).Msg("Failed to hash password")
	}
	// GenerateFromPassword returns a byte slice so we need to
	// convert the bytes to a string and return it
	user.Password = string(hash)

	return user
}

func ComparePasswords(hashedPwd string, plainPwd string) bool {
	byteHash := []byte(hashedPwd)
	bytePlain := []byte(plainPwd)

	err := bcrypt.CompareHashAndPassword(byteHash, bytePlain)
	if err != nil {
		log.Error().Err(err).Msg("Failed to compare passwords")
		return false
	}
	return true
}

func CreateUser(user User) error {
	collection := globals.MongoClient.Database(Database).Collection("Users")

	// hash password
	user = HashUserPassword(user)
	user.ID = bson.NewObjectID()
	user.RefreshToken = ""

	// set default empty arrays
	user.ContactIDs = []bson.ObjectID{}
	user.GroupIDs = []bson.ObjectID{}
	user.Contacts = []Contact{}
	user.GroupRequests = []GroupRequest{}

	// remove contacts from user object

	inserted, err := collection.InsertOne(context.TODO(), user)
	if err != nil {
		log.Error().Err(err).Msg("Failed to insert user")
	}
	log.Info().Msgf("Inserted user with ID %v", inserted.InsertedID)
	return err
}

func UpdateUserPassword(id string, password string) error {
	collection := globals.MongoClient.Database(Database).Collection("Users")

	// hash password
	user := User{Password: password}
	user = HashUserPassword(user)

	i, err := bson.ObjectIDFromHex(id)
	if err != nil {
		log.Error().Err(err).Msg("Failed to convert user ID")
	}
	filter := bson.D{{"_id", i}}
	update := bson.D{{"$set", bson.D{{"password", user.Password}}}}

	_, err = collection.UpdateOne(context.TODO(), filter, update)
	if err != nil {
		log.Error().Err(err).Msg("Failed to update user")
	}
	return err
}

func DeleteUser(id string) error {
	collection := globals.MongoClient.Database(Database).Collection("Users")

	i, _ := bson.ObjectIDFromHex(id)
	filter := bson.D{{"_id", i}}

	_, err := collection.DeleteOne(context.TODO(), filter)
	if err != nil {
		log.Error().Err(err).Msg("Failed to delete user")
	}
	return err
}

func enumRequestStatus(status int8) string {
	switch status {
	case 0:
		return "Pending"
	case 1:
		return "Accepted"
	case 2:
		return "Rejected"
	default:
		return "Unknown"
	}
}

func CreateGroupRequest(
	senderUserID string,
	receiverUserID string,
	groupID string,
	groupName string,
	message string,
) (string, error) {
	collection := globals.MongoClient.Database(Database).Collection("Users")

	i, err := bson.ObjectIDFromHex(senderUserID)
	if err != nil {
		log.Error().Err(err).Msg("Failed to convert sender user ID")
		return "", err
	}

	j, err := bson.ObjectIDFromHex(groupID)
	if err != nil {
		log.Error().Err(err).Msg("Failed to convert group ID")
		return "", err
	}

	k, err := bson.ObjectIDFromHex(receiverUserID)
	if err != nil {
		log.Error().Err(err).Msg("Failed to convert receiver user ID")
		return "", err
	}

	RequestID := bson.NewObjectID()
	request := GroupRequest{
		ID:            RequestID,
		SendUserID:    i,
		RecieveUserID: k,
		GroupID:       j,
		GroupName:     groupName,
		Message:       message,
		Timestamp:     bson.Timestamp{},
		Status:        0,
	}

	updateResult, err := collection.UpdateOne(
		context.Background(),
		bson.D{{Key: "_id", Value: k}},
		bson.D{{Key: "$push", Value: bson.D{{Key: "group_requests", Value: request}}}},
	)
	if err != nil {
		log.Error().Err(err).Msg("Failed to create new message group request")
		return "", err
	}
	if updateResult.MatchedCount == 0 || updateResult.ModifiedCount == 0 {
		log.Error().Msgf("Matched %d documents and modified %d documents", updateResult.MatchedCount, updateResult.ModifiedCount)
		return "", errors.New("Failed to create new message group request")
	}

	log.Info().Msgf("Inserted message group request for user: %v", receiverUserID)
	return RequestID.Hex(), err
}

func AcceptGroupRequest(userID string, groupID string) error {
	collection := globals.MongoClient.Database(Database).Collection("Users")

	i, _ := bson.ObjectIDFromHex(groupID)
	j, _ := bson.ObjectIDFromHex(userID)

	// Update Request Status
	filter := bson.D{{"group_id", i}, {"receive_user_id", j}}
	update := bson.D{{"$set", bson.D{{"status", 1}}}}

	_, err := collection.UpdateOne(context.TODO(), filter, update)
	if err != nil {
		log.Error().Err(err).Msg("Failed to accept message group request")
	}

	// Add user to group
	err = AddUserToGroup(groupID, userID)
	if err != nil {
		log.Error().Err(err).Msg("Failed to add user to group")
	}

	return err
}

func RejectGroupRequest(userID string, groupID string) error {
	collection := globals.MongoClient.Database(Database).Collection("Users")

	i, _ := bson.ObjectIDFromHex(groupID)
	j, _ := bson.ObjectIDFromHex(userID)

	// Update Request Status
	filter := bson.D{{"group_id", i}, {"receive_user_id", j}}
	update := bson.D{{"$set", bson.D{{"status", 2}}}}

	_, err := collection.UpdateOne(context.TODO(), filter, update)
	if err != nil {
		log.Error().Err(err).Msg("Failed to reject message group request")
	}

	return err
}

func UserExists(id string) bool {
	i, _ := bson.ObjectIDFromHex(id)
	filter := bson.D{{"_id", i}}

	collection := globals.MongoClient.Database(Database).Collection("Users")
	count, err := collection.CountDocuments(context.TODO(), filter)
	if err != nil {
		log.Error().Err(err).Msg("Failed to check if user exists")
	}
	return count > 0
}

func UpdateUserRefreshToken(user User) error {
	collection := globals.MongoClient.Database(Database).Collection("Users")

	i, _ := bson.ObjectIDFromHex(user.ID.Hex())
	filter := bson.D{{"_id", i}}
	update := bson.D{{"$set", bson.D{{"refresh_token", user.RefreshToken}}}}

	_, err := collection.UpdateOne(context.TODO(), filter, update)
	if err != nil {
		log.Error().Err(err).Msg("Failed to update user refresh token")
	}
	return err
}
