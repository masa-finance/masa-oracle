// welcome/welcome.go
package welcome

import (
	"fmt"
)

func DisplayWelcomeMessage(multiAddr string, ipAddr string) {
	borderLine := "#######################################"

	fmt.Println(borderLine)
	fmt.Println("#                                     #")
	fmt.Println("#                                     #")
	fmt.Println("#                                     #")
	fmt.Println("#                                     #")
	fmt.Println("#     __  __    _    ____    _        #")
	fmt.Println("#    |  \\/  |  / \\  / ___|  / \\       #")
	fmt.Println("#    | |\\/| | / _ \\ \\___ \\ / _ \\      #")
	fmt.Println("#    | |  | |/ ___ \\ ___) / ___ \\     #")
	fmt.Println("#    |_|  |_/_/   \\_\\____/_/   \\_\\    #")
	fmt.Println("#                                     #")
	fmt.Println("#                                     #")
	fmt.Println("#                                     #")
	fmt.Println("#                                     #")
	fmt.Println(borderLine)

	// Displaying the multi-address and IP address below the ASCII art
	fmt.Printf("Multiaddress: %s\n", multiAddr)
	fmt.Printf("IP Address:   %s\n", ipAddr)
}
