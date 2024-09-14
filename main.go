package safetrade

import (
	"crypto/hmac"
	"crypto/sha256"
	"encoding/hex"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"strconv"
	"time"

	_ "github.com/gorilla/websocket"
)

type FloatLike float64

func (b *FloatLike) UnmarshalJSON(data []byte) error {
	// Try to unmarshal as float64 or string then float
	var floatVal float64
	if err := json.Unmarshal(data, &floatVal); err == nil {
		*b = FloatLike(floatVal)
		return nil
	}
	var stringVal string
	if err := json.Unmarshal(data, &stringVal); err == nil {
		f, err := strconv.ParseFloat(stringVal, 64)
		if err != nil {
			return fmt.Errorf("Failed to parse float from string: %s", err)
		}
		*b = FloatLike(f)
		return nil
	}

	return fmt.Errorf("unable to unmarshal balance")
}

type SpotAccount struct {
	Balance          FloatLike `json:"balance"`
	Currency         string    `json:"currency"`
	DepositAddresses []struct {
		Address    string   `json:"address"`
		Currencies []string `json:"currencies"`
		Network    string   `json:"network"`
	} `json:"deposit_addresses"`
	Locked FloatLike `json:"locked"`
	Type   string    `json:"type"`
}

func NewRequest(key string, secret string, method string, path string, data []byte) (*http.Request, error) {
	nonce := fmt.Sprintf("%d", time.Now().UnixMilli())
	signature := nonce + key
	hmac_sign := hmac.New(sha256.New, []byte(secret))
	if _, err := hmac_sign.Write([]byte(signature)); err != nil {
		return nil, fmt.Errorf("Failed HMAC signing request: %s", err)
	}
	headers := http.Header{
		"X-Auth-Apikey":    {key},
		"X-Auth-Nonce":     {nonce},
		"X-Auth-Signature": {hex.EncodeToString(hmac_sign.Sum(nil))},
		"content-type":     {"application/json"},
	}
	req, err := http.NewRequest(method, "https://safe.trade/api/v2"+path, nil)
	if err != nil {
		return nil, fmt.Errorf("Failed to create request: %s", err)
	}
	req.Header = headers
	return req, err
}

func RequestMaker(key string, secret string) func(method string, path string, data []byte) (*http.Request, error) {
	return func(method string, path string, data []byte) (*http.Request, error) {
		return NewRequest(key, secret, method, path, data)
	}
}

func GetSpotAccounts(key string, secret string) ([]SpotAccount, error) {
	req, err := RequestMaker(key, secret)("GET", "/trade/account/balances/spot", []byte(`{}`))
	if err != nil {
		return nil, fmt.Errorf("Failed to build request: %s", err)
	}
	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		return nil, fmt.Errorf("Failed to send request: %s", err)
	}
	defer resp.Body.Close()
	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("Status code is not OK: %d", resp.StatusCode)
	}
	spotaccounts := []SpotAccount{}

	bodyBytes, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, fmt.Errorf("Failed to read response body: %s", err)
	}

	if err := json.Unmarshal(bodyBytes, &spotaccounts); err != nil {
		return nil, fmt.Errorf("Failed to unmarshal response body: %s", err)
	}
	return spotaccounts, nil
}
