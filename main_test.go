package main

import (
	"os"
	"testing"

	"github.com/stretchr/testify/require"
)

func TestCreateGame(t *testing.T) {
	if len(os.Getenv("RUNNING_GITHUB_ACTIONS")) > 0 {
		t.Skip("requires database, skipping github actions")
	}
	requestBody := `{
 "title": "Game Created in Test Mode"
}`

	res, bodyString := PostRequest(t, "http://localhost:8005/games", nil, requestBody)
	require.Contains(t, bodyString, `{"id":`)
	require.Equal(t, "201 Created", res.Status)
}
