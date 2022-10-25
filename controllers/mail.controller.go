package controllers

import (
	"crypto/tls"
	"fmt"
	"log"

	"github.com/jasonlvhit/gocron"
	"github.com/spf13/viper"

	"gopkg.in/gomail.v2"
)

func (sc ServerController) Periodically() {
	timeat := viper.GetString("SEND_MAIL_TIME")
	gocron.Every(1).Day().At(timeat).Do(sc.SendEmail)
	<-gocron.Start()
}

func (sc ServerController) SendEmail() {
	mail := viper.GetString("EMAIL_HOST_USER")
	passmail := viper.GetString("EMAIL_HOST_PASSWORD")

	totalServer, countOn, countOff := sc.CheckStatus()
	msg := fmt.Sprintf("Total number of server : %d \nSERVERS ON : %d \nSERVERS OFF : %d ", totalServer, countOn, countOff)

	m := gomail.NewMessage()
	m.SetHeader("From", "anhntpvcs@gmail.com")
	m.SetHeader("To", "anhntpvcs@gmail.com")

	m.SetBody("text/plain", msg)

	d := gomail.NewDialer("smtp.gmail.com", 587, mail, passmail)
	d.TLSConfig = &tls.Config{InsecureSkipVerify: true}

	if err := d.DialAndSend(m); err != nil {
		// ctx.JSON(http.StatusBadGateway, gin.H{"status": "error", "message": err.Error()})
		// return
		log.Fatal("Failed to Send Mail periodically ", err)
	}
	// ctx.JSON(http.StatusOK, gin.H{"status": "success"})
	fmt.Println("Completed to Send Mail periodically")
}
