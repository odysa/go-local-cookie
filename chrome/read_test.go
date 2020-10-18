package chrome

import (
	"testing"
)

func TestGetCookies(t *testing.T) {
	// test empty name
	res, err := getCookiesCore("google.com", "")
	if err != nil {
		t.Errorf("failed to load cookie")
	}
	if len(res) == 0 {
		t.Errorf("cookies should not be empty")
	}

	// test empty domain
	res, err = getCookiesCore("", "SID")
	if err != nil {
		t.Errorf("failed to load cookie,%v", err)
	}
	if len(res) == 0 {
		t.Errorf("cookies should not be empty")
	}
	res, err = getCookiesCore("google.com", "SID")
	if err != nil {
		t.Errorf("failed to load cookie,%v", err)
	}
	if len(res) == 0 {
		t.Errorf("cookie should not be empty,%v", res)
	}
}

func TestGetCookieByName(t *testing.T) {
	const testData = "SID"
	res, err := GetCookieByName(testData)
	if err != nil {
		t.Errorf("failed to load cookie,%v", err)
	}
	if res.Name != testData {
		t.Errorf("Name should be %s,but got %s", testData, res.Name)
	}
	if len(res.Value)==0{
		t.Errorf("failed to read value of cookie")
	}
}
func TestGetCookiesByDomain(t *testing.T) {
	res, err := GetCookiesByDomain("google.com")
	if err != nil {
		t.Errorf("failed to load cookie")
	}
	if len(res) == 0 {
		t.Errorf("cookies should not be empty")
	}
}

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
