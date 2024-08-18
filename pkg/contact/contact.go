package contact

import (
	"context"
	"fmt"
	"os"

	brevo "github.com/sendinblue/APIv3-go-library/v2/lib"
)

func SendResetEmail(email string, oneTimeToken string) error {
	emailContent := fmt.Sprintf(
		`<html>
		<body>
			<p>
				Your account (%s) made a request to reset your password. If you would like to do so please click this link:
			</p>
			<br>
			<a href="https://jasontasks.com/login/reset?id=%s">
				Click here
			</a>
		</body>
	</html>`,
		email, oneTimeToken)

	emailModel := brevo.SendSmtpEmail{
		Sender: &brevo.SendSmtpEmailSender{
			Name:  "Contact",
			Email: "contact@jasontasks.com",
		},
		To:          []brevo.SendSmtpEmailTo{{Email: email}},
		Subject:     "Jasontasks forgot password request",
		HtmlContent: emailContent,
	}

	var ctx context.Context
	cfg := brevo.NewConfiguration()

	cfg.AddDefaultHeader("api-key", os.Getenv("EMAIL_API_KEY"))

	sib := brevo.NewAPIClient(cfg)

	result, resp, err := sib.TransactionalEmailsApi.SendTransacEmail(ctx, emailModel)

	if err != nil {
		fmt.Println("Error when calling AccountApi->get_account: ", err.Error())
		fmt.Println("Result: ", result)
		fmt.Println("Response", resp)
	}
	return err
}
