package main

import (
	"context"
	"log"
	"net/http"
	"time"

	"go.mongodb.org/mongo-driver/bson/primitive"

	"github.com/jasonlvhit/gocron"
	"github.com/k0kubun/pp"
	"github.com/labstack/echo"
	"github.com/labstack/echo/middleware"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

// Log Log
type Log struct {
	ID         *primitive.ObjectID `json:"id" bson:"_id,omitempty"`
	Name       string              `json:"name"`
	Productive bool                `json:"productive"`
	StartAt    time.Time           `json:"start_at"`
}

// Worktime interval info
type Worktime struct {
	ID       *primitive.ObjectID `json:"id" bson:"_id,omitempty"`
	LeftEnd  time.Time           `json:"left_end" bson:"left_end"`
	RightEnd time.Time           `json:"right_end" bson:"right_end"`
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

	api.GET("/worktimes", getWorktimes)
	// start server
	e.Logger.Fatal(e.Start(":8754"))

	gocron.Every(5).Seconds().Do(addLog, db.Collection("logs"))
	// function Start start all the pending jobs
	<-gocron.Start()
}

func addLog(collection *mongo.Collection) {
	newOne := Log{
		Name:    "",
		StartAt: time.Now(),
	}
	irs, err := collection.InsertOne(context.Background(), newOne)
	if err != nil {
		log.Fatalf("job add log failed by error %v", err)
	}
	log.Printf("add log success with record %v", irs.InsertedID)
}

func bod(t time.Time) time.Time {
	year, month, day := t.Date()
	return time.Date(year, month, day, 0, 0, 0, 0, t.Location())
}

func seedData(db *mongo.Database) {

	logCollection := db.Collection("logs")
	worktimeCollection := db.Collection("worktimes")

	err := db.Drop(context.Background())
	if err != nil {
		log.Fatalf("drop db failed by error %v", err)
	}

	initData := []interface{}{}
	first := Log{
		Name:       "set mongo connection",
		Productive: true,
		StartAt:    time.Now().UTC().Add(time.Duration(-1) * time.Hour),
	}
	second := Log{
		Name:       "seed data",
		Productive: true,
		StartAt:    time.Now().UTC().Add(time.Duration(+1) * time.Hour),
	}

	third := Log{
		Name:       "test sort",
		Productive: false,
		StartAt:    time.Now().UTC().Add(time.Duration(0) * time.Hour),
	}

	initData = append(initData, first)
	initData = append(initData, second)
	initData = append(initData, third)

	insertManyRs, err := logCollection.InsertMany(context.Background(), initData)
	if err != nil {
		log.Fatalf("Seed logs failed by error %v", err)
	}

	pp.Println("inserted these: \n", insertManyRs.InsertedIDs)

	worktime := Worktime{
		LeftEnd:  bod(time.Now().UTC()).Add(time.Hour * time.Duration(8)),
		RightEnd: bod(time.Now().UTC()).Add(time.Hour * time.Duration(17)),
	}

	ion, err := worktimeCollection.InsertOne(context.Background(), worktime)
	if err != nil {
		log.Fatalf("Seed worktime failed by error %v", err)
	}
	pp.Println("inserted %v", ion.InsertedID)

}

func getWorktimes(c echo.Context) error {
	collection := db.Collection("worktimes")

	filter := bson.M{}
	findOptions := options.Find()
	findOptions.SetSort(bson.D{{"left_end", 1}})
	cursor, err := collection.Find(context.Background(), filter, findOptions)
	if err != nil {
		return err
	}

	result := []Worktime{}
	var worktimeRecord Worktime
	for cursor.Next(context.Background()) {
		err := cursor.Decode(&worktimeRecord)
		if err != nil {
			return err
		}

		result = append(result, worktimeRecord)
	}

	response := make(map[string]interface{})
	response["data"] = result

	return c.JSON(http.StatusOK, response)
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
	seedData(connection)
	// seed data
	return connection, nil
}
