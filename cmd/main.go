package main

import (
	"fmt"
	"github.com/ervitis/japanvisacovidbot/repo"
	"log"
	"time"

	"github.com/ervitis/japanvisacovidbot/jacrawler"
	"github.com/kelseyhightower/envconfig"
	tb "gopkg.in/tucnak/telebot.v2"
)

type (
	ApiSecrets struct {
		Token string `envconfig:"TOKEN"`
	}

	TelegramUser struct {
		ID int `envconfig:"TELEGRAM_USERID"`
	}
)

var (
	ApiSecretParameters ApiSecrets
	TelegUser           TelegramUser
)

func init() {
	envconfig.MustProcess("", &ApiSecretParameters)
	envconfig.MustProcess("", &TelegUser)

	repo.LoadDBConfig()
}

func main() {
	db := repo.New(&repo.DBConfig)

	if TelegUser.ID == 0 {
		panic("Must set telegram user ID")
	}

	embassies := []jacrawler.IEmbassyData{
		jacrawler.NewJapaneseEmbassy(),
		jacrawler.NewEnglishEmbassy(),
	}

	covidBot, err := tb.NewBot(tb.Settings{Token: ApiSecretParameters.Token, Poller: &tb.LongPoller{Timeout: 10 * time.Second}})
	if err != nil {
		panic(err)
	}

	user := &tb.User{ID: TelegUser.ID}

	// make a tick to execute this or cron every 2 hours
	for _, embassy := range embassies {
		crawler := jacrawler.NewCovidCrawler(embassy)
		data, err := crawler.CrawlPage()
		if err != nil {
			log.Printf("Error crawling data: %s", err)
			continue
		}

		if embassy.IsDateUpdated(data, db) {
			continue
		}

		if err := embassy.UpdateDate(data, db); err != nil {
			log.Printf("Error updating data: %s", err)
			continue
		}

		msg := fmt.Sprintf("There is an update in the embassy of %s, go to the web %s", embassy.GetISO(), embassy.GetUri())
		if _, err := covidBot.Send(user, msg); err != nil {
			log.Printf("Error sending message to telegram: %s", err)
			continue
		}
	}
}
