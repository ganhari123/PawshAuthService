package RouteHandler

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"

	model "../model"
)

func HandleHttpRoutes() {
	http.HandleFunc("/health", healthHandler)
	http.HandleFunc("/login", loginHandler)
}

func healthHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method == http.MethodGet {
		fmt.Fprintf(w, "Healthy")
	} else {
		http.Error(w, "Invalid request", 405)
	}
}

func loginHandler(w http.ResponseWriter, r *http.Request) {
	var user model.User
	if r.Method == http.MethodPost {
		//valid Login request
		jsonDecode := json.NewDecoder(r.Body)
		err := jsonDecode.Decode(&user)
		if err != nil {
			log.Fatal(err)
		}

		if ok, err := user.VerifyUserCredentials(); ok && err == nil {
			fmt.Println(user)
			//generate JWT token for user

		}
	} else {
		http.Error(w, "Invalid request", 405)
	}
}
