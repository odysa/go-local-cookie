package localCookie

import "testing"

func TestGetCookies(t *testing.T) {
	// test empty name
	res, err := GetCookies("google.com", "")
	if err != nil {
		t.Errorf("failed to load cookie")
	}
	if len(res) == 0 {
		t.Errorf("cookies should not be empty")
	}

	// test empty domain
	res, err = GetCookies("", "SID")
	if err != nil {
		t.Errorf("failed to load cookie,%v", err)
	}
	if len(res) == 0 {
		t.Errorf("cookies should not be empty")
	}
	res, err = GetCookies("google.com", "SID")
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