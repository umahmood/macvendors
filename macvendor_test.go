package macvendors

import (
	"io/ioutil"
	"net/http"
	"strings"
	"testing"
)

type mock struct {
	doFunc func(req *http.Request) (*http.Response, error)
}

func (m *mock) lookup(macAddress string) (*MAC, error) {
	return lookupHelper(m, macAddress)
}

func (m *mock) name(macAddress string) (vendorName string, err error) {
	return venorNameHelper(m, macAddress)
}

func (m *mock) do(req *http.Request) (*http.Response, error) {
	return m.doFunc(req)
}

// macsAreEqual compares two MAC instances
func macsAreEqual(a, b *MAC) (bool, string) {
	if a.Address != b.Address {
		return false, "Address"
	}
	if a.Company != b.Company {
		return false, "Company"
	}
	if a.Country != b.Country {
		return false, "Country"
	}
	if a.MacPrefix != b.MacPrefix {
		return false, "MacPrefix"
	}
	if a.StartHex != b.StartHex {
		return false, "StartHex"
	}
	if a.EndHex != b.EndHex {
		return false, "EndHex"
	}
	if a.Type != b.Type {
		return false, "Type"
	}
	return true, ""
}

func TestLookup(t *testing.T) {
	mock := &mock{
		doFunc: func(req *http.Request) (*http.Response, error) {
			body := `{"result":{"company":"Microsoft Corporation","mac_prefix":"28:18:78","address":"One Microsoft Way,Redmond  Washington  98052-6399,US","start_hex":"281878000000","end_hex":"281878FFFFFF","country":"US","type":"MA-L"}}`
			return &http.Response{
				StatusCode: 200,
				Body:       ioutil.NopCloser(strings.NewReader(body)),
			}, nil
		},
	}

	wantMAC := &MAC{
		&Vendor{
			Address:   "One Microsoft Way,Redmond  Washington  98052-6399,US",
			Company:   "Microsoft Corporation",
			Country:   "US",
			MacPrefix: "28:18:78",
			StartHex:  "281878000000",
			EndHex:    "281878FFFFFF",
			Type:      "MA-L",
		},
	}

	testMAC := "28:18:78:6D:64:42"

	gotMAC, err := mock.lookup(testMAC)
	if err != nil {
		t.Errorf("Fail: %v", err)
	}

	if equal, diff := macsAreEqual(gotMAC, wantMAC); !equal {
		t.Errorf("Fail: macs are not equal, diff field %s", diff)
	}
}

func TestVendorName(t *testing.T) {
	mock := &mock{
		doFunc: func(req *http.Request) (*http.Response, error) {
			body := "Microsoft Corporation"
			return &http.Response{
				StatusCode: 200,
				Body:       ioutil.NopCloser(strings.NewReader(body)),
			}, nil
		},
	}

	wantVendorName := "Microsoft Corporation"
	testMAC := "28:18:78:6D:64:42"

	gotVendorName, err := mock.name(testMAC)
	if err != nil {
		t.Errorf("Fail: %v", err)
	}

	if gotVendorName != wantVendorName {
		t.Errorf("Fail: got vendor name %s want %s", gotVendorName, wantVendorName)
	}
}

func TestBadMACVendorName(t *testing.T) {
	mock1 := &mock{
		doFunc: func(req *http.Request) (*http.Response, error) {
			body := "Please provide mac address"
			return &http.Response{
				StatusCode: 200,
				Body:       ioutil.NopCloser(strings.NewReader(body)),
			}, nil
		},
	}
	wantErr := ErrMACFormat
	testMAC := "xx"
	_, gotErr := mock1.name(testMAC)
	if gotErr != wantErr {
		t.Errorf("fail name got err %v want err %v", gotErr, wantErr)
	}
}

func TestBadMACLookup(t *testing.T) {
	mock := &mock{
		doFunc: func(req *http.Request) (*http.Response, error) {
			body := `{"result": {"error": "no result"}}`
			return &http.Response{
				StatusCode: 200,
				Body:       ioutil.NopCloser(strings.NewReader(body)),
			}, nil
		},
	}
	wantErr := ErrMACFormat
	testMAC := "xx"
	_, gotErr := mock.lookup(testMAC)
	if gotErr != wantErr {
		t.Errorf("fail lookup got err %v want err %v", gotErr, wantErr)
	}
}
