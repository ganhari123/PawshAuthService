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
	http.HandleFunc("/verifyRegistrationCode", verifyRegistrationCode)
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
			log.Println(err)
			http.Error(w, "Response json marshalling error", 400)
			return
		}

		ok, err := user.VerifyUserCredentials()
		if err != nil {
			log.Println(err)
			http.Error(w, "Response json marshalling error", 400)
			return
		}

		if ok == "User does not exist" {
			log.Println("User does not exist")
			http.Error(w, "User does not exist", 400)
			return
		}

		if ok == "User is not verified" {
			user, err := user.GetUserDetails()
			if err != nil {
				log.Println(err)
				http.Error(w, "Unable to obtain user", 400)
				return
			}
			log.Println(user.PhoneNumber)
			err = twilio.SendVerificationCode(user.Email, user.PhoneNumber)
			if err != nil {
				log.Println(err)
				http.Error(w, "twilio verification code error", 400)
				return
			}

			log.Println("User is not verified yet")
			response := util.HttpResponse{
				Body:    "User is not verified yet",
				Success: false,
				Error:   "User is not verified yet",
			}
			res, err := json.Marshal(response)
			if err != nil {
				http.Error(w, "Response json marshalling error", 400)
				return
			}

			fmt.Fprintf(w, string(res))
			return
		}

		if ok == "User is verified" {
			//generate JWT token for user
			token, err := jwt.GenerateJwtAccessToken(user.Email)
			if err != nil {
				http.Error(w, "JWT token generation error", 400)
				return
			}

			response := util.HttpResponse{
				Body:    token,
				Success: true,
				Error:   "None",
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
			log.Println(err)
			http.Error(w, "json decoder error", 400)
			return
		}

		success, err := user.AddUserToUserTable()
		if !success && (err != nil) {
			log.Println(err)
			http.Error(w, "User addition to table failed", 400)
			return
		}

		err = twilio.SendVerificationCode(user.Email, user.PhoneNumber)
		if err != nil {
			log.Println(err)
			http.Error(w, "twilio verification code error", 400)
			return
		}

		response := util.HttpResponse{
			Body:    "Registration successful.. Awaiting verification",
			Success: true,
			Error:   "None",
		}
		res, err := json.Marshal(response)
		if err != nil {
			log.Println(err)
			http.Error(w, "Response json marshalling error", 400)
			return
		}
		fmt.Fprintf(w, string(res))
		return
	}
	http.Error(w, "Invalid request", 405)
	return
}

func verifyRegistrationCode(w http.ResponseWriter, r *http.Request) {
	var user *model.User
	if r.Method == http.MethodPost {
		jsonDecode := json.NewDecoder(r.Body)
		err := jsonDecode.Decode(&user)
		if err != nil {
			log.Println(err)
			http.Error(w, "json decoder error", 400)
			return
		}

		isVerified, err := twilio.VerifyCode(user.Email, user.VerificationCode)
		if err != nil {
			log.Println(err)
			http.Error(w, "unable to verify code", 400)
			return
		}

		if isVerified {
			err = user.UpdateUserTableVerified()
			if err != nil {
				log.Println(err)
				http.Error(w, "unable to verify code", 400)
				return
			}

			token, err := jwt.GenerateJwtAccessToken(user.Email)
			if err != nil {
				log.Println(err)
				http.Error(w, "unable to generate JWT token", 400)
				return
			}
			response := util.HttpResponse{
				Body:    token,
				Success: true,
				Error:   "None",
			}
			res, err := json.Marshal(response)
			if err != nil {
				log.Println(err)
				http.Error(w, "Response json marshalling error", 400)
				return
			}
			fmt.Fprintf(w, string(res))
			return
		}

		response := util.HttpResponse{
			Body:    "Registration unsuccessful",
			Success: false,
			Error:   "User not verified correctly",
		}
		res, err := json.Marshal(response)
		if err != nil {
			http.Error(w, "Response json marshalling error", 400)
			return
		}
		fmt.Fprintf(w, string(res))
		return
	}
}
