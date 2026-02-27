package english

import (
	"log"
)

func examplesEnglish() {
	log.Print("запуск сервера")                   // want `log message should be in English only`
	log.Print("ошибка подключения к базе данных") // want `log message should be in English only`
	log.Print("starting server")
	log.Print("failed to connect to database")
	log.Print("Привет мир") // want `log message should be in English only`
}
