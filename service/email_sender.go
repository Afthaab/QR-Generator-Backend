package service

import (
	"encoding/base64"
	"fmt"
	"io/ioutil"
	"net/smtp"
	"os"
	"qrgen/service/model"
)

func SendEmail(imgPath string, studentData model.Student) error {
	from := os.Getenv("EMAIL")
	password := os.Getenv("EMAIL_PASSWORD")

	// toList is a list of email addresses that the email is to be sent to.
	toList := []string{studentData.Email}

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

	// Create the email headers
	subject := "Subject: Congratulations! Your Registration on Student Attendance Portal is Complete\n"
	mimeHeaders := "MIME-version: 1.0;\nContent-Type: multipart/related; boundary=\"boundary\"\n\n"

	// Create the HTML email body with the image embedded
	body := fmt.Sprintf(`
--boundary
Content-Type: text/html; charset="UTF-8"

<html>
  <body>
    <h1>Hello, %s!</h1>
    <p>Welcome to the <strong>Student Attendance Portal</strong>. Your registration was successful, and your student portal account is now active!</p>
    <p>Here’s what you can do with your new account:</p>
    <ul>
      <li>View your personal details</li>
      <li>Track your attendance</li>
    </ul>
    <p>To get started, please <a href="[http://localhost:8000/sign_in.html]" target="_blank">log in to the student portal</a> with your credentials provided below.</p>
	<ul>
      <li>Email: %s</li>
      <li>Password: %s</li>
    </ul>
	<p>Please <strong>Do Not Share</strong> your password and email with anyone. If you have any questions or need assistance, feel free to reach out to our support team</p>

	<p>Below there is also a <strong>QR CODE</strong> attached to the email, Please download the image to register your attendance in the college. The physical copy will be tagged along with your ID Card</p>
	
    <p>We’re excited to support your academic journey.</p>
    <br>
    <p>Best regards,</p>
    <p><strong>College Administration</strong></p>
    <hr>
    <p><em>This email was generated automatically. Please do not reply to this email.</em></p>
    <img src="cid:image_id" alt="Welcome Image">
  </body>
</html>

--boundary
Content-Type: image/jpeg
Content-Transfer-Encoding: base64
Content-ID: <image_id>

%s
--boundary--`, studentData.Name, studentData.Email, studentData.Password, imgBase64Str)

	// Combine headers and body into a single message
	msg := []byte(subject + mimeHeaders + body)

	// Authenticate with the SMTP server
	auth := smtp.PlainAuth("", from, password, host)

	// Send the email
	err = smtp.SendMail(host+":"+port, auth, from, toList, msg)
	if err != nil {
		fmt.Println("Error sending email:", err)
		return err
	}

	fmt.Println("Successfully sent mail to", studentData.Email)
	return nil
}
