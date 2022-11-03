package sendemail

import (
	"net/smtp"
)

func Sendemail(f string, pass string, t []string, b string, h string, p string, s string) error {

	//host := "smtp.gmail.com"
	//port := "587"
	address := h + ":" + p

	//toEmailAddress := "<paste the email address you want to send to>"

	// subject := "Subject: This is the subject of the mail\n"
	// body := "This is the body of the mail"
	message := []byte(s + b)

	auth := smtp.PlainAuth("", f, pass, h)

	err := smtp.SendMail(address, auth, f, t, message)
	if err != nil {
		return err
	}
	return err
}

// func ConfirmLink(e string) (string, error) {
// 	timeNow := time.Now().String()
// 	confirmLink, err := bcrypt.GenerateFromPassword([]byte(e+timeNow), 8)
// 	if err != nil {
// 		return "", err
// 	}

// 	return string(confirmLink), err

// }
