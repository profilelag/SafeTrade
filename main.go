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

func RequestMaker(key string, secret string) func (method string, path string, data []byte) (*http.Request, error) {
	return func (method string, path string, data []byte) (*http.Request, error) {
		nonce := fmt.Sprintf("%d", time.Now().UnixNano()/1000000)
		signature := nonce + key
		hmac_sign := hmac.New(sha256.New, []byte(signature))
		hmac_sign.Write([]byte(secret))
		headers := http.Header{
			"X-Auth-Apikey": {key},
			"X-Auth-Nonce": {nonce},
			"X-Auth-Signature": {hex.EncodeToString(hmac_sign.Sum(nil))},
			"content-type": {"application/json"},
		}
		req, err := http.NewRequest(method, "https://safe.trade/api/v2" + path, nil)
		req.Header = headers
		return req, err
	}
}