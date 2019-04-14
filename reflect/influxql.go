package main

import (
	"encoding/json"
	"errors"
)

// ParseIqlResults 解析influxdb查询结果，返回单个结果集
func ParseIqlResults(jsonData string, f func(map[string]interface{})) error {
	//	Partly JSON unmarshal into a map in Go.
	//	参考：https://stackoverflow.com/questions/11066946/partly-json-unmarshal-into-a-map-in-go
	var m map[string]*json.RawMessage
	err := json.Unmarshal([]byte(jsonData), &m)
	if err != nil {
		message, _ := json.Marshal(err.Error())
		return errors.New(string(message))
	}

	errMsg, existsError := m["error"]
	if existsError {
		return UnmarshalError(errMsg)
	}

	var results []*json.RawMessage
	_ = json.Unmarshal(*m["results"], &results)

	for _, r := range results {
		var result map[string]*json.RawMessage
		_ = json.Unmarshal(*r, &result)

		errMsg, existsError = result["error"]
		if existsError {
			return UnmarshalError(errMsg)
		}

		series, seriesOK := result["series"]
		if !seriesOK {
			continue
		}

		var seriesRows []*json.RawMessage
		_ = json.Unmarshal(*series, &seriesRows)

		for _, seriesRow := range seriesRows {
			var m map[string]*json.RawMessage
			if err := json.Unmarshal(*seriesRow, &m); err != nil {
				return err
			}

			var tags map[string]string
			json.Unmarshal(*m["tags"], &tags)

			var columns []string
			if err := json.Unmarshal(*m["columns"], &columns); err != nil {
				return err
			}

			var seriesValues []*json.RawMessage
			_ = json.Unmarshal(*m["values"], &seriesValues)

			for _, seriesValue := range seriesValues {
				var values []*json.RawMessage
				_ = json.Unmarshal(*seriesValue, &values)

				result := make(map[string]interface{})

				for tagKey, tagVal := range tags {
					result[tagKey] = tagVal
				}

				for i, row := range values {
					key := columns[i]
					var v interface{}
					json.Unmarshal(*row, &v)
					result[key] = v
				}

				f(result)
			}
		}
	}

	return nil
}
