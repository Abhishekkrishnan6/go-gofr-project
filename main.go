
package main

import (
	"context"
	"encoding/json"
	"log"
	
	"time"

	"gofr.dev/pkg/gofr"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"github.com/gofrs/uuid"
)

type Car struct {
	ID           string    `json:"id" bson:"_id"`
	LicensePlate string    `json:"licensePlate" bson:"licensePlate"`
	EntryTime    time.Time `json:"entryTime" bson:"entryTime"`
	Status       string    `json:"status" bson:"status"`
}

type User struct {
	ID       string `json:"id" bson:"_id"`
	Username string `json:"username" bson:"username"`
	Password string `json:"password" bson:"password"`
}

var userCollection *mongo.Collection
var collection *mongo.Collection

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
	userCollection = client.Database("car_garage").Collection("users")
}

func main() {
	app := gofr.New()
	app.POST("/addEntry", addEntryHandler)
	app.GET("/listCars", listCarsHandler)
	app.PUT("/updateEntry", updateEntryHandler)   
	app.POST("/deleteEntry", deleteEntryHandler) 
	app.POST("/register", registerHandler)
	app.POST("/login", loginHandler)
	port := 8000
	log.Printf("Server is running on :%d\n", port)
	app.Start()
}

func registerHandler(ctx *gofr.Context) (interface{}, error) {
	var user User
	err := json.NewDecoder(ctx.Request().Body).Decode(&user)
	if err != nil {
		return nil, err
	}

	existingUser := User{}
	err = userCollection.FindOne(context.Background(), bson.D{{"username", user.Username}}).Decode(&existingUser)
	if err == nil {
		return nil, err
	}

	userID, err := uuid.NewV4()
	if err != nil {
		return nil, err
	}
	user.ID = userID.String()

	user.Password = "hashed_password"

	_, err = userCollection.InsertOne(context.Background(), user)
	if err != nil {
		return nil, err
	}

	return map[string]string{"id": user.ID}, nil
}

func loginHandler(ctx *gofr.Context) (interface{}, error) {
	var loginData map[string]string
	err := json.NewDecoder(ctx.Request().Body).Decode(&loginData)
	if err != nil {
		return nil, err
	}

	username := loginData["username"]
	password := loginData["password"]

	user := User{}
	err = userCollection.FindOne(context.Background(), bson.D{{"username", username}}).Decode(&user)
	if err != nil {
		return nil, err
	}

	if user.Password != password {
		return nil, err
	}

	return map[string]string{"message": "Login successful"}, nil
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
