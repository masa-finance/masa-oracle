package main

import (
	"encoding/json"
	"fmt"
	"log"
	"os"
	"path/filepath"
	"strings"

	"github.com/ethereum/go-ethereum/accounts/abi/bind"
)

type Contract struct {
	ABI      []interface{} `json:"abi"`
	Bytecode string        `json:"bytecode"`
}

func AbiGen() error {
	if len(os.Args) < 2 {
		panic("Please provide the path to the JSON file as a command-line argument.")
	}
	inFile := os.Args[1]
	contractName := strings.TrimSuffix(filepath.Base(inFile), filepath.Ext(inFile))

	// Read the compiled contract JSON file
	data, err := os.ReadFile(inFile)
	if err != nil {
		return err
	}

	// Unmarshal the JSON data into a Contract struct
	var contract Contract
	if err := json.Unmarshal(data, &contract); err != nil {
		return err
	}

	// Marshal the ABI into JSON
	abi, err := json.Marshal(contract.ABI)
	if err != nil {
		return err
	}

	// Write the ABI JSON to a file
	if err := os.WriteFile(fmt.Sprintf("%s.abi.json", contractName), abi, 0644); err != nil {
		return err
	}

	// Write the bytecode to a file
	if err := os.WriteFile(fmt.Sprintf("%s.bin", contractName), []byte(contract.Bytecode), 0644); err != nil {
		return err
	}
	return nil
}

func AbiGenBindings(inPath string) error {
	if inPath == "" {
		if len(os.Args) < 2 {
			panic("Please provide the path to the JSON file as a command-line argument.")
		}
		inPath = os.Args[1]
	}
	contractName := strings.TrimSuffix(filepath.Base(inPath), filepath.Ext(inPath))
	// Compile the Solidity contract
	//contracts, err := compiler.ParseCombinedJSON("", string(sol))
	//if err != nil {
	//	log.Fatalf("Failed to compile Solidity contract: %v", err)
	//}

	abiData, err := os.ReadFile(inPath)
	if err != nil {
		log.Fatalf("[-] Failed to read ABI: %v", err)
	}

	var contract Contract
	if err := json.Unmarshal(abiData, &contract); err != nil {
		return err
	}

	// Generate the Go bindings
	code, err := bind.Bind(
		[]string{contractName},      // Names of the Go contracts
		[]string{string(abiData)},   // Ethereum contract ABIs
		[]string{contract.Bytecode}, // Ethereum contract bytecodes
		nil,                         // Ethereum contract fsigs
		"package_name",              // Name of the Go package
		bind.LangGo,                 // Target language
		nil,                         // Ethereum contract libraries
		nil,                         // Ethereum contract aliases
	)
	if err != nil {
		log.Fatalf("[-] Failed to generate Go bindings: %v", err)
	}

	// Write the Go bindings to a file
	err = os.WriteFile(fmt.Sprintf("%s.go", contractName), []byte(code), 0644)
	if err != nil {
		log.Fatalf("[-] Failed to write Go bindings to file: %v", err)
	}
	fmt.Println("[+] Go bindings generated successfully!")
	return nil
}
