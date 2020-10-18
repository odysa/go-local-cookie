package chrome

import (
	"encoding/base64"
	"github.com/tidwall/gjson"
	"io/ioutil"
)

const keyDirWin = `\AppData\Local\Google\Chrome\User Data\Local State`

func chromeDecrypt(encrypted []byte) (string, error) {
	return aesDecrypt(string(encrypted))
}

func aesDecrypt(encrypted string) (string, error) {
	key,err := getKeyFromChrome("win")
	if err != nil {
		return "", err
	}

	decodedKey, err := base64.StdEncoding.DecodeString(key)
	if err != nil {
		return "", err
	}
	decodedKey = decodedKey[5:]

	return string(decodedKey), nil
}

func readJsonFile(fileName string) (string, error) {
	bytes, err := ioutil.ReadFile(fileName)

	if err != nil {
		return "", err
	}

	return string(bytes), nil
}

func getKeyFromChrome(os string) (string, error) {
	var key string
	if os == "win" {
		jsn, err := readJsonFile(getUsrDir() + keyDirWin)
		if err != nil {
			return "", err
		}
		// get key
		key = gjson.Get(jsn, "os_crypt.encrypted_key").String()
	}
	return key, nil
}
