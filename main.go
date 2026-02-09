package main

import (
	"fmt"
	"os"
	"path/filepath"
	"strings"

	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/app"
	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/dialog"
	"fyne.io/fyne/v2/widget"
)

func main() {
	a := app.New()
	w := a.NewWindow("File rename tool")
	w.Resize(fyne.NewSize(400, 300))

	// UI Elements
	pathLabel := widget.NewLabel("No folder selected")
	var selectedPath string

	oldPrefixEntry := widget.NewEntry()
	oldPrefixEntry.SetPlaceHolder("Old Prefix (e.g., img_)")

	newPrefixEntry := widget.NewEntry()
	newPrefixEntry.SetPlaceHolder("New Prefix (e.g., photo_)")

	statusLabel := widget.NewLabel("Ready")

	// Browse Button
	browseBtn := widget.NewButton("Browse Folder", func() {
		dialog.ShowFolderOpen(func(uri fyne.ListableURI, err error) {
			if err != nil || uri == nil {
				return
			}
			selectedPath = uri.Path()
			pathLabel.SetText(selectedPath)
		}, w)
	})

	// Logic: Recursive Rename
	runBtn := widget.NewButton("Rename All", func() {
		if selectedPath == "" || oldPrefixEntry.Text == "" {
			statusLabel.SetText("Error: Select folder and Old Prefix.")
			return
		}

		count := 0
		oldP := oldPrefixEntry.Text
		newP := newPrefixEntry.Text

		err := filepath.WalkDir(selectedPath, func(path string, d os.DirEntry, err error) error {
			if err != nil {
				return err
			}
			if d.IsDir() {
				return nil // Skip directories themselves
			}

			// Check file name (not full path) for prefix
			fileName := d.Name()
			if strings.HasPrefix(fileName, oldP) {
				dir := filepath.Dir(path)
				newName := newP + strings.TrimPrefix(fileName, oldP)
				newPath := filepath.Join(dir, newName)

				renameErr := os.Rename(path, newPath)
				if renameErr != nil {
					fmt.Println("Error renaming:", renameErr)
				} else {
					count++
				}
			}
			return nil
		})

		if err != nil {
			statusLabel.SetText("Error: " + err.Error())
		} else {
			statusLabel.SetText(fmt.Sprintf("Success! Renamed %d files.", count))
		}
	})

	// Layout
	content := container.NewVBox(
		widget.NewLabel("Target Folder:"),
		browseBtn,
		pathLabel,
		widget.NewSeparator(),
		widget.NewLabel("Old Prefix:"), oldPrefixEntry,
		widget.NewLabel("New Prefix:"), newPrefixEntry,
		widget.NewSeparator(),
		runBtn,
		statusLabel,
	)

	w.SetContent(content)
	w.ShowAndRun()
}
