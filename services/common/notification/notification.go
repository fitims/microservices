package notification

import (
	"log"
	"os"
	"path"
	"path/filepath"
	"thaThrowdown/common/infrastructure"
)

type (
	RegistrationEmail struct {
		Name  string
		Email string
		Token string
	}

	ResetPasswordEmail struct {
		Name  string
		Email string
		Token string
	}

	NewPassword struct {
		Name        string
		NewPassword string
	}
)

// SendRegistrationEmail will compose and send a registration email to the provided user
func SendRegistrationEmail(name string, email string, token string) error {
	msg := infrastructure.MailMessage{
		From:      "Autoresponse <auto@thathrowdown.com>",
		To:        email,
		Subject:   "Thanks for registering !",
		PlainBody: "",
	}

	tplData := RegistrationEmail{
		Name:  name,
		Email: email,
		Token: token,
	}

	ex, err := os.Executable()
	exPath := path.Join(filepath.Dir(ex), "/templates/register_email.html")

	log.Println("Path = ", exPath)

	err = msg.SendTemplatedMail(exPath, tplData)
	if err != nil {
		log.Println("SendRegistrationEmail - error sending email : ", err.Error())
		return err
	}

	return nil
}

// SendResetPasswordEmail will compose and send a password request email to the provided user
func SendResetPasswordEmail(name string, email string, token string) error {
	msg := infrastructure.MailMessage{
		From:      "Autoresponse <auto@thathrowdown.com>",
		To:        email,
		Subject:   "Password Reset Request !",
		PlainBody: "",
	}

	tplData := ResetPasswordEmail{
		Name:  name,
		Email: email,
		Token: token,
	}

	ex, err := os.Executable()
	exPath := path.Join(filepath.Dir(ex), "/templates/reset_password_email.html")

	log.Println("Path = ", exPath)

	err = msg.SendTemplatedMail(exPath, tplData)
	if err != nil {
		log.Println("SendResetPasswordEmail - error sending email : ", err.Error())
		return err
	}

	return nil
}

// SendNewPasswordEmail will compose and send a new password email to the provided user
func SendNewPasswordEmail(name string, email string, newPassword string) error {
	msg := infrastructure.MailMessage{
		From:      "Autoresponse <auto@thathrowdown.com>",
		To:        email,
		Subject:   "New Password !",
		PlainBody: "",
	}

	tplData := NewPassword{
		Name:        name,
		NewPassword: newPassword,
	}

	ex, err := os.Executable()
	exPath := path.Join(filepath.Dir(ex), "/templates/new_password_email.html")

	log.Println("Path = ", exPath)

	err = msg.SendTemplatedMail(exPath, tplData)
	if err != nil {
		log.Println("SendNewPasswordEmail - error sending email : ", err.Error())
		return err
	}

	return nil
}
