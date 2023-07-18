// server.go
package main

import (
	"log"
	"net/http"
)

func main() {
	http.HandleFunc("/webhook", webhookHandler)
	log.Fatal(http.ListenAndServe(":8080", nil))
}

func webhookHandler(w http.ResponseWriter, r *http.Request) {
	// Call the function from mid.go
	err := callMintFunction()
	if err != nil {
		http.Error(w, "Failed to call mint function", http.StatusInternalServerError)
		return
	}

	w.Write([]byte("Mint function called successfully"))
}
