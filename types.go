package main

type CookieResult struct {
	Domain   string `db:"host_key"`
	Name     string `db:"name"`
	Path     string `db:"path"`
	Secure   bool   `db:"is_secure"`
	HttpOnly bool   `db:"is_httponly"`
	Expires  int64 `db:"expires_utc"`

	Value    string
}