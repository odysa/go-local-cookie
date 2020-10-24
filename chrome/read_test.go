package chrome

import (
	"testing"
)

func TestConnectDataBase(t *testing.T) {
	db, err := connectDatabase(winDir)
	defer db.Close()
	if err != nil {
		panic(err)
	}
	if db == nil {
		t.Errorf("db connect should not be nil!")
	}
}

func TestReadFromSqlite(t *testing.T) {
	db, err := connectDatabase(winDir)
	defer db.Close()
	if err != nil {
		panic(err)
	}
	res, err := readFromSqlite(db, "", "")
	if err != nil {
		panic(err)
	}
	if len(res) == 0 {
		t.Errorf("query data should not be 0!")
	}

}
