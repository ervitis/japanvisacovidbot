package scheduler

import (
	"context"
	"github.com/ervitis/japanvisacovidbot/japancovid"
	"github.com/ervitis/japanvisacovidbot/metrics"
	"github.com/ervitis/japanvisacovidbot/model"
	"github.com/ervitis/japanvisacovidbot/queue"
	"log"
)

func CovidDataFn(appMetrics metrics.IMetrics, covidService japancovid.IJapanCovidService) CovidJob {

	return CovidJob{
		Cron: "0 */2 * * *",
		Task: func() error {
			log.Println("executing task covidJob")
			ctx, cancel := context.WithCancel(context.Background())
			defer cancel()

			return appMetrics.ExecuteWithSegment(ctx, "covidJobTask", func(ctx context.Context) error {
				data := new(model.JapanCovidData)
				if err := covidService.GetLatest(ctx, data); err != nil {
					return err
				}

				dbData, err := covidService.GetDataByDateFromDB(ctx, data.Date)
				if err != nil {
					return err
				}

				if dbData.DateCovid.IsZero() {
					log.Println("saving new data of", data.Date)
					if err := covidService.SaveData(ctx, data); err != nil {
						return err
					}

					// send event
					{
						dayBefore := covidService.DateOneDayBefore(data)
						payload := map[string]interface{}{
							"dayBefore": covidService.DateToString(dayBefore),
							"dataNow":   data,
						}

						queue.Queue.Publish(queue.NewCovidEntryEvent, payload)
					}

					return nil
				}

				log.Println("updating data of", data.Date)
				if err := covidService.UpdateData(ctx, data); err != nil {
					return err
				}
				return nil
			})
		},
		TaskParams: nil,
	}
}
