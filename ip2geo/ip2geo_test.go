package ip2geo

import (
	"net"
	"testing"
)

func TestUnnamed0(t *testing.T) {
	reader, err := NewReader("")
	if err != nil {
		t.Fatal(err)
	}
	geo, err := reader.ReadIPInfo("27.69.27.62")
	if err != nil {
		t.Error(err)
	}
	t.Logf("geo: %#v", geo)
	if geo.CountryCode != "VN" {
		t.Errorf("error wrong CountryCode: %v", geo.CountryCode)
	}
	if geo.ISPName != "Viettel Group" {
		t.Errorf("error wrong ISPName: %v", geo.ISPName)
	}

	_, err = reader.ReadIPInfo("127.0.0.1")
	if err == nil {
		t.Error(err)
	}
}

func TestLookupIPFromHost(t *testing.T) {
	ip := LookupIPFromHost("lichess.org")
	if net.ParseIP(ip) == nil {
		t.Error("error unexpected nil IP")
	}
	ip1 := LookupIPFromHost("lichess.org11111")
	//t.Logf(ip1)
	if net.ParseIP(ip1) != nil {
		t.Error("error unexpected IP, should be err no such host")
	}
}
