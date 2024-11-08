package main

import (
	"crypto/ecdsa"
	"fmt"
	"math/big"
	"time"

	"github.com/fatih/color"
	"github.com/sirupsen/logrus"

	"github.com/masa-finance/masa-oracle/pkg/config"
	"github.com/masa-finance/masa-oracle/pkg/staking"
)

func handleStaking(privateKey *ecdsa.PrivateKey, cfg *config.AppConfig) error {
	// Staking logic
	// Convert the stake amount to the smallest unit, assuming 18 decimal places
	amountBigInt, ok := new(big.Int).SetString(cfg.StakeAmount, 10)
	if !ok {
		logrus.Fatal("Invalid stake amount")
	}
	amountInSmallestUnit := new(big.Int).Mul(amountBigInt, big.NewInt(1e18))

	stakingClient, err := staking.NewClient(privateKey)
	if err != nil {
		return err
	}
	// Function to start and stop a spinner with a message
	startSpinner := func(msg string, txHashChan <-chan string, done chan bool) {
		spinner := []string{"|", "/", "-", "\\"}
		i := 0
		var txHash string
		for {
			select {
			case txHash = <-txHashChan: // Receive the transaction hash
				// Do not print anything here, just update the txHash variable
			case <-done:
				fmt.Printf("\r%s\n", msg) // Print final message when done
				if txHash != "" {
					fmt.Println(txHash) // Print the transaction hash on a new line
				}
				return
			default:
				// Use carriage return `\r` to overwrite the spinner animation on the same line
				// Remove the newline character `\n` from the print statement
				if txHash != "" {
					fmt.Printf("\r%s %s - %s", spinner[i], msg, txHash)
				} else {
					fmt.Printf("\r%s %s", spinner[i], msg)
				}
				i = (i + 1) % len(spinner)
				time.Sleep(100 * time.Millisecond)
			}
		}
	}

	// Approve the staking contract to spend tokens on behalf of the user
	var approveTxHash string
	done := make(chan bool)
	txHashChan := make(chan string, 1) // Buffer of 1 to prevent blocking
	go startSpinner("Approving staking contract to spend tokens...", txHashChan, done)
	approveTxHash, err = stakingClient.Approve(amountInSmallestUnit)
	if err != nil {
		logrus.Error("[-] Failed to approve tokens for staking:", err)
		return err
	}
	txHashChan <- approveTxHash // Send the transaction hash to the spinner
	done <- true                // Stop the spinner
	color.Green("Approve transaction hash: %s", approveTxHash)

	// Stake the tokens after approval
	var stakeTxHash string
	done = make(chan bool)
	txHashChan = make(chan string, 1) // Buffer of 1 to prevent blocking
	go startSpinner("Staking tokens...", txHashChan, done)
	stakeTxHash, err = stakingClient.Stake(amountInSmallestUnit)
	if err != nil {
		logrus.Error("[-] Failed to stake tokens:", err)
		return err
	}
	txHashChan <- stakeTxHash // Send the transaction hash to the spinner
	done <- true              // Stop the spinner
	color.Green("Stake transaction hash: %s", stakeTxHash)

	return nil
}

func handleFaucet(privateKey *ecdsa.PrivateKey) error {
	faucetClient, err := staking.NewClient(privateKey)
	if err != nil {
		logrus.Error("[-] Failed to create staking client:", err)
		return err
	}

	startSpinner := func(msg string, txHashChan <-chan string, done chan bool) {
		spinner := []string{"|", "/", "-", "\\"}
		i := 0
		var txHash string
		for {
			select {
			case txHash = <-txHashChan: // Receive the transaction hash
				// Do not print anything here, just update the txHash variable
			case <-done:
				fmt.Printf("\r%s\n", msg) // Print final message when done
				if txHash != "" {
					fmt.Println(txHash) // Print the transaction hash on a new line
				}
				return
			default:
				// Use carriage return `\r` to overwrite the spinner animation on the same line
				// Remove the newline character `\n` from the print statement
				if txHash != "" {
					fmt.Printf("\r%s %s - %s", spinner[i], msg, txHash)
				} else {
					fmt.Printf("\r%s %s", spinner[i], msg)
				}
				i = (i + 1) % len(spinner)
				time.Sleep(100 * time.Millisecond)
			}
		}
	}

	// Run the faucet
	var faucetTxHash string
	done := make(chan bool)
	txHashChan := make(chan string, 1) // Buffer of 1 to prevent blocking
	go startSpinner("Requesting tokens from faucet...", txHashChan, done)
	faucetTxHash, err = faucetClient.RunFaucet()
	if err != nil {
		logrus.Error("[-] Failed to request tokens from faucet:", err)
		return err
	}
	txHashChan <- faucetTxHash // Send the transaction hash to the spinner
	done <- true               // Stop the spinner
	color.Green("[-] Faucet transaction hash: %s", faucetTxHash)

	return nil
}
