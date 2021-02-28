package util

import "net/url"

func extractValuesFromArrayMap(data []map[string]string, key string) []string {
	keys := []string{}
	for _, object := range data {
		keys = append(keys, object[key])
	}
	return keys
}

// ConvertQueryParamsToMapString - convert url query params to map
func ConvertQueryParamsToMapString(params url.Values) map[string]string {
	values := map[string]string{}
	for key, val := range params {
		values[key] = val[0]
	}
	return values
}
