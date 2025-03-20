package utils

import (
	"context"
	"fmt"

	"github.com/Pacerino/postal-go"
	"github.com/praction-networks/quantum-ISP365/webapp/src/config"
	"github.com/praction-networks/quantum-ISP365/webapp/src/logger"
)

func OTPOverPostalMail(email []string, otp int64) error {
	logger.Info("Initiating OTP request", "Email", email, "OTP", otp)

	PostalConfig, err := config.POSTALEnvGet()
	if err != nil {
		logger.Warn("Unable to get Postal Config")
		return err
	}

	// Dynamically insert the OTP into the email body
	emailBody := fmt.Sprintf(`
<!DOCTYPE html>
<html>
<head>
    <style>
        body {
            font-family: Arial, sans-serif;
            background-color: #f4f4f9;
            color: #333;
            margin: 0;
            padding: 20px;
        }
        .container {
            max-width: 600px;
            background: #ffffff;
            border-radius: 10px;
            box-shadow: 0px 4px 6px rgba(0, 0, 0, 0.1);
            margin: 0 auto;
            padding: 20px;
        }
        .header {
            text-align: center;
            border-bottom: 2px solid #5A67D8;
            padding-bottom: 10px;
        }
        .header img {
            max-width: 120px;
        }
        .content {
            margin-top: 20px;
        }
        .content h2 {
            color: #5A67D8;
        }
        .footer {
            margin-top: 30px;
            text-align: center;
            font-size: 12px;
            color: #888;
        }
        .social-links img {
            width: 24px;
            margin: 0 5px;
            vertical-align: middle;
        }
    </style>
</head>
<body>
    <div class="container">
        <div class="header">
            <img src="https://praction.in/assets/images/Praction-logo-original.svg" alt="Praction Networks" style="max-width: 300px;" class="CtoWUD">
        </div>
        <div class="content">
            <h2>One Time Password (OTP)</h2>
            <p>Dear User,</p>
            <p>Your Praction Verification Code is:</p>
            <h1 style="text-align: center; color: #5A67D8;">%d</h1>
            <p>Please do not share this code with anyone for your security.</p>
        </div>
        <div class="footer">
            <p>Best regards,</p>
            <p><strong>Praction Networks Team</strong></p>
            <p>Contact Us:</p>
            <p>Email: <a href="mailto:info@praction.in">info@praction.in</a> | Phone: +91-9312166166</p>
            <div class="social-links">
                <a href="https://www.facebook.com/practionnetworks" target="_blank"><img src="https://img.icons8.com/color/48/000000/facebook-new.png" alt="Facebook"></a>
                <a href="https://www.instagram.com/practionnetworks" target="_blank"><img src="https://img.icons8.com/color/48/000000/instagram-new.png" alt="Instagram"></a>
                <a href="https://www.linkedin.com/company/praction-networks/" target="_blank"><img src="https://img.icons8.com/color/48/000000/linkedin.png" alt="LinkedIn"></a>
            </div>
            <p>&copy; 2024 Praction Networks. All rights reserved.</p>
        </div>
    </div>
</body>
</html>
`, otp)

	client := postal.NewClient(PostalConfig.ServerURL, PostalConfig.ApiKey)

	msg := &postal.SendRequest{
		To:       email,
		From:     PostalConfig.EmailFrom,
		Sender:   "Praction Networks",
		Subject:  "Your Praction OTP Code",
		Tag:      "OTP",
		HTMLBody: emailBody,
	}

	resp, _, err := client.Send.Send(context.TODO(), msg)

	if err != nil {
		logger.Warn("Failed to send OTP to User over", "Email:", email)
		return err
	}

	logger.Info("OTP successfully sent to ", "Email:", email, "MessageID:", resp.MessageID)

	return nil
}
