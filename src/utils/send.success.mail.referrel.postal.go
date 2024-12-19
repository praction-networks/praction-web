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

// SendSuccessMailReferrel sends a success email to both the referrer and the referrals.
func SendSuccessMailReferrel(userDetail models.UserRefrence) error {
	logger.Info("Sending success emails for the referral process")

	// Get Postal configuration
	PostalConfig, err := config.POSTALEnvGet()
	if err != nil {
		logger.Error("Unable to retrieve Postal configuration", "Error", err)
		return err
	}

	// Generate the email body for the referrer

	logger.Info("Referel Code is ", "Couppon", userDetail.RefrelCoupon)
	emailBodyForReffredBy, err := generateEmailBodyForReffredBy(userDetail.ReferedBy.Name, userDetail.RefrelCoupon, userDetail.Referels)
	if err != nil {
		logger.Error("Failed to generate email body for referrer", "Error", err)
		return err
	}

	// Set up the postal client
	client := postal.NewClient(PostalConfig.ServerURL, PostalConfig.ApiKey)

	// Send email to the referrer
	msgForReffredBy := &postal.SendRequest{
		To:       []string{userDetail.ReferedBy.Email},
		From:     PostalConfig.EmailFrom,
		Sender:   "Praction Networks",
		Subject:  "Referral Request Successfully Received - Praction Networks",
		BCC:      []string{"info@praction.in"},
		Tag:      "Userreferrels",
		HTMLBody: emailBodyForReffredBy,
	}

	resp, _, err := client.Send.Send(context.TODO(), msgForReffredBy)
	if err != nil {
		logger.Error("Failed to send success email to referrer", "Email", userDetail.ReferedBy.Email, "Error", err)
		return err
	}
	logger.Info("Success email sent to referrer successfully", "Email", userDetail.ReferedBy.Email, "MessageID", resp.MessageID)

	// Now, send emails to the referrals
	for _, referrel := range userDetail.Referels {
		emailBodyForReferrel, err := generateEmailBodyForReferrel(referrel.Name, userDetail.ReferedBy.Name, userDetail.RefrelCoupon)
		if err != nil {
			logger.Error("Failed to generate email body for referral", "Error", err)
			continue // Continue to send emails for other referrals even if one fails
		}

		msgForReferrel := &postal.SendRequest{
			To:       []string{referrel.Email},
			From:     PostalConfig.EmailFrom,
			Sender:   "Praction Networks",
			Subject:  "Welcome to Praction Networks - Your Referral Details",
			BCC:      []string{"info@praction.in"},
			Tag:      "Userreferrels",
			HTMLBody: emailBodyForReferrel,
		}

		// Send email to the referral
		_, _, err = client.Send.Send(context.TODO(), msgForReferrel)
		if err != nil {
			logger.Error("Failed to send email to referral", "Email", referrel.Email, "Error", err)
			continue // Skip this referral and try with the next one
		}
		logger.Info("Success email sent to referral successfully", "Email", referrel.Email)
	}

	return nil
}

