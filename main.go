package main

import (
	"fmt"
	"io/ioutil"
	"path/filepath"
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
		// Retrieve the path for the selected favorite
		if path, ok := favoritesMap[mainText]; ok {
			updateTable(detailsTable, path)
		} else {
			detailsTable.Clear()
			detailsTable.SetCell(0, 0, tview.NewTableCell("Invalid favorite selected").SetTextColor(tview.Styles.PrimaryTextColor))
		}
	})

	// Favorites pane
	favoritesPane := tview.NewFlex()
	favoritesPane.SetTitle("Favorites")
	favoritesPane.SetDirection(tview.FlexRow)
	favoritesPane.SetBorder(true)
	favoritesPane.AddItem(leftList, 0, 1, true)

	// Details pane
	detailsPane := tview.NewFlex()
	detailsPane.SetTitle("Details")
	favoritesPane.SetDirection(tview.FlexRow)
	detailsPane.SetBorder(true)
	detailsPane.AddItem(detailsTable, 0, 1, true)

	// Create a flex layout to arrange the lists side-by-side.
	flex := tview.NewFlex().
		AddItem(favoritesPane, 0, 1, true).
		AddItem(detailsPane, 0, 3, false)

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
		list.AddItem(file.Name(), fullPath, 0, nil) // Display file name, store full path as secondary text
	}
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

// updateTable populates the details table with the contents of the specified directory.
func updateTable(table *tview.Table, dir string) {
	table.Clear()

	// Add header row
	headers := []string{"Name", "Date Modified", "Size", "Kind"}
	columnWidths := []int{30, 25, 15, 10} // Fixed widths for each column
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
		// File name
		table.SetCell(i+1, 0, tview.NewTableCell(formatCell(file.Name(), columnWidths[0])).
			SetTextColor(tview.Styles.PrimaryTextColor))

		// Date modified
		modTime := file.ModTime().Format(time.RFC1123)
		table.SetCell(i+1, 1, tview.NewTableCell(formatCell(modTime, columnWidths[1])).
			SetTextColor(tview.Styles.PrimaryTextColor))

		// File size
		size := fmt.Sprintf("%d bytes", file.Size())
		if file.IsDir() {
			size = "-"
		}
		table.SetCell(i+1, 2, tview.NewTableCell(formatCell(size, columnWidths[2])).
			SetTextColor(tview.Styles.PrimaryTextColor))

		// File kind
		kind := "File"
		if file.IsDir() {
			kind = "Directory"
		}
		table.SetCell(i+1, 3, tview.NewTableCell(formatCell(kind, columnWidths[3])).
			SetTextColor(tview.Styles.PrimaryTextColor))
	}

}

// formatCell truncates or pads a string to fit within a fixed width.
func formatCell(content string, width int) string {
	if len(content) > width {
		return content[:width-3] + "..." // Truncate and add ellipsis
	}
	return content + strings.Repeat(" ", width-len(content)) // Pad with spaces
}
