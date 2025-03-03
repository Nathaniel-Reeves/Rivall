package users

import (
	"Rivall-Backend/api/global"
	"context"
	"time"

	"github.com/davecgh/go-spew/spew"
	"github.com/rs/zerolog/log"
	"go.mongodb.org/mongo-driver/v2/bson"
	"golang.org/x/crypto/bcrypt"
)

/*
User struct represents a user in the system
{
	"_id": "5f8a0b9b0f1b5b1b3c1b1b1b",
	"firstName": "Nathaniel",
	"lastName": "Reeves",
	"email": "email@email.com",
	"password": "password",
	"image": "base64 encoded image",
	"groups": [<message_group>],
	"contacts": [<user>]
}
*/

const database string = "Rivall-DB"
const CollectionName string = "Users"

type User struct {
	ID           bson.ObjectID   `json:"_id"           bson:"_id"`
	FirstName    string          `json:"first_name"    bson:"first_name"`
	LastName     string          `json:"last_name"     bson:"last_name"`
	Email        string          `json:"email"         bson:"email"`
	Password     string          `json:"password" bson:"password"`
	AvatarImage  string          `json:"avatar_image"  bson:"avatar_image"`
	GroupIDs     []bson.ObjectID `bson:"groups_ids"`
	ContactIDs   []bson.ObjectID `bson:"contact_ids"`
	OTP          string          `json:"otp"`
	RefreshToken string          `json:"refresh_token" bson:"refresh_token"`
	// Contacts are not stored on the database, they are fetched from the contact_ids
	Contacts             []Contact             `json:"contacts"      bson:"omitempty"`
	MessageGroupRequests []MessageGroupRequest `json:"message_group_request_ids" bson:"message_group_request_ids"`
}

type MessageGroupRequest struct {
	ID            bson.ObjectID `json:"id"        bson:"id"`
	SendUserID    bson.ObjectID `json:"send_user_id" bson:"send_user_id"`
	RecieveUserID bson.ObjectID `json:"recieve_user_id" bson:"recieve_user_id"`
	GroupID       bson.ObjectID `json:"group_id"  bson:"group_id"`
	Message       string        `json:"message" bson:"message"`
	Timestamp     string        `json:"timestamp" bson:"timestamp"`
	Status        int8          `json:"status" bson:"status"`
}

// Status Codes
// 0 - Pending
// 1 - Accepted
// 2 - Rejected

