package utils

import "encoding/json"

// ToJSON interface to json string
func ToJSON(i interface{}) string {
	data, _ := json.Marshal(i)
	return string(data)
}

// ToRawMessage interface to raw json bytes
func ToRawMessage(i interface{}) json.RawMessage {
	data, _ := json.Marshal(i)
	return data
}

// ToJSONArray list of interfaces to json string
func ToJSONArray(i ...interface{}) string {
	return ToJSON(i)
}
