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

func GetMultiAddressesForHostQuiet(host host.Host) []multiaddr.Multiaddr {
	ma, err := GetMultiAddressesForHost(host)
	if err != nil {
		logrus.Fatal(err)
	}
	return ma
}

func GetPriorityAddress(addrs []multiaddr.Multiaddr) multiaddr.Multiaddr {
	var bestPublicAddr multiaddr.Multiaddr
	var bestPrivateAddr multiaddr.Multiaddr

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
	logrus.Debug("Best public address: ", bestPublicAddr)
	logrus.Debug("Best private address: ", bestPrivateAddr)
	logrus.Debug("Base address: ", baseAddr)
	gcpAddr := replaceGCPAddress(baseAddr)
	if gcpAddr != nil {
		baseAddr = gcpAddr
	}
	return baseAddr
}

func isPreferredAddress(addr multiaddr.Multiaddr) bool {
	// Check if the multiaddress contains the UDP protocol
	for _, p := range addr.Protocols() {
		if p.Code == multiaddr.P_UDP {
			return true
		}
	}
	return false
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

//func contains(slice []multiaddr.Multiaddr, item multiaddr.Multiaddr) bool {
//	for _, a := range slice {
//		if a.Equal(item) {
//			return true
//		}
//	}
//	return false
//}

func GetBootNodesMultiAddress(input string) ([]multiaddr.Multiaddr, error) {
	logrus.Infof("Getting bootnodes from %s", input)
	bootstrapPeers := strings.Split(input, ",")
	addrs := make([]multiaddr.Multiaddr, 0)
	for _, peerAddr := range bootstrapPeers {
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
	result := string(body)
	//check if the result is a valid IP
	if net.ParseIP(result) == nil {
		return false, ""
	}
	return true, result
}
