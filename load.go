package Browser

import (
	"fmt"
	"github.com/go-sqlite/sqlite3"
	"net/http"
	"os/user"
	"time"
)

const winDir = `\AppData\Local\Google\Chrome\User Data\default\Cookies`

func LoadCookieFromChrome(domain string) ([]*http.Cookie, error) {
	var cookies []*http.Cookie
	usr, _ := user.Current()
	cookiesFile := fmt.Sprintf(`%s%s`, usr.HomeDir, winDir)

	db, err := sqlite3.Open(cookiesFile)
	if err != nil {
		return nil, err
	}
	defer db.Close()

	return cookies, nil
}

func ReadFromSqlite(db *sqlite3.DbFile) error {
	return db.VisitTableRecords("cookies", func(rowId *int64, rec sqlite3.Record) error {
		if rowId == nil {
			return fmt.Errorf("unexpected nil RowID in Chrome sqlite database")
		}
		cookie := &http.Cookie{}

		if len(rec.Values) != 14 {
			return nil
		}

		domain, ok := rec.Values[1].(string)
		if !ok {
			return fmt.Errorf("expected column 2 (host_key) to to be string; got %T", rec.Values[1])
		}
		name, ok := rec.Values[2].(string)
		if !ok {
			return fmt.Errorf("expected column 3 (name) in cookie(domain:%s) to to be string; got %T", domain, rec.Values[2])
		}
		value, ok := rec.Values[3].(string)
		if !ok {
			return fmt.Errorf("expected column 4 (value) in cookie(domain:%s, name:%s) to to be string; got %T", domain, name, rec.Values[3])
		}
		path, ok := rec.Values[4].(string)
		if !ok {
			return fmt.Errorf("expected column 5 (path) in cookie(domain:%s, name:%s) to to be string; got %T", domain, name, rec.Values[4])
		}
		var expiresUtc int64
		switch i := rec.Values[5].(type) {
		case int64:
			expiresUtc = i
		case int:
			if i != 0 {
				return fmt.Errorf("expected column 6 (expires_utc) in cookie(domain:%s, name:%s) to to be int64 or int with value=0; got %T with value %v", domain, name, rec.Values[5], rec.Values[5])
			}
		default:
			return fmt.Errorf("expected column 6 (expires_utc) in cookie(domain:%s, name:%s) to to be int64 or int with value=0; got %T with value %v", domain, name, rec.Values[5], rec.Values[5])
		}
		encryptedValue, ok := rec.Values[12].([]byte)
		if !ok {
			return fmt.Errorf("expected column 13 (encrypted_value) in cookie(domain:%s, name:%s) to to be []byte; got %T", domain, name, rec.Values[12])
		}

		var expiry time.Time
		if expiresUtc != 0 {
			expiry = chromeCookieDate(expiresUtc)
		}
		//creation := chromeCookieDate(*rowId)

		cookie.Domain = domain
		cookie.Name = name
		cookie.Path = path
		cookie.Expires = expiry
		//cookie.Creation = creation
		cookie.Secure = rec.Values[6] == 1
		cookie.HttpOnly = rec.Values[7] == 1

		if len(encryptedValue) > 0 {
			decrypted, err := decryptValue(encryptedValue)
			if err != nil {
				return fmt.Errorf("decrypting cookie %v: %v", cookie, err)
			}
			cookie.Value = decrypted
		} else {
			cookie.Value = value
		}
		return nil
	})
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
