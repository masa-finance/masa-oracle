package network

import (
	"fmt"
	"io"
	"net"
	"net/http"
	"strings"

	"github.com/libp2p/go-libp2p/core/host"
	"github.com/libp2p/go-libp2p/core/peer"
	"github.com/multiformats/go-multiaddr"
	"github.com/sirupsen/logrus"
)

// The URL of the GCP metadata server for the external IP
const externalIPURL = "http://metadata.google.internal/computeMetadata/v1/instance/network-interfaces/0/access-configs/0/external-ip"

// GetMultiAddressesForHost returns the multiaddresses for the host
func GetMultiAddressesForHost(host host.Host) ([]multiaddr.Multiaddr, error) {
	peerInfo := peer.AddrInfo{
		ID:    host.ID(),
		Addrs: host.Addrs(),
	}
	multiaddrs, err := peer.AddrInfoToP2pAddrs(&peerInfo)
	if err != nil {
		return nil, err
	}
	logrus.Debug("Multiaddresses from AddrInfoToP2pAddrs: ", multiaddrs)
	addresses := make([]multiaddr.Multiaddr, 0)
	for _, addr := range multiaddrs {
		logrus.Debug(addr.String())
		// skip using localhost since it provides no value
		if !strings.Contains(addr.String(), "127.0.0.1") {
			addresses = append(addresses, addr)
		}
	}
	return addresses, nil
}

// GetMultiAddressesForHostQuiet returns the multiaddresses for the host without logging
// TODO rm
func GetMultiAddressesForHostQuiet(host host.Host) []multiaddr.Multiaddr {
	ma, err := GetMultiAddressesForHost(host)
	if err != nil {
		logrus.Fatal(err)
	}
	return ma
}

// getPublicMultiAddress returns the best public IP address (for some definition of "best")
// TODO: This is not guaranteed to work, and should not be necessary since we're using AutoNAT
func getPublicMultiAddress(addrs []multiaddr.Multiaddr) multiaddr.Multiaddr {
	ipBytes, err := Get("https://api.ipify.org?format=text", nil)
	externalIP := net.ParseIP(string(ipBytes))
	if err != nil {
		logrus.Warnf("[-] Failed to get public IP: %v", err)
		return nil
	}
	if externalIP == nil || externalIP.IsPrivate() {
		return nil
	}

	var addrToCopy multiaddr.Multiaddr
	if len(addrs) > 0 {
		addrToCopy = addrs[0]
	}
	publicMultiaddr, err := replaceIPComponent(addrToCopy, externalIP.String())
	if err != nil {
		logrus.Warnf("[-] Failed to create multiaddr with public IP: %v", err)
		return nil
	}
	return publicMultiaddr
}

// GetPriorityAddress returns the best public or private IP address
// TODO rm?
func GetPriorityAddress(addrs []multiaddr.Multiaddr) multiaddr.Multiaddr {
	var bestPrivateAddr multiaddr.Multiaddr
	bestPublicAddr := getPublicMultiAddress(addrs)

	for _, addr := range addrs {
		ipComponent, err := addr.ValueForProtocol(multiaddr.P_IP4)
		if err != nil {
			ipComponent, err = addr.ValueForProtocol(multiaddr.P_IP6)
			if err != nil {
				continue // Not an IP address
			}
		}
		ip := net.ParseIP(ipComponent)
		if ip == nil || ip.IsLoopback() {
			continue // Skip invalid or loopback addresses
		}
		if ip.IsPrivate() {
			// If it's the first private address found or if it's preferred over the current best, keep it
			if bestPrivateAddr == nil || isPreferredAddress(addr) {
				bestPrivateAddr = addr
			}
		} else {
			// If it's the first public address found or if it's preferred over the current best, keep it
			if bestPublicAddr == nil || isPreferredAddress(addr) {
				bestPublicAddr = addr
			}
		}
	}
	var baseAddr multiaddr.Multiaddr
	// Prefer public addresses over private ones
	if bestPublicAddr != nil {
		baseAddr = bestPublicAddr
	} else if bestPrivateAddr != nil {
		baseAddr = bestPrivateAddr
	} else {
		logrus.Warn("No address matches the priority criteria, returning the first entry")
		baseAddr = addrs[0]
	}
	logrus.Infof("Best public address: %s", bestPublicAddr)
	logrus.Debugf("Best private address: %s", bestPrivateAddr)
	logrus.Debugf("Base address: %s", baseAddr)
	gcpAddr := replaceGCPAddress(baseAddr)
	if gcpAddr != nil {
		baseAddr = gcpAddr
	}
	return baseAddr
}

