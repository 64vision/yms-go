package account

import (
	"fmt"
	"time"

	"gollux/email"
	u "gollux/utils"

	jwt "github.com/dgrijalva/jwt-go"
)

type Token struct {
	AccountID int
	jwt.StandardClaims
}
type Account struct {
	ID              int           `json:"id"`
	Username        string        `json:"username"`
	FirstName       string        `json:"first_name"`
	LastName        string        `json:"last_name"`
	Gender          string        `json:"gender"`
	BirthDate       string        `json:"birth_date"`
	Email           string        `json:"email"`
	MobileNo        string        `json:"mobile_no"`
	Password        string        `json:"password"`
	CreatedAt       time.Time     `json:"created_at"`
	Level           int           `json:"level"`
	Bank            Bank          `json:"bank"`
	UpdatedAt       time.Time     `json:"updated_at"`
	DeletedAt       *time.Time    `json:"deleted_at"`
	Status          string        `json:"status"`
	Media           []interface{} `json:"media"`
	Code            string        `json:"code"`
	Remarks         string        `json:"remarks"`
	Balance         float64       `json:"balance"`
	Referral        string        `json:"referral"`
	Commission      float64       `json:"commission"`
	BalanceUpdateAt time.Time     `json:"balance_update_at"`
	Agent           int           `json:"agent"`
	SubAgent        int           `json:"sub_agent"`
	Upline          int           `json:"upline"`
	SubUpline       int           `json:"sub_upline"`
	DirectUpline    int           `json:"direct_upline"`
	CompanyName     string        `json:"company_name"`
	CompanyAddress  string        `json:"company_address"`
	CompanyType     string        `json:"company_type"`
}

type Registration struct {
	FirstName string  `json:"first_name"`
	LastName  string  `json:"last_name"`
	Password  string  `json:"password"`
	BirthDate string  `json:"birth_date"`
	Address   Address `json:"address"`
	Value     string  `json:"value"`
	Referral  string  `json:"referral"`
	Type      string  `json:"type"`
	MobileNo  string  `json:"mobile_no"`
}

type VerifyReq struct {
	Code  string `json:"code"`
	Value string `json:"value"`
	Type  string `json:"type"`
}
type LogRequest struct {
	Username string `json:"username"`
	Password string `json:"password"`
	Lat      string `json:"lat"`
	Long     string `json:"long"`
}
type Location struct {
	ID        int       `json:"id"`
	Lat       string    `json:"lat"`
	Long      string    `json:"long"`
	AccountID int       `json:"account_id"`
	LoginAt   time.Time `json:"login_at"`
}
type Address struct {
	City     string `json:"city"`
	Province string `json:"province"`
	Country  string `json:"country"`
}

type Administrator struct {
	ID       int    `json:"id"`
	Username string `json:"username"`
	Fullname string `json:"fullname"`
	Password string `json:"password"`
	Role     string `json:"role"`
	Status   string `json:"status"`
}

func (login *LogRequest) Login() map[string]interface{} {

	var account Account
	fmt.Println("Username attemp: ", login.Username)
	fmt.Println("Pass attemp: ", login.Password)
	_, err := DBM.Query(&account, `SELECT * FROM accounts where username=?`, login.Username)
	if err != nil {
		//panic(err)
		return u.MessageCode(false, "Connection error. Please retry", 5001)
	}
	if account.ID == 0 {
		return u.MessageCode(false, "Username not found!", 5002)
	}
	if login.Password != account.Password {
		return u.MessageCode(false, "Invalid password. Please try again", 5003)
	}
	if account.Remarks == "For Verification" {
		return u.MessageCode(false, "This account is not yet verified. Please verify you account!", 5004)
	}

	account.UpdatedAt = time.Now()
	account.Code = u.GenNumCode(4)

	errdb := DBM.Update(&account)
	if errdb != nil {
		panic(errdb)
		return u.MessageCode(false, "Failed to create account, connection error", 5000)
	}
	//Worked! Logged In
	account.Password = ""

	//Create JWT token
	tk := Token{AccountID: account.ID}
	fmt.Println("TK:", tk)
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, tk)
	tokenString, _ := token.SignedString([]byte("token_mydog_app"))
	//account.Token = tokenString //Store the token in the response

	resp := u.MessageCode(true, "Logged In successful", 2000)
	resp["account"] = account

	if account.Agent != 0 {

		agent, err := GetAccountByID(account.Agent)
		if err != nil {
			fmt.Println("No agent")
		}
		resp["agent"] = map[string]interface{}{"id": agent.ID, "commission": agent.Commission, "mobile_no": agent.MobileNo, "first_name": agent.FirstName, "last_name": agent.LastName, "referral": agent.Referral}
	}

	resp["token"] = tokenString
	data, _err := ParseToken(tokenString)
	if _err != nil {
		panic(_err)
	} else {
		var Loc Location
		Loc.AccountID = account.ID
		Loc.LoginAt = time.Now()
		Loc.Lat = login.Lat
		Loc.Long = login.Long
		Loc.Log()

		fmt.Println("Parsed Token:", data)
	}

	return resp
}

