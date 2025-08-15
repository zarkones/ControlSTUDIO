package pages

import (
	"common/httpc"
	"ui/core"

	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/layout"
	"fyne.io/fyne/v2/widget"
)

func Settings() fyne.CanvasObject {
	c2Host := widget.NewEntry()
	c2Host.Text = httpc.BaseURL
	c2Host.OnChanged = func(s string) {
		httpc.BaseURL = s
		c2NodeLabel.SetText("C2: " + s)
	}

	return container.NewBorder(
		// Top.
		nil,

		// Bottom.
		nil,

		// Left.
		nil,

		// Right.
		nil,

		container.NewVBox(
			widget.NewLabel("C2 Host:"),
			c2Host,
			layout.NewSpacer(),
			widget.NewLabel("ControlSTUDIO "+core.VERSION),
		),
	)
}
