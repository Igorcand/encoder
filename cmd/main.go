package main

import (
	"encoder/application/services"
	"encoder/adapters/database"
	"encoder/adapters/queue"
	"log"
	"os"
	"strconv"

	"github.com/joho/godotenv"
	"github.com/streadway/amqp"
)

var db database.Database

func init(){
	err := godotenv.Load()
	if err != nil{
		log.Fatalf("error loading .env fil")
	}

	autoMigrateDb, err := strconv.ParseBool(os.Getenv("AUTO_MIGRATE_DB"))
	if err != nil{
		log.Fatalf("error parsing boolean env var")
	}
	debug, err := strconv.ParseBool(os.Getenv("DEBUG"))
	if err != nil{
		log.Fatalf("error parsing boolean env var")
	}

	db.AutoMigrateDb 	= autoMigrateDb
	db.Debug 			= debug
	db.DsnTest 			= os.Getenv("DSN_TEST")
	db.Dsn 				= os.Getenv("DSN")
	db.DbTypeTest 		= os.Getenv("DB_TYPE_TEST")
	db.DbType 			= os.Getenv("DB_TYPE")
	db.Env 				= os.Getenv("ENV")

}

func main(){
	messageChannel := make(chan amqp.Delivery)
	jobReturnChannel := make(chan services.JobWorkerResult)

	dbConnection, err := db.Connect()
	if err != nil{
		log.Fatalf("error connecting to db")
	}
	defer dbConnection.Close()

	rabbitMQ := queue.NewRabbitMQ()
	ch := rabbitMQ.Connect()
	defer ch.Close()

	rabbitMQ.Consume(messageChannel)

	jobManager := services.NewJobManager(dbConnection, rabbitMQ, jobReturnChannel, messageChannel)
	jobManager.Start(ch)


}