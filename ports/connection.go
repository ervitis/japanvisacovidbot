package ports

import (
	"context"
	"github.com/ervitis/japanvisacovidbot/model"
)

type (
	IConnection interface {
		Save(context.Context, *model.Embassy) error
		FetchLatestDateFromEmbassy(context.Context, *model.Embassy) (err error)
		SaveCovid(context.Context, *model.JapanCovidData, string) error
		UpdateCovid(context.Context, *model.JapanCovidData) error
		GetCovid(context.Context, *model.JapanCovidData) error
		GetAll(context.Context, *[]model.JapanCovidData, string) error
		Close() error
	}
)
