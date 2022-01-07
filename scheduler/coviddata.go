package scheduler

import (
	"context"
	"github.com/ervitis/japanvisacovidbot/japancovid"
	"github.com/ervitis/japanvisacovidbot/model"
	"github.com/ervitis/japanvisacovidbot/ports"
	"log"
)

func CovidDataFn(db ports.IConnection) CovidJob {
	dataCovid := japancovid.New(db)

	return CovidJob{
		Cron: "0 */1 * * *",
		Task: func() error {
			log.Println("executing task covidJob")
			ctx, cancel := context.WithCancel(context.Background())
			defer cancel()

			data := new(model.JapanCovidResponse)
			if err := dataCovid.GetLatest(ctx, data); err != nil {
				return err
			}

			log.Printf("data from api: %#v", data)

			covidData, err := dataCovid.GetData(ctx, data)
			if err != nil {
				return err
			}

			log.Printf("data from db: %#v", covidData)

			if covidData.Date == "" {
				log.Println("saving new data")
				if err := dataCovid.SaveData(ctx, data); err != nil {
					return err
				}
				return nil
			}

			log.Println("updating data")
			if err := dataCovid.UpdateData(ctx, data); err != nil {
				return err
			}
			return nil
		},
		TaskParams: nil,
	}
}
