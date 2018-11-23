package anydesk

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
	ID        string `json:"sid"`
	From      Client
	To        Client
	Active    bool
	StartTime int `json:"start-time"`
	EndTime   int `json:"end-time"`
	Duration  int
	Comment   string
}

// Sessions gets a list of sessions regarding a license's clients, or an individual client
// ClientID is optional; if set it will limit sessions to only that client
// Direction is one of these strings: "in" or "out"
// BeforeTime and AfterTime set the start/stop bounds on returned sessions
// Limit <= 0 means unlimited
// Sort is one of these strings:  "from.cid", "to.cid", "start-time", "end-time", "duration"
// Order is one of these variables: `anydesk.SortAsc` or `anydesk.SortDesc`
func (a *AnyDesk) Sessions() (*Sessions, error) {
	return nil, nil
}
