package main

import (
	"bufio"
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
	"runtime"
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

// handleClearScreen clears the terminal screen by executing the appropriate command based on the operating system.
// It uses "cls" command for Windows and "clear" for Unix-like systems.
func handleClearScreen() {
	var cmd *exec.Cmd
	switch runtime.GOOS {
	case "windows":
		cmd = exec.Command("cmd", "/c", "cls")
	default:
		cmd = exec.Command("clear") // Works on Unix-like systems
	}
	cmd.Stdout = os.Stdout
	cmd.Run()
}

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
			// handleOption(app, "3", output)
			output.SetText(appConfig.Model)
		}).
		AddItem("ChatGPT", "[darkgray]chat with a helpful assistant", '4', func() {
			// handleOption(app, "4", output)
			flex := tview.NewFlex().SetDirection(tview.FlexRow)
			var inputField *tview.InputField

			inputField = tview.NewInputField().
				SetLabel("> ").
				SetFieldWidth(100).
				SetDoneFunc(func(key tcell.Key) {
					if key == tcell.KeyEnter {
						// Handle input here, similar to the loop in your ChatGPT interaction
						output.SetText(inputField.GetText())
						// Switch back to the main menu or clear the input field as needed
						app.SetRoot(mainFlex, true) // Return to main view
					}
				})

			flex.AddItem(inputField, 3, 1, true)
			app.SetRoot(flex, true)

			// fmt.Print("type \\q to return to main menu\n")
			// for {
			// 	fmt.Print("> ")
			// 	reader := bufio.NewReader(os.Stdin)
			// 	userMessage, err := reader.ReadString('\n')
			// 	if err != nil {
			// 		logrus.Errorf("%v", err)
			// 	}
			// 	userMessage = strings.TrimSuffix(userMessage, "\n")

			// 	if userMessage == "\\q" {
			// 		break
			// 	}

			// 	prompt := os.Getenv("PROMPT")

			// 	resp, err := handleGPT(prompt, userMessage)
			// 	if err != nil {
			// 		logrus.Errorf("%v", err)
			// 	}
			// 	fmt.Println(resp)
			// if os.Getenv("ELAB_KEY") != "" || os.Getenv("ELAB_URL") != "" {
			// 	handleSpeak(resp)
			// }
			// }

		}).
		AddItem("Sentiment", "[darkgray]analyze sentiment from tweets", '5', func() {
			// handleOption(app, "5", output)
			output.SetText("Subscribe to the Sentiment Topic")
		})

	menu.AddItem("Quit", "[darkgray]press to exit", 'q', func() {
		app.Stop()
	}).SetBorder(true).SetBorderColor(tcell.ColorGray)

	return menu
}

// handleOption triggers actions based on user selection.
func handleOption(app *tview.Application, option string, output *tview.TextView) {

	scanner := bufio.NewScanner(os.Stdin)

	switch option {
	case "1":
		modalFlex := tview.NewFlex().SetDirection(tview.FlexRow)
		modalFlex.SetBorderPadding(1, 1, 1, 1)

		inputField := tview.NewInputField().
			SetLabel("Multiaddress: ").
			SetFieldWidth(60)

		modal := tview.NewModal().
			SetText("Enter the masa node multiaddress and click OK.").
			AddButtons([]string{"OK", "Cancel"}).
			SetDoneFunc(func(buttonIndex int, buttonLabel string) {
				if buttonLabel == "OK" {
					appConfig.Address = inputField.GetText()
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
				}
				app.SetRoot(mainFlex, true) // Return to main view
			})

		modalFlex.AddItem(inputField, 2, 1, true).
			AddItem(modal, 0, 2, false)

		app.SetRoot(modalFlex, true).SetFocus(inputField)
	case "2":
		radioButtons := NewRadioButtons([]string{"Claude", "GPT4"}, func(option string) {
			appConfig.Model = option
			output.SetText(fmt.Sprintf("Selected model: %s", option))
			app.SetRoot(mainFlex, true) // Return to main view after selection
		})

		radioButtons.SetBorder(true).
			SetTitle(" Choose LLM Model ").
			SetRect(85, 20, 30, 5) // centering on screen

		app.SetRoot(radioButtons, false)
	case "3":
		fmt.Print("Enter Twitter Username: ")
		scanner.Scan()
		appConfig.TwitterUser = scanner.Text()

		fmt.Print("Enter Twitter Password: ")
		scanner.Scan()
		appConfig.TwitterPassword = scanner.Text()

		fmt.Println("Credentials saved during this session only.")
	case "4":
		fmt.Print("type \\q to return to main menu\n")
		for {
			fmt.Print("> ")
			reader := bufio.NewReader(os.Stdin)
			userMessage, err := reader.ReadString('\n')
			if err != nil {
				logrus.Errorf("%v", err)
			}
			userMessage = strings.TrimSuffix(userMessage, "\n")

			if userMessage == "\\q" {
				break
			}

			prompt := os.Getenv("PROMPT")

			resp, err := handleGPT(prompt, userMessage)
			if err != nil {
				logrus.Errorf("%v", err)
			}
			fmt.Println(resp)
			if os.Getenv("ELAB_KEY") != "" || os.Getenv("ELAB_URL") != "" {
				handleSpeak(resp)
			}
		}
	case "5":
		if appConfig.Address == "" {
			fmt.Println("Please connect to a masa node and try again.")
			break
		}

		// node := struct {
		// 	*masa.OracleNode
		// }{}

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

		fmt.Println("Subscribed to Sentiment Topic. Type your query:")
		scanner.Scan()
		query := scanner.Text()
		count := 5
		queryData := fmt.Sprintf(`{"query":"%s","count":%d}`, query, count)

		uri := "http://" + handleIPAddress(appConfig.Address) + ":8080/analyzeSentiment"
		// uri := "http://" + "localhost" + ":8080/analyzeSentiment"
		resp, err := http.Post(uri, "application/json", strings.NewReader(queryData))
		if err != nil {
			logrus.Errorf("%v", err)
			return
		}
		var result map[string]interface{}
		if err := json.NewDecoder(resp.Body).Decode(&result); err != nil {
			logrus.Errorf("%v", err)
			return
		}
		handleClearScreen()
		for _, r := range result {
			fmt.Println("\n", r)
		}
	case "6":
		fmt.Println("Exiting...")
		os.Exit(0)
	default:
		fmt.Println("Invalid option, please select again")
	}
}
