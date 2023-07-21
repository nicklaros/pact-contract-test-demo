package common

import (
	"encoding/json"
	"io"
	"net/http"
)

func CallService(url string, jsonResp interface{}) {
	resp, _ := http.Get(url)

	body, _ := io.ReadAll(resp.Body)
	defer resp.Body.Close()

	json.Unmarshal(body, &jsonResp)
}
