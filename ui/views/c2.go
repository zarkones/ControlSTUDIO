package views

import (
	"common/httpc"
	"encoding/json"
	"fmt"
	"os"
	"ui/core"
	"ui/state"

	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/dialog"
	"fyne.io/fyne/v2/widget"
	profiles "github.com/zarkones/ControlPROFILE"
)

const C2_PANEL_WINDOW_WIDTH = 960
const C2_PANEL_WINDOW_HEIGHT = 640

var c2WindowOpen = false

func C2Window() {
	if c2WindowOpen {
		return
	}

	w := core.App.NewWindow("Command & Control Server Panel")
	w.Resize(fyne.NewSize(C2_PANEL_WINDOW_WIDTH, C2_PANEL_WINDOW_HEIGHT))
	w.CenterOnScreen()

	w.SetContent(C2Display(&w))

	c2WindowOpen = true
	w.SetOnClosed(func() {
		c2WindowOpen = false
	})

	w.Show()
}

func C2Display(w *fyne.Window) fyne.CanvasObject {
	profileCreation := widget.NewLabel("C2 profile creation")

	nameInput := widget.NewEntry()
	nameInput.SetPlaceHolder("Name")

	descInput := widget.NewMultiLineEntry()
	descInput.SetPlaceHolder("Description")

	loadProfileBtn := widget.NewButton("Load Profile", func() {
		importProfile(
			w,
			func(profile profiles.Profile) {
				if err := httpc.InsertProfile(nameInput.Text, descInput.Text, profile); err != nil {
					Alert(err.Error())
					return
				}
				Notify("Listener Creation Success", "Creation of a listener service was successful")
				nameInput.SetText("")
				descInput.SetText("")
				var err error
				state.Listeners, err = httpc.GetListeners()
				if err != nil {
					fmt.Println("error httpc.GetListeners:", err)
					return
				}
			})
	})

	return container.NewVBox(
		profileCreation,
		nameInput,
		descInput,
		loadProfileBtn,
	)
}

func importProfile(w *fyne.Window, callback func(profile profiles.Profile)) {
	var fd *dialog.FileDialog
	fd = dialog.NewFileOpen(func(reader fyne.URIReadCloser, err error) {
		if err != nil {
			Alert(err.Error())
			return
		}

		if reader == nil {
			return
		}
		if reader.URI() == nil {
			return
		}

		rawPipeline, err := os.ReadFile(reader.URI().Path())
		if err != nil {
			Alert(err.Error())
			return
		}

		var profile profiles.Profile

		if err := json.Unmarshal(rawPipeline, &profile); err != nil {
			Alert(err.Error())
			return
		}

		callback(profile)
	}, *w)

	fd.Show()
}
