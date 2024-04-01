package main

import (
	"os"
	"path/filepath"
	"runtime"

	"github.com/gdamore/tcell/v2"
	"github.com/joho/godotenv"
	"github.com/rivo/tview"
)

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
		SetText(" Welcome to the MASA Oracle Client ").
		SetTextAlign(tview.AlignCenter)

	content := Splash()

	mainFlex = tview.NewFlex().SetDirection(tview.FlexColumn).
		AddItem(content, 0, 1, false).
		AddItem(handleMenu(app, output), 0, 1, true).
		AddItem(output, 0, 3, false)

	output.SetBorder(true).SetBorderColor(tcell.ColorBlue)

	app.SetFocus(handleMenu(app, output))

	if err := app.SetRoot(mainFlex, true).EnableMouse(true).Run(); err != nil {
		panic(err)
	}
}
