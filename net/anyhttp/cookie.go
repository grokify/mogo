package anyhttp

import (
	"net/http"

	"github.com/valyala/fasthttp"
)

type Cookie struct {
	Name  string
	Value string
}

func (c *Cookie) ToNetHttp() *http.Cookie {
	cookie := http.Cookie{Name: c.Name, Value: c.Value}
	return &cookie
}

func (c *Cookie) ToFastHttp() *fasthttp.Cookie {
	var cookie fasthttp.Cookie
	cookie.SetKey(c.Name)
	cookie.SetValue(c.Value)
	return &cookie
}
