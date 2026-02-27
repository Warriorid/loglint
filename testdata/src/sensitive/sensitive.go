package sensitive

import (
	"log"
)

func examplesSensitive() {
	password := "secret123"
	log.Print("user password: " + password) // want `log message may contain sensitive data`

	apiKey := "mykey"
	log.Print("api_key=" + apiKey) // want `log message may contain sensitive data`

	token := "mytoken"
	log.Print("token: " + token) // want `log message may contain sensitive data`

	log.Print("user authenticated successfully")
	log.Print("api request completed")
	log.Print("token validated") // want `log message may contain sensitive data`
}
