package main

import (
	"fmt"
	"github.com/go-local-cookie/chrome"
)

func main() {
	_,err:= chrome.LoadCookieFromChrome(".zhaopin.com")
	if err!=nil{
		fmt.Println(err)
	}
}
