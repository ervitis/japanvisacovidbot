package japancovid

import (
	"context"
	"fmt"
	"github.com/ervitis/japanvisacovidbot/bots"
	"github.com/ervitis/japanvisacovidbot/model"
	"github.com/ervitis/japanvisacovidbot/ports"
	"github.com/ervitis/japanvisacovidbot/queue"
	"log"
	"time"
)

type (
	japanCovidService struct {
		covidEndpoint Endpointer
		db            ports.IConnection
		bots          []bots.IBot
	}

	IJapanCovidService interface {
		GetLatest(context.Context, *model.JapanCovidData) error
		SaveData(context.Context, *model.JapanCovidData) (err error)
		GetDataFromDB(context.Context) (*model.JapanCovidData, error)
		UpdateData(context.Context, *model.JapanCovidData) error
		DateOneDayBefore(*model.JapanCovidData) time.Time
		DateToString(time.Time) string
		CalculateDeltaBetweenDayBeforeAndToday(*queue.Message)
	}
)

const (
	dateLayout        = "20060102"
	dateLayoutMessage = "02 January 2006"
)

func New(db ports.IConnection, bots []bots.IBot, covidEndpoint Endpointer) IJapanCovidService {
	rc := NewRestClient()
	rc.R()
	return &japanCovidService{
		covidEndpoint: covidEndpoint,
		db:            db,
		bots:          bots,
	}
}

func (js *japanCovidService) GetLatest(ctx context.Context, covid *model.JapanCovidData) error {
	return js.covidEndpoint.GetData(ctx, covid)
}

func (js *japanCovidService) SaveData(ctx context.Context, data *model.JapanCovidData) (err error) {
	if err := js.db.SaveCovid(ctx, data, "coviddata"); err != nil {
		return err
	}
	return
}

func (js *japanCovidService) UpdateData(ctx context.Context, data *model.JapanCovidData) error {
	if err := js.db.UpdateCovid(ctx, data); err != nil {
		return err
	}
	return nil
}

func (js *japanCovidService) GetDataFromDB(ctx context.Context) (*model.JapanCovidData, error) {
	dbModel := new(model.JapanCovidData)

	if err := js.db.GetCovid(ctx, dbModel); err != nil {
		return nil, err
	}
	return dbModel, nil
}

func (js *japanCovidService) DateOneDayBefore(data *model.JapanCovidData) time.Time {
	t, _ := time.Parse(dateLayout, data.Date)
	return t.AddDate(0, 0, -1)
}

func (js *japanCovidService) DateToString(date time.Time) string {
	return date.Format(dateLayout)
}

func (js *japanCovidService) CalculateDeltaBetweenDayBeforeAndToday(message *queue.Message) {
	payload := message.Payload.(map[string]interface{})
	dayBefore := payload["dayBefore"].(string)
	dataDayBefore := new(model.JapanCovidData)
	dataDayBefore.Date = dayBefore

	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	if err := js.db.GetCovid(ctx, dataDayBefore); err != nil {
		log.Println("error trying get data of day before", err)
		return
	}

	if dataDayBefore.DateCovid.IsZero() {
		log.Println("there is no data of day before, so the notification won't be send")
		return
	}

	dataNow := payload["dataNow"].(*model.JapanCovidData)

	msg := `New cases on %s:

	💀 death: %d
	🚑 severe: %d
	🏥 hospitalized: %d
	😊 discharged: %d
	😷 positive: %d
`
	diffData := &model.JapanCovidData{
		Date:        dataNow.DateCovid.Format(dateLayout),
		DateCovid:   dataNow.DateCovid,
		Positive:    dataNow.Positive - dataDayBefore.Positive,
		Hospitalize: dataNow.Hospitalize - dataDayBefore.Hospitalize,
		Severe:      dataNow.Severe - dataDayBefore.Severe,
		Discharge:   dataNow.Discharge - dataDayBefore.Discharge,
		Death:       dataNow.Death - dataDayBefore.Death,
	}

	if err := js.db.SaveCovid(ctx, diffData, "diffdatacovid"); err != nil {
		log.Println("error saving diff data into diffdatacovid table", err)
	}

	msg = fmt.Sprintf(
		msg,
		diffData.DateCovid.Format(dateLayoutMessage),
		diffData.Death,
		diffData.Severe,
		diffData.Hospitalize,
		diffData.Discharge,
		diffData.Positive,
	)

	log.Println("Try sending notification to bot", msg)

	for _, bot := range js.bots {
		if err := bot.SendNotification(msg); err != nil {
			log.Println("error trying sending notification with data", err)
			return
		}
	}

}
