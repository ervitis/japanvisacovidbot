package japancovid

import (
	"context"
	"fmt"
	"github.com/ervitis/japanvisacovidbot/bots"
	"github.com/ervitis/japanvisacovidbot/model"
	"github.com/ervitis/japanvisacovidbot/ports"
	"github.com/ervitis/japanvisacovidbot/queue"
	"log"
	"strconv"
	"time"
)

type (
	japanCovidService struct {
		restClient IRestClient
		endpoint   string
		db         ports.IConnection
		bots       []bots.IBot
	}

	IJapanCovidService interface {
		GetLatest(context.Context, *model.JapanCovidResponse) error
		SaveData(context.Context, *model.JapanCovidResponse) (err error)
		GetData(context.Context, *model.JapanCovidResponse) (*model.JapanCovidData, error)
		UpdateData(context.Context, *model.JapanCovidResponse) error
		DateOneDayBefore(*model.JapanCovidResponse) time.Time
		DateToString(time.Time) string
		Transform(*model.JapanCovidResponse, *model.JapanCovidData) error
		CalculateDeltaBetweenDayBeforeAndToday(*queue.Message)
	}
)

const (
	dateLayout = "20060102"
)

func New(db ports.IConnection, bots []bots.IBot) IJapanCovidService {
	rc := NewRestClient()
	rc.R()
	return &japanCovidService{
		db:         db,
		restClient: NewRestClient(),
		endpoint:   `https://covid19-japan-web-api.now.sh/api/v1/total`,
		bots:       bots,
	}
}

func (js *japanCovidService) GetLatest(ctx context.Context, covid *model.JapanCovidResponse) error {
	resp, err := js.restClient.R().SetContext(ctx).SetResult(covid).Get(js.endpoint)
	if err != nil {
		return err
	}

	if resp.IsError() {
		return fmt.Errorf("response error: %d %s: %v", resp.StatusCode(), resp.Status(), resp.Error())
	}
	return nil
}

func (js *japanCovidService) SaveData(ctx context.Context, data *model.JapanCovidResponse) (err error) {
	dbModel := new(model.JapanCovidData)
	if err := js.Transform(data, dbModel); err != nil {
		return err
	}

	if err := js.db.SaveCovid(ctx, dbModel); err != nil {
		return err
	}
	return
}

func (js *japanCovidService) UpdateData(ctx context.Context, data *model.JapanCovidResponse) error {
	dbModel := new(model.JapanCovidData)

	if err := js.Transform(data, dbModel); err != nil {
		return err
	}

	if err := js.db.UpdateCovid(ctx, dbModel); err != nil {
		return err
	}
	return nil
}

func (js *japanCovidService) GetData(ctx context.Context, data *model.JapanCovidResponse) (*model.JapanCovidData, error) {
	dbModel := new(model.JapanCovidData)
	dbModel.Date = strconv.Itoa(data.Date)

	if err := js.db.GetCovid(ctx, dbModel); err != nil {
		return nil, err
	}
	return dbModel, nil
}

func (js *japanCovidService) DateOneDayBefore(data *model.JapanCovidResponse) time.Time {
	t, _ := time.Parse(dateLayout, strconv.Itoa(data.Date))
	return t.AddDate(0, 0, -1)
}

func (js *japanCovidService) DateToString(date time.Time) string {
	return date.Format(dateLayout)
}

func (js *japanCovidService) Transform(input *model.JapanCovidResponse, output *model.JapanCovidData) error {
	if output == nil {
		output = new(model.JapanCovidData)
	}

	output.Date = strconv.Itoa(input.Date)
	output.Pcr = input.Pcr
	output.Positive = input.Positive
	output.Symptom = input.Symptom
	output.Symptomless = input.Symptomless
	output.SymtomConfirming = input.SymtomConfirming
	output.Hospitalize = input.Hospitalize
	output.Mild = input.Mild
	output.Severe = input.Severe
	output.Confirming = input.Confirming
	output.Waiting = input.Waiting
	output.Discharge = input.Discharge
	output.Death = input.Death

	var err error
	output.DateCovid, err = time.Parse(dateLayout, output.Date)
	return err
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

	msg := `New cases:

	death: %d
	severe: %d
	hospitalized: %d
	discharged: %d
	positive: %d
`
	msg = fmt.Sprintf(
		msg,
		dataNow.Death-dataDayBefore.Death,
		dataNow.Severe-dataDayBefore.Severe,
		dataNow.Hospitalize-dataDayBefore.Hospitalize,
		dataNow.Discharge-dataDayBefore.Discharge,
		dataNow.Positive-dataDayBefore.Positive,
	)

	log.Println("Try sending notification to bot", msg)

	for _, bot := range js.bots {
		if err := bot.SendNotification(msg); err != nil {
			log.Println("error trying sending notification with data", err)
			return
		}
	}

}
