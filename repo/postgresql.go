package repo

import (
	"context"
	"database/sql"
	"fmt"
	"github.com/ervitis/japanvisacovidbot/model"
	"github.com/ervitis/japanvisacovidbot/ports"
	_ "github.com/lib/pq"
	"strings"
	"time"
)

type (
	postgresql struct {
		conn *sql.DB
	}
)

func New(parameters *DBConfigParameters) ports.IConnection {
	psqlconn := fmt.Sprintf("host=%s port=%d user=%s password=%s dbname=%s %s",
		parameters.Host, parameters.Port, parameters.User, parameters.Password, parameters.DB, parameters.Options)
	psqlconn = strings.TrimSpace(psqlconn)

	conn, err := sql.Open("postgres", psqlconn)
	if err != nil {
		panic(err)
	}

	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	if err := conn.PingContext(ctx); err != nil {
		panic(err)
	}

	return &postgresql{
		conn: conn,
	}
}

func (p *postgresql) Save(ctx context.Context, data *model.Embassy) error {
	query := `INSERT INTO embassydates (date, embassy) VALUES ($1, $2)`

	stmt, err := p.conn.PrepareContext(ctx, query)
	if err != nil {
		return err
	}

	_, err = stmt.ExecContext(ctx, data.UpdatedDate, data.EmbassyISO)
	if err != nil {
		return err
	}

	return nil
}

func (p *postgresql) FetchLatestDateFromEmbassy(ctx context.Context, data *model.Embassy) (err error) {
	query := `SELECT date, embassy FROM embassydates WHERE embassy=$1 ORDER BY date DESC LIMIT 1`

	stmt, err := p.conn.PrepareContext(ctx, query)
	if err != nil {
		return err
	}

	r, err := stmt.QueryContext(ctx, data.EmbassyISO)
	if err != nil {
		return err
	}

	defer func() {
		if err := r.Close(); err != nil {
			return
		}
	}()

	for r.Next() {
		err := r.Scan(&data.UpdatedDate, &data.EmbassyISO)
		if err != nil {
			return err
		}
	}
	return nil
}

func (p *postgresql) SaveCovid(ctx context.Context, data *model.JapanCovidData, tableName string) error {
	query := fmt.Sprintf(`INSERT INTO %s (datecovid, date, pcr, positive, symptom, symptomless, symtomConfirming, hospitalize, mild, severe, confirming, waiting, discharge, death) VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9, $10, $11, $12, $13, $14)`, tableName)

	stmt, err := p.conn.PrepareContext(ctx, query)
	if err != nil {
		return err
	}

	defer func() {
		_ = stmt.Close()
	}()

	_, err = stmt.ExecContext(ctx, data.DateCovid, data.Date, data.Pcr, data.Positive, data.Symptom, data.Symptomless, data.SymtomConfirming, data.Hospitalize, data.Mild, data.Severe, data.Confirming, data.Waiting, data.Discharge, data.Death)
	if err != nil {
		return err
	}
	return nil
}

func (p *postgresql) GetCovid(ctx context.Context, data *model.JapanCovidData) error {
	query := `SELECT datecovid, date, pcr, positive, symptom, symptomless, symtomConfirming, hospitalize, mild, severe, confirming, waiting, discharge, death FROM coviddata WHERE date = $1`

	stmt, err := p.conn.PrepareContext(ctx, query)
	if err != nil {
		return err
	}

	r, err := stmt.QueryContext(ctx, data.Date)
	if err != nil {
		return err
	}

	defer func() {
		if err := r.Close(); err != nil {
			return
		}
	}()

	for r.Next() {
		err := r.Scan(&data.DateCovid, &data.Date, &data.Pcr, &data.Positive, &data.Symptom, &data.Symptomless, &data.SymtomConfirming,
			&data.Hospitalize, &data.Mild, &data.Severe, &data.Confirming, &data.Waiting, &data.Discharge, &data.Death)
		if err != nil {
			return err
		}
	}

	return nil
}

func (p *postgresql) UpdateCovid(ctx context.Context, data *model.JapanCovidData) error {
	query := `UPDATE coviddata SET pcr = $2, positive = $3, symptom = $4, symptomless = $5, symtomConfirming = $6, hospitalize = $7, mild = $8, severe = $9, confirming = $10, waiting = $11, discharge = $12, death = $13 WHERE date = $1`

	stmt, err := p.conn.PrepareContext(ctx, query)
	if err != nil {
		return err
	}

	_, err = stmt.ExecContext(ctx, data.Date, data.Pcr, data.Positive, data.Symptom, data.Symptomless, data.SymtomConfirming,
		data.Hospitalize, data.Mild, data.Severe, data.Confirming, data.Waiting, data.Discharge, data.Death)
	if err != nil {
		return err
	}

	return nil
}

func (p *postgresql) Close() error {
	return p.conn.Close()
}
