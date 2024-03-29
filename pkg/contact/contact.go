package contact

import (
	"context"
	"fmt"

	brevo "github.com/sendinblue/APIv3-go-library/v2/lib"
)

func SendResetEmail(email string, oneTimeToken string) error {
	emailContect := fmt.Sprintf(
	`<html>
		<body>
			<p>
				Your account (%s) made a request to reset your password. If you would like to do so please click this link:
			</p>
			<br>
			<a>
				https://jasontaks.com/login/reset/%s
			</a>
		</body>
	</html>`, 
	email, oneTimeToken)
	emailModel := brevo.SendSmtpEmail{
		Sender: &brevo.SendSmtpEmailSender{
			Name:  "Kyler Bomhof",
			Email: "contact@jasontasks.com",
		},
		To:          []brevo.SendSmtpEmailTo{{Email: email}},
		Subject:     "Jasontasks forgot password request",
		HtmlContent: emailContect,
	}

	var ctx context.Context
	cfg := brevo.NewConfiguration()
	// TODO: Clear this and convert it to a kubernetes secret
	cfg.AddDefaultHeader("api-key", "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")

	//Configure API key authorization: partner-key
	// cfg.AddDefaultHeader("partner-key", "YOUR_API_KEY")

	sib := brevo.NewAPIClient(cfg)

	result, resp, err := sib.TransactionalEmailsApi.SendTransacEmail(ctx, emailModel)
	if err != nil {
		fmt.Println("Error when calling AccountApi->get_account: ", err.Error())
		return err
	}
	fmt.Println("GetAccount Object:", result, " GetAccount Response: ", resp)
	return err
}
