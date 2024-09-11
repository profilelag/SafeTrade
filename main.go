package safetrade

import (
	"crypto/hmac"
	"crypto/sha256"
	"encoding/hex"
	"fmt"
	"net/http"
	"time"

	_ "github.com/gorilla/websocket"
)

func NewRequest(key string, secret string, method string, path string, data []byte) (*http.Request, error) {
	nonce := fmt.Sprintf("%d", time.Now().UnixMicro()/1000)
	signature := nonce + key
	hmac_sign := hmac.New(sha256.New, []byte(secret))
	hmac_sign.Write([]byte(signature))
	headers := http.Header{
		"X-Auth-Apikey": {key},
		"X-Auth-Nonce": {nonce},
		"X-Auth-Signature": {hex.EncodeToString(hmac_sign.Sum(nil))},
		"content-type": {"application/json"},
	}
	req, err := http.NewRequest(method, "https://safe.trade/api/v2"+path, nil)
	req.Header = headers
	return req, err
}

func RequestMaker(key string, secret string) func(method string, path string, data []byte) (*http.Request, error) {
	return func(method string, path string, data []byte) (*http.Request, error) {
		return NewRequest(key, secret, method, path, data)
	}
}