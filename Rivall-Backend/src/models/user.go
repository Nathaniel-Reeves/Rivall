
import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"os"
	"time"
	"github.com/gorilla/mux"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

type User struct {
	ID int64 `bson: "_id,omitempty" json: "_id,omitempty"` 
	FullText string `bson: "full_text,omitempty" json: "full_text,omitempty"`
	Contacts struct {
		UserID int64 `bson: "user_id" json: "user_id"`
	} `bson: "contacts,omitempty" json: "contacts,omitempty"`
}
