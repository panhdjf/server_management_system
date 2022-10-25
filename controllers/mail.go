package controllers

import (
	"crypto/tls"
	"fmt"
	"log"

	"github.com/jasonlvhit/gocron"
	"github.com/panhdjf/server_management_system/models"
	"github.com/spf13/viper"
	"gorm.io/gorm"

	"gopkg.in/gomail.v2"
)

type MailController struct {
	DB *gorm.DB
}

func NewMailController(DB *gorm.DB) MailController {
	return MailController{DB}
}

func (mc MailController) Cron() {

	gocron.Every(1).Day().At("09:05:00").Do(mc.SendEmail)
	<-gocron.Start()
}

func (mc MailController) SendEmail() {
	// initializers.LoadConfig("app.env")
	mail := viper.GetString("EMAIL_HOST_USER")

	// mail := os.Getenv("EMAIL_HOST_USER")
	passmail := viper.GetString("EMAIL_HOST_PASSWORD")

	var servers []models.Server
	mc.DB.Find(&servers)

	countServerOn := 0
	countServerOff := 0

	for _, server := range servers {
		if server.Status == "online" {
			countServerOn++
		} else {
			countServerOff++
		}
	}

	msg := fmt.Sprintf("Total number of server : %d \nSERVERS ON : %d \nSERVERS OFF : %d ", len(servers), countServerOn, countServerOff)

	m := gomail.NewMessage()
	m.SetHeader("From", "anhntpvcs@gmail.com")
	m.SetHeader("To", "anhntpvcs@gmail.com")

	m.SetBody("text/plain", msg)

	d := gomail.NewDialer("smtp.gmail.com", 587, mail, passmail)
	d.TLSConfig = &tls.Config{InsecureSkipVerify: true}

	// time.Sleep(time.Second * 10)
	if err := d.DialAndSend(m); err != nil {
		// ctx.JSON(http.StatusBadGateway, gin.H{"status": "error", "message": err.Error()})
		// return
		log.Fatal("Error", err)
	}
	// ctx.JSON(http.StatusOK, gin.H{"status": "success"})
	fmt.Println("Send Mail ok")
}
