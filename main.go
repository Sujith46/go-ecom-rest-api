package main

import (
	"context"
	// "encoding/json"
	"log"
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/sujith46/ecom-rest-api/types"
	"github.com/sujith46/ecom-rest-api/utils"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

// Define a data model

var client *mongo.Client

// Create a new person
func createPerson(c *gin.Context) {
	var person types.Person

	if err := c.ShouldBindJSON(&person); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	person.CreatedAt = time.Now()

	collection := client.Database("test").Collection("people")
	result, err := collection.InsertOne(context.Background(), person)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	person.ID = result.InsertedID.(string)
	c.JSON(http.StatusCreated, person)
}

// Get all people
func getPeople(c *gin.Context) {
	var people []types.Person
	collection := client.Database("ecom").Collection("people")
	cursor, err := collection.Find(context.Background(), bson.M{})
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	defer cursor.Close(context.Background())

	for cursor.Next(context.Background()) {
		var person types.Person
		cursor.Decode(&person)
		people = append(people, person)
	}

	c.JSON(http.StatusOK, people)
}

func retirePerson(c *gin.Context) {
	personID := c.Param("id")

	collection := client.Database("test").Collection("people")

	filter := bson.M{"_id": personID}
	update := bson.M{"$set": bson.M{"retired": true}}

	result, err := collection.UpdateOne(context.Background(), filter, update)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	if result.ModifiedCount == 0 {
		c.JSON(http.StatusNotFound, gin.H{"error": "Person not found"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Person retired successfully"})
}

func main() {
	// Connect to MongoDB
	mongoURL := utils.ReadEnv("MONGODB_URI")
	clientOptions := options.Client().ApplyURI(mongoURL)
	var err error
	client, err = mongo.Connect(context.Background(), clientOptions)
	if err != nil {
		log.Fatal(err)
	}
	defer client.Disconnect(context.Background())

	// Initialize the Gin router
	router := gin.Default()

	// Define API routes
	router.POST("/people", createPerson)
	router.GET("/people", getPeople)
	router.GET("/people/:id/retire", retirePerson)

	// Start the server
	log.Println("Server is running on :8080...")
	log.Fatal(router.Run(":8080"))
}
