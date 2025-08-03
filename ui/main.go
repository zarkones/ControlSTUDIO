package main

import (
	"ui/core"
	"ui/layouts"

	"fyne.io/fyne/v2"
)

func main() {
	core.App.Settings().SetTheme(&mainTheme{})
	core.MainW.Resize(fyne.NewSize(core.WIN_WIDTH, core.WIN_HEIGHT))
	core.MainW.CenterOnScreen()
	core.MainW.SetContent(layouts.Default())
	core.MainW.ShowAndRun()
}
