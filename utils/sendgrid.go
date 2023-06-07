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
	plainTextContent := fmt.Sprintf("Klik tautan di bawah ini untuk mereset password Anda:\n%s", resetLink)
	htmlContent := fmt.Sprintf("<p><img src=\"https://firebasestorage.googleapis.com/v0/b/b-impact-7d077.appspot.com/o/b-impact-icon.jpeg?alt=media&\" alt=\"Logo\"></p><p>Klik tautan di bawah ini untuk mereset password Anda:</p><p><a href=\"%s\">Reset Password</a></p>", resetLink)

	message := mail.NewSingleEmail(from, subject, to, plainTextContent, htmlContent)
	client := sendgrid.NewSendClient(os.Getenv("YOUR_SENDGRID_API_KEY"))

	_, err := client.Send(message)
	if err != nil {
		return err
	}

	return nil
}
