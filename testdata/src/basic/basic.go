package basic

import (
	"log"
	"log/slog"
)

func examples() {
	log.Print("Starting server on port 8080")   // want `log message should start with a lowercase letter`
	log.Print("starting server on port 8080")
	slog.Error("Failed to connect to database") // want `log message should start with a lowercase letter`
	slog.Error("failed to connect to database")

	log.Print("запуск сервера")                   // want `log message should be in English only`
	log.Print("ошибка подключения к базе данных") // want `log message should be in English only`
	log.Print("starting server")

	log.Print("server started! 🚀") // want `log message should be in English only` `log message should not contain special characters or emojis`
	log.Print("connection failed!!!")              // want `log message should not contain special characters or emojis`
	log.Print("warning: something went wrong...")  // want `log message should not contain special characters or emojis`
	log.Print("server started")
	log.Print("connection failed")
	log.Print("something went wrong")

	password := "secret123"
	log.Print("user password: " + password) // want `log message may contain sensitive data`
	apiKey := "mykey"
	log.Print("api_key=" + apiKey) // want `log message may contain sensitive data`
	token := "mytoken"
	log.Print("token: " + token) // want `log message may contain sensitive data`
	log.Print("user authenticated successfully")
}