func (loc *Location) Log() {
	_, errdb := DBM.Model(loc).Insert()
	if errdb != nil {
		panic(errdb)
	}

}
func ParseToken(tokenString string) (*Token, error) {
	token, err := jwt.ParseWithClaims(tokenString, &Token{}, func(token *jwt.Token) (interface{}, error) {
		return []byte("token_mydog_app"), nil
	})

	if err != nil {
		return nil, err
	}

	// Extract and validate claims
	if claims, ok := token.Claims.(*Token); ok && token.Valid {
		return claims, nil
	}

	return nil, fmt.Errorf("invalid token")
}

func (verify *VerifyReq) Verify() map[string]interface{} {
	qryStr := "select * from accounts where mobile_no=? and code=?"
	if verify.Type == "email" {
		qryStr = "select * from accounts where email=? and code=?"
	}
	var p Account
	_, err := DBM.Query(&p, qryStr, verify.Value, verify.Code)
	if err != nil {
		panic(err)
		return u.Message(false, "Failed in verification")
	}
	if p.ID == 0 {
		return u.Message(false, "Invalid verification code!")
	} else {
		if p.Remarks != "For Verification" {
			return u.Message(false, "This account is already verified.")
		}
		_, err = DBM.Exec(`UPDATE accounts SET remarks=? where id=? `, "Verified", p.ID)
		if err != nil {
			panic(err)
			return u.Message(false, "Failed in verification update!")
		}
	}

	response := u.Message(true, "Verification Successful!")
	return response
}
func (acct *Account) Validate() (map[string]interface{}, bool) {
	var p Account
	_, err := DBM.Query(&p, `SELECT * FROM accounts where username=?`, acct.Username)
	if err != nil {
		panic(err)
		return u.Message(false, "Failed"), false
	}
	if p.ID != 0 {
		return u.Message(false, "Email already used."), false
	}
	// _, err = DBM.Query(&p, `SELECT * FROM accounts where email=?`, acct.Email)
	// if err != nil {
	// 	panic(err)
	// 	return u.Message(false, "Failed."), false
	// }
	// if p.ID != 0 {
	// 	return u.Message(false, "Email already used."), false
	// }
	//fmt.Println("MobileNo", acct.MobileNo)
	_, err = DBM.Query(&p, `SELECT * FROM accounts where mobile_no=?`, acct.MobileNo)
	if err != nil {
		panic(err)
		return u.Message(false, "Failed."), false
	}
	if p.ID != 0 {
		return u.Message(false, "Mobile No. already in used."), false
	}

	return u.Message(false, "Requirement passed"), true
}

