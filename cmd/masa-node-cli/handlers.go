package main

import (
	"bytes"
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"log"
	"net/http"
	"os"
	"os/exec"
	"strconv"
	"strings"

	"github.com/gdamore/tcell/v2"
	"github.com/libp2p/go-libp2p"
	pubsub "github.com/libp2p/go-libp2p-pubsub"
	"github.com/libp2p/go-libp2p/core/peer"
	quic "github.com/libp2p/go-libp2p/p2p/transport/quic"
	"github.com/multiformats/go-multiaddr"
	"github.com/rivo/tview"
	"github.com/sashabaranov/go-openai"
	"github.com/sirupsen/logrus"
)

// handleIPAddress takes a multi-address string as input and extracts the IP address from it.
// The function supports both "/ip4/" and "/dns/" multi-address formats.
// If the multi-address contains an IP address, it is returned directly.
// If the multi-address contains a DNS name, it is also returned as is.
// In case the multi-address does not follow the expected format or the IP address/DNS name is not found, an empty string is returned.
func handleIPAddress(multiAddr string) string {
	parts := strings.Split(multiAddr, "/")
	// Assuming the IP address is always after "/ip4/"
	for i, part := range parts {
		if part == "ip4" {
			return parts[i+1]
		} else if part == "dns" {
			return parts[i+1]
		}
	}
	return ""
}

// handleOpenFile reads the content of a file specified by the filename 'f' and returns it as a string.
// If the file cannot be read, the function logs a fatal error and exits the program.
// Parameters:
// - f: The name of the file to read.
// Returns:
// - A string containing the content of the file.
// func handleOpenFile(f string) string {
// 	dat, err := os.ReadFile(f)
// 	if err != nil {
// 		log.Print(err)
// 		return ""
// 	}
// 	return string(dat)
// }

// handleSaveFile writes the provided content to a file specified by the filename 'f'.
// It appends the content to the file if it already exists, or creates a new file with the content if it does not.
// The file is created with permissions set to 0755.
// Parameters:
// - f: The name of the file to which the content will be written.
// - content: The content to write to the file.
func handleSaveFile(f string, content string) {
	file, err := os.OpenFile(f, os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0755)
	file.WriteString(content + "\n")
	if err != nil {
		log.Println(err)
		return
	}
}

// handleGPT takes a prompt and a user message, sends them to the OpenAI API, and returns the generated response.
// It utilizes the OPENAI_API_KEY environment variable for authentication.
// Parameters:
// - prompt: A string containing the initial prompt for the AI.
// - user_message: A string containing the user's message to the AI.
// Returns:
// - A string containing the AI's response.
// - An error if the request to the OpenAI API fails.
func handleGPT(prompt string, userMessage string) (string, error) {
	key := os.Getenv("OPENAI_API_KEY")
	if key == "" {
		log.Println("OPENAI_API_KEY is not set. Please set the environment variable and try again.")
		return "", errors.New("OPENAI_API_KEY is not set")
	}
	client := openai.NewClient(key)
	resp, err := client.CreateChatCompletion(
		context.Background(),
		openai.ChatCompletionRequest{
			Model: openai.GPT4,
			Messages: []openai.ChatCompletionMessage{
				{
					Role:    openai.ChatMessageRoleSystem,
					Content: prompt,
				},
				{
					Role:    openai.ChatMessageRoleUser,
					Content: userMessage,
				},
			},
		},
	)
	if err != nil {
		log.Print(err)
		return "", err
	}
	return resp.Choices[0].Message.Content, nil
}

