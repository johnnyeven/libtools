package agent

import (
	"encoding/json"
)

func MarshalData(v interface{}) ([]byte, error) {
	return json.Marshal(v)
}

func UnmarshalData(data []byte, v interface{}) error {
	if len(data) == 0 {
		return nil
	}
	return json.Unmarshal(data, v)
}
