package main

import (
	"io/ioutil"
	"net/http"
	"strings"
	"testing"

	"github.com/stretchr/testify/require"
)

// PostRequest(t, url, nil, requestBody)
func PostRequest(t *testing.T, url string, h http.Header, requestBody string) (*http.Response, string) {
	req, err := http.NewRequest("POST", url, strings.NewReader(requestBody))
	if h != nil {
		req.Header = h
	}
	req.Header.Set("Content-Type", "application/json")
	res, err := http.DefaultClient.Do(req)
	if t != nil {
		require.NoError(t, err)
	} else {
		if err != nil {
			return nil, ""
		}
	}
	defer res.Body.Close()

	responseBody, err := ioutil.ReadAll(res.Body)
	if t != nil {
		require.NoError(t, err)
	}
	return res, string(responseBody)
}
