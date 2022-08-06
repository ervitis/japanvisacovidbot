package main

import (
	"context"
	"github.com/ervitis/japanvisacovidbot/japancovid"
	"github.com/ervitis/japanvisacovidbot/model"
	"github.com/ervitis/japanvisacovidbot/queue"
	"github.com/ervitis/japanvisacovidbot/repo"
	"log"
	"time"
)

func init() {
	repo.LoadDBConfig()
}

func main() {
	db := repo.New(&repo.DBConfig)

	dataCovid := japancovid.New(db, nil, japancovid.NewCovidEndpoint(japancovid.NewRestClient()))

	all := make([]model.JapanCovidData, 0)
	allDiff := make([]model.JapanCovidData, 0)

	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	if err := db.GetAll(ctx, &all, "coviddata"); err != nil {
		panic(err)
	}

	if err := db.GetAll(ctx, &allDiff, "diffdatacovid"); err != nil {
		panic(err)
	}

	if len(all) == len(allDiff) {
		log.Println("nothing to do")
		return
	}

	firstTime, err := time.Parse("20060102", "20220106")
	if err != nil {
		panic(err)
	}

	lastTime := allDiff[len(allDiff)-1].DateCovid

	for _, data := range all {
		if data.DateCovid.Equal(firstTime) || lastTime.Before(data.DateCovid) {
			continue
		}

		dr := &model.JapanCovidData{
			Date:             data.Date,
			Pcr:              data.Pcr,
			Positive:         data.Positive,
			Symptom:          data.Symptom,
			Symptomless:      data.Symptomless,
			SymtomConfirming: data.SymtomConfirming,
			Hospitalize:      data.Hospitalize,
			Mild:             data.Mild,
			Severe:           data.Severe,
			Confirming:       data.Confirming,
			Waiting:          data.Waiting,
			Discharge:        data.Discharge,
			Death:            data.Death,
		}
		dayBefore := dataCovid.DateOneDayBefore(dr)

		message := &queue.Message{
			Payload: map[string]interface{}{
				"dayBefore": dayBefore.Format("20060102"),
				"dataNow":   &data,
			},
		}
		dataCovid.CalculateDeltaBetweenDayBeforeAndToday(message)
	}

	log.Println("finished proccess")
}
