package main

import (
	"fmt"
	"github.com/go-local-cookie/chrome"
)

func main() {
	res,err:= chrome.GetCookiesByDomain("zhaopin.com")
	if err!=nil{
		fmt.Println(err)
	}
	fmt.Println(res)
}
