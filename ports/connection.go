package ports

import "github.com/ervitis/japanvisacovidbot/model"

type (
	IConnection interface {
		Save(*model.Embassy) error
		FetchLatestDateFromEmbassy(*model.Embassy) (err error)
	}
)