func (reg *Account) Add() map[string]interface{} {
	//	var acct Account
	fmt.Println("create account!")
	var verType string

	if reg.MobileNo == "" {
		if !u.IsValidPHMobileNumber(reg.MobileNo) {
			return u.Message(false, "Invalid Mobile number!")
		}
		return u.Message(true, "Lacking Mobile No. parameter")
	}
	if reg.Email == "" {
		return u.Message(true, "Lacking Email parameter")
	}

	if reg.FirstName == "" {
		return u.Message(true, "Lacking FirstName parameter")
	}
	if reg.LastName == "" {
		return u.Message(true, "Lacking LastName parameter")
	}
	if reg.CompanyName == "" {
		return u.Message(true, "Lacking Company Name parameter")
	}
	if reg.CompanyAddress == "" {
		return u.Message(true, "Lacking Company Address parameter")
	}
	// if reg.BirthDate == "" {
	// 	return u.Message(true, "Lacking BirthDate parameter")
	// }
	if reg.Password == "" {
		return u.Message(true, "Lacking Password parameter")
	}
	if resp, ok := reg.Validate(); !ok {
		return resp
	}
	// age, err := CalculateAge(reg.BirthDate)
	// if err != nil {
	// 	return u.Message(false, err.Error())
	// } else {
	// 	if age < 18 {
	// 		return u.Message(false, "To play you must be 18+ years old")
	// 	}
	// }
	// if reg.Referral != "" {
	// 	agent, _ := GetAccountByReferral(reg.Referral) // player as an agent

	// 	if agent.ID != 0 {
	// 		// if agent.Level == 2 {
	// 		// 	acct.SubAgent = agent.ID
	// 		// 	acct.Agent = agent.Agent
	// 		// } else {
	// 		// 	acct.Agent = agent.ID
	// 		// }
	// 		acct.DirectUpline = agent.ID
	// 		acct.SubUpline = agent.DirectUpline
	// 		acct.Upline = agent.SubUpline
	// 	} else {
	// 		return u.Message(false, "Invalid referral code!")
	// 	}
	// } else {
	// 	acct.Agent = 801500000
	// 	acct.DirectUpline = acct.Agent
	// 	acct.SubUpline = acct.Agent
	// 	acct.Upline = acct.Agent
	// }

	// acct.FirstName = reg.FirstName
	// acct.LastName = reg.LastName
	// acct.Password = u.Md5hash(reg.Password)
	// acct.BirthDate = reg.BirthDate

	code := u.GenNumCode(6)
	reg.Code = code
	reg.CreatedAt = time.Now()
	reg.UpdatedAt = time.Now()
	// acct.Level = 3
	// acct.Referral = u.GenCharCode(6)
	reg.Status = "Active"
	reg.Remarks = "For Verification"
	// smsmsg := "Your verification code is " + code + ". Enter this code to verify your account."
	// sms.Send(acct.MobileNo, smsmsg)

	_, errdb := DBM.Model(reg).Insert()
	if errdb != nil {
		panic(errdb)
		return u.Message(false, "Failed to create account, connection error")
	}

	reg.SendVerificationEmail(reg.Email, code)

	response := u.Message(true, "Account has been registered successfully. Please check your email for verification code.")
	response["account"] = map[string]interface{}{
		"username":   reg.Username,
		"email":      reg.Email,
		"mobile_no":  reg.MobileNo,
		"created_at": reg.CreatedAt,
		"status":     reg.Status,
		"remarks":    reg.Remarks,
	}
	return response
}

func (acct *Account) SendVerificationEmail(to_email string, code string) {
	body := "<p>Hi!,</p><p>Thank you for using ZERA Suite Yard Booking. Here is your verification code.</p>" +
		"<p style='padding: 15px;  color: #000000; font-weight: bold; font-sise: 28px; margin-bottom: 20px; margin-top: 20px; margin-left: 50px;'><b> Verification Code: " + code + "</b><p>" +
		"<p style='color: #999; margin-top: 50px;'>Thank you!<br /> ZERA Suite Team<p>"
	email.Send("Account Verification", body, to_email, "no-reply@zerasuite.com", "ZERA Suite")
}
func CalculateAge(birthdateStr string) (int, error) {
	// Define the date format (adjust if needed)
	layout := "2006-01-02" // YYYY-MM-DD format

	// Parse the string to time.Time
	birthdate, err := time.Parse(layout, birthdateStr)
	if err != nil {
		return 0, fmt.Errorf("invalid date format: %v", err)
	}

	// Get the current date
	now := time.Now()

	// Calculate age
	age := now.Year() - birthdate.Year()

	// Adjust if the birthday hasn't occurred yet this year
	if now.YearDay() < birthdate.YearDay() {
		age--
	}

	return age, nil
}

func GetAccountByReferral(referral string) (*Account, error) {
	var account Account
	_, err := DBM.Query(&account, `SELECT * FROM accounts where referral=?`, referral)
	if err != nil {
		//panic(err)
		return &account, err
	}
	return &account, nil
}
func GetAccountByID(id int) (*Account, error) {
	var account Account
	_, err := DBM.Query(&account, `SELECT * FROM accounts where id=?`, id)
	if err != nil {
		//panic(err)
		return &account, err
	}
	return &account, nil
}

func GetAccountByMobile(mobile string) (*Account, error) {
	var account Account
	_, err := DBM.Query(&account, `SELECT * FROM accounts where mobile_no=?`, mobile)
	if err != nil {
		//panic(err)
		return &account, err
	}

	return &account, nil
}
