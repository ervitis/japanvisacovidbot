package main

import (
	"fmt"
	"github.com/ervitis/japanvisacovidbot"
	"github.com/ervitis/japanvisacovidbot/bots"
	"github.com/ervitis/japanvisacovidbot/bots/telegram"
	"github.com/ervitis/japanvisacovidbot/ports"
	"github.com/ervitis/japanvisacovidbot/repo"
	"log"
	"time"

	"github.com/ervitis/japanvisacovidbot/jacrawler"
)

var (
	tickerCheckEmbassyPages = time.NewTicker(3 * time.Hour)
)

func init() {
	telegram.LoadTelegramConfig()
	repo.LoadDBConfig()
	japanvisacovidbot.LoadGlobalSignalHandler()
}

func main() {
	db := repo.New(&repo.DBConfig)
	covidBots := []bots.IBot{
		telegram.New(&telegram.TelegramConfig),
	}
	server := japanvisacovidbot.NewServer()

	go func(bots []bots.IBot, db ports.IConnection) {
		for {
			select {
			case <-japanvisacovidbot.GlobalSignalHandler.Signals:
				log.Println("cleaning servers...")
				tickerCheckEmbassyPages.Stop()
				for _, bot := range bots {
					bot.Close()
				}
				server.Close()
				return
			case t := <-tickerCheckEmbassyPages.C:
				log.Println("Executed ticker at", t)
				for _, bot := range bots {
					doCrawlerService(bot, db)
				}
			}
		}
	}(covidBots, db)

	for _, bot := range covidBots {
		go func(bot bots.IBot) {
			log.Fatal(bot.StartServer())
		}(bot)
	}
	server.StartServer()
	close(japanvisacovidbot.GlobalSignalHandler.Signals)
}

func doCrawlerService(covidBot bots.IBot, db ports.IConnection) {
	embassies := []jacrawler.IEmbassyData{
		jacrawler.NewJapaneseEmbassy(),
		jacrawler.NewEnglishEmbassy(),
		jacrawler.NewSpanishEmbassy(),
	}

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
		if err := covidBot.SendNotification(msg); err != nil {
			log.Printf("Error sending message to telegram: %s", err)
			continue
		}
	}
}
