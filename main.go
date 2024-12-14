package main

import (
	"fmt"
	"github.com/gdamore/tcell/v2"
	"io/ioutil"
	"strings"
	"time"

	"github.com/rivo/tview"
)

func main() {
	app := tview.NewApplication()
	leftList := tview.NewList().ShowSecondaryText(false)
	detailsTable := tview.NewTable()

	favoritesMap := make(map[string]string)
	defaultFavorites(leftList, &favoritesMap)

	// Handle left list selection changes
	leftList.SetChangedFunc(func(index int, mainText string, secondaryText string, shortcut rune) {
		if path, ok := favoritesMap[mainText]; ok {
			updateTable(detailsTable, path)
		} else {
			detailsTable.Clear()
			detailsTable.SetCell(0, 0, tview.NewTableCell("Invalid favorite selected").SetTextColor(tview.Styles.PrimaryTextColor))
		}
	})

	// Handle Enter key to change focus to the right pane
	leftList.SetSelectedFunc(func(index int, mainText string, secondaryText string, shortcut rune) {
		app.SetFocus(detailsTable)
	})

	// Add navigation and selection functionality to the table
	detailsTable.SetSelectable(true, false) // Rows selectable, not columns
	detailsTable.SetInputCapture(func(event *tcell.EventKey) *tcell.EventKey {
		switch event.Key() {
		case tcell.KeyEnter:
			row, _ := detailsTable.GetSelection()
			cell := detailsTable.GetCell(row, 0)
			if cell != nil {
				selectedItem := strings.TrimSpace(cell.Text)
				fmt.Printf("Selected item: %s\n", selectedItem)
			}
		case tcell.KeyEsc:
			app.SetFocus(leftList) // Return focus to the left pane
		}
		return event
	})

	// Favorites pane
	favoritesPane := tview.NewFlex()
	favoritesPane.SetTitle("Favorites")
	favoritesPane.SetDirection(tview.FlexRow)
	favoritesPane.SetBorder(true)
	favoritesPane.AddItem(leftList, 0, 1, true)

	// Details pane
	detailsPane := tview.NewFlex()
	detailsPane.SetDirection(tview.FlexRow)
	detailsPane.SetBorder(true)
	detailsPane.AddItem(detailsTable, 0, 1, true)

	// Create a flex layout to arrange the lists side-by-side.
	flex := tview.NewFlex().
		AddItem(favoritesPane, 0, 1, true).
		AddItem(detailsPane, 0, 3, false)

	// Start the application
	if err := app.SetRoot(flex, true).SetFocus(leftList).Run(); err != nil {
		panic(err)
	}
}

// updateTable populates the details table with the contents of the specified directory.
func updateTable(table *tview.Table, dir string) {
	table.Clear()

	// Add header row
	headers := []string{"Name", "Date Modified", "Size", "Kind"}
	columnWidths := []int{30, 25, 15, 10}
	for i, header := range headers {
		table.SetCell(0, i, tview.NewTableCell(formatCell(header, columnWidths[i])).
			SetTextColor(tview.Styles.SecondaryTextColor).
			SetSelectable(false).
			SetAlign(tview.AlignLeft))
	}

	// Read directory contents
	files, err := ioutil.ReadDir(dir)
	if err != nil {
		table.SetCell(1, 0, tview.NewTableCell(fmt.Sprintf("Error: %v", err)).
			SetTextColor(tview.Styles.PrimaryTextColor))
		return
	}

	// Populate table with file details
	for i, file := range files {
		table.SetCell(i+1, 0, tview.NewTableCell(formatCell(file.Name(), columnWidths[0])).
			SetTextColor(tview.Styles.PrimaryTextColor).
			SetAlign(tview.AlignLeft))
		modTime := file.ModTime().Format(time.RFC1123)
		table.SetCell(i+1, 1, tview.NewTableCell(formatCell(modTime, columnWidths[1])).
			SetTextColor(tview.Styles.PrimaryTextColor).
			SetAlign(tview.AlignLeft))
		size := fmt.Sprintf("%d bytes", file.Size())
		if file.IsDir() {
			size = "-"
		}
		table.SetCell(i+1, 2, tview.NewTableCell(formatCell(size, columnWidths[2])).
			SetTextColor(tview.Styles.PrimaryTextColor).
			SetAlign(tview.AlignRight))
		kind := "File"
		if file.IsDir() {
			kind = "Directory"
		}
		table.SetCell(i+1, 3, tview.NewTableCell(formatCell(kind, columnWidths[3])).
			SetTextColor(tview.Styles.PrimaryTextColor).
			SetAlign(tview.AlignLeft))
	}
}

// formatCell truncates or pads a string to fit within a fixed width.
func formatCell(content string, width int) string {
	if len(content) > width {
		return content[:width-3] + "..."
	}
	return content + strings.Repeat(" ", width-len(content))
}

// defaultFavorites initializes the favorites list and map.
func defaultFavorites(list *tview.List, favoritesMap *map[string]string) {
	list.Clear()
	(*favoritesMap)["Desktop"] = "/Users/dustin/Desktop"
	(*favoritesMap)["Documents"] = "/Users/dustin/Documents"
	(*favoritesMap)["Pictures"] = "/Users/dustin/Pictures"
	(*favoritesMap)["Dustin"] = "/Users/dustin"

	for favorite, path := range *favoritesMap {
		list.AddItem(favorite, path, 0, nil)
	}
}
