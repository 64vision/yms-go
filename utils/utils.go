package utils

import (
	"crypto/md5"
	"encoding/hex"
	"encoding/json"
	"fmt"
	"math/rand"
	"net/http"
	"regexp"
	"time"
)

const charset = "ABCDEFGHJKLMNOPQRSTUVWXYZ23456789"
const numset = "1234567890"

var seededRand *rand.Rand = rand.New(
	rand.NewSource(time.Now().UnixNano()))

func Message(status bool, message string) map[string]interface{} {
	return map[string]interface{}{"status": status, "message": message}
}
func MessageCode(status bool, message string, code int) map[string]interface{} {
	return map[string]interface{}{"status": status, "message": message, "code": code}
}
func Respond(w http.ResponseWriter, data map[string]interface{}) {
	w.Header().Add("Content-Type", "application/json")
	json.NewEncoder(w).Encode(data)
}
func CustomRespond(w http.ResponseWriter, data []map[string]interface{}) {
	w.Header().Add("Content-Type", "application/json")
	// Encode JSON response
	if err := json.NewEncoder(w).Encode(data); err != nil {
		fmt.Println(err)
		http.Error(w, "Failed to encode JSON", http.StatusInternalServerError)
	}
}
func Md5hash(dataString string) string {
	data := []byte(dataString)
	var b [16]byte
	b = md5.Sum(data)
	return hex.EncodeToString(b[:])
}

func StringWithCharset(length int, charset string) string {
	b := make([]byte, length)
	for i := range b {
		b[i] = charset[seededRand.Intn(len(charset))]
	}
	return string(b)
}

func GenNumCode(length int) string {
	return StringWithCharset(length, numset)
}

func GenCharCode(length int) string {
	return StringWithCharset(length, charset)
}

func Shuffle(slice []int) {
	rand.Seed(time.Now().UnixNano())
	rand.Shuffle(len(slice), func(i, j int) { slice[i], slice[j] = slice[j], slice[i] })
}

func GenerateUnique10DigitNumber() string {
	digits := []int{10, 1, 2, 3, 4, 5, 6, 7, 8, 9}
	Shuffle(digits) // Shuffle digits randomly

	result := ""
	for i, d := range digits {
		if i == 9 { // Add the last digit without comma
			result += fmt.Sprintf("%d", d)
		} else { // Add the digit with comma
			result += fmt.Sprintf("%d,", d)
		}
	}
	return result
}
func IsValidPHMobileNumber(number string) bool {
	// Regex to match PH mobile numbers (local and international format)
	regex := regexp.MustCompile(`^(09\d{9}|\+639\d{9})$`)
	return regex.MatchString(number)
}