type Contact struct {
	ID          bson.ObjectID `json:"_id"          bson:"_id"`
	FirstName   string        `json:"first_name"   bson:"first_name"`
	LastName    string        `json:"last_name"    bson:"last_name"`
	Email       string        `json:"email"        bson:"email"`
	AvatarImage string        `json:"avatar_image" bson:"avatar_image"`
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

func ReadByUserId(id string) User {
	var result User

	log.Debug().Msgf("Reading user with ID '%v'", id)
	i, _ := bson.ObjectIDFromHex(id)
	filter := bson.D{{"_id", i}}

	collection := global.MongoClient.Database(database).Collection(CollectionName)
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
	for _, requests := range result.MessageGroupRequests {
		request := ReadMessageGroupRequestById(requests.ID.Hex())
		if request.ID == bson.NilObjectID {
			log.Error().Msg("Failed to read message group request")
			continue
		}
		result.MessageGroupRequests = append(result.MessageGroupRequests, request)
	}

	return result
}

func ReadMessageGroupRequestById(id string) MessageGroupRequest {
	var result MessageGroupRequest

	log.Debug().Msgf("Reading message group request with ID '%v'", id)
	i, _ := bson.ObjectIDFromHex(id)
	filter := bson.D{{"_id", i}}

	collection := global.MongoClient.Database(database).Collection(CollectionName)
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

	collection := global.MongoClient.Database(database).Collection(CollectionName)
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

	collection := global.MongoClient.Database(database).Collection(CollectionName)
	err := collection.FindOne(context.TODO(), filter).Decode(&result)
	if err != nil {
		log.Error().Err(err).Msg("Failed to read user")
	}
	return result
}

func CreateUser(user User) error {
	collection := global.MongoClient.Database(database).Collection(CollectionName)

	if log.Debug().Enabled() {
		log.Debug().Msgf("Creating user...")
		spew.Dump(user)
	}

	// hash password
	user = HashUserPassword(user)
	user.ID = bson.NewObjectID()
	user.RefreshToken = ""

	// set default empty arrays
	user.ContactIDs = []bson.ObjectID{}
	user.GroupIDs = []bson.ObjectID{}

	// remove contacts from user object

	inserted, err := collection.InsertOne(context.TODO(), user)
	if err != nil {
		log.Error().Err(err).Msg("Failed to insert user")
	}
	log.Info().Msgf("Inserted user with ID %v", inserted.InsertedID)
	return err
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

func UpdateUserPassword(id string, password string) error {
	collection := global.MongoClient.Database(database).Collection(CollectionName)

	// hash password
	user := User{Password: password}
	user = HashUserPassword(user)

	i, _ := bson.ObjectIDFromHex(id)
	filter := bson.D{{"_id", i}}
	update := bson.D{{"$set", bson.D{{"password", user.Password}}}}

	_, err := collection.UpdateOne(context.TODO(), filter, update)
	if err != nil {
		log.Error().Err(err).Msg("Failed to update user")
	}
	return err
}

func DeleteUser(id string) error {
	collection := global.MongoClient.Database(database).Collection(CollectionName)

	i, _ := bson.ObjectIDFromHex(id)
	filter := bson.D{{"_id", i}}

	_, err := collection.DeleteOne(context.TODO(), filter)
	if err != nil {
		log.Error().Err(err).Msg("Failed to delete user")
	}
	return err
}

func CreateUserContact(userID string, contactID string) error {
	collection := global.MongoClient.Database(database).Collection(CollectionName)

	// TODO: check if contact already exists, dont add if it does
	i, _ := bson.ObjectIDFromHex(userID)
	filter := bson.D{{"_id", i}}
	update := bson.D{{"$push", bson.D{{"contact_ids", contactID}}}}

	_, err := collection.UpdateOne(context.TODO(), filter, update)
	if err != nil {
		log.Error().Err(err).Msg("Failed to add user contact")
	}
	return err
}

func CreateUserMessageGroupRequest(senderUserID string, recieverUserID string, groupID string, message string) error {
	collection := global.MongoClient.Database(database).Collection(CollectionName)

	i, _ := bson.ObjectIDFromHex(senderUserID)
	j, _ := bson.ObjectIDFromHex(groupID)
	k, _ := bson.ObjectIDFromHex(recieverUserID)
	request := MessageGroupRequest{
		SendUserID:    i,
		RecieveUserID: k,
		GroupID:       j,
		Message:       message,
		Timestamp:     time.Now().String(),
		Status:        0,
	}

	inserted, err := collection.InsertOne(context.Background(), request)
	if err != nil {
		log.Error().Err(err).Msg("Failed to create new message group request")
	}
	log.Info().Msgf("Inserted message group request with ID %v", inserted.InsertedID)
	return err
}

func AcceptUserMessageGroupRequest(userID string, groupID string) error {
	collection := global.MongoClient.Database(database).Collection(CollectionName)

	i, _ := bson.ObjectIDFromHex(groupID)
	j, _ := bson.ObjectIDFromHex(userID)

	// Update Request Status
	filter := bson.D{{"group_id", i}, {"recieve_user_id", j}}
	update := bson.D{{"$set", bson.D{{"status", 1}}}}

	_, err := collection.UpdateOne(context.TODO(), filter, update)
	if err != nil {
		log.Error().Err(err).Msg("Failed to accept message group request")
	}

	// Add user to group
	filter = bson.D{{"_id", i}}
	update = bson.D{{"$push", bson.D{{"group_ids", i}}}}

	_, err = collection.UpdateOne(context.TODO(), filter, update)
	if err != nil {
		log.Error().Err(err).Msg("Failed to add user to message group")
	}

	return err
}

func RejectUserMessageGroupRequest(userID string, groupID string) error {
	collection := global.MongoClient.Database(database).Collection(CollectionName)

	i, _ := bson.ObjectIDFromHex(groupID)
	j, _ := bson.ObjectIDFromHex(userID)

	// Update Request Status
	filter := bson.D{{"group_id", i}, {"recieve_user_id", j}}
	update := bson.D{{"$set", bson.D{{"status", 2}}}}

	_, err := collection.UpdateOne(context.TODO(), filter, update)
	if err != nil {
		log.Error().Err(err).Msg("Failed to reject message group request")
	}

	return err
}

func RemoveUserContact(userID string, contactID string) error {
	collection := global.MongoClient.Database(database).Collection(CollectionName)

	i, _ := bson.ObjectIDFromHex(userID)
	filter := bson.D{{"_id", i}}
	update := bson.D{{"$pull", bson.D{{"contact_ids", contactID}}}}

	_, err := collection.UpdateOne(context.TODO(), filter, update)
	if err != nil {
		log.Error().Err(err).Msg("Failed to delete user contact")
	}
	return err
}

func UserExists(id string) bool {
	i, _ := bson.ObjectIDFromHex(id)
	filter := bson.D{{"_id", i}}

	collection := global.MongoClient.Database(database).Collection(CollectionName)
	count, err := collection.CountDocuments(context.TODO(), filter)
	if err != nil {
		log.Error().Err(err).Msg("Failed to check if user exists")
	}
	return count > 0
}

func UpdateUserRefreshToken(user User) error {
	collection := global.MongoClient.Database(database).Collection(CollectionName)

	i, _ := bson.ObjectIDFromHex(user.ID.Hex())
	filter := bson.D{{"_id", i}}
	update := bson.D{{"$set", bson.D{{"refresh_token", user.RefreshToken}}}}

	_, err := collection.UpdateOne(context.TODO(), filter, update)
	if err != nil {
		log.Error().Err(err).Msg("Failed to update user refresh token")
	}
	return err
}
