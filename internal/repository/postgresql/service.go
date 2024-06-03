package postgresql

import (
	"OPIS/internal/models"
	"context"
	"github.com/sirupsen/logrus"

	"github.com/jackc/pgx/v5"
)

const (
	insertQuery = `INSERT INTO programs (
                      name, content_type, contract_start, contract_end, air_duration, air_frequency, time_types)
						VALUES 
    						($1, $2, $3, $4, $5, $6, $7)`

	selectQuery = `SELECT * FROM programs;`

	createQuery = `CREATE TABLE programs (
                          id SERIAL PRIMARY KEY,
                          name VARCHAR(255) NOT NULL,
                          content_type VARCHAR(100) NOT NULL,
                          contract_start TIMESTAMP NOT NULL,
                          contract_end TIMESTAMP NOT NULL,
                          air_duration INTERVAL NOT NULL,
                          air_frequency VARCHAR(50) NOT NULL,
                          time_types TEXT[] NOT NULL
);`
)

type DB struct {
	ctx  context.Context
	conn *pgx.Conn
	log  *logrus.Logger
}

func NewDB(ctx context.Context, conn *pgx.Conn, log *logrus.Logger) *DB {
	return &DB{ctx, conn, log}

}

func (db *DB) CreateTMPTable() error {
	_, err := db.conn.Exec(db.ctx, createQuery)
	if err != nil {
		db.log.Errorf("cat't create table: %v", err)
		return err
	}

	return nil
}

func (db *DB) InsertData() error {

	_, err := db.conn.Query(db.ctx, insertQuery)
	if err != nil {
		db.log.Errorf("cat't insert data to db: %v", err)
		return err
	}
	return nil
}

func (db *DB) SelectData() ([]models.ProgramsDTO, error) {

	rows, err := db.conn.Query(db.ctx, selectQuery)
	if err != nil {
		db.log.Errorf("cat't select data from db: %v", err)
		return nil, err
	}

	data := make([]models.ProgramsDTO, 0)

	tmp := models.ProgramsDTO{}
	for rows.Next() != false {

		err := rows.Scan(&tmp)

		if err != nil {
			db.log.Errorf("cat't scan data from db: %v", err)
			return nil, err
		}

		data = append(data, tmp)

	}
	return data, nil
}
