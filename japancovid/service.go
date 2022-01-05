package japancovid

import (
	"context"
	"fmt"
	"github.com/ervitis/japanvisacovidbot/model"
)

type (
	japanCovidService struct {
		restClient IRestClient
		endpoint   string
	}

	IJapanCovidService interface {
		GetLatest(context.Context, *model.JapanCovidResponse) error
	}
)

func New() IJapanCovidService {
	rc := NewRestClient()
	rc.R()
	return &japanCovidService{
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
