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
	var floatVal float64
	var stringVal string

	// Try to unmarshal as float64
	if err := json.Unmarshal(data, &floatVal); err == nil {
		*b = FloatLike(floatVal)
		return nil
	}

	// Try to unmarshal as string
	if err := json.Unmarshal(data, &stringVal); err == nil {
		// Convert string to float64
		f, err := strconv.ParseFloat(stringVal, 64)
		if err != nil {
			return err
		}
		*b = FloatLike(f)
		return nil
	}

	return fmt.Errorf("unable to unmarshal balance")
}

type SpotAccount struct {
	Balance           FloatLike `json:"balance"`
	Currency           string `json:"currency"`
	DepositAddresses []struct {
		Address    string `json:"address"`
		Currencies []string `json:"currencies"`
		Network    string `json:"network"`
	} `json:"deposit_addresses"`
	Locked FloatLike `json:"locked"`
	Type  string `json:"type"`
}

func NewRequest(key string, secret string, method string, path string, data []byte) (*http.Request, error) {
	nonce := fmt.Sprintf("%d", time.Now().UnixMicro()/1000)
	signature := nonce + key
	hmac_sign := hmac.New(sha256.New, []byte(secret))
	hmac_sign.Write([]byte(signature))
	headers := http.Header{
		"X-Auth-Apikey":    {key},
		"X-Auth-Nonce":     {nonce},
		"X-Auth-Signature": {hex.EncodeToString(hmac_sign.Sum(nil))},
		"content-type":     {"application/json"},
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

func GetSpotAccounts(key string, secret string) ([]SpotAccount, error) {
	req, err := RequestMaker(key, secret)("GET", "/trade/account/balances/spot", []byte(`{}`))
	if err != nil {
		return nil, err
	}
	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()
	if resp.StatusCode != 200 {
		return nil, fmt.Errorf("status code is not 200")
	}
	spotaccounts := []SpotAccount{}

	bodyBytes, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}

	// use bodyBytes to decode the response into a struct
	// Unmarshal the response into the slice of SpotAccount structs
	err = json.Unmarshal(bodyBytes, &spotaccounts)
	if err != nil {
		return nil, err
	}
	return spotaccounts, nil
}
