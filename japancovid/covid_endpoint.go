package japancovid

import (
	"context"
	"fmt"
	"github.com/ervitis/japanvisacovidbot/model"
	"github.com/pkg/errors"
	"time"
)

var (
	ErrCouldNotConvert           = errors.New("cannot convert response data")
	ErrCouldNotParseResponseDate = errors.New("cannot parse latest data date")
)

const (
	dateLayoutFormatResponse = "2006-01-02"
)

type (
	responseModel struct {
		Daily []struct {
			Confirmed                       int    `json:"confirmed"`
			ConfirmedCumulative             int    `json:"confirmedCumulative"`
			Deceased                        int    `json:"deceased"`
			DeceasedCumulative              int    `json:"deceasedCumulative"`
			ReportedDeceased                int    `json:"reportedDeceased"`
			ReportedDeceasedCumulative      int    `json:"reportedDeceasedCumulative"`
			Recovered                       int    `json:"recovered"`
			RecoveredCumulative             int    `json:"recoveredCumulative"`
			Critical                        int    `json:"critical"`
			CriticalCumulative              int    `json:"criticalCumulative"`
			Tested                          int    `json:"tested"`
			TestedCumulative                int    `json:"testedCumulative"`
			Active                          int    `json:"active"`
			ActiveCumulative                int    `json:"activeCumulative"`
			CruiseConfirmedCumulative       int    `json:"cruiseConfirmedCumulative"`
			CruiseDeceasedCumulative        int    `json:"cruiseDeceasedCumulative"`
			CruiseRecoveredCumulative       int    `json:"cruiseRecoveredCumulative"`
			CruiseTestedCumulative          int    `json:"cruiseTestedCumulative"`
			CruiseCriticalCumulative        int    `json:"cruiseCriticalCumulative"`
			Date                            string `json:"date"`
			ConfirmedAdjustment             int    `json:"confirmedAdjustment"`
			RecoveredAdjustment             int    `json:"recoveredAdjustment"`
			ConfirmedAvg3D                  int    `json:"confirmedAvg3d"`
			ConfirmedAvg7D                  int    `json:"confirmedAvg7d"`
			ConfirmedCumulativeAvg3D        int    `json:"confirmedCumulativeAvg3d"`
			ConfirmedCumulativeAvg7D        int    `json:"confirmedCumulativeAvg7d"`
			DeceasedAvg3D                   int    `json:"deceasedAvg3d"`
			DeceasedAvg7D                   int    `json:"deceasedAvg7d"`
			DeceasedCumulativeAvg3D         int    `json:"deceasedCumulativeAvg3d"`
			DeceasedCumulativeAvg7D         int    `json:"deceasedCumulativeAvg7d"`
			ReportedDeceasedAvg3D           int    `json:"reportedDeceasedAvg3d"`
			ReportedDeceasedAvg7D           int    `json:"reportedDeceasedAvg7d"`
			ReportedDeceasedCumulativeAvg3D int    `json:"reportedDeceasedCumulativeAvg3d"`
			ReportedDeceasedCumulativeAvg7D int    `json:"reportedDeceasedCumulativeAvg7d"`
			RecoveredAvg3D                  int    `json:"recoveredAvg3d"`
			RecoveredAvg7D                  int    `json:"recoveredAvg7d"`
			RecoveredCumulativeAvg3D        int    `json:"recoveredCumulativeAvg3d"`
			RecoveredCumulativeAvg7D        int    `json:"recoveredCumulativeAvg7d"`
			Deaths                          int    `json:"deaths,omitempty"`
		} `json:"daily"`
		Updated time.Time `json:"updated"`
	}

	covid19japanSummary struct { // TODO use a generic model
		responseModel
		url string
	}

	endpoint struct {
		covid19japanSummary
		client IRestClient
	}

	Endpointer interface {
		GetData(ctx context.Context, data *model.JapanCovidData) error
		TransformIntoModel(resp interface{}, data *model.JapanCovidData) error
	}
)

func newEndpoint() covid19japanSummary {
	return covid19japanSummary{
		url: "https://data.covid19japan.com/summary/latest.json",
	}
}

func NewCovidEndpoint(client IRestClient) Endpointer {
	return &endpoint{
		covid19japanSummary: newEndpoint(),
		client:              client,
	}
}

func (e endpoint) GetData(ctx context.Context, data *model.JapanCovidData) error {
	responseData := e.covid19japanSummary
	resp, err := e.client.R().SetContext(ctx).SetResult(responseData).Get(e.url)
	if err != nil {
		return err
	}

	if resp.IsError() {
		return fmt.Errorf("response error: %d %s: %v", resp.StatusCode(), resp.Status(), resp.Error())
	}
	return e.TransformIntoModel(&responseData, data)
}

func (e endpoint) TransformIntoModel(resp interface{}, data *model.JapanCovidData) error {
	respData, ok := resp.(responseModel)
	if !ok {
		return ErrCouldNotConvert
	}

	latest := respData.Daily[len(respData.Daily)-1]
	dateCovidLatest, err := time.Parse(dateLayoutFormatResponse, latest.Date)
	if err != nil {
		return ErrCouldNotParseResponseDate
	}

	data = &model.JapanCovidData{
		Date:        latest.Date,
		DateCovid:   dateCovidLatest,
		Pcr:         latest.TestedCumulative,
		Hospitalize: latest.CruiseCriticalCumulative,
		Discharge:   latest.RecoveredCumulative,
		Positive:    latest.ConfirmedCumulative,
		Severe:      latest.CriticalCumulative,
		Death:       latest.DeceasedCumulative,
	}
	return nil
}
