package util

import (
	"fmt"
	"io/ioutil"
	"net/http"
	"net/url"
	"strings"
)

func CallAPI(method, urlStr string, queryParams map[string]string) ([]byte, error) {
	if len(queryParams) > 0 {
		query := url.Values{}
		for key, value := range queryParams {
			query.Add(key, value)
		}
		urlStr += "?" + query.Encode()
	}

	req, err := http.NewRequest(strings.ToUpper(method), urlStr, nil)
	if err != nil {
		return nil, fmt.Errorf("error creating request: %v", err)
	}

	client := http.DefaultClient
	resp, err := client.Do(req)
	if err != nil {
		return nil, fmt.Errorf("error sending request: %v", err)
	}
	defer resp.Body.Close()

	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return nil, fmt.Errorf("error reading response body: %v", err)
	}

	return body, nil
}