func replaceGCPAddress(addr multiaddr.Multiaddr) multiaddr.Multiaddr {
	// After finding the best address, try to get the GCP external IP
	var err error
	var bestAddr multiaddr.Multiaddr
	gotExternalIP, externalIP := getGCPExternalIP()
	if gotExternalIP && externalIP != "" {
		bestAddr, err = replaceIPComponent(addr, externalIP)
		if err != nil {
			logrus.Warnf("Failed to replace IP component: %s", err)
			return nil
		}
	}
	logrus.Debug("Got external IP: ", gotExternalIP)
	logrus.Debug("Address after replacing IP component: ", bestAddr)
	return bestAddr
}

func replaceIPComponent(maddr multiaddr.Multiaddr, newIP string) (multiaddr.Multiaddr, error) {
	var components []multiaddr.Multiaddr
	for _, component := range multiaddr.Split(maddr) {
		if component.Protocols()[0].Code == multiaddr.P_IP4 || component.Protocols()[0].Code == multiaddr.P_IP6 {
			// Create a new IP component
			newIPComponent, err := multiaddr.NewMultiaddr(fmt.Sprintf("/ip4/%s", newIP))
			if err != nil {
				return nil, err
			}
			components = append(components, newIPComponent)
		} else {
			components = append(components, component)
		}
	}
	return multiaddr.Join(components...), nil
}

func getGCPExternalIP() (bool, string) {

	// Create a new HTTP client with a specific timeout
	client := &http.Client{}

	// Make a request to the metadata server
	req, err := http.NewRequest("GET", externalIPURL, nil)
	if err != nil {
		return false, ""
	}

	// GCP metadata server requires this specific header
	req.Header.Add("Metadata-Flavor", "Google")

	// Perform the request
	resp, err := client.Do(req)
	if err != nil {
		return false, ""
	}
	defer func(Body io.ReadCloser) {
		err := Body.Close()
		if err != nil {
			logrus.Error(err)
		}
	}(resp.Body)

	// Check if the metadata server returns a successful status code
	if resp.StatusCode != http.StatusOK {
		logrus.Debug("Metadata server response status: ", resp.StatusCode)
		return false, ""
	}
	//Read the external IP from the response
	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return true, ""
	}
	// Check that the response is a valid IP address
	if net.ParseIP(string(body)) == nil {
		return false, ""
	}
	logrus.Debug("External IP from metadata server: ", string(body))
	return true, string(body)
}

// isPreferredAddress checks if the multiaddress contains the UDP protocol
func isPreferredAddress(addr multiaddr.Multiaddr) bool {
	// Check if the multiaddress contains the UDP protocol
	for _, p := range addr.Protocols() {
		if p.Code == multiaddr.P_UDP {
			return true
		}
	}
	return false
}

// GetBootNodesMultiAddress returns the multiaddresses for the bootstrap nodes
func GetBootNodesMultiAddress(bootstrapNodes []string) ([]multiaddr.Multiaddr, error) {
	addrs := make([]multiaddr.Multiaddr, 0)
	for _, peerAddr := range bootstrapNodes {
		if peerAddr == "" {
			continue
		}
		addr, err := multiaddr.NewMultiaddr(peerAddr)
		if err != nil {
			return nil, err
		}
		addrs = append(addrs, addr)
	}
	return addrs, nil
}
