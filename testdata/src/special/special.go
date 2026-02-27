package special

import (
	"log"
)

func examplesSpecial() {
	log.Print("server started! 🚀") // want `log message should not contain special characters or emojis`
	log.Print("connection failed!!!")              // want `log message should not contain special characters or emojis`
	log.Print("warning: something went wrong...")  // want `log message should not contain special characters or emojis`
	log.Print("server started")
	log.Print("connection failed")
	log.Print("something went wrong")
	log.Print("request completed with status 200")
}
