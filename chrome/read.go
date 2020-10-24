package chrome

import (
	"database/sql"
	_ "github.com/mattn/go-sqlite3"
	"net/http"
	"os/user"
	"time"
)

const winDir = `\AppData\Local\Google\Chrome\User Data\default\Cookies`

const queryStr = `SELECT host_key,name,path,is_secure,is_httponly,expires_utc,encrypted_value FROM cookies where host_key like ?`

func GetCookiesCore(domain string, name string) ([]http.Cookie, error) {

	db, err := connectDatabase(winDir)
	if err != nil {
		return nil, err
	}
	defer db.Close()

	res, err := readFromSqlite(db, domain, name)
	if err != nil {
		return nil, err
	}
	return res, nil
}

func connectDatabase(location string) (*sql.DB, error) {

	cookiesFile := getUsrDir() + location

	db, err := sql.Open("sqlite3", cookiesFile)
	if err != nil {
		return nil, err
	}
	return db, err
}

// read cookies from chrome's sqlite3 data
func readFromSqlite(db *sql.DB, targetDomain string, targetName string) ([]http.Cookie, error) {
	var (
		domain, name, path, value string
		secure, httponly          bool
		expire                    int64
		result                    []http.Cookie
	)
	err := db.Ping()
	if err != nil {
		return nil, err
	}
	var query = queryStr
	if len(targetName) > 0 {
		query += `AND name is ?`
	}

	rows, err := db.Query(query, "%"+targetDomain+"%", targetName)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	for rows.Next() {
		err = rows.Scan(&domain, &name, &path, &secure, &httponly, &expire, &value)
		decodedValue, err := decrypt(value)
		if err != nil {
			panic(err)
		}
		result = append(result, http.Cookie{
			Domain:   domain,
			Name:     name,
			Path:     path,
			Secure:   secure,
			HttpOnly: httponly,
			Expires:  getChromeCookieDate(expire),
			Value:    decodedValue,
		})
	}
	return result, nil
}

// See https://cs.chromium.org/chromium/src/base/time/time.h?l=452&rcl=fceb9a030c182e939a436a540e6dacc70f161cb1
const windowsToUnixMicrosecondsOffset = 11644473600000000

// chromeCookieDate converts microseconds to a time.Time object,
// accounting for the switch to Windows epoch (Jan 1 1601).
func getChromeCookieDate(timestampUtc int64) time.Time {
	if timestampUtc > windowsToUnixMicrosecondsOffset {
		timestampUtc -= windowsToUnixMicrosecondsOffset
	}

	return time.Unix(timestampUtc/1000000, (timestampUtc%1000000)*1000)
}

func getUsrDir() string {
	usr, _ := user.Current()
	return usr.HomeDir
}
