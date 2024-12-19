package utils

import (
	"bytes"
	"context"
	"fmt"
	"html/template"
	"strings"

	"github.com/Pacerino/postal-go"
	"github.com/praction-networks/quantum-ISP365/webapp/src/config"
	"github.com/praction-networks/quantum-ISP365/webapp/src/logger"
	"github.com/praction-networks/quantum-ISP365/webapp/src/models"
)

// SendSuccessMailIntrestForPlan sends a success email to the user.
func SendSuccessMailIntrestForPlan(userDetail models.AvailableUserRequest, planDetails models.PlanSpecific) error {
	logger.Info("Sending success email to the user")

	// Get Postal configuration
	PostalConfig, err := config.POSTALEnvGet()
	if err != nil {
		logger.Error("Unable to retrieve Postal configuration", "Error", err)
		return err
	}

	// Extract user and plan details
	name := userDetail.FirstName
	mobile := userDetail.Mobile
	email := userDetail.Email
	planName := planDetails.Name
	planTextSpeed := fmt.Sprintf("%.2f %s", planDetails.Speed, planDetails.SpeedUnit)
	planPrice := planDetails.Price * float64(planDetails.Period)
	planPriceText := fmt.Sprintf("%.2f", planPrice)
	planTextValidity := fmt.Sprintf("%d %s", planDetails.Period, planDetails.PeriodUnit)
	specificFeatures, remainingOfferings := FilterOfferings(planDetails.Offering)

	// Generate the email body using a template
	emailBody, err := generateEmailBodyforPlan(name, mobile, email, planName, planTextSpeed, planTextValidity, planPriceText, specificFeatures, remainingOfferings)
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
		BCC:      []string{"info@praction.in", "sales@praction.in"},
		Tag:      "UserPlanRequest",
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

func generateEmailBodyforPlan(
	name, mobile, email, planName, planTextSpeed, planTextValidity, planPriceText string,
	specificFeatures, remainingOfferings map[string]bool,
) (string, error) {
	const emailTemplate = `
<!DOCTYPE html>
<html>
<head>
    <style>
        body {
            font-family: 'Arial', sans-serif;
            background-color: #f4f4f9;
            color: #333;
            margin: 0;
            padding: 20px;
        }
        .container {
            background-color: #ffffff;
            border-radius: 8px;
            box-shadow: 0 4px 8px rgba(0, 0, 0, 0.1);
            max-width: 600px;
            margin: 0 auto;
            overflow: hidden;
        }
        .header {
            background: linear-gradient(90deg, #6a11cb, #2575fc);
            color: white;
            padding: 20px;
            text-align: center;
        }
        .header h1 {
            margin: 0;
            font-size: 24px;
        }
        .logo {
            text-align: center;
            margin: 20px 0;
        }
        .logo img {
            width: 150px;
        }
        .details {
            padding: 20px;
        }
        .details p {
            margin: 10px 0;
            font-size: 16px;
        }
        .details .highlight {
            font-weight: bold;
            color: #6a11cb;
        }
        ul {
            list-style-type: none;
            padding: 0;
        }
        ul li {
            display: flex;
            align-items: center;
            margin-bottom: 10px;
            font-size: 16px;
        }
        ul li img {
            width: 16px;
            height: 16px;
            margin-right: 10px;
        }
        .cta {
            text-align: center;
            margin: 20px 0;
        }
        .cta a {
            background: linear-gradient(90deg, #6a11cb, #2575fc);
            color: white;
            padding: 12px 20px;
            border-radius: 5px;
            text-decoration: none;
            font-size: 16px;
            box-shadow: 0 4px 6px rgba(0, 0, 0, 0.1);
        }
        .cta a:hover {
            background: linear-gradient(90deg, #2575fc, #6a11cb);
        }
        .footer {
            background: #f4f4f9;
            text-align: center;
            padding: 10px;
            font-size: 12px;
            color: #888;
        }
        .social-links img {
            width: 24px;
            margin: 0 5px;
        }
    </style>
</head>
<body>
    <div class="container">
        <div class="header">
            <h1>Welcome to Praction Networks!</h1>
        </div>
        <div class="logo" style="text-align: center; margin-bottom: 20px;">
            <img src="https://praction.in/assets/images/logo.png" alt="Praction Networks" style="width: 200px; height: auto;">
        </div>
        <div class="details">
            <p>Hello <span class="highlight">{{.Name}}</span>,</p>
            <p>Thank you for showing interest in our services. Here's a summary of the plan you're interested in:</p>

            <p>Plan Details:</p>
                <ul style="list-style-type: none; padding: 0; margin: 0;">
                    <li style="display: flex; align-items: center; margin-bottom: 10px;">
                        <img src="https://img.icons8.com/color/48/000000/star--v1.png" alt="Plan" style="margin-right: 10px;">
                        <strong style="margin-right: 15px;">Plan:</strong>
                        <span>{{.PlanName}}</span>
                    </li>
                    <li style="display: flex; align-items: center; margin-bottom: 10px;">
                        <img src="https://img.icons8.com/color/48/000000/speedometer.png" alt="Speed" style="margin-right: 10px;">
                        <strong style="margin-right: 15px;">Speed:</strong>
                        <span>{{.PlanTextSpeed}}</span>
                    </li>
                    <li style="display: flex; align-items: center; margin-bottom: 10px;">
                        <img src="https://img.icons8.com/color/48/000000/rupee.png" alt="Price" style="margin-right: 10px;">
                        <strong style="margin-right: 15px;">Price:</strong>
                        <span>₹ {{.PlanPriceText}}</span>
                    </li>
                    <li style="display: flex; align-items: center; margin-bottom: 10px;">
                        <img src="https://img.icons8.com/color/48/000000/calendar.png" alt="Validity" style="margin-right: 10px;">
                        <strong style="margin-right: 15px;">Validity:</strong>
                        <span>{{.PlanTextValidity}}</span>
                    </li>
                </ul>
        </div>

        <div class="details">
            <p><strong>Plan Features:</strong></p>
            <ul style="margin-bottom: 20px;">
                {{range $key, $value := .RemainingOfferings}}
                {{if $value}}
                <li style="margin-bottom: 5px;">
                    <img src="https://img.icons8.com/color/48/000000/checked-checkbox.png" alt="Available">
                    <strong>{{$key}}:</strong> Wow! You got this feature included!
                </li>
                {{else}}
                <li style="margin-bottom: 5px;">
                    <img src="https://img.icons8.com/fluency/48/000000/close-window.png" alt="Missing">
                    <strong>{{$key}}:</strong> Oops! This feature isn’t included. Upgrade your plan to unlock it!
                </li>
                {{end}}
                {{end}}
            </ul>
        </div>

        <div class="details">
            <p><strong>Installation and Security:</strong></p>
            <ul style="margin-bottom: 20px;">
                {{range $key, $value := .SpecificFeatures}}
                {{if $value}}
                <li style="margin-bottom: 10px;">
                    <img src="https://img.icons8.com/color/48/000000/sad.png" alt="Required">
                    <strong>{{$key}}:</strong> Oh no! You need to pay for this amount with this plan. You Can chose diffrent plan we will cover this for you!
                </li>
                {{else}}
                <li style="margin-bottom: 10px;">
                    <img src="https://img.icons8.com/color/48/000000/happy.png" alt="Free">
                    <strong>{{$key}}:</strong> Wow! We cover security deposit and installation cost for you with this plan!
                </li>
                {{end}}
                {{end}}
            </ul>
        </div>

        <div class="cta">
            <a href="https://praction.in" target="_blank">Explore More Plans</a>
        </div>
            <!-- Footer Section -->
            <div class="footer">
                <p>Follow us for updates:</p>
                <div class="social-links">
                    <a href="https://www.facebook.com/practionnetworks">
                        <img src="https://img.icons8.com/color/64/000000/facebook.png" alt="Facebook">
                    </a>
                    <a href="https://www.instagram.com/practionnetworks">
                        <img src="https://img.icons8.com/color/64/000000/instagram-new.png" alt="Instagram">
                    </a>
                    <a href="https://www.linkedin.com/company/praction-networks/">
                        <img src="https://img.icons8.com/color/64/000000/linkedin.png" alt="LinkedIn">
                    </a>
                </div>
                <div class="contact-info">
                    <p><strong>Phone:</strong> <a href="tel:+919312166166" class="highlight">+91 93121 66166</a></p>
                    <p><strong>Email:</strong> <a href="mailto:info@praction.in" class="highlight">info@praction.in</a></p>
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
		Name               string
		Mobile             string
		Email              string
		PlanName           string
		PlanTextSpeed      string
		PlanTextValidity   string
		PlanPriceText      string
		RemainingOfferings map[string]bool
		SpecificFeatures   map[string]bool
	}{
		Name:               name,
		Mobile:             mobile,
		Email:              email,
		PlanName:           planName,
		PlanTextSpeed:      planTextSpeed,
		PlanTextValidity:   planTextValidity,
		PlanPriceText:      planPriceText,
		RemainingOfferings: remainingOfferings,
		SpecificFeatures:   specificFeatures,
	}

	if err := tmpl.Execute(&body, data); err != nil {
		return "", fmt.Errorf("failed to execute email template: %w", err)
	}

	return body.String(), nil
}

// FilterOfferings separates specific features into a separate map and keeps the rest.
func FilterOfferings(offerings map[string]bool) (map[string]bool, map[string]bool) {
	specificFeatures := map[string]bool{}
	remainingOfferings := map[string]bool{}

	for key, value := range offerings {
		if strings.Contains(key, "Security Deposit") || strings.Contains(key, "Installation Cost") || strings.Contains(key, "Installation Fee") {
			// Add to specific features map
			specificFeatures[key] = value
		} else {
			// Add to remaining offerings map
			remainingOfferings[key] = value
		}
	}

	return specificFeatures, remainingOfferings
}
