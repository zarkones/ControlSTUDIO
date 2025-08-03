package layouts

import (
	"ui/pages"

	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/container"
)

func Default() fyne.CanvasObject {
	tabs := container.NewAppTabs(
		container.NewTabItem("Agents", pages.Agents()),
		container.NewTabItem("Settings", pages.Settings()),
	)
	// tabs.SetTabLocation(container.TabLocationLeading)
	return tabs
}
