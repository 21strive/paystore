package paystore

import (
	"github.com/gofiber/fiber/v2"
	"os"
	"paystore/command"
	"paystore/config"
	"paystore/fetch"
	"paystore/lib/helper"
)

func main() {
	writeDB := helper.CreatePostgresConnection(
		os.Getenv("DB_WRITE_HOST"), os.Getenv("DB_WRITE_PORT"), os.Getenv("DB_WRITE_USER"),
		os.Getenv("DB_WRITE_PASSWORD"), os.Getenv("DB_WRITE_NAME"), os.Getenv("DB_WRITE_SSLMODE"))
	defer writeDB.Close()
	readDB := helper.CreatePostgresConnection(
		os.Getenv("DB_READ_HOST"), os.Getenv("DB_READ_PORT"), os.Getenv("DB_READ_USER"),
		os.Getenv("DB_READ_PASSWORD"), os.Getenv("DB_READ_NAME"), os.Getenv("DB_READ_SSLMODE"))
	defer readDB.Close()
	redis := helper.ConnectRedis(os.Getenv("REDIS_HOST"), os.Getenv("REDIS_USER"),
		os.Getenv("REDIS_PASS"), false)

	config := config.DefaultConfig(os.Getenv("VENDOR_TABLE_NAME"))

	client := command.New(writeDB, readDB, redis, config)
	httpClient := command.NewHTTPHandler(client)

	app := fiber.New()
	app.Post("/webhook/xendit", httpClient.ReceivePayment)

	app.Listen(":" + os.Getenv("PORT"))
}

type HTTPFetcherHandler struct {
	paystoreFetcher *fetch.PaystoreFetcher
}
