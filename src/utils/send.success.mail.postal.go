package utils

import (
	"bytes"
	"context"
	"fmt"
	"html/template"

	"github.com/Pacerino/postal-go"
	"github.com/praction-networks/quantum-ISP365/webapp/src/config"
	"github.com/praction-networks/quantum-ISP365/webapp/src/logger"
	"github.com/praction-networks/quantum-ISP365/webapp/src/models"
)

// SendSuccessMail sends a success email to the user.
func SendSuccessMailIntrest(userDetail models.UserInterest) error {
	logger.Info("Sending success email to the user")

	// Get Postal configuration
	PostalConfig, err := config.POSTALEnvGet()
	if err != nil {
		logger.Error("Unable to retrieve Postal configuration", "Error", err)
		return err
	}

	name := userDetail.Name
	mobile := userDetail.Mobile
	email := userDetail.Email
	pinCode := userDetail.PinCode
	// Extract user details based on type

	// Generate the email body using a template
	emailBody, err := generateEmailBody(name, mobile, email, pinCode)
	if err != nil {
		logger.Error("Failed to generate email body", "Error", err)
		return err
	}

	// Set up the postal client
	client := postal.NewClient(PostalConfig.ServerURL, PostalConfig.ApiKey)

	// Create the email message
	msg := &postal.SendRequest{
		To:       []string{email},
		From:     PostalConfig.EmailFrom,
		Sender:   "Praction Networks",
		Subject:  "Request Successfully Received - Praction Networks",
		BCC:      []string{"info@praction.in"},
		Tag:      "UserRequest",
		HTMLBody: emailBody, // Use the HTML body here
	}

	// Send the email
	resp, _, err := client.Send.Send(context.TODO(), msg)
	if err != nil {
		logger.Error("Failed to send success email", "Email", email, "Error", err)
		return err
	}

	// Log success
	logger.Info("Success email sent successfully", "Email", email, "MessageID", resp.MessageID)
	return nil
}

// generateEmailBody generates the email body using an HTML template.
func generateEmailBody(name, mobile, email string, pinCode int64) (string, error) {
	const emailTemplate = `
<!DOCTYPE html>
<html>
<head>
    <style>
        body {
            font-family: Arial, sans-serif;
            background-color: #f4f4f9;
            color: #333;
            padding: 20px;
            margin: 0;
        }
        .container {
            background-color: #ffffff;
            padding: 20px;
            border-radius: 8px;
            box-shadow: 0 2px 4px rgba(0, 0, 0, 0.1);
            max-width: 600px;
            margin: 0 auto;
        }
        h2 {
            color: #5A67D8;
        }
        .details {
            margin-top: 20px;
        }
        .details p {
            margin: 10px 0;
        }
        .footer {
            margin-top: 30px;
            font-size: 12px;
            color: #888;
        }
        .footer p {
            margin: 5px 0;
        }
    </style>
</head>
<body>
    <div class="container">
	    <div class="logo">
            <img src="https://praction.in/assets/images/logo.png" alt="Praction Networks" style="max-width: 300px;" class="CtoWUD">
        </div>
        <h2>Thank You for Your Interest!</h2>
        <p>Dear {{.Name}},</p>
        <p>Thank you for showing your interest in our services. Our representative will be in touch with you shortly.</p>
        <div class="details">
            <p><strong>Name:</strong> {{.Name}}</p>
            <p><strong>Mobile:</strong> {{.Mobile}}</p>
            <p><strong>Email:</strong> {{.Email}}</p>
            <p><strong>PinCode:</strong> {{.PinCode}}</p>
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
`

	tmpl, err := template.New("email").Parse(emailTemplate)
	if err != nil {
		return "", fmt.Errorf("failed to parse email template: %w", err)
	}

	var body bytes.Buffer
	data := struct {
		Name    string
		Mobile  string
		Email   string
		PinCode int64
	}{
		Name:    name,
		Mobile:  mobile,
		Email:   email,
		PinCode: pinCode,
	}

	if err := tmpl.Execute(&body, data); err != nil {
		return "", fmt.Errorf("failed to execute email template: %w", err)
	}

	return body.String(), nil
}
