// package main

// import (
// 	"fmt"
// 	"net/smtp"

// 	"github.com/gin-gonic/gin"
// 	"github.com/panhdjf/server_management_system/models"
// )

// func main() {
// 	var ctx *gin.Context
// 	currentUser := ctx.MustGet("currentUser").(models.User)

// 	// Sender data.
// 	from := currentUser.Email
// 	password := currentUser.Password

// 	// Receiver email address.
// 	to := []string{
// 		"phuonganh080701a3@gmail.com",
// 	}

// 	// smtp server configuration.
// 	smtpHost := "localhost"
// 	smtpPort := "8000"

// 	// Message.
// 	message := []byte("This is a test email message.")

// 	// Authentication.
// 	auth := smtp.PlainAuth("", from, password, smtpHost)

// 	// Sending email.
// 	err := smtp.SendMail(smtpHost+":"+smtpPort, auth, from, to, message)
// 	if err != nil {
// 		fmt.Println(err)
// 		return
// 	}
// 	fmt.Println("Email Sent Successfully!")
// }
