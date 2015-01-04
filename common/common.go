package common

import (
	"errors"
	"fmt"
	"html/template"
	"log"
	"net/http"
	"os"
	"regexp"

	"github.com/keighl/mandrill"

	"github.com/iangudger/basicGoAPI/database"
)

const (
	// Define templates
	serverErrorTemplateText = `<!DOCTYPE html>
<html>
	<head>
		<title>
			500 Internal Server Error
		</title>
	</head>
	<body>
		<h1>
			500 Internal Server Error
		</h1>
		<p>
			The server has encountered an unexpected condition that prevented it from fulfilling the request by your client or web browser.
		</p>
		<p>
			Detailed error message:
			<br />
			{{.Message}}
		</p>
	</body>
</html>
`
)

// Errors
var InvalidUsernameOrPassword = errors.New("Invalid username or password.")
var EmailFailed = errors.New("Sending email failed.")
var InvalidSessionID = errors.New("Invalid sessionid.")
var InvalidPassword = errors.New("Invalid password.")
var DatabaseError = errors.New("Unknown database error.")

// Misc constants
const (
	minPasswordLength = 6
)

// Shared template store
var templates = make(map[string]*template.Template)

// Mandrill
var mandrillKey = os.Getenv("MANDRILL_APIKEY")

const (
	fromEmail = "donotreply@example.com"
	fromName  = "Basic Go API"
)

// Regex
var emailRegex *regexp.Regexp
var priceRegex *regexp.Regexp

func init() {
	templates["serverError"] = template.Must(template.New("serverError").Parse(serverErrorTemplateText))
	emailRegex = regexp.MustCompile("^.+@.+\\..+$")
}

func InternalServerError(w http.ResponseWriter, message string) {
	w.Header().Set("Content-Type", "text/html; charset=utf-8")
	w.WriteHeader(http.StatusInternalServerError)
	ExecTemplate(templates["serverError"], w, map[string]interface{}{"Message": message})
}

func ExecTemplate(tmpl *template.Template, w http.ResponseWriter, pc map[string]interface{}) {
	if err := tmpl.Execute(w, pc); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}
}

func Login(db database.DB, username string, password string) (sessionid string, err error) {
	userid, err := db.GetUserID(username, password)
	if err != nil {
		log.Printf("Error while logging in user (%s): %s\n", username, err.Error())
		err = InvalidUsernameOrPassword
		return
	}

	sessionid, err = db.NewSession(userid)
	if err != nil {
		log.Printf("Error while creating session for user (%s): %s\n", username, err.Error())
		err = InvalidUsernameOrPassword
		return
	}

	return
}

func Register(db database.DB, username string) (sessionid string, err error) {
	password, err := db.RegisterUser(username)
	if err != nil {
		return
	}

	err = sendRegEmail(username, password)
	if err != nil {
		return
	}

	return Login(db, username, password)
}

func sendRegEmail(email, password string) error {
	emailText := fmt.Sprintf(`Thank you for registerering with the Basic Go API.

Your username is: %s
Your temporary password is: %s

Please change your password after logging in.
`, email, password)
	return sendEmail(email, "Welcome to the Basic Go API!", "", emailText)
}

func sendResetEmail(email, password string) error {
	emailText := fmt.Sprintf(`We are sorry that you forgot your password to the Basic Go API.

Your username is: %s
Your new temporary password is: %s

Please change your password after logging in.

If you did not request this password reset please contact support.
`, email, password)
	return sendEmail(email, "Basic Go API Password Reset", "", emailText)
}

func sendEmail(recipient, subject, html, text string) error {
	log.Printf("Sending email to: %s\n", recipient)
	log.Printf("Subject: %s\nText:\n%s\n", subject, text)

	client := mandrill.ClientWithKey(mandrillKey)

	message := &mandrill.Message{}
	message.AddRecipient(recipient, recipient, "to")
	message.FromEmail = fromEmail
	message.FromName = fromName
	message.Subject = subject
	message.HTML = html
	message.Text = text

	responses, apiError, err := client.MessagesSend(message)
	if err != nil || apiError != nil {
		if err != nil {
			log.Printf("Error: %s\n", err.Error())
		}
		if apiError != nil {
			log.Printf("Mandrill API Error: %+v\n", apiError)
		}
		return EmailFailed
	}
	log.Printf("Mandrill responses: %+v\n", responses)
	return nil
}

func ValidEmail(email string) bool {
	return emailRegex.Match([]byte(email))
}

func ValidPrice(price string) bool {
	return priceRegex.Match([]byte(price))
}

func ResetPassword(db database.DB, username string) error {
	password, err := db.ResetPassword(username)
	if err != nil {
		return err
	}
	return sendResetEmail(username, password)
}

func ChangePassword(db database.DB, sessionid, oldPassword, newPassword string) error {
	log.Printf("Looking up email with sessionid: %s\n", sessionid)

	// Get email from session
	email, err := db.GetEmail(sessionid)
	if err != nil {
		log.Printf("Error retrieving email from sessionid (%s): %s\n", sessionid, err.Error())
		return InvalidSessionID
	}
	log.Printf("Sessionid: %s is associated with the email: %s\n", sessionid, email)

	// Check old password
	_, err = db.GetUserID(email, oldPassword)
	if err != nil {
		log.Printf("Error validating old password while changing password for user (%s): %s\n", email, err.Error())
		return InvalidPassword
	}

	// Check new password meets requirements
	if len(newPassword) < minPasswordLength {
		log.Printf(
			"New password for user %s of length %d is too short. %d required.\n",
			email,
			len(newPassword),
			minPasswordLength,
		)
		return InvalidPassword
	}

	return db.ChangePassword(email, newPassword)
}

func Logout(db database.DB, sessionid string) error {
	return db.Logout(sessionid)
}

func GetEmail(db database.DB, sessionid string) (email string, err error) {
	// Update the session's last used date.
	// GetEmail is called when the app's main page is displayed.
	err = db.UpdateSession(sessionid)
	if err != nil {
		return
	}

	return db.GetEmail(sessionid)
}
