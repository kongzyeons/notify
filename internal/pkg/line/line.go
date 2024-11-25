package line

import (
	"fmt"
	"go_notify/internal/config"
	"log"
	"time"

	"github.com/juunini/simple-go-line-notify/notify"
)

type LineAPI interface {
	SendMessageJob()
	SendMessage(message string) error
}

type lineAPI struct {
	AccToken string
	loc      *time.Location
}

func NewLineAPI() LineAPI {
	cfg := config.InitConfig()

	loc, err := time.LoadLocation("Asia/Bangkok")
	if err != nil {
		log.Println(err)
		loc = time.UTC
	}

	return &lineAPI{
		AccToken: cfg.LineToken,
		loc:      loc,
	}
}

func (self *lineAPI) SendMessageJob() {
	log.Println("run start..")
	timeNow := time.Now().In(self.loc)
	layout := "2006-01-02 15:04:05"
	currentTimeString := timeNow.Format(layout)

	message := fmt.Sprintf("%s : hello", currentTimeString)

	if err := notify.SendText(self.AccToken, message); err != nil {
		log.Println(err)
		panic(err)
	}

	log.Println("send success")

}

func (self *lineAPI) SendMessage(messageJson string) error {
	timeNow := time.Now().In(self.loc)
	layout := "2006-01-02 15:04:05"
	currentTimeString := timeNow.Format(layout)

	message := fmt.Sprintf("%s : %s", currentTimeString, messageJson)

	if err := notify.SendText(self.AccToken, message); err != nil {
		return err
	}
	return nil
}
