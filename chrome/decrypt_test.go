package chrome

import (
	"testing"
)

func TestGetKeyFromChrome(t *testing.T) {
	key, _ := getKeyFromChrome("win")
	if len(key) == 0 {
		t.Errorf("load key failed")
	}
}