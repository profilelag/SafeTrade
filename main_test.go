package safetrade

import (
	"io"
	"net/http"
	"os"
	"testing"
	"time"

	"github.com/joho/godotenv"
)

func TestNewRequest(t *testing.T) {
	err := godotenv.Load()
	if err != nil {
		t.Fatal("Error loading .env file")
	}
	t.Log(time.Now().UnixMicro()/1000)
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
