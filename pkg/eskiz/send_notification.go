package eskiz

import (
	"bytes"
	"encoding/json"
	"errors"
	"net/http"
)

type SendSMSRequest struct {
	MobilePhone string `json:"mobile_phone"`
	Message     string `json:"message"`
	From        string `json:"from"`
}

type SendSMSResponse struct {
	ID      string `json:"id"`
	Message string `json:"message"`
	Status  string `json:"status"`
}

func SendNotification(token, phone, message string) error {
	url := "https://notify.eskiz.uz/api/message/sms/send"

	// Create request payload
	payload := SendSMSRequest{
		MobilePhone: phone,
		Message:     message,
		From:        "4546",
	}
	jsonPayload, _ := json.Marshal(payload)

	// Create HTTP request
	req, err := http.NewRequest("POST", url, bytes.NewBuffer(jsonPayload))
	if err != nil {
		return err
	}

	// Set headers
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Authorization", "Bearer "+token)

	// Send request
	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	// Decode response
	var smsResp SendSMSResponse
	if err := json.NewDecoder(resp.Body).Decode(&smsResp); err != nil {
		return err
	}

	// Check for success
	if smsResp.Status != "waiting" {
		return errors.New("failed to send SMS: " + smsResp.Message)
	}

	return nil
}

func main() {
	// Example usage
	token := "eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJleHAiOjE3NDA2NzU1MDUsImlhdCI6MTczODA4MzUwNSwicm9sZSI6InRlc3QiLCJzaWduIjoiYjFkMDEyZDg0YWE3MDU3YzMwYWRmZjVmZmRlNDNmY2U3MzJjODhhNmUwZWNlOGRmMmYyN2VkNDAzOTdkOWVmYiIsInN1YiI6Ijk2NzcifQ.5WHh4CIU-Qf46bIdYvXpMWQFKxLZHm5n-IljFVgs0rg"
	phone := "+998934628018"
	message := "Это тест от Eskiz"

	err := SendNotification(token, phone, message)
	if err != nil {
		panic(err)
	}
}
