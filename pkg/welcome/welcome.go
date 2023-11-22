// welcome/welcome.go
package welcome

import (
	"fmt"
)

func DisplayWelcomeMessage(multiAddr string, ipAddr string) {
	// ANSI escape code for yellow text
	yellow := "\033[33m"
	// ANSI escaoe code for blue text
	blue := "\033[34m"
	// ANSI escape code to reset color
	reset := "\033[0m"

	borderLine := "#######################################"

	fmt.Println(yellow + borderLine + reset)
	fmt.Println(yellow + "#                                     #" + reset)
	fmt.Println(yellow + "#                                     #" + reset)
	fmt.Println(yellow + "#                                     #" + reset)
	fmt.Println(yellow + "#                                     #" + reset)
	fmt.Println(yellow + "#     __  __    _    ____    _        #" + reset)
	fmt.Println(yellow + "#    |  \\/  |  / \\  / ___|  / \\       #" + reset)
	fmt.Println(yellow + "#    | |\\/| | / _ \\ \\___ \\ / _ \\      #" + reset)
	fmt.Println(yellow + "#    | |  | |/ ___ \\ ___) / ___ \\     #" + reset)
	fmt.Println(yellow + "#    |_|  |_/_/   \\_\\____/_/   \\_\\    #" + reset)
	fmt.Println(yellow + "#                                     #" + reset)
	fmt.Println(yellow + "#                                     #" + reset)
	fmt.Println(yellow + "#                                     #" + reset)
	fmt.Println(yellow + "#                                     #" + reset)
	fmt.Println(yellow + borderLine + reset)

	// Displaying the multi-address and IP address in yellow
	// Displaying the multi-address and IP address in blue
	fmt.Printf(blue+"Multiaddress: %s\n"+reset, multiAddr)
	fmt.Printf(blue+"IP Address:   %s\n"+reset, ipAddr)
}
