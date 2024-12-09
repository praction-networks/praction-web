package utils

import (
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"net/http"

	"github.com/praction-networks/quantum-ISP365/webapp/src/config"
	"github.com/praction-networks/quantum-ISP365/webapp/src/logger"
)

// MSG91ReSendOTP resends an OTP using the MSG91 service
func MSG91ReSendOTP(mobile int64, retryType string) error {
	logger.Info("Resending OTP to user via mobile", "Mobile", mobile)

	// Get MSG91 configuration
	MSG91Config, err := config.MSG91EnvGet()
	if err != nil {
		logger.Warn("Unable to get MSG91 configuration")
		return err
	}

	authKey := MSG91Config.AuthKey
	if authKey == "" {
		logger.Error("Missing configuration: authKey")
		return errors.New("missing SMS gateway configuration")
	}

	// Validate retryType
	if retryType != "text" && retryType != "voice" {
		return errors.New("invalid retryType, must be 'text' or 'voice'")
	}

	// Construct the retry URL
	msg91retryURL := fmt.Sprintf(
		"https://control.msg91.com/api/v5/otp/retry?authkey=%s&retrytype=%s&mobile=%d",
		authKey, retryType, mobile,
	)

	// Create the HTTP GET request
	req, err := http.NewRequest("GET", msg91retryURL, nil)
	if err != nil {
		logger.Error("Failed to generate HTTP request", "Error", err)
		return err
	}

	// Send the HTTP request
	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		logger.Error("Error sending HTTP request", "Error", err)
		return err
	}
	defer resp.Body.Close()

	// Read and parse the response body
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

	// Handle response based on MSG91 API documentation
	if resp.StatusCode == http.StatusOK && msg91Response.Type == "success" {
		logger.Info("OTP resend successful", "Mobile", mobile, "RequestID", msg91Response.RequestID)
		return nil
	}

	logger.Error("OTP resend failed", "Mobile", mobile, "StatusCode", resp.StatusCode, "Message", msg91Response.Message)
	return errors.New(msg91Response.Message)
}
