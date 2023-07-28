package dbconnector

import "testing"

// NewDBConnector returns a new DBConnector
func TestNewDbConnector(t *testing.T) {
	db, err := NewPostgreSQLConnector()

	if err != nil {
		t.Error(err)
	}
	err = DBStartupTask(db)

	if err != nil {
		t.Error(err)
	}
}