// handleSpeak sends a given text response to the ElevenLabs API for text-to-speech conversion,
// then plays the resulting audio. It uses the ELAB_KEY environment variable for API authentication.
func handleSpeak(response string) {
	key := os.Getenv("ELAB_KEY")

	data := SpeakRequest{
		Text: response,
		VoiceSettings: struct {
			Stability       float64 "json:\"stability\""
			SimilarityBoost float64 "json:\"similarity_boost\""
		}{
			Stability:       0.6,
			SimilarityBoost: 0.85,
		},
	}

	buf, _ := json.Marshal(data)

	req, err := http.NewRequest(http.MethodPost, os.Getenv("ELAB_URL"), bytes.NewBuffer(buf))
	if err != nil {
		log.Print(err)
		return
	}
	req.Header.Set("accept", "*/*")
	req.Header.Set("xi-api-key", key)
	req.Header.Set("Content-Type", "application/json")

	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		log.Print(err)
		return
	}
	defer resp.Body.Close()

	if resp.StatusCode == http.StatusOK {
		bodyBytes, err := io.ReadAll(resp.Body)
		if err != nil {
			log.Print(err)
			return
		}
		file, err := os.Create("output.mp3")
		file.Write(bodyBytes)
		if err != nil {
			log.Print(err)
			return
		} else {
			cmd := exec.Command("afplay", "output.mp3")
			go cmd.Run()
			go handleTranscribe("output.mp3", "transcription.txt")
		}
	}
}

// handleTranscribe takes an audio file and a target text file as input.
// It uses the OpenAI API to handleTranscribe the audio to text, then saves the text to the specified text file.
func handleTranscribe(audioFile string, txtFile string) error {
	key := os.Getenv("OPENAI_API_KEY")
	if key == "" {
		log.Println("OPENAI_API_KEY is not set. Please set the environment variable and try again.")
		return errors.New("OPENAI_API_KEY is not set")
	}
	client := openai.NewClient(key)
	ctx := context.Background()
	req := openai.AudioRequest{
		Model:    openai.Whisper1,
		FilePath: audioFile,
	}
	resp, err := client.CreateTranscription(ctx, req)
	if err != nil {
		fmt.Printf("Transcription error: %v\n", err)
		return err
	} else {
		handleSaveFile(txtFile, resp.Text)
	}
	return nil
}

// HandleMessage implement subscription handler here
func (handler *SubscriptionHandler) HandleMessage(message *pubsub.Message) {
	fmt.Println("Received a message")
	var gossip Gossip
	err := json.Unmarshal(message.Data, &gossip)
	if err != nil {
		logrus.Errorf("Failed to unmarshal message: %v", err)
		return
	}

	handler.mu.Lock()
	handler.Gossips = append(handler.Gossips, gossip)
	handler.mu.Unlock()

	logrus.Infof("added: %+v", gossip)
}

