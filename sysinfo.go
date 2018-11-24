package anydesk

import (
	"encoding/json"
	"errors"
)

// SysInfo is general license info and AnyDesk system information
type SysInfo struct {
	Name       string
	APIVersion string `json:"api-ver"`
	License    struct {
		Name           string
		Expires        int
		MaxClients     int `json:"max-clients"`
		MaxSessions    int `json:"max-sessions"`
		MaxSessionTime int `json:"max-session-time"`
		Namespaces     []struct {
			Name string
			Size int
		}
		LicenseID   string `json:"license-id"`
		LicenseKey  string `json:"license-key"`
		APIPassword string `json:"api-password"`
	}
	Clients struct {
		Total  int
		Online int
	}
	Sessions struct {
		Total  int
		Online int
	}
	Standalone bool
}

// SysInfo gets general license info and AnyDesk system information
func (a *AnyDesk) SysInfo() (*SysInfo, error) {
	req, err := a.makeRequest("GET", "sysinfo", nil, nil)
	if err != nil {
		return nil, err
	}
	resp, err := a.httpClient.Do(req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()
	var data *SysInfo
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
