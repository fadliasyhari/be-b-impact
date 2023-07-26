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
	htmlContent := fmt.Sprintf(`
	<div>
  <div style="align-items: flex-start; background-color: var(--baseblue-50); border: 1px none; padding: 20.0px 30px; width: 350px;">
    <div style="align-items: center; background-color: var(--basewhite); gap: 40px; margin-top: 0; padding: 10px; width: 480px;">
      <img src="https://anima-uploads.s3.amazonaws.com/projects/64ba02d41ae0244a5013366d/releases/64bf8ef1001762aedcd3ff05/img/group-10@2x.png" alt="Group 10" style="height: 48px; position: absolute; top: %s; left: %s; width: 146.45px; padding-bottom: 20px;" />
      <div style="align-items: flex-start; align-self: stretch; flex: 0 0 auto; gap: 32px; width: %s;">
        <div style="align-items: flex-start; align-self: stretch; flex: 0 0 auto; gap: 40px; width: %s;">
          <p style="align-self: stretch; color: var(--black); font-weight: 400; line-height: 22px; margin-top: -1.00px; position: relative;">
            Someone (hopefully you) has requested a password reset for your B-IMPACT account. Click the button below to set a new password:
          </p>
          <div style="align-self: stretch; color: var(--foundation-blueblue-500); font-weight: 400; line-height: 22px; position: relative; text-decoration: underline;">
            %s
          </div>
          <p style="align-self: stretch; color: var(--black); font-weight: 400; line-height: 22px; position: relative;">
            The link will expire in 3 hours. If you don&#39;t wish to reset your password, disregard this email and no action will be taken.<br /><br />Thanks,<br />The B-IMPACT Team
          </p>
        </div>
				<hr style="align-self: stretch; border-bottom: 0.1px solid #5F737F; position: relative; width: %s;" />
        <div style="align-items: center; align-self: stretch; flex: 0 0 auto; gap: 4px; width: %s;">
          <p style="align-self: stretch; color: #5F737F; font-weight: 400; line-height: 22px; margin-top: -1.00px; position: relative; text-align: center;">
            You&#39;re receiving this email because a password reset was requested for your account.
          </p>
          <div style="align-self: stretch; color: #5F737F; font-weight: 400; line-height: 22px; position: relative; text-align: center; white-space: nowrap;">
            © 2023 B-IMPACT
          </div>
        </div>
      </div>
    </div>
  </div>
</div>

	`, "50%", "50%", "100%", "100%", resetLink, "100%", "100%")

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
	<div>
	<div style="align-items: flex-start; background-color: var(--baseblue-50); border: 1px none; padding: 20.0px 30px; width: 400px;">
		<div style="align-items: center; background-color: var(--basewhite); flex-direction: column; gap: 40px; margin-top: 0; padding: 10px; position: relative; width: %s;">
			<img
				style="height: 48px; align-self: stretch; position: absolute; top: %s; left: %s; width: 146.45px; padding-bottom: 20px; margin: 0 auto;"
				src="https://anima-uploads.s3.amazonaws.com/projects/64ba02d41ae0244a5013366d/releases/64bf8ef1001762aedcd3ff05/img/group-10@2x.png"
				alt="Group 10"
			/>
			<div style="align-items: flex-start; align-self: stretch; flex-direction: column; gap: 32px; position: relative; width: %s;">
				<div style="align-items: flex-start; align-self: stretch; flex-direction: column; gap: 40px; position: relative; width: %s;">
					<div style="align-self: stretch; color: var(--black2); font-weight: 800; line-height: 22px; position: relative; font-size: 18px;">Hi,</div>
					<p style="align-self: stretch; color: var(--black2); font-weight: 400; line-height: 22px; position: relative;">
						Welcome to B-IMPACT. To verify your account, please use the following OTP code to complete your
						registration.
					</p>
					<div style="align-self: stretch; height: 52px; position: absolute; top: %s; left: %s; width: %s; margin: 0 auto;">
						<div style="align-items: center; background-color: #EFF1F2; border: 1px solid;  border-color: #B5BFC4; border-radius: 8px;  display: inline-flex; gap: 10px; justify-content: center; left: 161px; padding: 10px; position: relative;">
							<div style="color: var(--black); font-family: 'Inter'; font-size: 18px; font-weight: 700; letter-spacing: 6.00px; line-height: 32px; margin-top: -1.00px; position: relative; text-align: center; white-space: nowrap; width: fit-content; display: flex; flex-direction: column; justify-content: center;">%s</div>
						</div>
					</div>
					<p style="align-self: stretch; color: var(--black2); font-weight: 400; line-height: 22px; position: relative;">
						The OTP code is only valid for 5 minutes. Please do not show this OTP to anyone.<br /><br />Thanks,<br />The
						B-IMPACT Team
					</p>
				</div>
				<hr style="align-self: stretch; border: 0.01px solid #5F737F; position: relative; width: %s;" />
				<div style="align-items: center; align-self: stretch; flex-direction: column; gap: 4px; position: relative; width: %s;">
					<p style="align-self: stretch; color: #5F737F; font-weight: 400; line-height: 22px; margin-top: -1.00px; position: relative; text-align: center;">
						You&#39;re receiving this email because an account was registered with this e-mail.
					</p>
					<div style="align-self: stretch; color: #5F737F; font-weight: 400; line-height: 22px; position: relative; text-align: center; white-space: nowrap;">© 2023 B-IMPACT</div>
				</div>
			</div>
		</div>
	</div>
</div>


				`, "100%", "50%", "50%", "100%", "100%", "100%", "50%", "50%", otp, "100%", "100%")

	message := mail.NewSingleEmail(from, subject, to, "", htmlContent)
	client := sendgrid.NewSendClient(os.Getenv("YOUR_SENDGRID_API_KEY"))

	_, err := client.Send(message)
	if err != nil {
		return err
	}

	return nil
}
