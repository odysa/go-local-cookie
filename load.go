package main

import (
	"database/sql"
	"fmt"
	_ "github.com/mattn/go-sqlite3"
	"net/http"
	"os/user"
	"time"
)

const winDir = `\AppData\Local\Google\Chrome\User Data\default\Cookies`

func LoadCookieFromChrome(domain string) ([]*http.Cookie, error) {
	var cookies []*http.Cookie
	db, err := ConnectDatabase(winDir)
	defer db.Close()
	if err != nil {
		fmt.Println(err)
		return nil, err
	}

	ReadFromSqlite(db, domain)

	return cookies, nil
}

func ConnectDatabase(location string) (*sql.DB, error) {
	usr, _ := user.Current()
	cookiesFile := fmt.Sprintf(`%s%s`, usr.HomeDir, location)
	db, err := sql.Open("sqlite3", cookiesFile)
	if err != nil {
		return nil, err
	}
	return db, err
}

const retrieveQ = `SELECT host_key,name,path,is_secure,is_httponly,expires_utc,value FROM cookies where host_key like ?`

func ReadFromSqlite(db *sql.DB, targetDomain string) ([]http.Cookie, error) {
	var (
		domain, name, path, value string
		secure, httponly          bool
		expire                    int64
		result                    []http.Cookie
	)
	err := db.Ping()
	if err != nil {
		fmt.Println(err)
	}
	rows, err := db.Query(retrieveQ, "%"+targetDomain+"%")
	defer rows.Close()
	if err != nil {
		return nil, err
	}

	for rows.Next() {
		err = rows.Scan(&domain, &name, &path, &secure, &httponly, &expire, &value)
		result = append(result, http.Cookie{
			Domain:   domain,
			Name:     name,
			Path:     path,
			Secure:   secure,
			HttpOnly: httponly,
			Expires:  chromeCookieDate(expire),
			Value:    decrypt(value),
		})
	}

	fmt.Printf("%v\n", result)
	return result, nil
}

// See https://cs.chromium.org/chromium/src/base/time/time.h?l=452&rcl=fceb9a030c182e939a436a540e6dacc70f161cb1
const windowsToUnixMicrosecondsOffset = 11644473600000000

// chromeCookieDate converts microseconds to a time.Time object,
// accounting for the switch to Windows epoch (Jan 1 1601).
func chromeCookieDate(timestampUtc int64) time.Time {
	if timestampUtc > windowsToUnixMicrosecondsOffset {
		timestampUtc -= windowsToUnixMicrosecondsOffset
	}

	return time.Unix(timestampUtc/1000000, (timestampUtc%1000000)*1000)
}
