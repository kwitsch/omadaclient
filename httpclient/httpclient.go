package httpclient

import (
	"crypto/tls"
	"net/http"
	"net/http/cookiejar"

	"golang.org/x/net/publicsuffix"
)

func New(skipVerify bool) (*http.Client, error) {
	options := cookiejar.Options{
		PublicSuffixList: publicsuffix.List,
	}

	cookies, err := cookiejar.New(&options)
	if err != nil {
		return nil, err
	}

	http := http.Client{
		Jar: cookies,
		Transport: &http.Transport{
			TLSClientConfig: &tls.Config{
				InsecureSkipVerify: skipVerify,
			},
		},
	}

	return &http, nil
}
