package health

import (
	"net/http"
)

// Read godoc
//
//	@summary		Read health
//	@description	Read health
//	@tags			health
//	@success		200
//	@router			/../livez [get]
func Read(w http.ResponseWriter, _ *http.Request) {

	// Send a ping to confirm a healthy db connection
	// var result bson.M
	// if err := global.MongoClient.Database("health").RunCommand(context.TODO(), bson.D{{"ping", 1}}).Decode(&result); err != nil {
	// 	log.Error().Err(err).Msg("Failed to ping MongoDB")
	// 	w.WriteHeader(http.StatusInternalServerError)
	// 	w.Write([]byte("Failed to ping MongoDB, API is unhealthy"))
	// 	return
	// }
	// log.Info().Msg("Pinged your deployment. MongoDB is connected.  API is Healthy.")

	w.Write([]byte("API is Healthy"))
}
