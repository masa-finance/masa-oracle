package config

import (
	"fmt"
)

func DisplayWelcomeMessage(multiAddr, ipAddr, publicKeyHex string, isStaked bool, isValidator bool, isTwitterScraper bool, isTelegramScraper bool, isDiscordScraper bool, isWebScraper bool, version string) {
	// ANSI escape code for yellow text
	yellow := "\033[33m"
	blue := "\033[34m"
	reset := "\033[0m"

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

	fmt.Printf(blue+"%-20s %s\n"+reset, "Version:", yellow+version)
	fmt.Printf(blue+"%-20s %s\n"+reset, "Multiaddress:", multiAddr)
	fmt.Printf(blue+"%-20s %s\n"+reset, "IP Address:", ipAddr)
	fmt.Printf(blue+"%-20s %s\n"+reset, "Public Key:", publicKeyHex)
	fmt.Printf(blue+"%-20s %t\n"+reset, "Is Staked:", isStaked)
	fmt.Printf(blue+"%-20s %t\n"+reset, "Is Validator:", isValidator)
	fmt.Printf(blue+"%-20s %t\n"+reset, "Is TwitterScraper:", isTwitterScraper)
	fmt.Printf(blue+"%-20s %t\n"+reset, "Is DiscordScraper:", isDiscordScraper)
	fmt.Printf(blue+"%-20s %t\n"+reset, "Is TelegramScraper:", isTelegramScraper)
	fmt.Printf(blue+"%-20s %t\n"+reset, "Is WebScraper:", isWebScraper)
	fmt.Println("")
}
