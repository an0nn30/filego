package main

import (
	"fmt"
	"io/ioutil"
	"path/filepath"

	"github.com/rivo/tview"
)

func main() {
	app := tview.NewApplication()
	leftList := tview.NewList().ShowSecondaryText(false)
	rightList := tview.NewList().ShowSecondaryText(false)

	// Get initial directory contents
	updateList(leftList, ".")

	// Handle left list selection changes
	leftList.SetChangedFunc(func(index int, mainText string, secondaryText string, shortcut rune) {
		updateList(rightList, mainText)
	})

	// Create flex layouts for the boxes (to hold the list and the title)
	leftFlex := tview.NewFlex().SetDirection(tview.FlexRow).
		AddItem(tview.NewTextView().SetText("Directories"), 1, 0, false). // Title
		AddItem(leftList, 0, 1, true) // List

	rightFlex := tview.NewFlex().SetDirection(tview.FlexRow).
		AddItem(tview.NewTextView().SetText("Files"), 1, 0, false). // Title
		AddItem(rightList, 0, 1, true) // List

	// Create a flex layout to arrange the boxes side-by-side.
	flex := tview.NewFlex().
		AddItem(leftFlex.SetBorder(true), 0, 1, true). // Left box gets 1/4 of the space
		AddItem(rightFlex.SetBorder(true), 0, 3, false) // Right box gets 3/4 of the space

	// Start the application
	if err := app.SetRoot(flex, true).Run(); err != nil {
		panic(err)
	}
}

// updateList populates the given list with the contents of the specified directory.
func updateList(list *tview.List, dir string) {
	list.Clear()
	files, err := ioutil.ReadDir(dir)
	if err != nil {
		list.AddItem(fmt.Sprintf("Error: %v", err), "", 0, nil)
		return
	}

	for _, file := range files {
		fullPath := filepath.Join(dir, file.Name())
		list.AddItem(fullPath, "", 0, nil)
	}
}
