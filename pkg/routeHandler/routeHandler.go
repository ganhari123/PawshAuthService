package RouteHandler

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"

	jwt "../jwt"
	model "../model"
	util "../util"
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
			//generate JWT token for user
			token, err := jwt.GenerateJwtToken(&user)
			if err != nil {
				log.Fatal(err)
				return
			}

			response := util.HttpResponse{
				Body:    token,
				Success: true,
				Error:   "",
			}

			res, err := json.Marshal(response)
			if err != nil {
				log.Fatal(err)
				return
			}
			fmt.Fprintf(w, string(res))
			return
		}

		response := util.HttpResponse{
			Body:    "",
			Success: false,
			Error:   "Username or password was incorrect",
		}

		_, err = json.Marshal(response)
		if err != nil {
			log.Fatal(err)
			return
		}
		http.Error(w, "Invalid request", 405)
	} else {
		http.Error(w, "Invalid request", 405)
	}
}
