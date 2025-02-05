package users

import (
	"Rivall-Backend/api/global"
	"context"

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
const collectionName string = "Users"

type User struct {
	ID          bson.ObjectID   `json:"_id"           bson:"_id"`
	FirstName   string          `json:"first_name"    bson:"first_name"`
	LastName    string          `json:"last_name"     bson:"last_name"`
	Email       string          `json:"email"         bson:"email"`
	Password    string          `json:"password"      bson:"password"`
	AvatarImage string          `json:"avatar_image"  bson:"avatar_image"`
	Groups      []bson.ObjectID `json:"groups"        bson:"groups"`
	Contacts    []bson.ObjectID `json:"contacts"      bson:"contacts"`
}

func ReadId(id string) User {
	var result User

	log.Debug().Msgf("Reading user with ID '%v'", id)
	i, _ := bson.ObjectIDFromHex(id)
	filter := bson.D{{"_id", i}}

	collection := global.MongoClient.Database(database).Collection(collectionName)
	err := collection.FindOne(context.TODO(), filter).Decode(&result)
	if err != nil {
		log.Error().Err(err).Msg("Failed to read user")

	}
	return result
}

func ReadEmail(email string) User {
	var result User

	log.Debug().Msgf("Reading user with email '%v'", email)
	filter := bson.D{{"email", email}}

	collection := global.MongoClient.Database(database).Collection(collectionName)
	err := collection.FindOne(context.TODO(), filter).Decode(&result)
	if err != nil {
		log.Error().Err(err).Msg("Failed to read user")
	}
	return result
}

func ReadIdPopulateContacts(id string) User {
	var result User

	log.Debug().Msgf("Reading user with ID '%v'", id)
	i, _ := bson.ObjectIDFromHex(id)
	filter := bson.D{{"_id", i}}

	collection := global.MongoClient.Database(database).Collection(collectionName)
	err := collection.FindOne(context.TODO(), filter).Decode(&result)
	if err != nil {
		log.Error().Err(err).Msg("Failed to read user")
	}

	return result
}

func CreateUser(user User) error {
	collection := global.MongoClient.Database(database).Collection(collectionName)

	if log.Debug().Enabled() {
		log.Debug().Msgf("Creating user...")
		spew.Dump(user)
	}

	// hash password
	user = HashUserPassword(user)
	user.ID = bson.NewObjectID()

	// TODO: Need to make the default values for groups and contacts
	// and empty array, currently the are set to null at new user registration
	//
	// user.Contacts = bson.TypeArray
	// user.Groups = bson.TypeArray

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
	collection := global.MongoClient.Database(database).Collection(collectionName)

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
	collection := global.MongoClient.Database(database).Collection(collectionName)

	i, _ := bson.ObjectIDFromHex(id)
	filter := bson.D{{"_id", i}}

	_, err := collection.DeleteOne(context.TODO(), filter)
	if err != nil {
		log.Error().Err(err).Msg("Failed to delete user")
	}
	return err
}

func CreateUserContact(userID string, contactID string) error {
	collection := global.MongoClient.Database(database).Collection(collectionName)

	i, _ := bson.ObjectIDFromHex(userID)
	filter := bson.D{{"_id", i}}
	update := bson.D{{"$push", bson.D{{"contacts", contactID}}}}

	_, err := collection.UpdateOne(context.TODO(), filter, update)
	if err != nil {
		log.Error().Err(err).Msg("Failed to add user contact")
	}
	return err
}
