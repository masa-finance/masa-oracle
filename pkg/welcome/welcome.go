// welcome/welcome.go
package welcome

import (
	"fmt"
)

func DisplayWelcomeMessage(multiAddr string, ipAddr string) {
	multiAddrLine := fmt.Sprintf("#  Multiaddress: %-70s #", multiAddr)
	ipAddrLine := fmt.Sprintf("#  IP Address: %-70s #", ipAddr)

	fmt.Println("################################################################################")
	fmt.Println("#                                                                              #")
	fmt.Println("#  Welcome to the Masa Oracle Node                                             #")
	fmt.Println(multiAddrLine)
	fmt.Println(ipAddrLine)
	fmt.Println("#                                                                              #")
	fmt.Println("################################################################################")
}
