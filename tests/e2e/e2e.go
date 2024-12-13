// Blockchain integration test
package masa_test

import (
	"bytes"
	"fmt"
	"io"
	"net/http"
	"os"

	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
)

func getData(url string) ([]byte, error) {
	// Create a new GET request
	req, err := http.NewRequest("GET", url, nil)
	if err != nil {
		fmt.Println("Error creating request:", err)
		return nil, err
	}

	// Send the request using http.DefaultClient
	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		fmt.Println("Error sending request:", err)
		return nil, err
	}
	defer resp.Body.Close()

	// Read the response body
	responseBody, err := io.ReadAll(resp.Body)
	if err != nil {
		fmt.Println("Error reading response:", err)
		return nil, err
	}

	fmt.Printf("Response status: %s\n", resp.Status)
	fmt.Printf("Response body: %s\n", string(responseBody))

	return responseBody, nil
}

func postData(url string, requestBody []byte) ([]byte, error) {
	// Create a new POST request
	req, err := http.NewRequest("POST", url, bytes.NewBuffer(requestBody))
	if err != nil {
		fmt.Println("Error creating request:", err)
		return nil, err
	}

	// Set the headers
	req.Header.Set("accept", "application/json")
	req.Header.Set("Content-Type", "application/json")

	// Send the request using http.DefaultClient
	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		fmt.Println("Error sending request:", err)
		return nil, err
	}
	defer resp.Body.Close()

	// Read the response body
	responseBody, err := io.ReadAll(resp.Body)
	if err != nil {
		fmt.Println("Error reading response:", err)
		return nil, err
	}

	fmt.Printf("Response status: %s\n", resp.Status)
	fmt.Printf("Response body: %s\n", string(responseBody))

	return responseBody, nil
}

var _ = Describe("E2E tests", func() {

	// Make sure we don't run e2e tests accidentally
	BeforeEach(func() {
		if os.Getenv("E2E") != "true" {
			Skip("Skipping E2E tests")
		}

		Eventually(func() bool {
			data, err := getData("http://localhost:9092/api/v1/peers")
			if err != nil {
				return false
			}

			fmt.Println(string(data))
			// {"data":[{"peerId":"16Uiu2HAmQTSU51k4kx2hjGoXfNKJJjFMKvq2xWmVYgDSujgHpqoR"},{"peerId":"16Uiu2HAm74aXSU3CeBD2XWSb42UR1J6pGGKGFt9dLkAqYCaaRpY5"}],"success":true,"totalCount":2}

			return bytes.Contains(data, []byte(`"totalCount":2`))
		}, "1m").Should(BeTrue())
	})

	Context("can use the API", func() {
		It("scrapes the web", func() {
			response, err := postData("http://localhost:9092/api/v1/data/web", []byte(`{ "depth": 1, "url": "https://www.google.com"}`))

			Expect(err).ToNot(HaveOccurred())
			Expect(response).ToNot(BeNil())
			Expect(string(response)).To(ContainSubstring("google"))
		})
	})
})
