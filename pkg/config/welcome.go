package config

import (
	"fmt"

	"github.com/multiformats/go-multiaddr"
	"github.com/sirupsen/logrus"
)

func DisplayWelcomeMessage(multiAddrs []multiaddr.Multiaddr, publicKeyHex string, isStaked bool, isValidator bool, isTwitterScraper bool, isTelegramScraper bool, isDiscordScraper bool, isWebScraper bool, version, protocolVersion string) {
	// ANSI escape code for yellow text
	yellow := "\033[33m"
	blue := "\033[34m"
	reset := "\033[0m"

	var maddrs, ips string

	for _, ma := range multiAddrs {
		ip, err := ma.ValueForProtocol(multiaddr.P_IP4) // Get the IP address
		if err != nil {
			logrus.Errorf("[-] Error while parsing getting IP address for %v: %v", ma, err)
			continue
		}

		maddrs = fmt.Sprintf("%s %s", maddrs, ma)
		ips = fmt.Sprintf("%s %s", ips, ip)
	}

	fmt.Println("")
	borderLine := "#######################################"

	fmt.Println(yellow + borderLine + reset)
	fmt.Println(yellow + "#     __  __    _    ____    _        #" + reset)
	fmt.Println(yellow + "#    |  \\/  |  / \\  / ___|  / \\       #" + reset)
	fmt.Println(yellow + "#    | |\\/| | / _ \\ \\___ \\ / _ \\      #" + reset)
	fmt.Println(yellow + "#    | |  | |/ ___ \\ ___) / ___ \\     #" + reset)
	fmt.Println(yellow + "#    |_|  |_/_/   \\_\\____/_/   \\_\\    #" + reset)
	fmt.Println(yellow + "#                                     #" + reset)
	fmt.Println(yellow + borderLine + reset)
	fmt.Println("")

	fmt.Printf(blue+"%-20s %s\n"+reset, "Application Version:", yellow+version)
	fmt.Printf(blue+"%-20s %s\n"+reset, "Protocol Version:", yellow+protocolVersion)
	fmt.Printf(blue+"%-20s %s\n"+reset, "Multiaddresses:", maddrs)
	fmt.Printf(blue+"%-20s %s\n"+reset, "IP Addresses:", ips)
	fmt.Printf(blue+"%-20s %s\n"+reset, "Public Key:", publicKeyHex)
	fmt.Printf(blue+"%-20s %t\n"+reset, "Is Staked:", isStaked)
	fmt.Printf(blue+"%-20s %t\n"+reset, "Is Validator:", isValidator)
	fmt.Printf(blue+"%-20s %t\n"+reset, "Is TwitterScraper:", isTwitterScraper)
	fmt.Printf(blue+"%-20s %t\n"+reset, "Is DiscordScraper:", isDiscordScraper)
	fmt.Printf(blue+"%-20s %t\n"+reset, "Is TelegramScraper:", isTelegramScraper)
	fmt.Printf(blue+"%-20s %t\n"+reset, "Is WebScraper:", isWebScraper)
	fmt.Println("")
}
