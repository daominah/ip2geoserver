package ip2geo

import (
	_ "embed"
	"fmt"
	"net"
	"os"
	"path/filepath"
	"strings"

	geoip2 "github.com/oschwald/geoip2-golang"
)

type Reader struct {
	CityDatabase *geoip2.Reader
	ASNDatabase  *geoip2.Reader
}

// NewReader loads IP info database (MaxMind) from files in dir dataPath,
// DEPRECATED
func NewReader(dataPath string) (*Reader, error) {
	var err error
	ret := &Reader{}
	ret.CityDatabase, err = geoip2.Open(
		filepath.Join(dataPath, "GeoLite2-City.mmdb"))
	if err != nil {
		return nil, err
	}
	ret.ASNDatabase, err = geoip2.Open(
		filepath.Join(dataPath, "GeoLite2-ASN.mmdb"))
	if err != nil {
		return nil, err
	}
	return ret, nil
}

//go:embed GeoLite2-City.mmdb
var DatabaseCity []byte

//go:embed GeoLite2-ASN.mmdb
var DatabaseASN []byte

// NewReader loads IP info database (MaxMind) with Go embed
func NewDefaultReader() (*Reader, error) {
	var err error
	ret := &Reader{}

	tmpFileCity, err := os.CreateTemp("", "tmpFileCity")
	if err != nil {
		return nil, err
	}
	_, err = tmpFileCity.Write(DatabaseCity)
	if err != nil {
		return nil, err
	}
	tmpFileASN, err := os.CreateTemp("", "tmpFileASN")
	if err != nil {
		return nil, err
	}
	_, err = tmpFileASN.Write(DatabaseASN)
	if err != nil {
		return nil, err
	}

	ret.CityDatabase, err = geoip2.Open(tmpFileCity.Name())
	if err != nil {
		return nil, err
	}
	ret.ASNDatabase, err = geoip2.Open(tmpFileASN.Name())
	if err != nil {
		return nil, err
	}

	os.Remove(tmpFileCity.Name())
	os.Remove(tmpFileASN.Name())

	return ret, nil
}

func (r Reader) ReadIPInfo(addrIP string) (IPInfo, error) {
	ret := IPInfo{IP: addrIP}
	ip0 := net.ParseIP(addrIP)
	if ip0 == nil {
		return ret, fmt.Errorf("invalid IP format")
	}
	if CheckIsPrivateIP(ip0) {
		return ret, fmt.Errorf("input is a private IP")
	}
	row, err := r.CityDatabase.City(ip0)
	if err != nil {
		return ret, fmt.Errorf("read city: %v", err)
	}
	ret.Continent = row.Continent.Names[LangEN]
	ret.ContinentCode = row.Continent.Code
	ret.Country = row.Country.Names[LangEN]
	ret.CountryCode = row.Country.IsoCode
	ret.City = row.City.Names[LangEN]
	ret.TimeZoneName = row.Location.TimeZone
	row2, err := r.ASNDatabase.ASN(ip0)
	if err != nil {
		return ret, fmt.Errorf("read ASN: %v", err)
	}
	ret.ISPName = row2.AutonomousSystemOrganization
	return ret, nil
}

type IPInfo struct {
	IP            string
	Continent     string
	ContinentCode string
	Country       string
	CountryCode   string
	City          string
	TimeZoneName  string
	ISPName       string
}

const LangEN = "en" // language english

// Continent codes
const (
	ContinentAfrica       = "AF"
	ContinentAntarctica   = "AN"
	ContinentAsia         = "AS"
	ContinentEurope       = "EU"
	ContinentNorthAmerica = "NA"
	ContinentOceania      = "OC"
	ContinentSouthAmerica = "SA"
)

var privateIPBlocks []*net.IPNet

func init() {
	for _, cidr := range []string{
		"127.0.0.0/8",    // IPv4 loopback
		"10.0.0.0/8",     // RFC1918
		"172.16.0.0/12",  // RFC1918
		"192.168.0.0/16", // RFC1918
		"169.254.0.0/16", // RFC3927 link-local
		"::1/128",        // IPv6 loopback
		"fe80::/10",      // IPv6 link-local
		"fc00::/7",       // IPv6 unique local addr
	} {
		_, block, err := net.ParseCIDR(cidr)
		if err != nil {
			fmt.Printf("error init privateIPBlocks: %v", err)
		}
		privateIPBlocks = append(privateIPBlocks, block)
	}
}

// CheckIsPrivateIP returns true if input is a private IP
func CheckIsPrivateIP(ip net.IP) bool {
	if ip.IsLoopback() || ip.IsLinkLocalUnicast() || ip.IsLinkLocalMulticast() {
		return true
	}
	for _, block := range privateIPBlocks {
		if block.Contains(ip) {
			return true
		}
	}
	return false
}

// GetIpFromAddress removes port from the address
func GetIpFromAddress(hostPort string) string {
	colonIdx := strings.Index(hostPort, ":")
	if colonIdx == -1 {
		return hostPort
	}
	return hostPort[:colonIdx]
}

// LookupIPFromHost returns the first IPv4.
// This func can return an IPv6 or an error string if not found.
func LookupIPFromHost(hostname string) string {
	addrs, err := net.LookupHost(hostname)
	if err != nil {
		return err.Error()
	}
	var ipv6 net.IP
	for _, addr := range addrs {
		ip := net.ParseIP(addr)
		if ip == nil {
			continue
		}
		if ip.To4() == nil {
			ipv6 = ip
			continue
		}
		return ip.String()
	}
	return ipv6.String()
}
