package twilio

import (
	"encoding/json"
	"errors"
	"fmt"
	"math/rand"
	"net/http"
	"net/url"
	"strconv"
	"strings"
	"time"
)

const (
	accountSID = "ACc143b38a773f712b06bbe43cec93e13b"
	authToken  = "1eab48aae29b79143de8a9d110b1a4fc"
	urlStr     = "https://api.twilio.com/2010-04-01/Accounts/" + accountSID + "/Messages.json"
	codeLength = 6
)

func SendVerificationCode(toPhoneNumber string) (string, error) {
	msgData := url.Values{}
	msgData.Set("To", toPhoneNumber)
	msgData.Set("From", "+16013361157")
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
		return "", err
	} else if resp.StatusCode >= 200 && resp.StatusCode < 300 {
		var data map[string]interface{}
		decoder := json.NewDecoder(resp.Body)
		err := decoder.Decode(&data)
		if err == nil {
			return code, err
		} else {
			return "", err
		}
	} else {
		return "", errors.New("Response status is out of bounds")
	}

	return "", nil
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
