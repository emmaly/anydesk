package anydesk

import (
	"encoding/json"
	"errors"
	"fmt"
)

// Clients is a response for a Clients method call
type Clients struct {
	Count      int
	Selected   int
	Offset     int
	Limit      int
	OnlineOnly bool     `json:"online"`
	Clients    []Client `json:"list"`
}

// Client is an individual AnyDesk-powered client attached to this license
type Client struct {
	ID             int `json:"cid"`
	Alias          string
	ClientVersion  string `json:"client-version"`
	Online         bool
	OnlineTime     int       `json:"online-time"`
	RecentSessions []Session `json:"last-sessions"`
}

// Clients gets a list of individual AnyDesk-powered clients attached to this license
// IncludeOffline is a bool
// Limit <= 0 means unlimited
// Sort is one of these strings:  "cid", "alias", "online"
// Order is one of these variables: `anydesk.SortAsc` or `anydesk.SortDesc`
func (a *AnyDesk) Clients(includeOffline bool, offset, limit int, sort string, order bool) (*Clients, error) {
	onlineOnlyStr := ""
	if !includeOffline {
		onlineOnlyStr = "&online"
	}
	if limit <= 0 {
		limit = -1
	}
	req, err := a.makeRequest("GET", fmt.Sprintf("/clients?offset=%d&limit=%d&sort=%s&order=%s%s", offset, limit, sort, orderString(order), onlineOnlyStr), "")
	if err != nil {
		return nil, err
	}
	resp, err := a.httpClient.Do(req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()
	var data *Clients
	j := json.NewDecoder(resp.Body)
	err = j.Decode(&data)
	if err != nil {
		return nil, err
	}
	if resp.StatusCode != 200 {
		return data, errors.New(resp.Status)
	}
	return data, nil
}

// Client returns data about a specific client
func (a *AnyDesk) Client(id int) (*Client, error) {
	req, err := a.makeRequest("GET", fmt.Sprintf("/clients/%d", id), "")
	if err != nil {
		return nil, err
	}
	resp, err := a.httpClient.Do(req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()
	var data *Client
	j := json.NewDecoder(resp.Body)
	err = j.Decode(&data)
	if err != nil {
		return nil, err
	}
	if resp.StatusCode != 200 {
		return data, errors.New(resp.Status)
	}
	return data, nil
}
