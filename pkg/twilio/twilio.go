package twilio

import (
	"encoding/json"
	"errors"
	"fmt"
	"log"
	"math/rand"
	"net/http"
	"net/url"
	"strconv"
	"strings"
	"time"

	model "../model"
)

const (
	accountSID = ""
	authToken  = ""
	urlStr     = "https://api.twilio.com/2010-04-01/Accounts/" + accountSID + "/Messages.json"
	fromNumber = ""
	codeLength = 6
	timeFormat = "20060102150405"
)

var userCodeMap = model.UserCodeMapping{
	UserCodeMap: make(map[string]model.TwilioCode),
}

// SendVerificationCode - This function makes the request to twilio api and sends a code via text message to the user
func SendVerificationCode(email string, toPhoneNumber string) error {
	var twilioCode model.TwilioCode
	msgData := url.Values{}
	msgData.Set("To", toPhoneNumber)
	msgData.Set("From", fromNumber)
	code := generateRandomCode()
	msgData.Set("Body", code)
	msgDataReader := *strings.NewReader(msgData.Encode())

	client := &http.Client{}
	req, _ := http.NewRequest("POST", urlStr, &msgDataReader)
	req.SetBasicAuth(accountSID, authToken)
	req.Header.Add("Accept", "application/json")
	req.Header.Add("Content-Type", "application/x-www-form-urlencoded")

	resp, err := client.Do(req)
	if err != nil {
		return err
	} else if resp.StatusCode >= 200 && resp.StatusCode < 300 {
		var data map[string]interface{}
		decoder := json.NewDecoder(resp.Body)
		err := decoder.Decode(&data)
		if err == nil {
			userCodeMap.Mux.Lock()
			twilioCode.Code = code
			t := time.Now().UTC()
			twilioCode.TimeStamp = t.Format(timeFormat)
			userCodeMap.UserCodeMap[email] = twilioCode
			userCodeMap.Mux.Unlock()
			return err
		}
		return err
	}

	return errors.New("Response status is out of bounds")
}

// VerifyCode - Verifies the code sent to the user
func VerifyCode(email string, code string) (bool, error) {
	var twilioCode model.TwilioCode
	userCodeMap.Mux.Lock()
	twilioCode = userCodeMap.UserCodeMap[email]
	delete(userCodeMap.UserCodeMap, email)
	userCodeMap.Mux.Unlock()
	origTime, err := time.Parse(timeFormat, twilioCode.TimeStamp)
	if err != nil {
		log.Println(err)
		return false, err
	}
	tSince := time.Since(origTime).Seconds()
	if tSince > 100 {
		return false, errors.New("Code expired, request a new code")
	}
	if twilioCode.Code != code {
		return false, errors.New("Code incorrect please request a new code")
	}
	return true, nil
}

func generateRandomCode() string {
	rand.Seed(time.Now().UnixNano())
	var str = ""
	for i := 0; i < codeLength; i++ {
		num := rand.Intn(10)
		tempStr := strconv.Itoa(num)
		str = str + tempStr
	}

	fmt.Println(str)
	return str
}
