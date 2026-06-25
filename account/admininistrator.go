package account

import (
	"fmt"

	jwt "github.com/dgrijalva/jwt-go"
	u "hyperball.com/utils"
)

func (login *Administrator) Login() map[string]interface{} {

	var account Administrator

	_, err := DBM.Query(&account, `SELECT * FROM administrators where username=?`, login.Username)
	if err != nil {
		panic(err)
		return u.Message(false, "Connection error. Please retry")
	}
	if account.ID == 0 {
		return u.Message(false, "Username not found!")
	}
	if login.Password != account.Password {
		return u.Message(false, "Invalid password. Please try again")
	}

	//Worked! Logged In
	account.Password = ""
	//Create JWT token
	tk := Token{AccountID: account.ID}
	fmt.Println("TK:", tk)
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, tk)
	tokenString, _ := token.SignedString([]byte("token_mydog_app"))
	//account.Token = tokenString //Store the token in the response

	resp := u.Message(true, "Logged In successful")
	resp["account"] = account
	resp["token"] = tokenString
	data, _err := ParseToken(tokenString)
	if _err != nil {
		panic(_err)
	} else {
		fmt.Println("Parsed Token:", data)
	}

	return resp
}
