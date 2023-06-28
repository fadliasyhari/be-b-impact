package utils

import (
	"fmt"
	"os"

	"github.com/sendgrid/sendgrid-go"
	"github.com/sendgrid/sendgrid-go/helpers/mail"
)

func SendResetEmail(email string, resetLink string) error {
	from := mail.NewEmail("B-Impact", "geikisen@gmail.com")
	to := mail.NewEmail("", email)
	subject := "Reset Password"
	plainTextContent := fmt.Sprintf("Click on the link below to set your password:\n%s", resetLink)
	htmlContent := fmt.Sprintf("<p><img src=\"https://firebasestorage.googleapis.com/v0/b/b-impact-7d077.appspot.com/o/b-impact-icon.jpeg?alt=media&\" alt=\"Logo\"></p><p>Click on the link below to set your password:</p><p><a href=\"%s\">Reset Password</a></p>", resetLink)

	message := mail.NewSingleEmail(from, subject, to, plainTextContent, htmlContent)
	client := sendgrid.NewSendClient(os.Getenv("YOUR_SENDGRID_API_KEY"))

	_, err := client.Send(message)
	if err != nil {
		return err
	}

	return nil
}

func SendOTP(email string, otp string) error {
	from := mail.NewEmail("B-Impact", "geikisen@gmail.com")
	to := mail.NewEmail("", email)
	subject := "OTP Verification"
	// plainTextContent := fmt.Sprintf("Kode OTP anda: %s", otp)
	htmlContent := fmt.Sprintf(`
	<div style="font-family: Helvetica,Arial,sans-serif;min-width:1000px;overflow:auto;line-height:2">
  <div style="margin:20px 50px;width: %s;">
	<p><img src="https://firebasestorage.googleapis.com/v0/b/b-impact-7d077.appspot.com/o/b-impact-icon.jpeg?alt=media&" alt="Logo" style="width:%s;"></p>
    <div style="border-bottom:1px solid #eee">
      <a href="" style="font-size:1.4em;color: #00466a;text-decoration:none;font-weight:600">B-Impact</a>
    </div>
    <h3>Hi,</h3>
		<h2>Welcome to the club.</h2>
    <p>Thank you for joining B-Impact. Use the following OTP to complete your Sign Up procedures.<br> OTP is valid for 5 minutes</p>
    <h2 style="background: #00466a;margin: 20px 50px;width: max-content;padding: 0 10px;color: #fff;border-radius: 4px;">%s</h2>
    <p style="font-size:0.9em;">Regards,<br />B-Impact</p>
    <hr style="border:none;border-top:1px solid #eee" />
  </div>
</div>
				`, "70%", "30%", otp)

	message := mail.NewSingleEmail(from, subject, to, "", htmlContent)
	client := sendgrid.NewSendClient(os.Getenv("YOUR_SENDGRID_API_KEY"))

	_, err := client.Send(message)
	if err != nil {
		return err
	}

	return nil
}
