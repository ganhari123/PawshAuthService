package model

type User struct {
	Username    string `json:"username"`
	Password    string `json:"password"`
	FullName    string `json:"full_name"`
	Address     string `json:"address"`
	PhoneNumber string `json:"phone_number"`
}

func (u *User) VerifyUserCredentials() (bool, error) {
	if u.Username == "pawsh" {
		if u.Password == "password" {
			return true, nil
		}
	}
	return false, nil
}
