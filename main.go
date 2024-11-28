package main

import (
	"context"
	"fmt"
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
	"go.mongodb.org/mongo-driver/v2/bson"
	"go.mongodb.org/mongo-driver/v2/mongo"
	"go.mongodb.org/mongo-driver/v2/mongo/options"
	"go.mongodb.org/mongo-driver/v2/mongo/readpref"
)

type todo struct {
	ID        	bson.ObjectID `json:"id,omitempty" bson:"_id,omitempty"`
	Item    		string `json:"item"`
	Completed 	bool   `json:"completed"`
}

var (
	collection *mongo.Collection
)

func createTodo(ginContext *gin.Context) {
	var todo todo
	if err := ginContext.ShouldBindJSON(&todo); err != nil {
		ginContext.JSON(http.StatusBadRequest, gin.H{"error": "invalid request body"})
		return
	}

	if todo.Item == "" {
		ginContext.JSON(http.StatusBadRequest, gin.H{"error": "item is required"})
		return
	}

	insertResult, err := collection.InsertOne(context.Background(), todo)
	
	if err != nil {
		ginContext.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
	}

	todo.ID = insertResult.InsertedID.(bson.ObjectID)

	ginContext.IndentedJSON(http.StatusCreated, todo)
}

func getTodos(ginContext *gin.Context) {
	var todos = []todo{}

	cursor, err := collection.Find(context.Background(), bson.M{})

	if err != nil {
		ginContext.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	defer cursor.Close(context.Background())

	for cursor.Next(context.Background()) {
		var todo todo
		if err := cursor.Decode(&todo); err != nil {
			ginContext.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			return
		}
		todos = append(todos, todo)
	}

	ginContext.IndentedJSON(http.StatusOK, todos)
}

func getTodo(ginContext *gin.Context) {
	id := ginContext.Param("id")
	objectId, err := bson.ObjectIDFromHex(id)

	if err != nil {
		ginContext.JSON(http.StatusBadRequest, gin.H{"error": "invalid id format"})
		return
	}

	filter := bson.M{"_id": objectId}

	var todoLocal todo

	if err := collection.FindOne(context.Background(), filter).Decode(&todoLocal); err != nil {
		if err == mongo.ErrNoDocuments {
			ginContext.JSON(http.StatusNotFound, gin.H{"error": err.Error()})
			return
		}
		ginContext.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	ginContext.IndentedJSON(http.StatusOK, todoLocal)
}

func toggleTodoStatus(ginContext *gin.Context) {
	id := ginContext.Param("id")
	objectId, err := bson.ObjectIDFromHex(id)

	if err != nil {
		ginContext.JSON(http.StatusBadRequest, gin.H{"error": "invalid id format"})
		return
	}

	filter := bson.M{"_id": objectId}
	var todo todo

	if err := collection.FindOne(context.Background(), filter).Decode(&todo); err != nil {
		if err == mongo.ErrNoDocuments {
			ginContext.JSON(http.StatusNotFound, gin.H{"error": "document not found"})
			return
		} 
		
		ginContext.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	update := bson.M{"$set": bson.M{"completed": !todo.Completed}}

	_, err = collection.UpdateOne(context.Background(), filter, update)

	if err != nil {
		ginContext.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	todo.Completed = !todo.Completed

	ginContext.IndentedJSON(http.StatusOK, todo)
}

func updateTodo(ginContext *gin.Context) {
	id := ginContext.Param("id")
	objectId, err := bson.ObjectIDFromHex(id)

	if err != nil {
		ginContext.JSON(http.StatusBadRequest, gin.H{"error": "invalid id format"})
		return
	}

	var todoData todo
	if err := ginContext.ShouldBindJSON(&todoData); err != nil {
		ginContext.JSON(http.StatusBadRequest, gin.H{"error": "invalid request body"})
		return
	}

	filter := bson.M{"_id": objectId}
	update := bson.M{"$set": bson.M{"item": todoData.Item}}

	var updatedTodo todo
	err = collection.FindOneAndUpdate(
		context.Background(),
		filter,
		update,
		options.FindOneAndUpdate().SetReturnDocument(options.After),
	).Decode(&updatedTodo)

	if err != nil {
		if err == mongo.ErrNoDocuments {
			ginContext.JSON(http.StatusNotFound, gin.H{"error": "todo not found"})
			return
		}
			
		ginContext.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	ginContext.IndentedJSON(http.StatusOK, updatedTodo)
}



func deleteTodo(ginContext *gin.Context) {
	id := ginContext.Param("id")
	objectId, err := bson.ObjectIDFromHex(id)

	if err != nil {
		ginContext.JSON(http.StatusBadRequest, gin.H{"error": "invalid id format"})
		return
	}

	filter := bson.M{"_id": objectId}

	var deletedTodo todo
	err = collection.FindOneAndDelete(context.Background(), filter).Decode(&deletedTodo)

	if err != nil {
		if err == mongo.ErrNoDocuments {
			ginContext.JSON(http.StatusNotFound, gin.H{"error": "document not found"})
			return
		}
		ginContext.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	ginContext.IndentedJSON(http.StatusOK, deletedTodo)
}

func main() {
	client, _ := mongo.Connect(options.Client().ApplyURI("mongodb://admin:adminpassword@localhost:27017/"))

	defer func() {
    if err := client.Disconnect(context.Background()); err != nil {
        panic(err)	
    }
	}()

	ctx, cancel := context.WithTimeout(context.Background(), 2*time.Second)
	defer cancel()

	if err := client.Ping(ctx, readpref.Primary()); err != nil {
		panic(err)
	}

	fmt.Println("Connected to MONGODB")

	collection = client.Database("todos_db").Collection("todos")

	router := gin.Default()
	router.GET("/todos", getTodos)
	router.POST("/todos", createTodo)
	router.GET("/todos/:id", getTodo)
	router.PATCH("/todos/:id", toggleTodoStatus)
	router.PUT("/todos/:id", updateTodo)
	router.DELETE("/todos/:id", deleteTodo)
	router.Run("localhost:9191")
}
