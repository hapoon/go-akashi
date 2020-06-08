package akashi

import (
	"context"
	"encoding/json"
	"net/http"
	"net/http/httputil"
	"testing"

	"github.com/stretchr/testify/assert"
)

type testClient struct {
	ID   int    `json:"id"`
	Name string `json:"name"`
}

func TestClientGet(t *testing.T) {
	cli := NewClient("", "")
	testCase := map[string]struct {
		endpointURL string
		statusCode  int
		err         bool
		foo1        string
		foo2        string
		timeout     bool
	}{
		"ok(args not exist)": {
			endpointURL: "https://postman-echo.com/get",
			statusCode:  http.StatusOK,
			err:         false,
		},
		"ok(args exist)": {
			endpointURL: "https://postman-echo.com/get?foo1=abc&foo2=1",
			statusCode:  http.StatusOK,
			err:         false,
			foo1:        "abc",
			foo2:        "1",
		},
		"ok(not found)": {
			endpointURL: "https://postman-echo.com/gett",
			statusCode:  http.StatusNotFound,
			err:         false,
		},
		"ng": {
			endpointURL: "https://postman-echo.com/get",
			statusCode:  http.StatusNotFound,
			err:         true,
			timeout:     true,
		},
	}

	for scenario, test := range testCase {
		endpointURL = test.endpointURL
		ctx := context.Background()
		if test.timeout {
			ctx, _ = context.WithTimeout(ctx, 0)
		}
		res, err := cli.Get(ctx, "")
		var resBody struct {
			Args struct {
				Foo1 string `json:"foo1"`
				Foo2 string `json:"foo2"`
			} `json:"args"`
		}
		if test.err {
			assert.Error(t, err, scenario)
		} else {
			assert.NoError(t, err, scenario)
			assert.Equal(t, test.statusCode, res.StatusCode, scenario)
			b, _ := httputil.DumpResponse(res, true)
			t.Log(string(b))
			if res.StatusCode == http.StatusOK {
				json.NewDecoder(res.Body).Decode(&resBody)
				assert.Equal(t, test.foo1, resBody.Args.Foo1, scenario)
				assert.Equal(t, test.foo2, resBody.Args.Foo2, scenario)
			}
		}
	}
}

func TestClientPost(t *testing.T) {
	cli := NewClient("", "")
	ctx := context.Background()
	tests := map[string]struct {
		endpointURL string
		statusCode  int
		body        testClient
	}{
		"ok": {
			endpointURL: "https://postman-echo.com/post",
			statusCode:  http.StatusOK,
			body: testClient{
				ID:   1,
				Name: "Bob",
			},
		},
	}

	for scenario, test := range tests {
		endpointURL = test.endpointURL
		res, err := cli.Post(ctx, "", test.body)
		var resBody struct {
			Args struct{}   `json:"args"`
			Data testClient `json:"data"`
		}
		assert.NoError(t, err, scenario)
		assert.Equal(t, test.statusCode, res.StatusCode, scenario)
		b, _ := httputil.DumpResponse(res, true)
		t.Log(string(b))
		json.NewDecoder(res.Body).Decode(&resBody)
		assert.Equal(t, test.body, resBody.Data, scenario)
	}
}
