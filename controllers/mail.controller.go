package controllers

import (
	"crypto/tls"
	"fmt"
	"log"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/jasonlvhit/gocron"
	"github.com/spf13/viper"

	"gopkg.in/gomail.v2"
)

func (sc ServerController) DailyReport() {
	timeat := viper.GetString("SEND_MAIL_TIME")
	gocron.Every(1).Day().At(timeat).Do(sc.SendEmail)
	<-gocron.Start()
}

func (sc ServerController) SendEmail() {
	mail := viper.GetString("EMAIL_HOST_USER")
	passmail := viper.GetString("EMAIL_HOST_PASSWORD")
	host := viper.GetString("EMAIL_HOST")
	emailPort := viper.GetInt("EMAIL_PORT")

	totalServer, countOn, countOff, avgUptime := sc.CheckStatusServer()
	msg := fmt.Sprintf("Total number of server : %d \nSERVERS ON : %d \nSERVERS OFF : %d \nAverage Uptime of Servers: %fs", totalServer, countOn, countOff, avgUptime)

	m := gomail.NewMessage()
	m.SetHeader("From", mail)
	m.SetHeader("To", "anhntpvcs@gmail.com")
	m.SetHeader("Subject", "Mail Daily Report Automatically")
	m.SetBody("text/plain", msg)

	d := gomail.NewDialer(host, emailPort, mail, passmail)
	d.TLSConfig = &tls.Config{InsecureSkipVerify: true}

	if err := d.DialAndSend(m); err != nil {
		log.Fatal("Failed to Send Mail periodically ", err)
	}
	fmt.Println("Completed to Send Mail periodically")
}

func (sc ServerController) ReportManually(ctx *gin.Context) {
	mail := viper.GetString("EMAIL_HOST_USER")
	passmail := viper.GetString("EMAIL_HOST_PASSWORD")
	host := viper.GetString("EMAIL_HOST")
	emailPort := viper.GetInt("EMAIL_PORT")

	totalServer, countOn, countOff, avgUptime := sc.CheckStatusServer()
	msg := fmt.Sprintf("Total number of server : %d \nSERVERS ON : %d \nSERVERS OFF : %d \nAverage Uptime of Servers: %fs", totalServer, countOn, countOff, avgUptime)

	m := gomail.NewMessage()
	m.SetHeader("From", mail)
	m.SetHeader("To", "anhntpvcs@gmail.com")
	m.SetHeader("Subject", "Mail Report Manually")
	m.SetBody("text/plain", msg)

	d := gomail.NewDialer(host, emailPort, mail, passmail)
	d.TLSConfig = &tls.Config{InsecureSkipVerify: true}

	if err := d.DialAndSend(m); err != nil {
		ctx.JSON(http.StatusBadGateway, gin.H{"status": "error", "message": err.Error()})
		return
	}
	ctx.JSON(http.StatusOK, gin.H{"status": "success"})

}
