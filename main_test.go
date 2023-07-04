package main

import (
	"testing"

	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/test"
)

// test the inline structure and function.
func Test_makeUI(t *testing.T) {
	var testCfg Config
	edit, preview := testCfg.makeUI()

	test.Type(edit, "Hello")

	if preview.String() != "Hello" {
		t.Error("Failed -- did not find expected value in preview")
	}
}

func Test_RunAPP(t *testing.T) {
	var testCfg Config
	testApp := test.NewApp()
	testWin := testApp.NewWindow("Test Markdown")

	// make UI
	edit, preview := testCfg.makeUI()

	// create menu items
	testCfg.createMenuItems(testWin)

	// show the window
	testWin.SetContent(container.NewHSplit(edit, preview))

	testApp.Run()

	test.Type(edit, "Some Test")
	if preview.String() != "Some Test" {
		t.Error("failed")
	}
}
