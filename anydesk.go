package anydesk

import (
	"crypto/hmac"
	"crypto/sha1"
	"encoding/base64"
	"errors"
	"fmt"
	"net/http"
	"strings"
	"time"
)

const anyDeskBaseURL = "https://v1.api.anydesk.com:8081"

// Errors
var (
	ErrMissingAPIKey    = errors.New("missing API Key")
	ErrMissingLicenseID = errors.New("missing License ID")
)

// AnyDesk is an AnyDesk API client
type AnyDesk struct {
	apiKey     string
	licenseID  string
	httpClient *http.Client
	baseURL    string
}

// Options are optional options for an AnyDesk API client
type Options struct {
	HTTPClient *http.Client
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
		baseURL:    o.BaseURL,
	}

	if a.baseURL == "" {
		a.baseURL = anyDeskBaseURL
	}

	if a.httpClient == nil {
		a.httpClient = &http.Client{
			Timeout: time.Second * 5,
		}
	}

	return a, nil
}

func (a *AnyDesk) makeRequest(method, resource, body string) (*http.Request, error) {
	timestamp := time.Now().Unix()

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
	req.Header.Set("Authorization", authHeader)
	return req, nil
}
