package utils

import "encoding/json"

func PrettyJSON(a any) string {
	jsonB, err := json.MarshalIndent(a, "", "  ")
	if err != nil {
		return err.Error()
	}
	return string(jsonB)
}