// generateEmailBodyForReffredBy generates the email body for the referrer.
func generateEmailBodyForReffredBy(referrerName, refrelCoupon string, referrels []models.UserType) (string, error) {
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
        <h2>Thank You for Your Referral, {{.ReferrerName}}!</h2>
        <p>Dear {{.ReferrerName}},</p>
        <p>Thank you for referring the following individuals to us:</p>
        <div class="details">
            {{range .Referrels}}
            <p><strong>Name:</strong> {{.Name}}</p>
            <p><strong>Mobile:</strong> {{.Mobile}}</p>
            <p><strong>Email:</strong> {{.Email}}</p>
            <p><strong>PinCode:</strong> {{.PinCode}}</p>
            <hr>
            {{end}}
           <p><strong>Your Referral Coupon Code:</strong> {{.RefrelCoupon}}</p>
			<p><strong>How and When to Use the Coupon Code:</strong></p>
			<ul>
			    <li><strong>Valid for 30 days</strong> from the date of receipt.</li>
			    <li>Applicable for plans of <strong>6 months or more</strong> and with a value above <strong>3500 INR</strong>.</li>
			    <li>Once your referred friend successfully completes their <strong>6-month subscription</strong>, you (the referrer) will receive <strong>1 month</strong> added to your subscription.</li>
			    <li>Your referred friend will also receive <strong>one additional month</strong> added to their bundle upon successful completion of the 6-month subscription.</li>
			    <li>This offer can be cancelled at any time with prior notice. All decisions regarding this offer are at the sole discretion of the company.</li>
			</ul>

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
	tmpl, err := template.New("referrerEmail").Parse(emailTemplate)
	if err != nil {
		return "", fmt.Errorf("failed to parse email template for referrer: %w", err)
	}

	var body bytes.Buffer
	data := struct {
		ReferrerName string
		Referrels    []models.UserType
		RefrelCoupon string
	}{
		ReferrerName: referrerName,
		Referrels:    referrels,
		RefrelCoupon: refrelCoupon,
	}

	if err := tmpl.Execute(&body, data); err != nil {
		return "", fmt.Errorf("failed to execute email template for referrer: %w", err)
	}
	logger.Info("Generated email body for referrer: ", "EmailBody", body.String())
	return body.String(), nil
}

// generateEmailBodyForReferrel generates the email body for each referral.
func generateEmailBodyForReferrel(referrelName, referrerName, refrelCoupon string) (string, error) {
	const emailTemplate = `
<!DOCTYPE html>
<html>
<head>
    <style>
        body { font-family: Arial, sans-serif; background-color: #f4f4f9; color: #333; padding: 20px; margin: 0; }
        .container { background-color: #ffffff; padding: 20px; border-radius: 8px; box-shadow: 0 2px 4px rgba(0, 0, 0, 0.1); max-width: 600px; margin: 0 auto; }
        h2 { color: #5A67D8; }
        .details { margin-top: 20px; }
        .details p { margin: 10px 0; }
        .footer { margin-top: 30px; font-size: 12px; color: #888; }
    </style>
</head>
<body>
    <div class="container">
		<div class="logo">
            <img src="https://praction.in/assets/images/logo.png" alt="Praction Networks">
        </div>
        <h2>Welcome, {{.ReferrelName}}!</h2>
        <p>Dear {{.ReferrelName}},</p>
        <p>You have been referred to Praction Networks by {{.ReferrerName}}.</p>
        <p>Here are the details:</p>
        <div class="details">
            <p><strong>Referred By:</strong> {{.ReferrerName}}</p>
            <p><strong>Your Coupon Code:</strong> {{.RefrelCoupon}}</p>
			<p><strong>How and When to Use the Coupon Code:</strong></p>
			<ul>
			    <li><strong>Valid for 30 days</strong> from the date of receipt.</li>
			    <li>Applicable for plans of <strong>6 months or more</strong> and with a value above <strong>3500 INR</strong>.</li>
			    <li>Upon successful completion of your <strong>6-month subscription</strong>, you will receive <strong>one additional month</strong> added to your bundle.</li>
			    <li>Your referrer will also receive <strong>1 month</strong> once you complete your 6-month subscription.</li>
			    <li>This offer can be cancelled at any time with prior notice. All decisions regarding this offer are at the sole discretion of the company.</li>
			</ul>
            </div>
        <div class="footer">
            <p>Best regards,</p>
            <p>Praction Networks</p>
            <p>Email: info@praction.in</p>
            <p>Mobile: +91-9312166166</p>
        </div>
    </div>
</body>
</html>
`
	tmpl, err := template.New("referrelEmail").Parse(emailTemplate)
	if err != nil {
		return "", fmt.Errorf("failed to parse email template for referral: %w", err)
	}

	var body bytes.Buffer
	data := struct {
		ReferrelName string
		ReferrerName string
		RefrelCoupon string
	}{
		ReferrelName: referrelName,
		ReferrerName: referrerName,
		RefrelCoupon: refrelCoupon,
	}

	if err := tmpl.Execute(&body, data); err != nil {
		return "", fmt.Errorf("failed to execute email template for referral: %w", err)
	}

	return body.String(), nil
}
