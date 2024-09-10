package service

import (
	"bytes"
	"encoding/base64"
	"fmt"
	"io/ioutil"
	"mime/multipart"
	"net/smtp"
	"os"
)

func SendEmail(senderEmail string, imgPath string) error {
	from := os.Getenv("EMAIL")
	password := os.Getenv("EMAIL_PASSWORD")

	// toList is a list of email addresses that the email is to be sent to.
	toList := []string{senderEmail}

	// SMTP server configuration
	host := "smtp.gmail.com"
	port := "587"

	// Read the image file
	imgBytes, err := ioutil.ReadFile(imgPath)
	if err != nil {
		fmt.Println("Error reading image file:", err)
		return err
	}

	// Encode the image in base64
	imgBase64Str := base64.StdEncoding.EncodeToString(imgBytes)

	// Create the HTML email body with the image embedded
	var body bytes.Buffer
	writer := multipart.NewWriter(&body)

	// Set up the MIME headers
	mimeHeaders := fmt.Sprintf("MIME-version: 1.0;\nContent-Type: multipart/related; boundary=%s;\n\n", writer.Boundary())
	body.Write([]byte(mimeHeaders))

	// Write the HTML content with the embedded image
	htmlPart := fmt.Sprintf(`
--%s
Content-Type: text/html; charset="UTF-8"

<html>
  <body>
    <h1>Hello geeks!!!</h1>
    <p>This is an image embedded in the email body:</p>
    <img src="cid:image_id">
  </body>
</html>
`, writer.Boundary())

	body.Write([]byte(htmlPart))

	// Write the image part
	imagePart := fmt.Sprintf(`
--%s
Content-Type: image/jpeg
Content-Transfer-Encoding: base64
Content-ID: <image_id>

%s
--%s--
`, writer.Boundary(), imgBase64Str, writer.Boundary())

	body.Write([]byte(imagePart))

	// Convert body to bytes
	emailBody := body.Bytes()

	// PlainAuth uses the given username and password to authenticate to host and act as identity.
	auth := smtp.PlainAuth("", from, password, host)

	// SendMail uses TLS connection to send the mail
	err = smtp.SendMail(host+":"+port, auth, from, toList, emailBody)
	if err != nil {
		fmt.Println("Error sending email:", err)
		return err
	}

	fmt.Println("Successfully sent mail to", senderEmail)
	return nil
}
