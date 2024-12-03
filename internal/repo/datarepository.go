package repo

import (
	"context"
	"database/sql"
	"errors"
	"fmt"
	"log"
	"time"

	_ "github.com/lib/pq"

	sq "github.com/Masterminds/squirrel"
	"github.com/georgysavva/scany/pgxscan"
	"github.com/jackc/pgx/v4"

	"github.com/s-buhar0v/demoapp/internal/metrics"
)

const (
	dataTable      = "data"
	testDataColumn = "test_data"
	delayProcedure = "delay()"
)

func (p *Postgres) GetData(ctx context.Context, long bool) ([]string, error) {
	now := time.Now()

	if long {
		// Вызываем процедуру для нагрузки
		err := p.RunStoredProcedure()
		elapsedTime := time.Since(now).Seconds()
		metrics.DatabaseQuesriesTotal.WithLabelValues("CALL", delayProcedure).Inc()
		metrics.DbQueryDuration.WithLabelValues("CALL", delayProcedure).Observe(elapsedTime)
		if err != nil {
			return nil, err
		}
	}

	query := sq.Select(testDataColumn).From(dataTable)

	sqlQuery, args, err := query.ToSql()
	if err != nil {
		log.Println("Failed to create query for operation Get test_data from table data")
		return nil, err
	}

	var data []string
	err = pgxscan.Select(ctx, p.db, &data, sqlQuery, args...)

	elapsedTime := time.Since(now).Seconds()
	metrics.DatabaseQuesriesTotal.WithLabelValues("SELECT", dataTable).Inc()
	metrics.DbQueryDuration.WithLabelValues("SELECT", dataTable).Observe(elapsedTime)

	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return nil, fmt.Errorf("%s: %w", "data record not found", err)
		}

		log.Println("Failed to Get test_data from table data")
		return nil, err
	}

	return data, nil
}

func (p *Postgres) SaveData(ctx context.Context, newData []string) error {
	query := sq.Insert(dataTable).
		Columns(testDataColumn)

	for _, val := range newData {
		query = query.Values(val)
	}

	query = query.PlaceholderFormat(sq.Dollar)

	sqlQuery, args, err := query.ToSql()
	if err != nil {
		log.Println("Failed to create query for operation Save to table data")
		return err
	}

	result, err := p.db.Exec(ctx, sqlQuery, args...)
	if err != nil {
		log.Println("Failed to Insert test data rows to table data")
		return err
	}

	if result.RowsAffected() == 0 {
		return fmt.Errorf("new data was not inserted")
	}

	return err
}

func (p *Postgres) RunStoredProcedure() error {
	// Настройки подключения к базе данных
	connString := fmt.Sprintf("host=%s port=%s user=%s password=%s dbname=%s sslmode=disable",
		"postgres-db", "5432", "postgres", "postgres", "postgres")
	db, err := sql.Open("postgres", connString)
	if err != nil {
		log.Fatal(err)
	}
	defer db.Close()

	_, err = db.Exec(fmt.Sprintf("%s %s", "CALL", delayProcedure))
	if err != nil {
		return err
	}
	return nil
}
