package handlers

import "encoding/json"

func reason(s string) []byte {
	res, _ := json.Marshal(map[string]any{"success": false, "reason": s})
	return res
}

func success(data any, status int8) []byte {
	res, _ := json.Marshal(map[string]any{"success": true, "data": data, "status": status})
	return res
}
