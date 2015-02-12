package urlutil

import (
	"io"
	"io/ioutil"
	"net/http"
	"net/url"
)

func BuildUrl(sUrlBase string, dParams map[string]string) string {
	if len(dParams) < 1 {
		return sUrlBase
	}
	vals := url.Values{}
	for key, val := range dParams {
		vals.Set(key, val)
	}
	qryString := vals.Encode()
	sUrlFull := sUrlBase + "?" + qryString
	return sUrlFull
}

func GetUrlBody(sUrl string) ([]byte, error) {
	req, err := http.NewRequest("GET", sUrl, nil)
	cli := &http.Client{}
	res, err := cli.Do(req)
	if err != nil {
		return []byte{}, err
	}
	defer res.Body.Close()
	return ioutil.ReadAll(res.Body)
}

func GetUrlPostBody(url string, bodyType string, reqBody io.Reader) ([]byte, error) {
	client := &http.Client{}
	res, err := client.Post(url, bodyType, reqBody)
	if err != nil {
		return []byte{}, err
	}
	defer res.Body.Close()
	return ioutil.ReadAll(res.Body)
}
