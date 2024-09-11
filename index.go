package main

import (
	"github.com/profilelag/safetrade"
	"net/http"
)

func main() {
	key := "YOUR"
	secret := ""
	requestMaker := safetrade.RequestMaker(key, secret)
	req, err := requestMaker("GET", "/markets", nil)
	if err != nil {
		panic(err)
	}
	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		panic(err)
	}
	defer resp.Body.Close()
}