package router

import (
	"net/netip"
	"strconv"

	"github.com/xtls/xray-core/common/net"
	"go4.org/netipx"
)

type GeoIPMatcher struct {
	countryCode  string
	reverseMatch bool
	ip4          *netipx.IPSet
	ip6          *netipx.IPSet
}

func (m *GeoIPMatcher) Init(cidrs []*CIDR) error {
	var builder4, builder6 netipx.IPSetBuilder

	for _, cidr := range cidrs {
		ip := net.IP(cidr.GetIp())
		ipPrefixString := ip.String() + "/" + strconv.Itoa(int(cidr.GetPrefix()))
		ipPrefix, err := netip.ParsePrefix(ipPrefixString)
		if err != nil {
			return err
		}

		switch len(ip) {
		case net.IPv4len:
			builder4.AddPrefix(ipPrefix)
		case net.IPv6len:
			builder6.AddPrefix(ipPrefix)
		}
	}

	if ip4, err := builder4.IPSet(); err != nil {
		return err
	} else {
		m.ip4 = ip4
	}

	if ip6, err := builder6.IPSet(); err != nil {
		return err
	} else {
		m.ip6 = ip6
	}

	return nil
}

func (m *GeoIPMatcher) SetReverseMatch(isReverseMatch bool) {
	m.reverseMatch = isReverseMatch
}

func (m *GeoIPMatcher) match4(ip net.IP) bool {
	nip, ok := netipx.FromStdIP(ip)
	if !ok {
		return false
	}

	return m.ip4.Contains(nip)
}

func (m *GeoIPMatcher) match6(ip net.IP) bool {
	nip, ok := netipx.FromStdIP(ip)
	if !ok {
		return false
	}

	return m.ip6.Contains(nip)
}

func (m *GeoIPMatcher) GetCountryCode() string {
	return m.countryCode
}

// Match returns true if the given ip is included by the GeoIP.
func (m *GeoIPMatcher) Remove(ips []string) (error, bool) {
	ip4 := netipx.IPSetBuilder{}
	ip6 := netipx.IPSetBuilder{}
	ip4.AddSet(m.ip4)
	ip6.AddSet(m.ip6)
	isMatched4 := false
	isMatched6 := false
	for _, rip := range ips {

		nip := net.ParseIP(rip)
		ip, ok := netipx.FromStdIP(nip)
		if !ok {
			continue
		}
		
		if ip.Is4() && m.match4(nip) {
			ip4.Remove(ip)
			isMatched4 = true
		} else if ip.Is6() && m.match6(nip) {
			ip6.Remove(ip)
			isMatched6 = true
		}
	}

	if isMatched4 {
		ip4 , err := ip4.IPSet()
		if err != nil {
			return err, true
		}
		m.ip4 = ip4
	}
	if isMatched6 {
		ip6 , err := ip6.IPSet()
		if err != nil {
			return err, true
		}
		m.ip6 = ip6
	}
	return nil, true
}

// Match returns true if the given ip is included by the GeoIP.
func (m *GeoIPMatcher) Add(ips []string) (error, bool) {
	ip4 := netipx.IPSetBuilder{}
	ip6 := netipx.IPSetBuilder{}
	ip4.AddSet(m.ip4)
	ip6.AddSet(m.ip6)
	is4 := false
	is6 := false
	for _, rip := range ips {

		nip := net.ParseIP(rip)
		ip, ok := netipx.FromStdIP(nip)
		if !ok {
			continue
		}
		
		if ip.Is4(){
			ip4.Add(ip)
			is4 = true
		} else if ip.Is6(){
			ip6.Add(ip)
			is6 = true
		}
	}

	if is4 {
		ip4 , err := ip4.IPSet()
		if err != nil {
			return err, true
		}
		m.ip4 = ip4
	}
	if is6 {
		ip6 , err := ip6.IPSet()
		if err != nil {
			return err, true
		}
		m.ip6 = ip6
	}
	return nil, true
}

// Match returns true if the given ip is included by the GeoIP.
func (m *GeoIPMatcher) Match(ip net.IP) bool {
	isMatched := false
	switch len(ip) {
	case net.IPv4len:
		isMatched = m.match4(ip)
	case net.IPv6len:
		isMatched = m.match6(ip)
	}
	if m.reverseMatch {
		return !isMatched
	}
	return isMatched
}

// GeoIPMatcherContainer is a container for GeoIPMatchers. It keeps unique copies of GeoIPMatcher by country code.
type GeoIPMatcherContainer struct {
	matchers []*GeoIPMatcher
}

// Add adds a new GeoIP set into the container.
// If the country code of GeoIP is not empty, GeoIPMatcherContainer will try to find an existing one, instead of adding a new one.
func (c *GeoIPMatcherContainer) Add(geoip *GeoIP) (*GeoIPMatcher, error) {
	if len(geoip.CountryCode) > 0 {
		for _, m := range c.matchers {
			if m.countryCode == geoip.CountryCode && m.reverseMatch == geoip.ReverseMatch {
				return m, nil
			}
		}
	}

	m := &GeoIPMatcher{
		countryCode:  geoip.CountryCode,
		reverseMatch: geoip.ReverseMatch,
	}
	if err := m.Init(geoip.Cidr); err != nil {
		return nil, err
	}
	if len(geoip.CountryCode) > 0 {
		c.matchers = append(c.matchers, m)
	}
	return m, nil
}

var globalGeoIPContainer GeoIPMatcherContainer
