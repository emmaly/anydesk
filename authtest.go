package anydesk

import (
	"encoding/json"
	"errors"
)

// AuthTest tests to see if auth is working
func (a *AnyDesk) AuthTest() (*GenericResult, error) {
	req, err := a.makeRequest("GET", "auth", nil, nil)
	if err != nil {
		return nil, err
	}
	resp, err := a.httpClient.Do(req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()
	var data *GenericResult
	j := json.NewDecoder(resp.Body)
	err = j.Decode(&data)
	if err != nil {
		return nil, err
	}
	if resp.StatusCode < 200 || resp.StatusCode >= 300 {
		data.Success = false
		return data, errors.New(resp.Status)
	}
	data.Success = true
	return data, nil
}
