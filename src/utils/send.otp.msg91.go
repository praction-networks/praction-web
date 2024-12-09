package utils

import (
	"bytes"
	"encoding/json"
	"errors"
	"io"
	"net/http"
	"time"

	"github.com/praction-networks/quantum-ISP365/webapp/src/config"
	"github.com/praction-networks/quantum-ISP365/webapp/src/logger"
)

// MSG91Response represents the response from the MSG91 API
type MSG91Response struct {
	Type      string `json:"type"`
	RequestID string `json:"request_id"`
	Message   string `json:"message"`
}

// MSG91SendOTP sends an OTP using the MSG91 service
func MSG91SendOTP(mobile int64, otp int64) error {
	logger.Info("Initiating Mail for Sending User INtrest Details request", "Mobile", mobile, "OTP", otp)

	MSG91Config, err := config.MSG91EnvGet()

	if err != nil {
		logger.Warn("Unable to get MSG91 Config")
	}

	url := "https://control.msg91.com/api/v5/otp"
	authKey := MSG91Config.AuthKey
	templateID := MSG91Config.TemplateID
	if authKey == "" || templateID == "" {
		logger.Error("Missing configuration: authKey or templateID")
		return errors.New("missing SMS gateway configuration")
	}

	// Prepare request payload
	payload := map[string]interface{}{
		"template_id": templateID,
		"mobile":      mobile,
		"authkey":     authKey,
		"otp_length":  6,
		"otp":         otp,
	}

	payloadBytes, err := json.Marshal(payload)
	if err != nil {
		logger.Error("Error marshalling payload", "Error", err)
		return err
	}

	// Send HTTP request
	client := &http.Client{Timeout: 10 * time.Second}
	req, err := http.NewRequest("POST", url, bytes.NewBuffer(payloadBytes))
	if err != nil {
		logger.Error("Error creating HTTP request", "Error", err)
		return err
	}
	req.Header.Set("Content-Type", "application/json")

	resp, err := client.Do(req)
	if err != nil {
		logger.Error("Error sending HTTP request", "Error", err)
		return err
	}
	defer resp.Body.Close()

	// Read and parse response
	body, err := io.ReadAll(resp.Body)
	if err != nil {
		logger.Error("Error reading response body", "Error", err)
		return err
	}

	var msg91Response MSG91Response
	if err := json.Unmarshal(body, &msg91Response); err != nil {
		logger.Error("Error unmarshalling response", "Error", err)
		return err
	}

	// Handle response
	if resp.StatusCode == http.StatusOK && msg91Response.Type == "success" {
		logger.Info("OTP request successful", "Mobile", mobile, "RequestID", msg91Response.RequestID)
		return nil
	}

	logger.Error("OTP request failed", "Mobile", mobile, "StatusCode", resp.StatusCode, "Message", msg91Response.Message)
	return errors.New(msg91Response.Message)
}
