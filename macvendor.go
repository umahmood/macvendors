package macvendors

import (
	"bytes"
	"encoding/json"
	"errors"
	"fmt"
	"io/ioutil"
	"net/http"
)

// ErrMACFormat invalid formatted MAC address determined by macvendors.co
var ErrMACFormat = errors.New("invalid mac format supported formats hex ':' bit '-' dot '.'")

const (
	macEndPoint        = "https://macvendors.co/api/%s/json"
	vendorNameEndPoint = "https://macvendors.co/api/vendorname/%s"
)

// doer implemented by any type which can do HTTP requests
type doer interface {
	do(*http.Request) (*http.Response, error)
}

// API instance
type API struct{}

// New returns a new API instance
func New() *API {
	return &API{}
}

// Lookup vendor details by MAC address
func (m *API) Lookup(macAddress string) (*MAC, error) {
	return lookupHelper(&MAC{}, macAddress)
}

// Name lookup for vendor name only.
func (m *API) Name(macAddress string) (vendorName string, err error) {
	return venorNameHelper(&MAC{}, macAddress)
}

// Vendor instance
type Vendor struct {
	Address   string `json:"address"`
	Company   string `json:"company"`
	Country   string `json:"country"`
	Type      string `json:"type"`
	MacPrefix string `json:"mac_prefix"`
	StartHex  string `json:"start_hex"`
	EndHex    string `json:"end_hex"`
}

// MAC instance
type MAC struct {
	*Vendor `json:"result"`
}

// do performs a http request
func (m *MAC) do(req *http.Request) (*http.Response, error) {
	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		return nil, err
	}
	return resp, nil
}

// venorNameHelper helper for VendorName
func venorNameHelper(d doer, macAddress string) (string, error) {
	url := fmt.Sprintf(vendorNameEndPoint, macAddress)
	req, err := http.NewRequest("GET", url, nil)
	if err != nil {
		return "", err
	}
	req.Header.Set("User-Agent", "API Browser")
	rsp, err := d.do(req)
	if rsp != nil {
		defer rsp.Body.Close()
	}
	if err != nil {
		return "", err
	}
	body, err := ioutil.ReadAll(rsp.Body)
	if err != nil {
		return "", err
	}
	if bytes.Equal(body, []byte("Please provide mac address")) {
		return "", ErrMACFormat
	}
	return string(body), nil
}

// lookupHelper helper for Lookup
func lookupHelper(d doer, macAddress string) (*MAC, error) {
	url := fmt.Sprintf(macEndPoint, macAddress)
	req, err := http.NewRequest("GET", url, nil)
	if err != nil {
		return nil, err
	}
	req.Header.Set("User-Agent", "API Browser")
	rsp, err := d.do(req)
	if rsp != nil {
		defer rsp.Body.Close()
	}
	if err != nil {
		return nil, err
	}
	body, err := ioutil.ReadAll(rsp.Body)
	if err != nil {
		return nil, err
	}
	if bytes.Contains(body, []byte("error")) {
		return nil, ErrMACFormat
	}
	mac := &MAC{}
	err = json.Unmarshal(body, &mac)
	if err != nil {
		return nil, err
	}
	return mac, err
}
