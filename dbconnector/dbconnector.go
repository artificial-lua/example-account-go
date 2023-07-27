package dbconnector

import (
	"database/sql"
	"fmt"

	_ "github.com/lib/pq"
)

func SetClinetEncoding(conn *sql.DB) (sql.Result, error) {
	return conn.Exec("SET client_encoding TO 'UTF8'")
}

func checkAndMakeTable(conn *sql.DB, tableName string, query string) (sql.Result, error) {
	return conn.Exec(fmt.Sprintf(`CREATE TABLE IF NOT EXISTS %s (%s)`, tableName, query))
}

func makeUserTable(conn *sql.DB) (sql.Result, error) {
	defaultTableName := "account"
	return checkAndMakeTable(conn, defaultTableName, `email VARCHAR(255) PRIMARY KEY,
	hash VARCHAR(255) NOT NULL,
	salt VARCHAR(255),
	name VARCHAR(255) NOT NULL,
	birth DATE NOT NULL,
	gender VARCHAR(1) NOT NULL`)
}

// NewDBConnector returns a new DBConnector
func NewPostgreSQLConnector(username string, password string, db string) (*sql.DB, error) {
	dbinfo := fmt.Sprintf("user=%s password=%s dbname=%s sslmode=disable", username, password, db)

	conn, err := sql.Open("postgres", dbinfo)
	if err != nil {
		return nil, err
	}

	_, err = SetClinetEncoding(conn)
	if err != nil {
		return nil, err
	}

	return conn, nil
}

func DBStartupTask(conn *sql.DB) error {
	_, err := SetClinetEncoding(conn)
	if err != nil {
		return err
	}

	_, err = makeUserTable(conn)
	if err != nil {
		return err
	}

	return nil
}
