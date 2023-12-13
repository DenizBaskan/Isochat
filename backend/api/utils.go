package api

import (
	"encoding/json"
	"net/http"
	"net/url"
	"bytes"
	"os"
)

func reason(msg string) map[string]string {
	return map[string]string{"reason": msg}
}

func isHCapValid(token string) (bool, error) {
	values := url.Values{}

	values.Add("secret", os.Getenv("HCAP_SECRET"))
	values.Add("response", token)

	r, err := http.NewRequest("POST", "https://hcaptcha.com/siteverify", bytes.NewBuffer([]byte(values.Encode())))
	if err != nil {
		return false, err
	}

	r.Header.Add("content-type", "application/x-www-form-urlencoded")

	client := &http.Client{}
	res, err := client.Do(r)
	if err != nil {
		return false, err
	}

	defer res.Body.Close()

	var data map[string]any

	if err = json.NewDecoder(res.Body).Decode(&data); err != nil {
		return false, err
	}

	if data["success"].(bool) {
		return true, nil
	}

	return false, nil
}