// showMenu creates and returns the menu component.
func handleMenu(app *tview.Application, output *tview.TextView) *tview.List {
	menu := tview.NewList().
		AddItem("Connect", "[darkgray]connect to an oracle node", '1', func() {
			handleOption(app, "1", output)
		}).
		AddItem("LLM Model", "[darkgray]select an llm model to use", '2', func() {
			handleOption(app, "2", output)
		}).
		AddItem("Twitter", "[darkgray]set your Twitter/X credentials", '3', func() {
			handleOption(app, "3", output)
		}).
		AddItem("ChatGPT", "[darkgray]chat with a helpful assistant", '4', func() {
			handleOption(app, "4", output)
		}).
		AddItem("Sentiment", "[darkgray]analyze sentiment from tweets", '5', func() {
			handleOption(app, "5", output)
		}).
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
						logrus.Errorf("%v", err)
					}

					// Create a libp2p host to connect to the Masa node
					host, err := libp2p.New(libp2p.NoSecurity, libp2p.Transport(quic.NewTransport))
					if err != nil {
						logrus.Errorf("%v", err)
					}

					// Extract the peer ID from the multiaddress
					peerInfo, err := peer.AddrInfoFromP2pAddr(maddr)
					if err != nil {
						logrus.Errorf("%v", err)
					}

					// Connect to the peer
					if err := host.Connect(context.Background(), *peerInfo); err != nil {
						logrus.Errorf("%v", err)
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
	case "2":
		// @todo add more models here then pass that to the sentiment analyzer
		radioButtons := NewRadioButtons([]string{"claude-3-opus-20240229", "claude-3-sonnet-20240229", "claude-3-haiku-20240307", "GPT4"}, func(option string) {
			appConfig.Model = option
			output.SetText(fmt.Sprintf("Selected model: %s", option)).SetTextAlign(tview.AlignLeft)
			app.SetRoot(mainFlex, true) // Return to main view after selection
		})

		radioButtons.SetBorder(true).
			SetTitle(" Choose LLM Model ").
			SetRect(85, 20, 30, 10) // centering on screen

		app.SetRoot(radioButtons, false)
	case "3":
		// @todo used session stored twitter creds, right now pulling from .env
		modalFlex := tview.NewFlex().SetDirection(tview.FlexRow)
		modalFlex.SetBorderPadding(1, 1, 1, 1)

		var form *tview.Form

		// Create a new form
		form = tview.NewForm().
			AddInputField("Username", "", 60, nil, nil).
			AddPasswordField("Password", "", 60, '*', nil).
			AddButton("OK", func() {
				inputValue := form.GetFormItemByLabel("Username").(*tview.InputField).GetText()
				passValue := form.GetFormItemByLabel("Password").(*tview.InputField).GetText()
				appConfig.TwitterUser = inputValue
				appConfig.TwitterPassword = passValue

				if appConfig.TwitterUser == "" {
					output.SetText("A Twitter username was not entered. Please enter your username and try again.")
					return
				} else if appConfig.TwitterPassword == "" {
					output.SetText("A Twitter password was not entered. Please enter your password and try again.")
					return
				} else {
					output.SetText("Twitter credentials saved for this session only!")
				}
				app.SetRoot(mainFlex, true) // Return to main view
			}).
			AddButton("Cancel", func() {
				output.SetText("Cancelled storing Twitter credentials.")
				app.SetRoot(mainFlex, true) // Return to main view
			})

		form.SetBorder(true).SetBorderColor(tcell.ColorBlue).SetTitle(" Twitter Credentials ")

		modalFlex.AddItem(form, 0, 1, true)

		app.SetRoot(modalFlex, true).SetFocus(form)
	case "4":
		// Create the input field for user messages.
		inputField := tview.NewInputField().
			SetLabel("> ").
			SetFieldWidth(100)

		// Create the text view to display responses.
		textView := tview.NewTextView().
			SetDynamicColors(true).
			SetRegions(true).
			SetWordWrap(true)

		content := Splash()

		// Add an event listener to the input field.
		inputField.SetDoneFunc(func(key tcell.Key) {
			userMessage := inputField.GetText()
			// Clear the input field for the next message.
			inputField.SetText("")

			if userMessage == "\\q" {
				output.SetText(" Welcome to the MASA Oracle Client ")

				mainFlex.Clear().
					AddItem(content, 0, 1, false).
					AddItem(handleMenu(app, output), 0, 1, true).
					AddItem(output, 0, 3, false)

				app.SetRoot(mainFlex, true) // Return to main view
				return
			}

			prompt := os.Getenv("PROMPT")

			resp, err := handleGPT(prompt, userMessage)
			if err != nil {
				output.SetText(fmt.Sprintf("%v", err))
				return
			}
			if os.Getenv("ELAB_KEY") != "" || os.Getenv("ELAB_URL") != "" {
				handleSpeak(resp)
			}

			// Display the response in the text view.
			fmt.Fprintf(textView, "%s\n", resp)

		})

		inputField.Autocomplete().SetFieldWidth(0)

		// Arrange the input field and text view in a layout.
		flex := tview.NewFlex().
			SetDirection(tview.FlexRow).
			AddItem(textView, 0, 1, false).
			AddItem(inputField, 3, 1, true) // Set input field to be focused.

		mainFlex.Clear().
			AddItem(content, 0, 1, false).
			AddItem(handleMenu(app, output), 0, 1, true).
			AddItem(flex, 0, 3, false)

		flex.SetBorder(true).SetBorderColor(tcell.ColorBlue).
			SetTitle(" Start typing in the input box and press enter or type \\q to exit ")

		app.SetRoot(mainFlex, true).SetFocus(flex)
	case "5":
		// @todo use gossip topic instead of api to allow all staked nodes to participate in this analysis

		// node := &masa.OracleNode{}

		// gossipStatusHandler := &SubscriptionHandler{}
		// err := node.PubSubManager.Subscribe(config.TopicWithVersion(config.NodeGossipTopic), gossipStatusHandler)
		// if err != nil {
		// 	fmt.Println("Failed to subscribe to Sentiment Topic:", err)
		// 	return
		// }

		// err = node.PubSubManager.Publish(config.TopicWithVersion(config.NodeGossipTopic), []byte(message))
		// if err != nil {
		// 	fmt.Println("Failed to publish message:", err)
		// 	return
		// }

		if appConfig.Address == "" {
			output.SetText("Please connect to a masa node and try again.")
			break
		}

		var countMessage string
		var userMessage string

		// Create the input field for user messages.
		inputCountField := tview.NewInputField().
			SetLabel("# of Tweets to analyze ").
			SetFieldWidth(10)

		inputCountField.SetText("5")

		// Create the input field for user messages.
		inputField := tview.NewInputField().
			SetLabel("> ").
			SetFieldWidth(100)

		// Create the text view to display responses.
		textView := tview.NewTextView().
			SetDynamicColors(true).
			SetRegions(true).
			SetWordWrap(true)

		content := Splash()

		inputCountField.SetDoneFunc(func(key tcell.Key) {
			countMessage = inputCountField.GetText()
		})

		// Add an event listener to the input field.
		inputField.SetDoneFunc(func(key tcell.Key) {
			userMessage = inputField.GetText()
			// Clear the input field for the next message.
			inputField.SetText("")

			if userMessage == "\\q" {
				output.SetText(" Welcome to the MASA Oracle Client ")

				mainFlex.Clear().
					AddItem(content, 0, 1, false).
					AddItem(handleMenu(app, output), 0, 1, true).
					AddItem(output, 0, 3, false)

				app.SetRoot(mainFlex, true) // Return to main view
				return
			}

			count, _ := strconv.Atoi(countMessage)
			queryData := fmt.Sprintf(`{"query":"%s","count":%d}`, userMessage, count)

			uri := "http://" + handleIPAddress(appConfig.Address) + ":8080/analyzeSentiment"
			resp, err := http.Post(uri, "application/json", strings.NewReader(queryData))
			if err != nil {
				output.SetText(fmt.Sprintf("%v", err))
				return
			}
			var result map[string]interface{}
			if err := json.NewDecoder(resp.Body).Decode(&result); err != nil {
				output.SetText(fmt.Sprintf("%v", err))
				return
			}
			for _, r := range result {
				// Display the response in the text view.
				fmt.Fprintf(textView, "\n%s\n", r)
			}

		})

		inputField.Autocomplete().SetFieldWidth(0)

		// Arrange the input field and text view in a layout.
		flex := tview.NewFlex().
			SetDirection(tview.FlexRow).
			AddItem(textView, 0, 1, false).
			AddItem(inputCountField, 2, 1, false).
			AddItem(inputField, 3, 1, true)

		mainFlex.Clear().
			AddItem(content, 0, 1, false).
			AddItem(handleMenu(app, output), 0, 1, true).
			AddItem(flex, 0, 3, false)

		flex.SetBorder(true).SetBorderColor(tcell.ColorBlue).
			SetTitle(" Start typing in the input box and press enter or type \\q to exit ")

		app.SetRoot(mainFlex, true).SetFocus(flex)
	case "6":
		content := Splash()

		table := tview.NewTable().SetBorders(true).SetFixed(1, 0)

		// Set header titles
		table.SetCell(0, 0, tview.NewTableCell("Address").SetTextColor(tcell.ColorYellow).SetAlign(tview.AlignCenter))
		table.SetCell(0, 1, tview.NewTableCell("IsStaked").SetTextColor(tcell.ColorYellow).SetAlign(tview.AlignCenter))
		table.SetCell(0, 2, tview.NewTableCell("IsWriter").SetTextColor(tcell.ColorYellow).SetAlign(tview.AlignCenter))

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
