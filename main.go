package main
/*
В этой части кода импортируются необходимые пакеты для работы с базой данных PostgreSQL (database/sql), 
для работы с HTTP-сервером (net/http) и для драйвера базы данных PostgreSQL (github.com/lib/pq).
*/
import (
	"database/sql"
	"fmt"
	"net/http"

	_ "github.com/lib/pq"
)
/*
Эти константы содержат информацию о подключении к базе данных PostgreSQL: хост (host), 
порт (port), имя пользователя (user), пароль (password) и имя базы данных (dbname).
*/
const (
	host     = "localhost"
	port     = 5432
	user     = "admin"
	password = "admin"
	dbname   = "programs_db"
)
/*
Эта функция открывает соединение с базой данных PostgreSQL и выполняет проверку подключения (Ping). 
Если подключение успешно, функция возвращает объект базы данных (*sql.DB).
*/
func setupDB() (*sql.DB, error) {
	psqlInfo := fmt.Sprintf("host=%s port=%d user=%s password=%s dbname=%s sslmode=disable",
		host, port, user, password, dbname)

	db, err := sql.Open("postgres", psqlInfo)
	if err != nil {
		return nil, err
	}

	err = db.Ping()
	if err != nil {
		return nil, err
	}

	fmt.Println("Successfully connected to database!")
	return db, nil
}
/*
Эта функция создает таблицы в базе данных. Она выполняет SQL-запросы для создания таблиц программ (programs), 
типов программ (program_content_type, program_time_type и т.д.) и связующей таблицы (program_time).
*/
func createTables(db *sql.DB) {
	queries := []string{
		`CREATE TABLE IF NOT EXISTS program_content_type (
            id SERIAL PRIMARY KEY,
            name VARCHAR(50) UNIQUE NOT NULL
        );`,
		`CREATE TABLE IF NOT EXISTS program_time_type (
            id SERIAL PRIMARY KEY,
            name VARCHAR(50) UNIQUE NOT NULL
        );`,
		`CREATE TABLE IF NOT EXISTS contract_duration_unit (
            id SERIAL PRIMARY KEY,
            name VARCHAR(50) UNIQUE NOT NULL
        );`,
		`CREATE TABLE IF NOT EXISTS airtime_duration_unit (
            id SERIAL PRIMARY KEY,
            name VARCHAR(50) UNIQUE NOT NULL
        );`,
		`CREATE TABLE IF NOT EXISTS periodicity_unit (
            id SERIAL PRIMARY KEY,
            name VARCHAR(50) UNIQUE NOT NULL
        );`,
		`CREATE TABLE IF NOT EXISTS programs (
            id SERIAL PRIMARY KEY,
            name VARCHAR(255) NOT NULL,
            content_type_id INT REFERENCES program_content_type(id),
            contract_duration_unit_id INT REFERENCES contract_duration_unit(id),
            airtime_duration_unit_id INT REFERENCES airtime_duration_unit(id),
            periodicity_unit_id INT REFERENCES periodicity_unit(id),
            periodicity_value INT NOT NULL
        );`,
		`CREATE TABLE IF NOT EXISTS program_time (
            program_id INT REFERENCES programs(id),
            time_type_id INT REFERENCES program_time_type(id),
            PRIMARY KEY (program_id, time_type_id)
        );`,
	}

	for _, query := range queries {
		_, err := db.Exec(query)
		if err != nil {
			fmt.Println("Error creating table:", err)
		} else {
			fmt.Println("Table created or already exists")
		}
	}
}
/*
Эта функция вставляет начальные данные в созданные таблицы базы данных.
*/
func insertInitialData(db *sql.DB) {
	inserts := []string{
		`INSERT INTO program_content_type (name) VALUES 
        ('Реклама'), 
        ('Развлекательная'), 
        ('Образовательная'), 
        ('Детская'), 
        ('Для взрослых'), 
        ('Новости') 
        ON CONFLICT (name) DO NOTHING;`,

		`INSERT INTO program_time_type (name) VALUES 
        ('Утренняя'), 
        ('Дневная'), 
        ('Вечерняя'), 
        ('Ночная') 
        ON CONFLICT (name) DO NOTHING;`,

		`INSERT INTO contract_duration_unit (name) VALUES 
        ('Дни'), 
        ('Месяцы'), 
        ('Годы') 
        ON CONFLICT (name) DO NOTHING;`,

		`INSERT INTO airtime_duration_unit (name) VALUES 
        ('Секунды'), 
        ('Минуты'), 
        ('Часы') 
        ON CONFLICT (name) DO NOTHING;`,

		`INSERT INTO periodicity_unit (name) VALUES 
        ('Каждые n минут'), 
        ('Каждые n часов'), 
        ('Каждые n дней') 
        ON CONFLICT (name) DO NOTHING;`,

		`INSERT INTO programs (name, content_type_id, contract_duration_unit_id, airtime_duration_unit_id, periodicity_unit_id, periodicity_value) VALUES
        ('И дым отечества нам сладок и приятен. Parliament', 1, 3, 1, 1, 20),
        ('Утренняя гимнастика!', 2, 3, 2, 3, 24),
        ('Галилео', 3, 2, 3, 3, 3),
        ('Ежедневный мультсериал', 4, 2, 3, 3, 24),
        ('Политобзор', 6, 1, 2, 3, 7),
        ('Время', 6, 3, 2, 3, 24)
        ON CONFLICT (name) DO NOTHING;`,

		`INSERT INTO program_time (program_id, time_type_id) VALUES
        ((SELECT id FROM programs WHERE name = 'И дым отечества нам сладок и приятен. Parliament'), (SELECT id FROM program_time_type WHERE name = 'Утренняя')),
        ((SELECT id FROM programs WHERE name = 'И дым отечества нам сладок и приятен. Parliament'), (SELECT id FROM program_time_type WHERE name = 'Дневная')),
        ((SELECT id FROM programs WHERE name = 'И дым отечества нам сладок и приятен. Parliament'), (SELECT id FROM program_time_type WHERE name = 'Вечерняя')),
        ((SELECT id FROM programs WHERE name = 'И дым отечества нам сладок и приятен. Parliament'), (SELECT id FROM program_time_type WHERE name = 'Ночная')),
        ((SELECT id FROM programs WHERE name = 'Утренняя гимнастика!'), (SELECT id FROM program_time_type WHERE name = 'Утренняя')),
        ((SELECT id FROM programs WHERE name = 'Галилео'), (SELECT id FROM program_time_type WHERE name = 'Дневная')),
        ((SELECT id FROM programs WHERE name = 'Ежедневный мультсериал'), (SELECT id FROM program_time_type WHERE name = 'Дневная')),
        ((SELECT id FROM programs WHERE name = 'Ежедневный мультсериал'), (SELECT id FROM program_time_type WHERE name = 'Вечерняя')),
        ((SELECT id FROM programs WHERE name = 'Политобзор'), (SELECT id FROM program_time_type WHERE name = 'Ночная')),
        ((SELECT id FROM programs WHERE name = 'Время'), (SELECT id FROM program_time_type WHERE name = 'Утренняя')),
        ((SELECT id FROM programs WHERE name = 'Время'), (SELECT id FROM program_time_type WHERE name = 'Дневная')),
        ((SELECT id FROM programs WHERE name = 'Время'), (SELECT id FROM program_time_type WHERE name = 'Вечерняя'))
        ON CONFLICT DO NOTHING;`,
	}

	for _, insert := range inserts {
		_, err := db.Exec(insert)
		if err != nil {
			fmt.Println("Error inserting data:", err)
		} else {
			fmt.Println("Data inserted or already exists")
		}
	}
}
/*
Эта функция является точкой входа в программу. В ней происходит настройка базы данных, создание таблиц, 
вставка начальных данных, настройка обработчиков HTTP-запросов и запуск HTTP-сервера.
*/
func main() {
	db, err := setupDB()
	if err != nil {
		fmt.Println(err)
		return
	}
	defer db.Close()

	createTables(db)
	insertInitialData(db)

	http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		fmt.Fprintf(w, "Welcome to the Home Page!")
	})
	//Этот обработчик выполняет запрос к таблице programs в базе данных и выводит список программ на веб-страницу.
	http.HandleFunc("/programs", func(w http.ResponseWriter, r *http.Request) {
		rows, err := db.Query("SELECT id, name FROM programs")
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
		defer rows.Close()

		var programs []struct {
			ID   int
			Name string
		}

		for rows.Next() {
			var program struct {
				ID   int
				Name string
			}
			if err := rows.Scan(&program.ID, &program.Name); err != nil {
				http.Error(w, err.Error(), http.StatusInternalServerError)
				return
			}
			programs = append(programs, program)
		}

		for _, program := range programs {
			fmt.Fprintf(w, "ID: %d, Name: %s\n", program.ID, program.Name)
		}
	})

	fmt.Println("Starting server at port 8080")
	if err := http.ListenAndServe(":8080", nil); err != nil {
		fmt.Println(err)
	}
}
