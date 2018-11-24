package anydesk

import (
	"encoding/json"
	"errors"
	"net/url"
	"strconv"
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

// ClientsOptions are options for client list query
type ClientsOptions struct {
	IncludeOffline bool   // Default is false.
	Offset         int    // The index of the first item to be returned.
	Limit          int    // Defaults to unlimited.
	Sort           string // `anydesk.SortClientID`, `anydesk.SortAlias`, or `anydesk.SortOnline`
	Order          bool   // `false` is descending (default), `true` is ascending
}

// Clients gets a list of individual AnyDesk-powered clients attached to this license
func (a *AnyDesk) Clients(opts *ClientsOptions) (*Clients, error) {
	q := make(url.Values)
	if opts != nil {
		if !opts.IncludeOffline {
			q.Set("online", "true")
		}
		if opts.Offset > 0 {
			q.Set("offset", strconv.Itoa(opts.Offset))
		}
		if opts.Limit > 0 {
			q.Set("limit", strconv.Itoa(opts.Limit))
		}
		if opts.Sort != "" {
			q.Set("sort", opts.Sort)
		}
		if opts.Order {
			q.Set("order", "asc")
		} else {
			q.Set("order", "desc")
		}
	}
	req, err := a.makeRequest("GET", "clients", &q, nil)
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
	if resp.StatusCode < 200 || resp.StatusCode >= 300 {
		return data, errors.New(resp.Status)
	}
	return data, nil
}

// Client returns data about a specific client
func (a *AnyDesk) Client(id int) (*Client, error) {
	req, err := a.makeRequest("GET", makeResource("clients", id), nil, nil)
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
	if resp.StatusCode < 200 || resp.StatusCode >= 300 {
		return data, errors.New(resp.Status)
	}
	return data, nil
}

// ClientAlias sets the alias for a specific client
func (a *AnyDesk) ClientAlias(id int, alias string) error {
	var value *string
	if alias != "" {
		value = &alias
	}
	req, err := a.makeRequest("PATCH", makeResource("clients", id), nil, map[string]*string{
		"alias": value, // if it's null, then it will erase the alias instead of just being empty
	})
	if err != nil {
		return err
	}
	resp, err := a.httpClient.Do(req)
	if err != nil {
		return err
	}
	defer resp.Body.Close()
	if resp.StatusCode < 200 || resp.StatusCode >= 300 {
		return errors.New(resp.Status)
	}
	return nil
}
