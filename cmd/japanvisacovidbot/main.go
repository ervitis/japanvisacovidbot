package main

import (
	"context"
	"fmt"
	"github.com/ervitis/japanvisacovidbot"
	"github.com/ervitis/japanvisacovidbot/bots"
	"github.com/ervitis/japanvisacovidbot/bots/telegram"
	"github.com/ervitis/japanvisacovidbot/email"
	"github.com/ervitis/japanvisacovidbot/japancovid"
	"github.com/ervitis/japanvisacovidbot/metrics"
	"github.com/ervitis/japanvisacovidbot/queue"
	"github.com/ervitis/japanvisacovidbot/repo"
	"github.com/ervitis/japanvisacovidbot/scheduler"
	"log"
	"time"

	"github.com/ervitis/japanvisacovidbot/jacrawler"
)

var (
	tickerCheckEmbassyPages = time.NewTicker(1 * time.Hour)
)

func init() {
	telegram.LoadTelegramConfig()
	repo.LoadDBConfig()
	japanvisacovidbot.LoadGlobalSignalHandler()
	email.LoadConfig()

	createTopics()
}

func createTopics() {
	for _, topic := range queue.AllTopics() {
		if err := queue.Queue.CreateTopic(topic); err != nil {
			panic(err)
		}
	}
}

func main() {
	db := repo.New(&repo.DBConfig)
	covidBots := []bots.IBot{
		telegram.New(&telegram.Config),
	}
	server := japanvisacovidbot.NewServer()

	cron := scheduler.New()

	appMetrics := metrics.New()

	if err := cron.ExecuteJob([]scheduler.CovidJob{
		scheduler.CovidDataFn(db, covidBots, appMetrics),
	}...); err != nil {
		log.Fatal("error executing job", err)
	}

	dataCovid := japancovid.New(db, covidBots)

	queue.Queue.Subscribe(queue.NewCovidEntryEvent, dataCovid.CalculateDeltaBetweenDayBeforeAndToday)

	embassies := []jacrawler.IEmbassyData{
		jacrawler.NewJapaneseEmbassy(db),
		jacrawler.NewEnglishEmbassy(db),
		jacrawler.NewSpanishEmbassy(db),
	}

	go func(bots []bots.IBot) {
		for {
			select {
			case <-japanvisacovidbot.GlobalSignalHandler.Signals:
				log.Println("cleaning servers...")
				tickerCheckEmbassyPages.Stop()
				cron.Stop()
				appMetrics.Stop()
				for _, bot := range bots {
					bot.Close()
				}
				server.Close()
				log.Println(db.Close())
				return
			case t := <-tickerCheckEmbassyPages.C:
				log.Println("Executed ticker at", t)
				for _, bot := range bots {
					doCrawlerService(bot, embassies, appMetrics)
				}
			}
		}
	}(covidBots)

	if err := appMetrics.Start(); err != nil {
		log.Println("Metrics cannot start", err)
	}

	for _, bot := range covidBots {
		go func(bot bots.IBot) {
			log.Fatal(bot.StartServer())
		}(bot)
	}
	server.StartServer()
	close(japanvisacovidbot.GlobalSignalHandler.Signals)
}

func doCrawlerService(covidBot bots.IBot, embassies []jacrawler.IEmbassyData, appMetrics metrics.IMetrics) {
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	if err := appMetrics.ExecuteWithSegment(ctx, "scrapWebEmbassies", func(ctx context.Context) error {
		for _, embassy := range embassies {
			crawler := jacrawler.NewCovidCrawler(embassy)
			data, err := crawler.CrawlPage()
			if err != nil {
				log.Printf("Error crawling data: %s", err)
				continue
			}

			if embassy.IsDateUpdated(ctx, data) {
				continue
			}

			if err := embassy.UpdateDate(ctx, data); err != nil {
				log.Printf("Error updating data: %s", err)
				continue
			}

			msg := fmt.Sprintf("There is an update in the embassy of %s, go to the web %s", embassy.GetISO(), embassy.GetUri())
			if err := covidBot.SendNotification(msg); err != nil {
				log.Printf("Error sending message to telegram: %s", err)
				continue
			}
		}
		return nil
	}); err != nil {
		log.Println("error scrapping web pages", err)
	}

}
