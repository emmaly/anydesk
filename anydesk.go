package anydesk

import (
	"crypto/hmac"
	"crypto/sha1"
	"encoding/base64"
	"encoding/json"
	"errors"
	"fmt"
	"net/http"
	"net/url"
	"regexp"
	"strconv"
	"strings"
	"time"
)

const (
	anyDeskBaseURL  = "https://v1.api.anydesk.com:8081"
	userAgentString = "https://github.com/emmaly/anydesk"
)

// Errors
var (
	ErrMissingAPIKey    = errors.New("missing API Key")
	ErrMissingLicenseID = errors.New("missing License ID")
	ErrBadResourceField = errors.New("bad resource field")
)

// AnyDesk is an AnyDesk API client
type AnyDesk struct {
	apiKey     string
	licenseID  string
	httpClient *http.Client
	userAgent  string
	baseURL    string
}

// Options are optional options for an AnyDesk API client
type Options struct {
	HTTPClient *http.Client
	UserAgent  string
	BaseURL    string
}

// GenericResult is a generic result struct
type GenericResult struct {
	Success     bool
	Error       string
	Code        string
	Method      string
	Resource    string
	RequestTime string `json:"request-time"`
	ContentHash string `json:"content-hash"`
	Result      string
	LicenseID   string
}

// Sort constants
var (
	SortClientID     = "cid"
	SortAlias        = "alias"
	SortOnline       = "online"
	SortClientIDFrom = "from.cid"
	SortClientIDTo   = "to.cid"
	SortTimeStart    = "start-time"
	SortTimeEnd      = "end-time"
	SortDuration     = "duration"
)

// Direction constants
const (
	DirectionAny = ""
	DirectionIn  = "in"
	DirectionOut = "out"
)

func orderString(b bool) string {
	if b {
		return "desc"
	}
	return "asc"
}

// New returns a new AnyDesk API client
func New(apiKey, licenseID string, o *Options) (*AnyDesk, error) {
	if apiKey == "" {
		return nil, ErrMissingAPIKey
	}

	if licenseID == "" {
		return nil, ErrMissingLicenseID
	}

	if o == nil {
		o = &Options{}
	}

	a := &AnyDesk{
		apiKey:     apiKey,
		licenseID:  licenseID,
		httpClient: o.HTTPClient,
		userAgent:  o.UserAgent,
		baseURL:    o.BaseURL,
	}

	if a.httpClient == nil {
		a.httpClient = &http.Client{
			Timeout: time.Second * 5,
		}
	}

	if a.userAgent == "" {
		a.userAgent = userAgentString
	}

	if a.baseURL == "" {
		a.baseURL = anyDeskBaseURL
	} else {
		a.baseURL = strings.TrimRight(a.baseURL, "/")
	}

	return a, nil
}

var cleanupResourceQuery = regexp.MustCompile("(^|&)(online)=[^&]*(&|$)")

func (a *AnyDesk) makeRequest(method, resource string, query *url.Values, data interface{}) (*http.Request, error) {
	timestamp := time.Now().Unix()

	// build the URL
	resource = "/" + strings.TrimLeft(resource, "/")
	if query != nil {
		q := query.Encode()
		q = cleanupResourceQuery.ReplaceAllString(q, "$1$2$3")
		if q != "" {
			resource += "?" + q
		}
	}

	// build the body
	var contentType string
	var body string
	if data != nil {
		bodyBytes, err := json.Marshal(data)
		if err != nil {
			return nil, err
		}
		body = string(bodyBytes)
		contentType = "application/json"
	}

	// build the bodyHash
	bh := sha1.New()
	bh.Write([]byte(body))
	bodyHash := base64.StdEncoding.EncodeToString(bh.Sum(nil))

	// build the requestString
	requestString := fmt.Sprintf("%s\n%s\n%d\n%s", method, resource, timestamp, bodyHash)

	// sign the requestString with the apiKey to generate a token
	h := hmac.New(sha1.New, []byte(a.apiKey))
	h.Write([]byte(requestString))
	token := base64.StdEncoding.EncodeToString(h.Sum(nil))

	// build the Authorization header
	authHeader := fmt.Sprintf("AD %s:%d:%s", a.licenseID, timestamp, token)

	// build and return the request
	req, err := http.NewRequest(method, a.baseURL+resource, strings.NewReader(body))
	if err != nil {
		return nil, err
	}
	req.Header.Set("User-Agent", a.userAgent)
	req.Header.Set("Authorization", authHeader)
	if contentType != "" {
		req.Header.Set("Content-Type", contentType)
	}
	return req, nil
}

func makeResource(r ...interface{}) string {
	part := make([]string, len(r))
	for i, v := range r {
		switch v.(type) {
		case string:
			part[i] = v.(string)
		case int:
			part[i] = strconv.Itoa(v.(int))
		}
	}
	return strings.Join(part, "/")
}
