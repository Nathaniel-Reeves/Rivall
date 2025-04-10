package db

import "go.mongodb.org/mongo-driver/v2/bson"

type Message struct {
	ID          bson.ObjectID   `json:"_id"           bson:"_id"`
	SenderID    bson.ObjectID   `json:"user"          bson:"user"`
	MessageData string          `json:"message_data"  bson:"message_data"`
	Timestamp   string          `json:"timestamp"     bson:"timestamp"`
	MessageType string          `json:"message_type"  bson:"message_type"`
	SeenBy      []bson.ObjectID `json:"seen_by"     bson:"seen_by"`
}
