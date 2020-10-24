package localCookie

import (
	"github.com/go-local-cookie/chrome"
	"net/http"
)

func GetCookies(domain string, name string) ([]http.Cookie, error) {
	return chrome.GetCookiesCore(domain, name)
}
func GetCookiesByDomain(domain string) ([]http.Cookie, error) {
	return chrome.GetCookiesCore(domain, "")
}
func GetCookieByName(name string) (http.Cookie, error) {
	res, err := chrome.GetCookiesCore("", name)
	if err != nil {
		return http.Cookie{}, err
	}
	return res[0], nil
}
