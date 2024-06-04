package main

import (
	"fmt"
	"strings"

	"github.com/gdamore/tcell/v2"
	"github.com/rivo/tview"
)

// NewInputBox returns a new inputbox primitive.
func NewInputBox() *InputBox {
	textView := tview.NewTextView().SetDynamicColors(true).SetRegions(true)
	return &InputBox{
		Box:      tview.NewBox().SetBorder(true).SetTitle("Input"),
		input:    make(chan rune),
		textView: textView,
	}
}

// InputHandler returns a function that processes keyboard input events for the InputBox.
// It listens for rune input (character keys) and sends the rune (character) to the input channel of the InputBox.
func (i *InputBox) InputHandler() func(event *tcell.EventKey) *tcell.EventKey {
	return func(event *tcell.EventKey) *tcell.EventKey {
		if event.Key() == tcell.KeyRune {
			i.input <- event.Rune()
		}
		return event
	}
}

// Draw renders the InputBox on the provided screen.
func (i *InputBox) Draw(screen tcell.Screen) {
	i.Box.DrawForSubclass(screen, i.Box)
	x, y, width, height := i.GetInnerRect()
	i.textView.SetRect(x, y, width, height)
	i.textView.Draw(screen)
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

const logo = ` 
  _____ _____    ___________   
 /     \\__  \  /  ___/\__  \  
|  Y Y  \/ __ \_\___ \  / __ \_
|__|_|  (____  /____  >(____  /
      \/     \/     \/      \/ 
`

const (
	subtitle   = `masa oracle client`
	navigation = `[yellow]use keys or mouse to navigate`
	mouse      = `[green]v0.0.6-beta`
)

// Splash shows the app info
func Splash() (content tview.Primitive) {
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
			// nothing todo
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
			AddItem(logoBox, logoWidth, 1, false).
			AddItem(tview.NewBox(), 0, 1, false), logoHeight, 1, false).
		AddItem(frame, 0, 10, false)

	return flex
}
