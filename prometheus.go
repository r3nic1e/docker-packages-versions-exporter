package main

import (
	"encoding/json"
	"net/http"
	"net/url"
)

type promResult struct {
	Metric map[string]string
	Value  []interface{}
}

type promResults struct {
	Status string
	Data   struct {
		ResultType string
		Result     []promResult
	}
}

func getRunningImages(promURL *url.URL) ([]string, error) {
	query := "max by (image) (container_start_time_seconds)"

	req, err := http.NewRequest("GET", promURL.String(), nil)
	if err != nil {
		return []string{}, err
	}

	q := req.URL.Query()
	q.Add("query", query)
	req.URL.RawQuery = q.Encode()

	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		return []string{}, err
	}
	defer resp.Body.Close()

	decoder := json.NewDecoder(resp.Body)
	var res promResults
	decoder.Decode(&res)

	var result []string

	for _, res := range res.Data.Result {
		result = append(result, res.Metric["image"])
	}

	return result, nil
}
