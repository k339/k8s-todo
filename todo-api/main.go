package main

import (
	"context"
	"fmt"
	"log"
	"net/http"

	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

type Todo struct {
	ID   string `json:"id"`
	Name string `json:"name" binding:"required"`
}

const connectionString = "mongodb://mongo-service"
const dbName = "test"
const collName = "todolist"

var collection *mongo.Collection
var requestCount = 0

func main() {
	connectDB()
	r := gin.Default()
	r.Use(cors.Default())
	r.GET("/api/health", health)
	r.GET("/api/todo", listTodo)
	r.POST("/api/todo", addTodo)
	r.Run()
}

func connectDB() {
	clientOptions := options.Client().ApplyURI(connectionString)
	client, err := mongo.Connect(context.TODO(), clientOptions)

	if err != nil {
		log.Fatal(err)
	}

	err = client.Ping(context.TODO(), nil)

	if err != nil {
		log.Fatal(err)
	}

	fmt.Println("Connected to MongoDB!")
	collection = client.Database(dbName).Collection(collName)
	fmt.Println("Collection instance created!")
}

func health(c *gin.Context) {
	c.Writer.WriteString("ready")
}

func listTodo(c *gin.Context) {
	requestCount++
	cur, err := collection.Find(context.Background(), bson.D{})
	if err != nil {
		log.Fatal(err)
	}
	defer cur.Close(context.Background())
	var res = []Todo{}
	for cur.Next(context.Background()) {
		var result bson.D
		err := cur.Decode(&result)
		if err != nil {
			log.Fatal(err)
		}
		m := result.Map()
		res = append(res, Todo{ID: fmt.Sprintf("%v", m["uuid"]), Name: fmt.Sprintf("%v", m["name"])})
	}
	c.JSON(http.StatusOK, res)
	if requestCount == 100 {
		log.Fatal("limit 100 request")
	}
}

func addTodo(c *gin.Context) {
	var t Todo
	if err := c.ShouldBind(&t); err != nil {
		c.Writer.WriteHeader(http.StatusBadRequest)
		return
	}
	t.ID = string(uuid.New().String())
	_, err := collection.InsertOne(context.Background(), bson.D{{Key: "uuid", Value: t.ID}, {Key: "name", Value: t.Name}})
	if err != nil {
		panic(err)
	}
	c.Writer.WriteHeader(http.StatusCreated)
}
