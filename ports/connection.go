package ports

import (
	"context"
	"github.com/ervitis/japanvisacovidbot/model"
)

type (
	IConnection interface {
		Save(*model.Embassy) error
		FetchLatestDateFromEmbassy(*model.Embassy) (err error)
		SaveCovid(context.Context, *model.JapanCovidData) error
		UpdateCovid(context.Context, *model.JapanCovidData) error
		GetCovid(context.Context, *model.JapanCovidData) error
	}
)
