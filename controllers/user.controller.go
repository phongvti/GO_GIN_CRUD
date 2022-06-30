package controllers

import (
	"context"
	"net/http"
	"start-gin/configs"
	"start-gin/models"
	"start-gin/responses"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/go-playground/validator/v10"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	// "go.mongodb.org/mongo-driver/mongo/options"
)


var userCollection *mongo.Collection= configs.GetCollection(configs.DB, "users")
var validate= validator.New()


func CreateUser() gin.HandlerFunc{
	return func(c *gin.Context){
		ctx, cancel:= context.WithTimeout(context.Background(), 10*time.Second)
		var user models.User
		defer cancel()

		//validate the requets body
		if err:= c.BindJSON(&user); err!=nil{
			c.JSON(http.StatusBadRequest, responses.UserResponse{Status: http.StatusBadRequest, Message: "Error", Data: map[string]interface{}{"data": err.Error()}})
			return 
		}

		if validationErr:= validate.Struct(&user); validationErr!=nil{
			c.JSON(http.StatusBadRequest, responses.UserResponse{Status: http.StatusBadRequest, Message: "Error",Data: map[string]interface{}{"data": validationErr.Error()}})
			return 
		}
		
		newUser := models.User{
			Id: primitive.NewObjectID(),
			Name: user.Name,
			Location: user.Location,
			Title: user.Title,
		}

		result, err:= userCollection.InsertOne(ctx, newUser)
		if err!=nil{
			c.JSON(http.StatusInternalServerError, responses.UserResponse{Status: http.StatusInternalServerError, Message: "Error", Data: map[string]interface{}{"data": err.Error()}})
			return
		}

		c.JSON(http.StatusCreated, responses.UserResponse{Status: http.StatusCreated, Message: "Successfull", Data: map[string]interface{}{"data": result}})
		
	}
}

func GetAnUser() gin.HandlerFunc{
	return func (c *gin.Context){
		ctx, cancel:= context.WithTimeout(context.Background(), 10*time.Second)
		userId:= c.Param("userId")

		var user models.User
		defer cancel()

		objId, _ := primitive.ObjectIDFromHex(userId)

		err:= userCollection.FindOne(ctx, bson.M{"id": objId}).Decode(&user)

		if err!=nil{
			c.JSON(http.StatusInternalServerError, responses.UserResponse{Status: http.StatusInternalServerError, Message: "error", Data: map[string]interface{}{"data": err.Error()}})
			return
		}
		c.JSON(http.StatusOK, responses.UserResponse{Status: http.StatusOK, Message: "Successful", Data: map[string]interface{}{"data": user}})
	}
}

func GetAllUser() gin.HandlerFunc{
	return func(c *gin.Context) {
		ctx, cancel:= context.WithTimeout(context.Background(), 10*time.Second)
		var users []models.User
		defer cancel()

		results, err:= userCollection.Find(ctx, bson.M{})

		if err!=nil{
			c.JSON(http.StatusInternalServerError, responses.UserResponse{Status: http.StatusInternalServerError, Message: "error", Data: map[string]interface{}{"data": err.Error()}})
			return
		}

		//reading from the db in an optimal way
		defer results.Close(ctx)
		for results.Next(ctx){
			var singleUser models.User
			if err= results.Decode(&singleUser); err!=nil{
				c.JSON(http.StatusInternalServerError, responses.UserResponse{Status: http.StatusInternalServerError, Message: "Error", Data: map[string]interface{}{"data": err.Error()}})
			}

			users= append(users, singleUser)
		}

		c.JSON(http.StatusOK, responses.UserResponse{Status: http.StatusOK, Message: "Successful", Data: map[string]interface{}{"data": users}})
	}
}

func UpdateAnUser() gin.HandlerFunc{
	return func (c *gin.Context){
		ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
        userId := c.Param("userId")
        var user models.User
        defer cancel()
        objId, _ := primitive.ObjectIDFromHex(userId)

        //validate the request body
        if err := c.BindJSON(&user); err != nil {
            c.JSON(http.StatusBadRequest, responses.UserResponse{Status: http.StatusBadRequest, Message: "error", Data: map[string]interface{}{"data": err.Error()}})
            return
        }

		
        update := bson.M{}

		if user.Location!=""{
			update["location"]= user.Location
		}
		if user.Name!=""{
			update["name"]= user.Name
		}
		if user.Title!=""{
			update["title"]= user.Title
		}

        result, err := userCollection.UpdateOne(ctx, bson.M{"id": objId}, bson.M{"$set": update})
        if err != nil {
            c.JSON(http.StatusInternalServerError, responses.UserResponse{Status: http.StatusInternalServerError, Message: "error", Data: map[string]interface{}{"data": err.Error()}})
            return
        }

        //get updated user details
        var updatedUser models.User
        if result.MatchedCount == 1 {
            err := userCollection.FindOne(ctx, bson.M{"id": objId}).Decode(&updatedUser)
            if err != nil {
                c.JSON(http.StatusInternalServerError, responses.UserResponse{Status: http.StatusInternalServerError, Message: "error", Data: map[string]interface{}{"data": err.Error()}})
                return
            }
        }

        c.JSON(http.StatusOK, responses.UserResponse{Status: http.StatusOK, Message: "success", Data: map[string]interface{}{"data": updatedUser}})
   
	}
}

func DeleteAnUser() gin.HandlerFunc{
	return func (c *gin.Context)  {
		ctx, cancel:= context.WithTimeout(context.Background(), 10*time.Second)
		userId:= c.Param("userId")

		defer cancel()

		objId, _ := primitive.ObjectIDFromHex(userId)

		result, err:= userCollection.DeleteOne(ctx, bson.M{"id": objId})

		if err!=nil{
			c.JSON(http.StatusInternalServerError, responses.UserResponse{Status: http.StatusInternalServerError, Message: "Error delete", Data: map[string]interface{}{"data": err.Error()}})
			return 
		}

		if result.DeletedCount<1{
			c.JSON(http.StatusNotFound, responses.UserResponse{Status: http.StatusNotFound, Message: "Error: Not found id", Data: map[string]interface{}{"data": "User with id: "+ userId+ " not found"}})
			return
		}

		c.JSON(http.StatusOK, 
			responses.UserResponse{
				Status: http.StatusOK, 
				Message: "User successfully deleted", 
				Data: map[string]interface{}{"data": result},
			})
	}
}

