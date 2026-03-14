package httputilmore

import "net/http"

type Cookies []*http.Cookie

func (c Cookies) String() string {
	req, err := http.NewRequest(http.MethodGet, "https://example.com", nil)
	if err != nil {
		panic(err)
	}

	for _, cookie := range c {
		req.AddCookie(cookie)
	}

	return req.Header.Get("Cookie")
}
