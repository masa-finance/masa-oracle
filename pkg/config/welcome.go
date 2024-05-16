package config

import (
	"fmt"
)

func DisplayWelcomeMessage(multiAddr, ipAddr, publicKeyHex string, isStaked bool, isWriterNode bool, isTwitterScraper bool, isDiscordScraper bool, isWebScraper bool, version string) {
	// ANSI escape code for yellow text
	green := "\033[32m"
	yellow := "\033[33m"
	blue := "\033[34m"
	// red := "\033[31m"
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
	fmt.Printf(green+"Version:		%s\n"+reset, version)
	// Displaying the multi-address and IP address in blue
	fmt.Printf(blue+"Multiaddress:		%s\n"+reset, multiAddr)
	fmt.Printf(blue+"IP Address:		%s\n"+reset, ipAddr)
	fmt.Printf(blue+"Public Key:   		%s\n"+reset, publicKeyHex)
	fmt.Printf(blue+"Is Staked:    		%t\n"+reset, isStaked)
	fmt.Printf(blue+"Is Writer:    		%t\n"+reset, isWriterNode)
	fmt.Printf(blue+"Is TwitterScraper:	%t\n"+reset, isTwitterScraper)
	fmt.Printf(blue+"Is DiscordScraper:	%t\n"+reset, isDiscordScraper)
	fmt.Printf(blue+"Is WebScraper:   	%t\n"+reset, isWebScraper)
	fmt.Println("")
}
