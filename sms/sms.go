package sms

import (
	"bytes"
	"fmt"
	"io/ioutil"
	"mime/multipart"
	"net/http"
	"strconv"
)

const (
	PROJECT_API_KEY = "key here" //""
	URL             = "https://*****/messaging/v1/sms/push"
)

func Send(_mobile_no string, _message string) (bool, string) {
	fmt.Println("Sending SMS")
	method := "POST"
	payload := &bytes.Buffer{}
	writer := multipart.NewWriter(payload)
	_ = writer.WriteField("message", _message)
	_ = writer.WriteField("number", _mobile_no)
	err := writer.Close()
	if err != nil {
		fmt.Println(err)
	}
	client := &http.Client{
		CheckRedirect: func(req *http.Request, via []*http.Request) error {
			return http.ErrUseLastResponse
		},
	}
	req, err := http.NewRequest(method, URL, payload)
	if err != nil {
		fmt.Println(err)
	}
	req.Header.Add("X-TXTBOX-Auth", PROJECT_API_KEY)
	req.Header.Set("Content-Type", writer.FormDataContentType())
	res, err := client.Do(req)
	if err != nil {
		panic(err)
		return false, "Failed"
	}
	defer res.Body.Close()
	body, err := ioutil.ReadAll(res.Body)
	if err != nil {
		panic(err)
		return false, "Failed"
	} else {
		fmt.Println(string(body))
		return true, string(body)
	}
	return false, "Failed"
}

func SendWinNotification(accountID int) {
	url := "https://hyperballrace-default-rtdb.asia-southeast1.firebasedatabase.app/kahyper/winnoti/"
	account := strconv.Itoa(accountID)
	url = url + account + ".json"
	// JSON payload (data to update)
	jsonData := []byte(`{
	"iswinner": 1,
	"canplay": 0,
	"announcement":0,
	"message": ""
	}`)

	// Create a new PATCH request
	req, err := http.NewRequest("PATCH", url, bytes.NewBuffer(jsonData))
	if err != nil {
		fmt.Println("Error creating request:", err)
		return
	}

	// Set headers
	req.Header.Set("Content-Type", "application/json")

	// Send request using http.DefaultClient
	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		fmt.Println("Error sending request:", err)
		return
	}
	defer resp.Body.Close()

	// Read response
	body, _ := ioutil.ReadAll(resp.Body)
	fmt.Println("Response:", string(body))
}

func ClearWinNotification() {
	url := "https://hyperballrace-default-rtdb.asia-southeast1.firebasedatabase.app/kahyper/winnoti.json"

	// JSON payload (data to update)
	jsonData := []byte(`{
	"clear": 1
	}`)

	// Create a new PATCH request
	req, err := http.NewRequest("DELETE", url, bytes.NewBuffer(jsonData))
	if err != nil {
		fmt.Println("Error creating request:", err)
		return
	}

	// Set headers
	req.Header.Set("Content-Type", "application/json")

	// Send request using http.DefaultClient
	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		fmt.Println("Error sending request:", err)
		return
	}
	defer resp.Body.Close()

	// Read response
	body, _ := ioutil.ReadAll(resp.Body)
	fmt.Println("Response:", string(body))
}
