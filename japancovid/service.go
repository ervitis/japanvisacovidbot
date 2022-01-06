package japancovid

import (
	"context"
	"fmt"
	"github.com/ervitis/japanvisacovidbot/model"
	"github.com/ervitis/japanvisacovidbot/ports"
	"strconv"
	"time"
)

type (
	japanCovidService struct {
		restClient IRestClient
		endpoint   string
		db         ports.IConnection
	}

	IJapanCovidService interface {
		GetLatest(context.Context, *model.JapanCovidResponse) error
		SaveData(context.Context, *model.JapanCovidResponse) (err error)
		GetData(context.Context, *model.JapanCovidResponse) (*model.JapanCovidData, error)
		UpdateData(context.Context, *model.JapanCovidResponse) error
	}
)

const (
	dateLayout = "20060102"
)

func New(db ports.IConnection) IJapanCovidService {
	rc := NewRestClient()
	rc.R()
	return &japanCovidService{
		db:         db,
		restClient: NewRestClient(),
		endpoint:   `https://covid19-japan-web-api.now.sh/api/v1/total`,
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
	if err := js.transform(data, dbModel); err != nil {
		return err
	}

	if err := js.db.SaveCovid(ctx, dbModel); err != nil {
		return err
	}
	return
}

func (js *japanCovidService) UpdateData(ctx context.Context, data *model.JapanCovidResponse) error {
	dbModel := new(model.JapanCovidData)

	if err := js.transform(data, dbModel); err != nil {
		return err
	}

	if err := js.db.UpdateCovid(ctx, dbModel); err != nil {
		return err
	}
	return nil
}

func (js *japanCovidService) GetData(ctx context.Context, data *model.JapanCovidResponse) (*model.JapanCovidData, error) {
	dbModel := new(model.JapanCovidData)
	if err := js.transform(data, dbModel); err != nil {
		return nil, err
	}

	if err := js.db.GetCovid(ctx, dbModel); err != nil {
		return nil, err
	}
	return dbModel, nil
}

func (js *japanCovidService) transform(input *model.JapanCovidResponse, output *model.JapanCovidData) error {

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
