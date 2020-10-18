package chrome

import (
	"crypto/aes"
	"crypto/cipher"
	"encoding/base64"
	"encoding/hex"
	"fmt"
	dpapi "github.com/go-local-cookie/chrome/win"
	"github.com/tidwall/gjson"
	"io/ioutil"
)

const keyDirWin = `\AppData\Local\Google\Chrome\User Data\Local State`

func chromeDecrypt(encrypted string) (string, error) {
	return aesDecrypt(encrypted)
}

func aesDecrypt(encrypted string) (string, error) {
	key, err := getKeyFromChrome("win")
	if err != nil {
		return "", err
	}

	decodedKey, err := base64.StdEncoding.DecodeString(key)
	if err != nil {
		return "", err
	}
	decodedKey, err = dpapi.DecryptBytes(decodedKey[5:])
	var nonce string
	fmt.Println(encrypted)
	if len(encrypted) > 15 {
		nonce = encrypted[3:15]
	}

	block, err := aes.NewCipher(decodedKey)
	if err != nil {
		return "", err
	}

	aesGcm, err := cipher.NewGCM(block)
	deNonce,_ := hex.DecodeString(nonce)
	deText,_ := hex.DecodeString(encrypted)
	text, err := aesGcm.Open(nil, deNonce, deText, nil)
	if err != nil {
		return "", err
	}
	return string(text), nil
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
