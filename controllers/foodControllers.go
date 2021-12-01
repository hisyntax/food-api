package controllers

import (
	"context"
	"net/http"
	"os"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/go-playground/validator/v10"
	"github.com/hisyntax/food-api/database"
	"github.com/hisyntax/food-api/helpers"
	"github.com/hisyntax/food-api/models"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
)

//connect to the database and open a food collection
var foodCollection *mongo.Collection = database.OpenCollection(database.Client, os.Getenv("DB_COLLECTION_NAME"))

//create a validator objact
var validate = validator.New()

//this function creates a new food item into the database
func CreateFood(c *gin.Context) {
	//open a connection to the database cluster
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	//declear a variable of type food struct
	var food models.Food

	//bind the object that comes in with the decleared variable
	//thorw an error if one occurs
	if err := c.BindJSON(&food); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	//use the validator package to verify that all items comming in meets
	//the requirements of the struct
	validatorErr := validate.Struct(food)
	if validatorErr != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": validatorErr.Error()})
		return
	}

	//assigning the time stamps upon creation
	food.Created_at, _ = time.Parse(time.RFC3339, time.Now().Format(time.RFC3339))
	food.Updated_at, _ = time.Parse(time.RFC3339, time.Now().Format(time.RFC3339))

	//generate a new ID for the object to be created
	food.ID = primitive.NewObjectID()

	//assign the auto generated ID to the primary key attribute
	food.Food_id = food.ID.Hex()
	var num = helpers.ToFixed(*food.Price, 2)
	food.Price = &num

	result, insertErr := foodCollection.InsertOne(ctx, food)
	if insertErr != nil {
		msg := "Food item was not created"
		c.JSON(http.StatusInternalServerError, gin.H{"error": msg})
		return
	}

	//return the id of the just created object to the frontend
	c.JSON(http.StatusOK, gin.H{"message": result})
}
