package RouteHandler

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"

	jwt "../jwt"
	model "../model"
	twilio "../twilio"
	util "../util"
)

type userCode map[string]string

func HandleHttpRoutes() {
	http.HandleFunc("/health", healthHandler)
	http.HandleFunc("/login", loginHandler)
	http.HandleFunc("/register", registerHandler)
	//http.HandleFunc("/verifyRegistrationCode", verifyRegistrationCode)
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
				http.Error(w, "JWT token generation error", 400)
				return
			}

			response := util.HttpResponse{
				Body:    token,
				Success: true,
				Error:   "",
			}
			res, err := json.Marshal(response)
			if err != nil {
				http.Error(w, "Response json marshalling error", 400)
				return
			}

			fmt.Fprintf(w, string(res))
			return
		}

		http.Error(w, "Invalid user credentials", 400)
		return
	}
	http.Error(w, "Invalid request", 405)
	return
}

func registerHandler(w http.ResponseWriter, r *http.Request) {
	var user *model.User
	if r.Method == http.MethodPost {
		jsonDecode := json.NewDecoder(r.Body)
		err := jsonDecode.Decode(&user)
		if err != nil {
			http.Error(w, "json decoder error", 400)
			return

		}

		success, err := user.AddUserToUserTable()
		if !success && (err != nil) {
			http.Error(w, "User addition to table failed", 400)
			return
		}

		code, err := twilio.SendVerificationCode(user.PhoneNumber)
		if err != nil {
			http.Error(w, "twilio verification code error", 400)
			return
		}
		fmt.Println(code)
	}
}
