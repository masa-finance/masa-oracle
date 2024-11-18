package main

import (
	"context"
	"encoding/json"
	"fmt"

	"github.com/gdamore/tcell/v2"
	"github.com/libp2p/go-libp2p"
	pubsub "github.com/libp2p/go-libp2p-pubsub"
	"github.com/libp2p/go-libp2p/core/peer"
	quic "github.com/libp2p/go-libp2p/p2p/transport/quic"
	"github.com/multiformats/go-multiaddr"
	"github.com/rivo/tview"
	"github.com/sirupsen/logrus"
)

// HandleMessage implement subscription handler here
func (handler *SubscriptionHandler) HandleMessage(message *pubsub.Message) {
	fmt.Println("Received a message")
	var gossip Gossip
	err := json.Unmarshal(message.Data, &gossip)
	if err != nil {
		logrus.Errorf("[-] Failed to unmarshal message: %v", err)
		return
	}

	handler.mu.Lock()
	handler.Gossips = append(handler.Gossips, gossip)
	handler.mu.Unlock()

	// logrus.Infof("added: %+v", gossip)
}

// showMenu creates and returns the menu component.
func handleMenu(app *tview.Application, output *tview.TextView) *tview.List {
	menu := tview.NewList().
		AddItem("Connect", "[darkgray]connect to an oracle node", '1', func() {
			handleOption(app, "1", output)
		}).
		// Items 2-5 removed - LLM functionality eliminated in	https://github.com/masa-finance/masa-oracle/pull/626
		AddItem("Oracle Nodes", "[darkgray]view active nodes", '6', func() {
			handleOption(app, "6", output)
		})

	menu.AddItem("Quit", "[darkgray]press to exit", 'q', func() {
		handleOption(app, "7", output)
	}).SetBorder(true).SetBorderColor(tcell.ColorGray)

	return menu
}

// handleOption triggers actions based on user selection.
func handleOption(app *tview.Application, option string, output *tview.TextView) {
	switch option {
	case "1":
		modalFlex := tview.NewFlex().SetDirection(tview.FlexRow)
		modalFlex.SetBorderPadding(1, 1, 1, 1)

		var form *tview.Form

		// Create a new form
		form = tview.NewForm().
			AddInputField("Node Multiaddress", "", 60, nil, nil).
			AddButton("OK", func() {
				inputValue := form.GetFormItemByLabel("Node Multiaddress").(*tview.InputField).GetText()
				appConfig.Address = inputValue

				if appConfig.Address == "" {
					output.SetText("A multiaddress was not entered. Please enter the masa node multiaddress and try again.")
				} else {
					output.SetText(fmt.Sprintf("Connecting to: %s", appConfig.Address))

					maddr, err := multiaddr.NewMultiaddr(appConfig.Address)
					if err != nil {
						logrus.Errorf("[-] %v", err)
					}

					// Create a libp2p host to connect to the Masa node
					host, err := libp2p.New(libp2p.NoSecurity, libp2p.Transport(quic.NewTransport))
					if err != nil {
						logrus.Errorf("[-] %v", err)
					}

					// Extract the peer ID from the multiaddress
					peerInfo, err := peer.AddrInfoFromP2pAddr(maddr)
					if err != nil {
						logrus.Errorf("[-] %v", err)
					}

					// Connect to the peer
					if err := host.Connect(context.Background(), *peerInfo); err != nil {
						logrus.Errorf("[-] %v", err)
					}

					output.SetText(fmt.Sprintf("Successfully connected to node: %s", appConfig.Address))
				}
				app.SetRoot(mainFlex, true) // Return to main view
			}).
			AddButton("Cancel", func() {
				output.SetText("Cancelled entering a masa node multiaddress.")
				app.SetRoot(mainFlex, true) // Return to main view
			})

		form.SetBorder(true).SetBorderColor(tcell.ColorBlue)

		modalFlex.AddItem(form, 0, 1, true)

		app.SetRoot(modalFlex, true).SetFocus(form)
	case "6":
		content := Splash()

		table := tview.NewTable().SetBorders(true).SetFixed(1, 0)

		// Set header titles
		table.SetCell(0, 0, tview.NewTableCell("Address").SetTextColor(tcell.ColorYellow).SetAlign(tview.AlignCenter))
		table.SetCell(0, 1, tview.NewTableCell("IsStaked").SetTextColor(tcell.ColorYellow).SetAlign(tview.AlignCenter))
		table.SetCell(0, 2, tview.NewTableCell("IsValidator").SetTextColor(tcell.ColorYellow).SetAlign(tview.AlignCenter))

		// Set cell values for each row
		for i := 1; i <= 10; i++ {
			table.SetCell(i, 0, tview.NewTableCell("/ip4/127.0.0.1/udp/4001/quic-v1/p2p/16Uiu2HAmVRNDAZ6J1eHTV8twU6VaX8vqhe7VehPBNrCzDrHB9aQn"))
			table.SetCell(i, 1, tview.NewTableCell("false"))
			table.SetCell(i, 2, tview.NewTableCell("false"))
		}

		table.Select(0, 0).SetFixed(1, 1).SetDoneFunc(func(key tcell.Key) {
			if key == tcell.KeyEscape {
				mainFlex.Clear().
					AddItem(content, 0, 1, false).
					AddItem(handleMenu(app, output), 0, 1, true).
					AddItem(output, 0, 3, false)

				app.SetRoot(mainFlex, true) // Return to main view
				return
			}
			if key == tcell.KeyEnter {
				table.SetSelectable(true, true)
			}
		}).SetSelectedFunc(func(row int, column int) {
			table.GetCell(row, column).SetTextColor(tcell.ColorRed)
			table.SetSelectable(false, false)
		})

		flex := tview.NewFlex().
			SetDirection(tview.FlexRow).
			AddItem(table, 0, 1, false)

		mainFlex.Clear().
			AddItem(content, 0, 1, false).
			AddItem(handleMenu(app, output), 0, 1, false).
			AddItem(flex, 0, 3, true)

		flex.SetBorder(true).SetBorderColor(tcell.ColorBlue).
			SetTitle(" Masa Oracle Nodes, press esc to return to menu ")

		app.SetRoot(mainFlex, true).SetFocus(table)
	case "7":
		modalFlex := tview.NewFlex().SetDirection(tview.FlexRow)
		modalFlex.SetBorderPadding(1, 1, 1, 1)

		modal := tview.NewModal().
			SetText("Are you sure you want to quit?").
			AddButtons([]string{"Yes", "No"}).
			SetDoneFunc(func(buttonIndex int, buttonLabel string) {
				if buttonLabel == "Yes" {
					app.Stop()
				}
				app.SetRoot(mainFlex, true) // Return to main view
			})

		modalFlex.AddItem(modal, 0, 1, true)

		app.SetRoot(modalFlex, true)
	default:
		break
	}
}
