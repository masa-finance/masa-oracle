package main

import (
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"
	"strings"
)

type Contract struct {
	ABI      []interface{} `json:"abi"`
	Bytecode string        `json:"bytecode"`
}

func AbiGen() {
	if len(os.Args) < 2 {
		panic("Please provide the path to the JSON file as a command-line argument.")
	}
	inFile := os.Args[1]
	contractName := strings.TrimSuffix(filepath.Base(inFile), filepath.Ext(inFile))

	// Read the compiled contract JSON file
	data, err := os.ReadFile(inFile)
	if err != nil {
		panic(err)
	}

	// Unmarshal the JSON data into a Contract struct
	var contract Contract
	if err := json.Unmarshal(data, &contract); err != nil {
		panic(err)
	}

	// Marshal the ABI into JSON
	abi, err := json.Marshal(contract.ABI)
	if err != nil {
		panic(err)
	}

	// Write the ABI JSON to a file
	if err := os.WriteFile(fmt.Sprintf("%s.abi.json", contractName), abi, 0644); err != nil {
		panic(err)
	}

	// Write the bytecode to a file
	if err := os.WriteFile(fmt.Sprintf("%s.bin", contractName), []byte(contract.Bytecode), 0644); err != nil {
		panic(err)
	}
}
