package main
import (
	"context"
	"encoding/json"
	"github.com/gofrs/uuid"
	"gofr.dev/pkg/gofr"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"log"
	"time"
)
type Car struct {
	ID           string    `json:"id" bson:"_id"`
	LicensePlate string    `json:"licensePlate" bson:"licensePlate"`
	EntryTime    time.Time `json:"entryTime" bson:"entryTime"`
	Status       string    `json:"status" bson:"status"`
}
var (
	client *mongo.Client
	collection *mongo.Collection
)

func init() {
	clientOptions := options.Client().ApplyURI("mongodb://localhost:27017")
	client, err := mongo.Connect(context.Background(), clientOptions)
	if err != nil {
		log.Fatal(err)
	}

	err = client.Ping(context.Background(), nil)
	if err != nil {
		log.Fatal(err)
	}

	collection = client.Database("car_garage").Collection("cars")
}

func main() {
	app := gofr.New()
	app.POST("/addEntry", addEntryHandler)
	app.GET("/listCars", listCarsHandler)
	app.POST("/updateEntry", updateEntryHandler)
	app.POST("/deleteEntry", deleteEntryHandler)
	port := 8080
	log.Printf("Server is running on :%d\n", port)
	app.Start()
	
}

func addEntryHandler(ctx *gofr.Context) (interface{}, error) {

	id, err := uuid.NewV4()
	if err != nil {
		
		return nil, err
	}

	var car Car
	err = json.NewDecoder(ctx.Request().Body).Decode(&car)
	if err != nil {
		
		return nil, err
	}

	car.EntryTime = time.Now()
	car.ID = id.String()
	car.Status = "In Garage"

	_, err = collection.InsertOne(context.Background(), car)
	if err != nil {
		
		return nil, err
	}

	return map[string]string{"id": car.ID}, nil
}

func listCarsHandler(ctx *gofr.Context) (interface{}, error) {

	var cars []Car
	cursor, err := collection.Find(context.Background(), bson.D{})
	if err != nil {
		
		return nil, err
	}
	defer cursor.Close(context.Background())

	err = cursor.All(context.Background(), &cars)
	if err != nil {
		
		return nil, err
	}

	return cars, nil
}

func updateEntryHandler(ctx *gofr.Context) (interface{}, error) {

	var updateData map[string]interface{}
	err := json.NewDecoder(ctx.Request().Body).Decode(&updateData)
	if err != nil {
		
		return nil, err
	}

	carID, ok := updateData["id"].(string)
	if !ok {
		
		return nil, err
	}

	filter := bson.D{{"_id", carID}}
	update := bson.D{{"$set", updateData}}
	_, err = collection.UpdateOne(context.Background(), filter, update)
	if err != nil {
		
		return nil, err
	}

	return map[string]string{"message": "Car entry updated successfully"}, nil
}

func deleteEntryHandler(ctx *gofr.Context) (interface{}, error) {

	var deleteData map[string]interface{}
	err := json.NewDecoder(ctx.Request().Body).Decode(&deleteData)
	if err != nil {
		
		return nil, err
	}

	carID, ok := deleteData["id"].(string)
	if !ok {
		
		return nil, err
	}

	filter := bson.D{{"_id", carID}}
	_, err = collection.DeleteOne(context.Background(), filter)
	if err != nil {
		
		return nil, err
	}
	return map[string]string{"message": "Car entry deleted successfully"}, nil
}
