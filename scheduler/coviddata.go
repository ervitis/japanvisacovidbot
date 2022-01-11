package scheduler

import (
	"context"
	"github.com/ervitis/japanvisacovidbot/bots"
	"github.com/ervitis/japanvisacovidbot/japancovid"
	"github.com/ervitis/japanvisacovidbot/model"
	"github.com/ervitis/japanvisacovidbot/ports"
	"github.com/ervitis/japanvisacovidbot/queue"
	"log"
)

func CovidDataFn(db ports.IConnection, bot []bots.IBot) CovidJob {
	dataCovid := japancovid.New(db, bot)

	return CovidJob{
		Cron: "*/2 * * * *",
		Task: func() error {
			log.Println("executing task covidJob")
			ctx, cancel := context.WithCancel(context.Background())
			defer cancel()

			data := new(model.JapanCovidResponse)
			if err := dataCovid.GetLatest(ctx, data); err != nil {
				return err
			}

			covidData, err := dataCovid.GetData(ctx, data)
			if err != nil {
				return err
			}

			if covidData.DateCovid.IsZero() {
				log.Println("saving new data of", data.Date)
				if err := dataCovid.SaveData(ctx, data); err != nil {
					return err
				}

				// send event
				{
					dayBefore := dataCovid.DateOneDayBefore(data)
					t := new(model.JapanCovidData)
					_ = dataCovid.Transform(data, t)
					payload := map[string]interface{}{
						"dayBefore": dataCovid.DateToString(dayBefore),
						"dataNow":   t,
					}

					queue.Queue.Publish(queue.NewCovidEntryEvent, payload)
				}

				return nil
			}

			log.Println("updating data of", data.Date)
			if err := dataCovid.UpdateData(ctx, data); err != nil {
				return err
			}
			return nil
		},
		TaskParams: nil,
	}
}
