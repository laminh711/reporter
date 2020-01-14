package main

import (
	"context"
	"fmt"
	"log"
	"net/http"
	"time"

	"github.com/k0kubun/pp"
	"github.com/labstack/echo"
	"github.com/labstack/echo/middleware"
	"github.com/spf13/viper"

	"github.com/laminh711/reporter/models"
	_wlh "github.com/laminh711/reporter/worklog/delivery/http"
	_wlr "github.com/laminh711/reporter/worklog/repository"
	_wlu "github.com/laminh711/reporter/worklog/usecase"
	_wh "github.com/laminh711/reporter/worktime/delivery/http"
	_wr "github.com/laminh711/reporter/worktime/repository"
	_wu "github.com/laminh711/reporter/worktime/usecase"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"

	"golang.org/x/crypto/bcrypt"
)

var db *mongo.Database

func init() {
	viper.SetConfigFile(`config.json`)
	if err := viper.ReadInConfig(); err != nil {
		panic(err)
		// TODO
		// if _, ok := err.(viper.ConfigFileNotFoundError); ok {

		// } else {

		// }
	}

	if viper.GetBool(`debug`) {
		fmt.Println("Running on DEBUG mode")
	}
}

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

	worktimeRepository := _wr.NewMongoWorktimeRepository(db, "worktimes")
	// TODO find out what this timeout use
	worktimeUsecase := _wu.NewWorktimeUsecase(worktimeRepository, 30)
	_wh.NewWorktimeHandler(e, worktimeUsecase)

	worklogRepository := _wlr.NewWorklogRepository(db, "worklogs")
	worklogUsecase := _wlu.NewWorklogUsecase(worklogRepository)
	_wlh.NewWorklogHandler(e, worklogUsecase)

	//
	_ = db.Collection("users").Drop(context.Background())
	password, _ := bcrypt.GenerateFromPassword([]byte("admin"), bcrypt.DefaultCost)
	adminUser := models.User{
		Username: "admin",
		Password: string(password),
	}
	_, err = db.Collection("Users").InsertOne(context.Background(), adminUser)
	if err != nil {
		pp.Printf("insert admin error by %v", err)
	}
	// modify later
	e.POST("api/auth", func(ec echo.Context) error {
		type LoginRequest struct {
			Username string
			Password string
		}
		type LoginSuccess struct {
			Message string `json:"message"`
		}
		type LoginFailed struct {
			Message string `json:"message"`
		}
		request := new(LoginRequest)
		if err := ec.Bind(request); err != nil {
			return ec.JSON(http.StatusBadRequest, LoginFailed{err.Error()})
		}

		// db

		// if bcrypt.CompareHashAndPassword()

		if request.Username == "admin" && request.Password == "admin" {
			return ec.JSON(http.StatusOK, LoginSuccess{"ok"})
		}

		return ec.JSON(http.StatusOK, LoginFailed{"Wrong credentials"})
	})

	// start server
	e.Logger.Fatal(e.Start(":8754"))

	// gocron.Every(5).Seconds().Do(addLog, db.Collection("logs"))
	// function Start start all the pending jobs
	// <-gocron.Start()
}

func addLog(collection *mongo.Collection) {
	newOne := models.Worklog{
		Name:       "",
		FinishedAt: time.Now(),
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

// func seedData(db *mongo.Database) {

// 	logCollection := db.Collection("logs")
// 	worktimeCollection := db.Collection("worktimes")

// 	err := db.Drop(context.Background())
// 	if err != nil {
// 		log.Fatalf("drop db failed by error %v", err)
// 	}

// 	initData := []interface{}{}
// 	first := models.Worklog{
// 		Name:       "set mongo connection",
// 		Productive: true,
// 		StartAt:    time.Now().UTC().Add(time.Duration(-1) * time.Hour),
// 	}
// 	second := models.Worklog{
// 		Name:       "seed data",
// 		Productive: true,
// 		StartAt:    time.Now().UTC().Add(time.Duration(+1) * time.Hour),
// 	}

// 	third := models.Worklog{
// 		Name:       "test sort",
// 		Productive: false,
// 		StartAt:    time.Now().UTC().Add(time.Duration(0) * time.Hour),
// 	}

// 	initData = append(initData, first)
// 	initData = append(initData, second)
// 	initData = append(initData, third)

// 	insertManyRs, err := logCollection.InsertMany(context.Background(), initData)
// 	if err != nil {
// 		log.Fatalf("Seed logs failed by error %v", err)
// 	}

// 	pp.Println("inserted these: \n", insertManyRs.InsertedIDs)

// 	worktime := models.Worktime{
// 		LeftEnd:  bod(time.Now().UTC()).Add(time.Hour * time.Duration(8)),
// 		RightEnd: bod(time.Now().UTC()).Add(time.Hour * time.Duration(17)),
// 	}

// 	ion, err := worktimeCollection.InsertOne(context.Background(), worktime)
// 	if err != nil {
// 		log.Fatalf("Seed worktime failed by error %v", err)
// 	}
// 	pp.Println("inserted %v", ion.InsertedID)

// }

func getWorktimes(c echo.Context) error {
	collection := db.Collection("worktimes")

	filter := bson.M{}
	findOptions := options.Find()
	findOptions.SetSort(bson.D{{"left_end", 1}})
	cursor, err := collection.Find(context.Background(), filter, findOptions)
	if err != nil {
		return err
	}

	result := []models.Worktime{}
	var worktimeRecord models.Worktime
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

	result := []models.Worklog{}
	var logRecord models.Worklog
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
	rec := new(models.Worklog)
	if err := c.Bind(rec); err != nil {
		return err
	}
	rec.FinishedAt = time.Now().UTC()

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

func getMongoConnection() (*mongo.Database, error) {
	dbHost := viper.GetStringSlice("database.host")
	dbName := viper.GetString("database.db")
	dbUsername := viper.GetString("database.username")
	dbPassword := viper.GetString("database.password")
	ctxTimeout := viper.GetInt("context.timeout")
	connectTimeout := time.Duration(ctxTimeout) * time.Second

	config := &options.ClientOptions{
		Hosts:          dbHost,
		ConnectTimeout: &connectTimeout,
		Auth: &options.Credential{
			Username: dbUsername,
			Password: dbPassword,
		},
	}

	client, err := mongo.Connect(context.Background(), options.MergeClientOptions(config))
	if err != nil {
		pp.Printf("mongo: could not connect to mongodb on %s by err %v", dbHost[0], err)
	}
	connection := client.Database(dbName)
	pp.Printf("mongo: connected to mongdb %s", dbHost[0])

	// seed data
	// seedData(connection)
	// seed data
	return connection, nil
}
