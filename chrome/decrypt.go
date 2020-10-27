package chrome

import (
	"crypto/aes"
	"crypto/cipher"
	"encoding/base64"
	dpapi "github.com/odysa/go-local-cookie/chrome/win"
	"github.com/tidwall/gjson"
	"io/ioutil"
)

const keyDirWin = `\AppData\Local\Google\Chrome\User Data\Local State`

func decrypt(encrypted string) (string, error) {
	return aesDecrypt(encrypted)
}

func aesDecrypt(encrypted string) (string, error) {
	key, err := getKeyFromChrome("win")
	if err != nil {
		return "", err
	}

	// decode key of aes
	decodedKey, err := base64.RawStdEncoding.DecodeString(key)
	if err != nil {
		return "", err
	}
	decodedKey, err = dpapi.DecryptBytes(decodedKey[5:])

	// get nonce and text of value
	nonce, text := parseText(encrypted)

	block, err := aes.NewCipher(decodedKey)
	if err != nil {
		return "", err
	}

	aesGcm, err := cipher.NewGCM(block)
	// decipher
	res, err := aesGcm.Open(nil, nonce, text, nil)

	if err != nil {
		return "", err
	}
	return string(res), nil
}

// parse new algorithm
func parseText(encrypted string) (nonce []byte, text []byte) {
	nonce = []byte(encrypted[3:15])
	text = []byte(encrypted[15:])
	return
}

func readJsonFile(fileName string) (string, error) {
	bytes, err := ioutil.ReadFile(fileName)

	if err != nil {
		return "", err
	}

	return string(bytes), nil
}

// get key from chrome storage
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
