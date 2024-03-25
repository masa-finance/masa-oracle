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
	"path/filepath"
	"runtime"
	"strings"
	"sync"

	"github.com/gdamore/tcell/v2"
	"github.com/joho/godotenv"
	"github.com/multiformats/go-multiaddr"
	"github.com/rivo/tview"

	"github.com/libp2p/go-libp2p"
	pubsub "github.com/libp2p/go-libp2p-pubsub"
	"github.com/libp2p/go-libp2p/core/peer"
	quic "github.com/libp2p/go-libp2p/p2p/transport/quic"
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

// extractIPAddress takes a multi-address string as input and extracts the IP address from it.
// The function supports both "/ip4/" and "/dns/" multi-address formats.
// If the multi-address contains an IP address, it is returned directly.
// If the multi-address contains a DNS name, it is also returned as is.
// In case the multi-address does not follow the expected format or the IP address/DNS name is not found, an empty string is returned.
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

// openFile reads the content of a file specified by the filename 'f' and returns it as a string.
// If the file cannot be read, the function logs a fatal error and exits the program.
// Parameters:
// - f: The name of the file to read.
// Returns:
// - A string containing the content of the file.
// func openFile(f string) string {
// 	dat, err := os.ReadFile(f)
// 	if err != nil {
// 		log.Print(err)
// 		return ""
// 	}
// 	return string(dat)
// }

// saveFile writes the provided content to a file specified by the filename 'f'.
// It appends the content to the file if it already exists, or creates a new file with the content if it does not.
// The file is created with permissions set to 0755.
// Parameters:
// - f: The name of the file to which the content will be written.
// - content: The content to write to the file.
// func saveFile(f string, content string) {
// 	file, err := os.OpenFile(f, os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0755)
// 	file.WriteString(content + "\n")
// 	if err != nil {
// 		log.Println(err)
// 		return
// 	}
// }

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
		}
	}
}

// transcribe takes an audio file and a target text file as input.
// It uses the OpenAI API to transcribe the audio to text, then saves the text to the specified text file.
// func transcribe(audioFile string, txtFile string) error {
// 	key := os.Getenv("OPENAI_API_KEY")
// 	if key == "" {
// 		log.Println("OPENAI_API_KEY is not set. Please set the environment variable and try again.")
// 		return errors.New("OPENAI_API_KEY is not set")
// 	}
// 	client := openai.NewClient(key)
// 	ctx := context.Background()
// 	req := openai.AudioRequest{
// 		Model:    openai.Whisper1,
// 		FilePath: audioFile,
// 	}
// 	resp, err := client.CreateTranscription(ctx, req)
// 	if err != nil {
// 		fmt.Printf("Transcription error: %v\n", err)
// 		return err
// 	} else {
// 		saveFile(txtFile, resp.Text)
// 	}
// 	return nil
// }

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

// RadioButtons implements a simple primitive for radio button selections.
type RadioButtons struct {
	*tview.Box
	options       []string
	currentOption int
	onSelect      func(option string)
}

// NewRadioButtons returns a new radio button primitive.
func NewRadioButtons(options []string, onSelect func(option string)) *RadioButtons {
	return &RadioButtons{
		Box:      tview.NewBox(),
		options:  options,
		onSelect: onSelect,
	}
}

// Draw draws this primitive onto the screen.
func (r *RadioButtons) Draw(screen tcell.Screen) {
	r.Box.DrawForSubclass(screen, r)
	x, y, width, height := r.GetInnerRect()

	for index, option := range r.options {
		if index >= height {
			break
		}
		radioButton := "\u25ef" // Unchecked.
		if index == r.currentOption {
			radioButton = "\u25c9" // Checked.
		}
		line := fmt.Sprintf(`%s[white]  %s`, radioButton, option)
		tview.Print(screen, line, x, y+index, width, tview.AlignLeft, tcell.ColorYellow)
	}
}

// InputHandler returns the handler for this primitive.
func (r *RadioButtons) InputHandler() func(event *tcell.EventKey, setFocus func(p tview.Primitive)) {
	return r.WrapInputHandler(func(event *tcell.EventKey, setFocus func(p tview.Primitive)) {
		switch event.Key() {
		case tcell.KeyUp:
			r.currentOption--
			if r.currentOption < 0 {
				r.currentOption = 0
			}
		case tcell.KeyDown:
			r.currentOption++
			if r.currentOption >= len(r.options) {
				r.currentOption = len(r.options) - 1
			}
		case tcell.KeyEnter:
			if r.onSelect != nil {
				r.onSelect(r.options[r.currentOption]) // Call the onSelect callback with the selected option
			}
		}
	})
}

// MouseHandler returns the mouse handler for this primitive.
func (r *RadioButtons) MouseHandler() func(action tview.MouseAction, event *tcell.EventMouse, setFocus func(p tview.Primitive)) (consumed bool, capture tview.Primitive) {
	return r.WrapMouseHandler(func(action tview.MouseAction, event *tcell.EventMouse, setFocus func(p tview.Primitive)) (consumed bool, capture tview.Primitive) {
		x, y := event.Position()
		_, rectY, _, _ := r.GetInnerRect()
		if !r.InRect(x, y) {
			return false, nil
		}

		if action == tview.MouseLeftClick {
			setFocus(r)
			index := y - rectY
			if index >= 0 && index < len(r.options) {
				r.currentOption = index
				consumed = true
				if r.onSelect != nil {
					r.onSelect(r.options[r.currentOption]) // Call the callback with the selected option
					// Logic to close the RadioButtons view goes here
				}
			}
		}
		return
	})
}

