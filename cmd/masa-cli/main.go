package main

// masa cli v2
// do not use this code, it's a WIP

import (
	"fmt"
	"log"

	"github.com/charmbracelet/bubbles/textinput"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
	masa "github.com/masa-finance/masa-oracle/pkg"
)

type Styles struct {
	BorderColor lipgloss.Color
	InputField  lipgloss.Style
}

func DefaultStyles() *Styles {
	return &Styles{
		BorderColor: lipgloss.Color("36"),
		InputField: lipgloss.NewStyle().
			Foreground(lipgloss.Color("#ffcc00")).
			BorderStyle(lipgloss.NormalBorder()).
			Padding(1).
			Width(80),
	}
}

type model struct {
	index       int
	questions   []Questions
	width       int
	height      int
	answerField textinput.Model
	styles      *Styles
}

type Questions struct {
	question string
	answer   string
}

func NewQuestion(question string) Questions {
	return Questions{question: question}
}

func New(questions []Questions) *model {
	answerField := textinput.New()
	answerField.Placeholder = "Enter your answer"
	answerField.Focus()
	return &model{questions: questions, answerField: answerField, styles: DefaultStyles()}
}

func (m model) Init() tea.Cmd {
	return nil
}

func (m model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	var cmd tea.Cmd
	current := &m.questions[m.index]
	switch msg := msg.(type) {
	case tea.WindowSizeMsg:
		m.width = msg.Width
		m.height = msg.Height
	case tea.KeyMsg:
		switch msg.String() {
		case "enter":
			current.answer = m.answerField.Value()
			m.Next()
			m.answerField.SetValue("")
			return m, nil
		case "q", "ctrl+c":
			return m, tea.Quit
		}
	}
	m.answerField, cmd = m.answerField.Update(msg)
	return m, cmd
}

func (m model) View() string {
	if m.width == 0 {
		return "loading..."
	}
	return lipgloss.Place(
		m.width, m.height,
		lipgloss.Center,
		lipgloss.Center,
		lipgloss.JoinVertical(
			lipgloss.Center,
			m.questions[m.index].question,
			m.styles.InputField.Render(m.answerField.View()),
		),
	)
}

func (m *model) Next() {
	if m.index < len(m.questions)-1 {
		m.index++
	} else {
		m.index = 0
	}
}

func main() {

	// this needs to also be a node to use all the functionality
	var node *masa.OracleNode
	fmt.Println(node)

	questions := []Questions{
		NewQuestion("Connect to Oracle Node"),
		NewQuestion("Select LLM Model"),
		NewQuestion("Set Twitter Credentials"),
		NewQuestion("Analyze Sentiment from Tweets"),
		NewQuestion("Analyze Sentiment from Website"),
		NewQuestion("Oracle Nodes"),
	}
	m := New(questions)
	p := tea.NewProgram(m, tea.WithAltScreen())
	if _, err := p.Run(); err != nil {
		log.Fatal(err)
	}
}
