package conv

import "encoding/json"

func StructToMap(i interface{}) map[string]interface{} {
	m := make(map[string]interface{})
	j, _ := json.Marshal(i)
	_ = json.Unmarshal(j, &m)
	return m
}
