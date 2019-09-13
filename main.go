package main

import (
	"context"
	"log"
	"net/http"
	"time"

	"github.com/k0kubun/pp"
	"github.com/labstack/echo"
	"github.com/labstack/echo/middleware"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

// Log Log
type Log struct {
	Name       string    `json:"name"`
	Productive bool      `json:"productive"`
	StartAt    time.Time `json:"start_at"`
}

// type RequestBody struct {
// 	Name    string `json:"name"`
// 	StartAt string `json:"start_at"`
// }

var db *mongo.Database

func main() {
	var err error
	db, err = getMongoConnection()
	if err != nil {
		pp.Println("cant setup db by %v", err)
	}

	// Echo instance
	e := echo.New()

	// middleware
	e.Use(middleware.Logger())
	e.Use(middleware.Recover())
	// CORS restricted
	// Allows requests from any `https://labstack.com` or `https://labstack.net` origin
	// wth GET, PUT, POST or DELETE method.
	e.Use(middleware.CORSWithConfig(middleware.CORSConfig{
		AllowOrigins: []string{"*"},
		AllowMethods: []string{http.MethodGet, http.MethodPut, http.MethodPost, http.MethodDelete},
	}))

	// routes
	api := e.Group("api")
	api.GET("/logs", showlog)
	api.POST("/logs", writelog)

	// start server
	e.Logger.Fatal(e.Start(":8754"))
}

func seedData(collection *mongo.Collection) {
	all := bson.M{}
	deleteManyRs, err := collection.DeleteMany(context.Background(), all)
	if err != nil {
		log.Fatal("Can not remove logs")
	}
	pp.Println("deleted ", deleteManyRs.DeletedCount, " records")

	initData := []interface{}{}
	first := Log{
		Name:       "set mongo connection",
		Productive: true,
		StartAt:    time.Now().Add(time.Duration(-1) * time.Hour),
	}
	second := Log{
		Name:       "seed data",
		Productive: true,
		StartAt:    time.Now().Add(time.Duration(+1) * time.Hour),
	}

	third := Log{
		Name:       "test sort",
		Productive: false,
		StartAt:    time.Now().Add(time.Duration(0) * time.Hour),
	}

	initData = append(initData, first)
	initData = append(initData, second)
	initData = append(initData, third)

	insertManyRs, err := collection.InsertMany(context.Background(), initData)
	if err != nil {
		log.Fatal("Seed failed")
	}
	pp.Println("inserted these: \n", insertManyRs.InsertedIDs)
}

func showlog(c echo.Context) error {
	collection := db.Collection("logs")

	filter := bson.M{}
	findOptions := options.Find()
	findOptions.SetSort(bson.D{{"startat", -1}})
	cursor, err := collection.Find(context.Background(), filter, findOptions)
	if err != nil {
		return err
	}

	result := []Log{}
	var logRecord Log
	for cursor.Next(context.Background()) {
		err := cursor.Decode(&logRecord)
		if err != nil {
			return err
		}

		result = append(result, logRecord)
	}

	response := make(map[string]interface{})
	response["data"] = result

	return c.JSON(http.StatusOK, response)
}

func writelog(c echo.Context) error {
	rec := new(Log)
	if err := c.Bind(rec); err != nil {
		return err
	}
	rec.StartAt = time.Now().UTC()

	collection := db.Collection("logs")
	insertRs, err := collection.InsertOne(context.Background(), rec)
	if err != nil {
		return err
	}
	log.Println("inserted one record ", insertRs.InsertedID)

	response := make(map[string]interface{})
	response["data"] = rec

	return c.JSON(http.StatusOK, response)
}

var (
	HOSTS    = []string{"localhost:27017"}
	DBNAME   = "devdb"
	USERNAME = "who"
	PASSWORD = "dat"
)

func getMongoConnection() (*mongo.Database, error) {
	config := &options.ClientOptions{
		Hosts: HOSTS,
		// ConnectTimeout: time.Duration(10 * time.Second),
		Auth: &options.Credential{
			Username: USERNAME,
			Password: PASSWORD,
		},
	}

	client, err := mongo.Connect(context.Background(), options.MergeClientOptions(config))
	if err != nil {
		pp.Println("mongo: could not connect to mongodb on %s", HOSTS)
	}

	connection := client.Database(DBNAME)
	pp.Println("mongo: connected to mongdb %s", HOSTS)

	// seed data
	seedData(connection.Collection("logs"))
	// seed data
	return connection, nil
}