// showMenu creates and returns the menu component.
func showMenu(app *tview.Application, output *tview.TextView) *tview.List {
	menu := tview.NewList().
		AddItem("Connect", "Connect to an oracle node.", '1', func() {
			handleOption(app, "1", output)
		}).
		AddItem("LLM Model", "Select an llm model to use.", '2', func() {
			handleOption(app, "2", output)
		}).
		AddItem("Twitter", "Set your Twitter/X credentials.", '3', func() {
			// handleOption(app, "3", output)
			output.SetText(appConfig.Model)
		}).
		AddItem("ChatGPT", "Chat with a helpful assistant.", '4', func() {

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

			// 	resp, err := gpt(prompt, userMessage)
			// 	if err != nil {
			// 		logrus.Errorf("%v", err)
			// 	}
			// 	fmt.Println(resp)
			// 	if os.Getenv("ELAB_KEY") != "" {
			// 		speak(resp)
			// 	}
			// }

		}).
		AddItem("Sentiment", "Analyze sentiment from tweets.", '5', func() {
			// handleOption(app, "2", output)
			output.SetText("Subscribe to the Sentiment Topic")
		})

	menu.AddItem("Quit", "Press to exit", 'q', func() {
		app.Stop()
	}).SetBorder(true).SetBorderColor(tcell.ColorBlue)

	return menu
}

// handleOption triggers actions based on user selection.
func handleOption(app *tview.Application, option string, output *tview.TextView) {

	scanner := bufio.NewScanner(os.Stdin)

	switch option {
	case "1":
		modalFlex := tview.NewFlex().SetDirection(tview.FlexRow)
		inputField := tview.NewInputField().
			SetLabel("Multiaddress: ").
			SetFieldWidth(20)

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
			SetRect(0, 0, 30, 5)

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
}

const logo = ` 
  _____ _____    ___________   
 /     \\__  \  /  ___/\__  \  
|  Y Y  \/ __ \_\___ \  / __ \_
|__|_|  (____  /____  >(____  /
      \/     \/     \/      \/ 
`

const (
	subtitle = `masa oracle client`
	// navigation = `[yellow] Up    [yellow]: Down    [yellow]Q[-]: Quit`
	navigation = `[yellow]use keys to navigate the menu`
	mouse      = `[yellow]or use your mouse`
)

// Splash shows the app info
func Splash() (content tview.Primitive) {
	// What's the size of the logo?
	lines := strings.Split(logo, "\n")
	logoWidth := 0
	logoHeight := len(lines)
	for _, line := range lines {
		if len(line) > logoWidth {
			logoWidth = len(line)
		}
	}
	logoBox := tview.NewTextView().
		SetTextColor(tcell.ColorGreen).
		SetDoneFunc(func(key tcell.Key) {
			// nextSlide()
		})
	fmt.Fprint(logoBox, logo)

	frame := tview.NewFrame(tview.NewBox()).
		SetBorders(0, 0, 0, 0, 0, 0).
		AddText(subtitle, true, tview.AlignCenter, tcell.ColorWhite).
		AddText("", true, tview.AlignCenter, tcell.ColorWhite).
		AddText(navigation, true, tview.AlignCenter, tcell.ColorDarkMagenta).
		AddText(mouse, true, tview.AlignCenter, tcell.ColorDarkMagenta)

	// Create a Flex layout that centers the logo and subtitle.
	flex := tview.NewFlex().
		SetDirection(tview.FlexRow).
		AddItem(tview.NewBox(), 0, 7, false).
		AddItem(tview.NewFlex().
			AddItem(tview.NewBox(), 0, 1, false).
			AddItem(logoBox, logoWidth, 1, true).
			AddItem(tview.NewBox(), 0, 1, false), logoHeight, 1, true).
		AddItem(frame, 0, 10, false)

	return flex
}

func main() {
	var err error
	_, b, _, _ := runtime.Caller(0)
	rootDir := filepath.Join(filepath.Dir(b), "../..")
	if _, _ = os.Stat(rootDir + "/.env"); !os.IsNotExist(err) {
		_ = godotenv.Load()
	}

	app := tview.NewApplication()

	output := tview.NewTextView().
		SetDynamicColors(true).
		SetText(" Welcome to the MASA Oracle Client ")

	content := Splash()

	mainFlex = tview.NewFlex().SetDirection(tview.FlexColumn).
		AddItem(content, 0, 1, true).
		AddItem(showMenu(app, output), 0, 1, true).
		AddItem(output, 0, 3, false)

	output.SetBorder(true).SetBorderColor(tcell.ColorBlue)

	app.SetFocus(showMenu(app, output))

	if err := app.SetRoot(mainFlex, true).EnableMouse(true).Run(); err != nil {
		panic(err)
	}
}
