package main

import (
	"bufio"
	"bytes"
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"github.com/gdamore/tcell/v2"
	"github.com/joho/godotenv"
	"github.com/rivo/tview"
	"io"
	"log"
	"net/http"
	"os"
	"os/exec"
	"path/filepath"
	"runtime"
	"strings"
	"sync"

	pubsub "github.com/libp2p/go-libp2p-pubsub"
	"github.com/sashabaranov/go-openai"
	"github.com/sirupsen/logrus"
)

type AppConfig struct {
	Address         string
	Model           string
	TwitterUser     string
	TwitterPassword string
}

var appConfig = AppConfig{}

var mainFlex *tview.Flex

type Gossip struct {
	Content  string
	Metadata map[string]string
}

type SpeakRequest struct {
	Text          string `json:"text"`
	VoiceSettings struct {
		Stability       float64 `json:"stability"`
		SimilarityBoost float64 `json:"similarity_boost"`
	} `json:"voice_settings"`
}

type SubscriptionHandler struct {
	Gossips     []Gossip
	GossipTopic *pubsub.Topic
	mu          sync.Mutex
}

// clearScreen clears the terminal screen by executing the appropriate command based on the operating system.
// It uses "cls" command for Windows and "clear" for Unix-like systems.
func clearScreen() {
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

func extractIPAddress(multiAddr string) string {
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

// openfile reads the content of a file specified by the filename 'f' and returns it as a string.
// If the file cannot be read, the function logs a fatal error and exits the program.
// Parameters:
// - f: The name of the file to read.
// Returns:
// - A string containing the content of the file.
func openfile(f string) string {
	dat, err := os.ReadFile(f)
	if err != nil {
		log.Print(err)
		return ""
	}
	return string(dat)
}

// saveFile writes the provided content to a file specified by the filename 'f'.
// It appends the content to the file if it already exists, or creates a new file with the content if it does not.
// The file is created with permissions set to 0755.
// Parameters:
// - f: The name of the file to which the content will be written.
// - content: The content to write to the file.
func saveFile(f string, content string) {
	file, err := os.OpenFile(f, os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0755)
	file.WriteString(content + "\n")
	if err != nil {
		log.Println(err)
		return
	}
}

// gpt takes a prompt and a user message, sends them to the OpenAI API, and returns the generated response.
// It utilizes the OPENAI_API_KEY environment variable for authentication.
// Parameters:
// - prompt: A string containing the initial prompt for the AI.
// - user_message: A string containing the user's message to the AI.
// Returns:
// - A string containing the AI's response.
// - An error if the request to the OpenAI API fails.
func gpt(prompt string, userMessage string) (string, error) {
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

// speak sends a given text response to the ElevenLabs API for text-to-speech conversion,
// then plays the resulting audio. It uses the ELAB_KEY environment variable for API authentication.
func speak(response string) {
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

	req, err := http.NewRequest(http.MethodPost, "https://api.elevenlabs.io/v1/text-to-speech/ErXwobaYiN019PkySvjV/stream", bytes.NewBuffer(buf))
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
		}
	}
}

// transcribe takes an audio file and a target text file as input.
// It uses the OpenAI API to transcribe the audio to text, then saves the text to the specified text file.
func transcribe(audioFile string, txtFile string) error {
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
		saveFile(txtFile, resp.Text)
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

// showMenu displays the secretive menu options to the user.
func showMenuOld() {
	// Colors
	yellow := "\033[33m"
	// blue := "\033[34m"
	// red := "\033[31m"
	green := "\033[32m"
	reset := "\033[0m"

	asciiArt := ` 
  _____ _____    ___________   
 /     \\__  \  /  ___/\__  \  
|  Y Y  \/ __ \_\___ \  / __ \_
|__|_|  (____  /____  >(____  /
      \/     \/     \/      \/ 
`
	fmt.Println(green + asciiArt + reset)

	fmt.Println(yellow + "MASA ORACLE NODE CLIENT" + reset)

	fmt.Println("-----------------------")
	fmt.Println("1. Connect to the Masa Oracle Network")
	fmt.Println("2. Choose your model")
	fmt.Println("3. Set your Twitter/X account credentials")
	fmt.Println("4. GPTChat")
	fmt.Println("5. Subscribe to the Sentiment Topic")
	fmt.Println("6. Quit")
	fmt.Println("")
	fmt.Println("Select an option:")
}

// showMenu creates and returns the menu component.
func showMenu(app *tview.Application, output *tview.TextView) *tview.List {
	menu := tview.NewList().
		AddItem("Connect to the Masa Oracle Network", "", 'c', func() {
			handleOption(app, "1", output)
		}).
		AddItem("Choose your model", "", 'm', func() {
			// handleOption(app, "2", output)
			output.SetText("Choose a model Claude or GPT4")
		}).
		AddItem("Quit", "Press to exit", 'q', func() {
			app.Stop()
		})

	// Example of how to handle menu selection that updates the output.
	// This is a placeholder to demonstrate interaction between the menu and the output.
	menu.AddItem("Show Message", "Press to show a message in the output", 's', func() {
		output.SetText("This is a message shown in the output panel.")
	}).SetBorder(true).SetBorderColor(tcell.ColorGoldenrod)

	return menu
}

// handleOption triggers actions based on user selection.
func handleOption(app *tview.Application, option string, output *tview.TextView) {

	scanner := bufio.NewScanner(os.Stdin)

	switch option {
	case "1":
		// Use a modal dialog to get the Masa node multiaddress
		inputField := tview.NewInputField().SetLabel("Enter the Masa node multiaddress: ")
		modal := tview.NewModal().
			SetText("Enter the Masa node multiaddress:").
			AddButtons([]string{"OK", "Cancel"}).
			SetDoneFunc(func(buttonIndex int, buttonLabel string) {
				if buttonLabel == "OK" {
					multiaddress := inputField.GetText()
					if multiaddress == "" {
						output.SetText("A multiaddress was not entered. Please enter the masa node multiaddress and try again.")
					} else {
						// Here you would normally handle the connection, but for simplicity, we just update the output
						output.SetText(fmt.Sprintf("Connecting to: %s", multiaddress))
						// Simulate connection logic...
					}
				}
				app.SetRoot(mainFlex, true) // Return to main view
			})
		app.SetRoot(modal, false) // Show modal dialog

		//fmt.Print("Enter the Masa node multiaddress: ")
		//scanner.Scan()
		//if scanner.Text() == "" {
		//	fmt.Println("A multiaddress was not entered. Please enter the masa node multiaddress and try again.")
		//	break
		//}
		//
		//appConfig.Address = scanner.Text()
		//
		//maddr, err := multiaddr.NewMultiaddr(appConfig.Address)
		//if err != nil {
		//	logrus.Errorf("%v", err)
		//}
		//
		//// Create a libp2p host to connect to the Masa node
		//host, err := libp2p.New(libp2p.NoSecurity, libp2p.Transport(quic.NewTransport))
		//if err != nil {
		//	logrus.Errorf("%v", err)
		//}
		//
		//// Extract the peer ID from the multiaddress
		//peerInfo, err := peer.AddrInfoFromP2pAddr(maddr)
		//if err != nil {
		//	logrus.Errorf("%v", err)
		//}
		//
		//// Connect to the peer
		//if err := host.Connect(context.Background(), *peerInfo); err != nil {
		//	logrus.Errorf("%v", err)
		//}
		//fmt.Println("Successfully connected to node")
	case "2":
		fmt.Println("Choose your model:\n1. Claude\n2. GPT4")
		scanner.Scan()
		modelChoice := scanner.Text()
		switch modelChoice {
		case "1":
			appConfig.Model = "Claude"
		case "2":
			appConfig.Model = "GPT4"
		default:
			fmt.Println("Invalid model selected.")
		}
		fmt.Printf("You selected: %s\n", appConfig.Model)
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

			resp, err := gpt(prompt, userMessage)
			if err != nil {
				logrus.Errorf("%v", err)
			}
			fmt.Println(resp)
			if os.Getenv("ELAB_KEY") != "" {
				speak(resp)
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

		uri := "http://" + extractIPAddress(appConfig.Address) + ":8080/analyzeSentiment"
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
		clearScreen()
		for _, r := range result {
			fmt.Println("\n", r)
		}
	case "6":
		fmt.Println("Exiting...")
		os.Exit(0)
	default:
		fmt.Println("Invalid option, please select again")
	}
	return
}

func main() {
	var err error
	_, b, _, _ := runtime.Caller(0)
	rootDir := filepath.Join(filepath.Dir(b), "../..")
	if _, _ = os.Stat(rootDir + "/.env"); !os.IsNotExist(err) {
		_ = godotenv.Load()
	}

	// Colors
	// yellow := "\033[33m"
	// blue := "\033[34m"
	// red := "\033[31m"
	// green := "\033[32m"
	// reset := "\033[0m"

	//asciiArt := `
	// _____ _____    ___________
	///     \\__  \  /  ___/\__  \
	//|  Y Y  \/ __ \_\___ \  / __ \_
	//|__|_|  (____  /____  >(____  /
	//     \/     \/     \/      \/
	//`
	// fmt.Println(green + asciiArt + reset)

	app := tview.NewApplication()

	output := tview.NewTextView().
		SetDynamicColors(true).
		SetText("Welcome to the MASA Oracle Node Client.\nSelect an option from the menu.")

	// Create the grid layout.
	//grid := tview.NewGrid().
	//	SetRows(0).
	//	SetColumns(30, 0).SetBorders(true).
	//	AddItem(showMenu(app, output), 0, 0, 1, 1, 0, 0, true).
	//	AddItem(output, 0, 1, 1, 1, 0, 0, false)
	//
	//mainFlex = tview.NewFlex().SetDirection(tview.FlexRow).
	//	AddItem(grid, 0, 1, false)

	output.SetBorder(true).SetBorderColor(tcell.ColorGoldenrod)

	mainFlex = tview.NewFlex().SetDirection(tview.FlexColumn).
		AddItem(showMenu(app, output), 0, 1, true).
		AddItem(output, 0, 3, false)

	// mainFlex.SetBorder(true).SetTitle(" MASA ORACLE CLIENT ").SetBorderColor(tcell.ColorGoldenrod)

	app.SetFocus(showMenu(app, output))

	if err := app.SetRoot(mainFlex, true).EnableMouse(true).Run(); err != nil {
		panic(err)
	}

	//reader := bufio.NewReader(os.Stdin)
	//for {
	//	clearScreen()
	//	showMenu()
	//	option, _ := reader.ReadString('\n')
	//	option = option[:len(option)-1]
	//	handleOption(option)
	//	fmt.Println("Press 'Enter' to continue...")
	//	_, _ = reader.ReadString('\n')
	//}
}
