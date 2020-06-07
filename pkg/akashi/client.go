package akashi

import (
	"context"
	"log"
	"net/http"
)

var (
	endpointURL = "https://atnd.ak4.jp/api/cooperation"
)

// Client is client interface
type Client interface {
	Get(ctx context.Context, url string) (*http.Response, error)
}

type client struct {
	hc *http.Client
	cc string // company_code
	t  string // access_token
}

func (c client) Get(ctx context.Context, url string) (*http.Response, error) {
	endpoint := endpointURL + url
	log.Println("GET", endpoint)
	req, err := http.NewRequestWithContext(ctx, http.MethodGet, endpoint, nil)
	if err != nil {
		return nil, err
	}
	res, err := c.hc.Do(req)
	if err != nil {
		return nil, err
	}
	return res, nil
}

// NewClient is constructor
func NewClient(companyCode, token string) Client {
	return &client{
		hc: &http.Client{},
		cc: companyCode,
		t:  token,
	}
}
