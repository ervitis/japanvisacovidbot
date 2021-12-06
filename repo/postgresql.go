package repo

import (
	"context"
	"database/sql"
	"fmt"
	"github.com/ervitis/japanvisacovidbot/model"
	"github.com/ervitis/japanvisacovidbot/ports"
	_ "github.com/lib/pq"
	"time"
)

type (
	postgresql struct {
		conn *sql.DB
	}
)

func New(parameters *DBConfigParameters) ports.IConnection {
	psqlconn := fmt.Sprintf("host=%s port=%d user=%s password=%s dbname=%s sslmode=disable",
		parameters.Host, parameters.Port, parameters.User, parameters.Password, parameters.DB)

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

func (p *postgresql) Save(data *model.Embassy) error {
	query := `INSERT INTO embassydates (date, embassy) VALUES ($1, $2)`
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

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

func (p *postgresql) FetchLatestDateFromEmbassy(data *model.Embassy) (err error) {
	query := `SELECT date, embassy FROM embassydates WHERE embassy=$1 ORDER BY date DESC LIMIT 1`
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

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
