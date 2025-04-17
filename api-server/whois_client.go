package main

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
)

func AnalyzeDomain(domain string) (map[string]interface{}, error) {
	url := fmt.Sprintf("http://phishing-analyzer:5001/analyze?domain=%s", domain)
	resp, err := http.Get(url)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	body, _ := ioutil.ReadAll(resp.Body)

	var result map[string]interface{}
	if err := json.Unmarshal(body, &result); err != nil {
		return nil, err
	}
	return result, nil
}
