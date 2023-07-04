package main

import (
	"io"
	"strings"

	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/app"
	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/dialog"
	"fyne.io/fyne/v2/storage"
	"fyne.io/fyne/v2/widget"
)

type Config struct {
	EditWidget    *widget.Entry
	PreviewWidget *widget.RichText
	CurrentFile   fyne.URI // use uri to represent the current file
	SaveMenuItem  *fyne.MenuItem
}

var cfg Config

// a simple markdown editor
func main() {
	// create fyne app
	a := app.New()
	// create the window for app
	win := a.NewWindow("Markdown Editor")
	// get the user interface
	edit, preview := cfg.makeUI()

	// create menu
	cfg.createMenuItems(win)

	// set the content for the window
	win.SetContent(container.NewHSplit(edit, preview))
	win.Resize(fyne.NewSize(800.0, 500.0))
	win.CenterOnScreen()

	// how to set the utf-8 encoding

	// show the window with the app
	win.ShowAndRun()
}

func (app *Config) makeUI() (*widget.Entry, *widget.RichText) {
	edit := widget.NewMultiLineEntry()
	preview := widget.NewRichTextFromMarkdown("")

	app.EditWidget = edit
	app.PreviewWidget = preview

	edit.OnChanged = preview.ParseMarkdown

	return edit, preview
}

func (app *Config) createMenuItems(win fyne.Window) {
	openMenuItem := fyne.NewMenuItem("Open...", app.openFunc(win))
	saveMenuItem := fyne.NewMenuItem("Save", app.saveFunc(win))

	// before you open a file, the save menu item is disabled
	app.SaveMenuItem = saveMenuItem
	app.SaveMenuItem.Disabled = true
	saveAsMenuItem := fyne.NewMenuItem("Save as...", app.saveAsFunc(win))

	fileMenu := fyne.NewMenu("File", openMenuItem, saveMenuItem, saveAsMenuItem)
	menu := fyne.NewMainMenu(fileMenu)

	win.SetMainMenu(menu)
}

var filter = storage.NewExtensionFileFilter([]string{".md", ".MD"})

func (app *Config) saveFunc(win fyne.Window) func() {
	return func() {
		// set app.CurrentFile when call the openFunc
		if app.CurrentFile != nil {
			write, err := storage.Writer(app.CurrentFile)
			if err != nil {
				dialog.ShowError(err, win)
				return
			}
			defer write.Close()

			_, err = write.Write([]byte(app.EditWidget.Text))
			if err != nil {
				dialog.ShowError(err, win)
				return
			}
		}
	}
}

func (app *Config) openFunc(win fyne.Window) func() {
	return func() {
		openDialog := dialog.NewFileOpen(func(read fyne.URIReadCloser, err error) {
			if err != nil {
				dialog.ShowError(err, win)
				return
			}

			// user cancel
			if read == nil {
				return
			}

			data, err := io.ReadAll(read)
			if err != nil {
				dialog.ShowError(err, win)
				return
			}

			defer read.Close()

			app.EditWidget.SetText(string(data))

			app.CurrentFile = read.URI()

			win.SetTitle(win.Title() + " - " + read.URI().Name())
		}, win)

		openDialog.SetFilter(filter)
		openDialog.Show()
	}
}

func (app *Config) saveAsFunc(win fyne.Window) func() {
	return func() {
		saveDialog := dialog.NewFileSave(func(write fyne.URIWriteCloser, err error) {
			if err != nil {
				dialog.ShowError(err, win)
				return
			}

			if write == nil {
				// user canceled
				return
			}

			if !strings.HasSuffix(strings.ToLower(write.URI().String()), ".md") {
				dialog.ShowInformation("Error", "You should save as a .md extension", win)
				return
			}

			write.Write([]byte(app.EditWidget.Text))
			app.CurrentFile = write.URI()

			defer write.Close()

			win.SetTitle(win.Title() + " - " + write.URI().Name())
			app.SaveMenuItem.Disabled = false
		}, win)

		// show the dialog
		saveDialog.SetFileName("untitled.md")
		saveDialog.SetFilter(filter)
		saveDialog.Show()
	}
}
