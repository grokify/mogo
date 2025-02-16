package httpsimple

import (
	"context"
	"fmt"
	"io"
	"net/http"
	"net/url"
	"strings"

	"github.com/grokify/mogo/net/http/httputilmore"
)

// CLI executes a HTTP request from CLI params. Attribute tags support usage
// with `github.com/jessevdk/go-flags`.`
type CLI struct {
	Method string   `short:"M" long:"method" description:"Request Method"`
	URL    string   `short:"U" long:"url" description:"Reaquest URL"`
	Header []string `short:"H" long:"header" description:"Request Header"`
	Body   string   `short:"B" long:"body" description:"Request Body"`
}

func (cli CLI) Request() (Request, error) {
	req := Request{
		Method:  cli.Method,
		URL:     cli.URL,
		Query:   url.Values{},
		Headers: http.Header{},
		Body:    cli.Body,
	}
	for _, hi := range cli.Header {
		if hi == "" {
			continue
		}
		parts := strings.SplitN(hi, ":", 2)
		if len(parts) != 2 {
			return req, fmt.Errorf("header malformed (%s)", hi)
		}
		req.Headers.Add(strings.TrimSpace(parts[0]), strings.TrimSpace(parts[1]))
	}
	return req, nil
}

func (cli CLI) Do(ctx context.Context, w io.Writer) error {
	req, err := cli.Request()
	if err != nil {
		return err
	}
	if ctx == nil {
		ctx = context.Background()
	}

	resp, err := req.Do(ctx)
	if err != nil {
		return err
	}

	if w != nil {
		if _, err := w.Write([]byte(fmt.Sprintf("Response Status Code: %d\n", resp.StatusCode))); err != nil {
			return err
		}
		if _, err := w.Write([]byte(fmt.Sprintf("Response ContentType: %s\n", resp.Header.Get(httputilmore.HeaderContentType)))); err != nil {
			return err
		}
		b, err := httputilmore.ResponseBodyMore(resp, "", "  ")
		if err != nil {
			return err
		}
		if _, err := w.Write([]byte(fmt.Sprintf("Body:\n%s", string(b)))); err != nil {
			return err
		}
	}

	return nil
}
