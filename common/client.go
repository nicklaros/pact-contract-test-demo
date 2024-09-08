package common

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
)

func GetServiceFullUrl(url string, productId string) string {
	return fmt.Sprintf("%s/?id=%s", url, productId)
}

func CallService(url string, jsonResp interface{}) {
	resp, _ := http.Get(url)

	body, _ := io.ReadAll(resp.Body)
	defer resp.Body.Close()

	json.Unmarshal(body, &jsonResp)
}
