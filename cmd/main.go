package main

/*
В этой части кода импортируются необходимые пакеты для работы с базой данных PostgreSQL (database/sql),
для работы с HTTP-сервером (net/http) и для драйвера базы данных PostgreSQL (github.com/lib/pq).
*/
import (
	"OPIS/internal/config"
	"OPIS/internal/repository/postgresql"
	"context"
	"fmt"
	"github.com/jackc/pgx/v5"

	"github.com/sirupsen/logrus"
)

func main() {
	ctx := context.Background()
	cfg := config.Read()
	log := logrus.New()

	url := fmt.Sprintf(
		"postgresql://%s:%s@%s:%d/%s",
		cfg.DB.User,
		cfg.DB.Password,
		cfg.DB.Host,
		cfg.DB.Port,
		cfg.DB.Name,
	)

	connector, err := pgx.Connect(ctx, url)
	if err != nil {
		log.Errorf("cat't connect to db: %v", err)
	}

	if err := connector.Ping(ctx); err != nil {
		log.Errorf("cat't ping db: %v", err)
	}

	db := postgresql.NewDB(context.Background(), connector, log)

	fmt.Println(db.SelectData())

}
