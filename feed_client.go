package main

import (
	"fmt"
	"net"
	"net/http"
	"net/netip"
	"syscall"
	"time"
)

var feedHTTPClient = newFeedHTTPClient()

var restrictedFeedNetworks = []netip.Prefix{
	netip.MustParsePrefix("0.0.0.0/8"),
	netip.MustParsePrefix("100.64.0.0/10"),
	netip.MustParsePrefix("192.0.0.0/24"),
	netip.MustParsePrefix("192.0.2.0/24"),
	netip.MustParsePrefix("192.88.99.0/24"),
	netip.MustParsePrefix("198.18.0.0/15"),
	netip.MustParsePrefix("198.51.100.0/24"),
	netip.MustParsePrefix("203.0.113.0/24"),
	netip.MustParsePrefix("240.0.0.0/4"),
	netip.MustParsePrefix("2001::/23"),
	netip.MustParsePrefix("2001:db8::/32"),
	netip.MustParsePrefix("2002::/16"),
	netip.MustParsePrefix("3fff::/20"),
}

var publicIPv6Network = netip.MustParsePrefix("2000::/3")

func newFeedHTTPClient() *http.Client {
	transport := http.DefaultTransport.(*http.Transport).Clone()
	// A proxy could resolve the destination itself and bypass the dialer's address check.
	transport.Proxy = nil
	dialer := &net.Dialer{
		Timeout:   30 * time.Second,
		KeepAlive: 30 * time.Second,
		Control:   checkFeedAddress,
	}
	transport.DialContext = dialer.DialContext
	return &http.Client{Transport: transport}
}

// The dialer supplies the resolved address before opening each connection,
// including connections made while following redirects.
func checkFeedAddress(_, address string, _ syscall.RawConn) error {
	target, err := netip.ParseAddrPort(address)
	if err != nil {
		return fmt.Errorf("parse feed destination: %w", err)
	}
	ip := target.Addr().Unmap()
	if !isPublicFeedAddress(ip) {
		return fmt.Errorf("feed destination %s is not a public IP address", ip)
	}
	return nil
}

func isPublicFeedAddress(ip netip.Addr) bool {
	if !ip.IsGlobalUnicast() || ip.IsPrivate() || (ip.Is6() && !publicIPv6Network.Contains(ip)) {
		return false
	}
	for _, network := range restrictedFeedNetworks {
		if network.Contains(ip) {
			return false
		}
	}
	return true
}
