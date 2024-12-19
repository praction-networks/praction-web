package utils

import (
	"fmt"
	"strconv"

	"github.com/praction-networks/quantum-ISP365/webapp/src/logger"
)

// SendOTP sends an OTP via SMS and email.
func SendOTP(email string, mobile string, otp int64, retryType string, resend bool) error {
	if resend && retryType == "voice" {
		// Prepend country code "91" to the mobile number
		mobileWithCode := "91" + mobile

		// Convert the mobile number to int64
		mobileNum, err := strconv.ParseInt(mobileWithCode, 10, 64)
		if err != nil {
			logger.Error("Error converting mobile number to int64", "Mobile", mobileWithCode, "Error", err)
			return fmt.Errorf("failed to convert mobile string to int64")
		}

		// Send OTP via SMS
		errMsg91 := MSG91ReSendOTP(mobileNum, retryType, otp)

		if errMsg91 != nil {
			logger.Warn("Failed to send OTP via SMS; user can use email OTP for verification", "Email", email, "Mobile", mobileWithCode, "Error", errMsg91)
			return nil
		}

		// Log success
		logger.Info("OTP sent successfully via both SMS and email", "Email", email, "Mobile", mobileWithCode)
		return nil
	}

	if resend && retryType == "text" {
		// Prepend country code "91" to the mobile number
		mobileWithCode := "91" + mobile

		// Convert the mobile number to int64
		mobileNum, err := strconv.ParseInt(mobileWithCode, 10, 64)
		if err != nil {
			logger.Error("Error converting mobile number to int64", "Mobile", mobileWithCode, "Error", err)
			return fmt.Errorf("failed to convert mobile string to int64")
		}

		// Send OTP via SMS
		errMsg91 := MSG91ReSendOTP(mobileNum, retryType, otp)

		// Send OTP via email
		errPostal := OTPOverPostalMail([]string{email}, otp)

		// Handle failure cases
		if errMsg91 != nil && errPostal != nil {
			logger.Error("Failed to resend OTP via both SMS and email", "Email", email, "Mobile", mobileWithCode)
			return fmt.Errorf("failed to resend OTP via both SMS and email")
		}

		if errMsg91 != nil {
			logger.Warn("Failed to resend OTP via SMS; user can use email OTP for verification", "Email", email, "Mobile", mobileWithCode, "Error", errMsg91)
			return nil
		}

		if errPostal != nil {
			logger.Warn("Failed to resend OTP via email; user can use SMS OTP for verification", "Email", email, "Mobile", mobileWithCode, "Error", errPostal)
			return nil
		}

		// Log success
		logger.Info("OTP resent successfully via both SMS and email", "Email", email, "Mobile", mobileWithCode)
		return nil
	}

	// Prepend country code "91" to the mobile number
	mobileWithCode := "91" + mobile

	// Convert the mobile number to int64
	mobileNum, err := strconv.ParseInt(mobileWithCode, 10, 64)
	if err != nil {
		logger.Error("Error converting mobile number to int64", "Mobile", mobileWithCode, "Error", err)
		return fmt.Errorf("failed to convert mobile string to int64")
	}

	// Send OTP via SMS
	errMsg91 := MSG91SendOTP(mobileNum, otp)

	// Send OTP via email
	errPostal := OTPOverPostalMail([]string{email}, otp)

	// Handle failure cases
	if errMsg91 != nil && errPostal != nil {
		logger.Error("Failed to send OTP via both SMS and email", "Email", email, "Mobile", mobileWithCode)
		return fmt.Errorf("failed to send OTP via both SMS and email")
	}

	if errMsg91 != nil {
		logger.Warn("Failed to send OTP via SMS; user can use email OTP for verification", "Email", email, "Mobile", mobileWithCode, "Error", errMsg91)
		return nil
	}

	if errPostal != nil {
		logger.Warn("Failed to send OTP via email; user can use SMS OTP for verification", "Email", email, "Mobile", mobileWithCode, "Error", errPostal)
		return nil
	}

	// Log success
	logger.Info("OTP sent successfully via both SMS and email and Valid For 30 Minutes", "Email", email, "Mobile", mobileWithCode)
	return nil

}
