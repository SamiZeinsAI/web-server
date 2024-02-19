package database

import (
	"encoding/json"
	"net/http"
)

func UnmarshalRequestBody(r *http.Request, params interface{}) error {
	err := json.NewDecoder(r.Body).Decode(&params)
	if err != nil {
		return err
	}
	return nil
}
