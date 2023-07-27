package dbconnector

import "testing"

const (
	DB_USER = "postgres"
	DB_PASS = "postgres"
	DB_NAME = "postgres"
)

// NewDBConnector returns a new DBConnector
func TestNewDbConnector(t *testing.T) {
	db, err := NewPostgreSQLConnector(DB_USER, DB_PASS, DB_NAME)

	if err != nil {
		t.Error(err)
	}
	err = DBStartupTask(db)

	if err != nil {
		t.Error(err)
	}
}
