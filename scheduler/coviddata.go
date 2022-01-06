package scheduler

import (
	"context"
	"github.com/ervitis/japanvisacovidbot/japancovid"
	"github.com/ervitis/japanvisacovidbot/model"
	"github.com/ervitis/japanvisacovidbot/ports"
)

func CovidDataFn(db ports.IConnection) CovidJob {
	dataCovid := japancovid.New(db)

	return CovidJob{
		Cron: "0 0/2 * * *",
		Task: func() error {
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

			if covidData.Date == "" {
				if err := dataCovid.SaveData(ctx, data); err != nil {
					return err
				}
				return nil
			}

			if err := dataCovid.UpdateData(ctx, data); err != nil {
				return err
			}
			return nil
		},
		TaskParams: nil,
	}
}
