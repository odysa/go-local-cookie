package main

import "fmt"

func main() {
	_,err:=LoadCookieFromChrome(".zhaopin.com")
	if err!=nil{
		fmt.Println(err)
	}
}
