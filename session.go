package anydesk

import (
	"encoding/json"
	"errors"
	"net/url"
	"strconv"
	"time"
)

// Sessions list all sessions regarding a license's clients
type Sessions struct {
	Count     int
	Selected  int
	Offset    int
	Limit     int
	ClientID  int
	Direction string
	Sessions  []Session `json:"list"`
}

// Session is detailed information about an individual session
type Session struct {
	ID           string `json:"sid"`
	ClientIDFrom Client
	ClientIDTo   Client
	Active       bool
	TimeStart    int `json:"start-time"` // time.Time?
	TimeEnd      int `json:"end-time"`   // time.Time?
	Duration     int
	Comment      string
}

// SessionsOptions are options for session list query
type SessionsOptions struct {
	ClientID   int       // If set it will limit returned sessions to that client.
	Direction  string    // `anydesk.DirectionAny` (default), `anydesk.DirectionIn`, or `anydesk.DirectionOut`
	TimeAfter  time.Time // Only return sessions after this time.
	TimeBefore time.Time // Only return sessions before this time.
	Offset     int       // The index of the first item to be returned.
	Limit      int       // Defaults to unlimited.
	Sort       string    // `anydesk.SortClientIDFrom`, `anydesk.SortClientIDTo`, `anydesk.SortTimeStart`, `anydesk.SortTimeEnd`, or `anydesk.SortDuration`
	Order      bool      // `false` is descending (default), `true` is ascending
}

// Sessions gets a list of sessions regarding a license's clients, or an individual client
func (a *AnyDesk) Sessions(opts *SessionsOptions) (*Sessions, error) {
	q := make(url.Values)
	if opts != nil {
		if opts.ClientID > 0 {
			q.Set("cid", strconv.Itoa(opts.ClientID))
		}
		if opts.Direction != "" {
			q.Set("direction", opts.Direction)
		}
		if !opts.TimeAfter.IsZero() {
			q.Set("to", strconv.FormatInt(opts.TimeAfter.Unix(), 10))
		}
		if !opts.TimeBefore.IsZero() {
			q.Set("from", strconv.FormatInt(opts.TimeBefore.Unix(), 10))
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
	req, err := a.makeRequest("GET", "sessions", &q, nil)
	if err != nil {
		return nil, err
	}
	resp, err := a.httpClient.Do(req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()
	var data *Sessions
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

// Session returns data about a specific session
func (a *AnyDesk) Session(id int) (*Session, error) {
	req, err := a.makeRequest("GET", makeResource("sessions", id), nil, nil)
	if err != nil {
		return nil, err
	}
	resp, err := a.httpClient.Do(req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()
	var data *Session
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

// SessionClose closes the specified currently open session
func (a *AnyDesk) SessionClose(id int) error {
	req, err := a.makeRequest("POST", makeResource("sessions", id, "action"), nil, map[string]string{
		"action": "close",
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

// SessionComment sets the comment on a session, overwriting the existing comment if there was one
func (a *AnyDesk) SessionComment(id int, comment string) error {
	var text *string
	if comment != "" {
		text = &comment
	}
	req, err := a.makeRequest("PATCH", makeResource("sessions", id), nil, map[string]*string{
		"comment": text, // if it's null, then it will erase the comment instead of just being empty
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
