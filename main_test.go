//go:build integration
// +build integration

package safetrade

import (
	"io"
	"net/http"
	"os"
	"testing"

	_ "github.com/gorilla/websocket"
	"github.com/joho/godotenv"
)

func TestNewRequest(t *testing.T) {
	err := godotenv.Load()
	if err != nil {
		t.Fatal("Error loading .env file")
	}
	key := os.Getenv("KEY")
	secret := os.Getenv("SECRET")
	method := "GET"
	path := "/trade/account/members/me"
	data := []byte(`{}`)
	req, err := NewRequest(key, secret, method, path, data)
	if err != nil {
		t.Fatal(err)
	}
	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		t.Fatal(err)
	}
	if resp.StatusCode != 200 {
		bodyBytes, err := io.ReadAll(resp.Body)
		if err != nil {
			t.Fatal(err)
		}
		bodyString := string(bodyBytes)
		t.Log(bodyString)
		t.Fatal("Status code is not 200")
	}
	t.Log("Status code is 200")
}

func TestGetSpotAccounts(t *testing.T) {
	err := godotenv.Load()
	if err != nil {
		t.Fatal("Error loading .env file")
		return
	}
	key := os.Getenv("KEY")
	secret := os.Getenv("SECRET")
	spotaccounts, err := GetSpotAccounts(key, secret)
	if err != nil {
		t.Fatal(err)
		return
	}
	t.Fatal(spotaccounts)
	for _, spotaccount := range spotaccounts {
		if spotaccount.Balance != 0 {
			t.Fatal(spotaccount.Balance)
		}
	}
}
