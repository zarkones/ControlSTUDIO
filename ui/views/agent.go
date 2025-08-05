package views

import (
	"bytes"
	"c2/models"
	"common/httpc"
	"fmt"
	"image/png"
	"time"
	"ui/core"
	"ui/static"

	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/canvas"
	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/layout"
	"fyne.io/fyne/v2/theme"
	"fyne.io/fyne/v2/widget"
)

const AGENT_WINDOW_WIDTH = 960
const AGENT_WINDOW_HEIGHT = 640
const FILES_DIR = "xenalang-scripts"

type ModifiedRichText struct {
	fyne.Container
	Node *widget.RichText
}

func (mv *ModifiedRichText) Tapped(e *fyne.PointEvent) {
	menuItems := make([]*fyne.MenuItem, 1)
	menuItems[0] = fyne.NewMenuItem("Copy", func() {
	})
	m := fyne.NewMenu(
		"MSG_MENU",
		menuItems...,
	)
	widget.ShowPopUpMenuAtPosition(
		m,
		fyne.CurrentApp().Driver().CanvasForObject(mv),
		e.AbsolutePosition,
	)
}

func AgentWindow(agent models.Agent) {
	w := core.App.NewWindow("XENA: " + agent.Hostname)
	w.Resize(fyne.NewSize(AGENT_WINDOW_WIDTH, AGENT_WINDOW_HEIGHT))
	w.CenterOnScreen()

	w.SetContent(AgentDisplay(agent, &w))

	w.Show()
}

func AgentDisplay(agent models.Agent, w *fyne.Window) fyne.CanvasObject {
	cmdInput := widget.NewEntry()
	cmdInput.SetPlaceHolder("Enter Command Here")

	messagesTxt := widget.NewRichText()
	messagesTxt.Wrapping = fyne.TextWrapWord

	messagesTxtScroll := container.NewVScroll(messagesTxt)

	lastMsgsCount := 0
	updatedOnRespID := ""
	updateMsg := func() {
		msgCtx, err := httpc.GetMessages(agent.ID, nil, nil, nil)
		if err != nil {
			fmt.Println("error httpc.GetMessages:", err)
			return
		}

		fyne.Do(func() {
			messagesTxt.Segments = make([]widget.RichTextSegment, len(msgCtx.Messages)*2)

			i := 0
			for _, msg := range msgCtx.Messages {
				messagesTxt.Segments[i] = &widget.TextSegment{
					Style: widget.RichTextStyleSubHeading,
					Text:  "_> " + msg.Request,
				}
				i++

				messagesTxt.Segments[i] = &widget.TextSegment{
					Style: widget.RichTextStyleCodeInline,
					Text:  msg.Response + "\n",
				}
				i++
			}

			// Scroll to bottom if latest message got a response.
			if len(msgCtx.Messages) != 0 && msgCtx.Messages[len(msgCtx.Messages)-1].ID != updatedOnRespID && len(msgCtx.Messages[len(msgCtx.Messages)-1].Response) != 0 {
				updatedOnRespID = msgCtx.Messages[len(msgCtx.Messages)-1].ID
				// messagesTxt.Refresh()
				messagesTxtScroll.ScrollToBottom()
			} else {
				// Don't refresh the messages text if no new messages are present and no response was received to the latest message.
				if lastMsgsCount != len(msgCtx.Messages) {
					// messagesTxt.Refresh()
				}
			}
			// Scroll to bottom if there is new message.
			if lastMsgsCount != len(msgCtx.Messages) {
				messagesTxtScroll.ScrollToBottom()
				lastMsgsCount = len(msgCtx.Messages)
			}

			messagesTxt.Refresh()
		})
	}

	sendBtn := widget.NewButtonWithIcon("SEND", theme.MailSendIcon(), func() {
		go func() {
			if cmdInput.Text == "" {
				return
			}

			if agent.ID == "" {
				Alert("exception: Agent ID Empty")
				return
			}

			if err := httpc.InsertMessage(agent.ID, cmdInput.Text); err != nil {
				Alert("Failed To Send Message, exception:" + err.Error())
			} else {
				fyne.Do(func() {
					updateMsg()
					cmdInput.SetText("")
				})
			}
		}()
	})

	go func() {
		updateMsg()
		messagesTxtScroll.ScrollToBottom()
		for range time.Tick(time.Second * 4) {
			updateMsg()
		}
	}()

	// inspector := ToolsInspector()
	// toolsTable := NewToolsTable(agent, &inspector)

	img, _ := png.Decode(bytes.NewReader(static.XenaAvatar))
	avatar := canvas.NewImageFromImage(img)
	avatar.FillMode = canvas.ImageFillOriginal

	return container.NewBorder(
		// Top.
		container.NewStack(
			container.NewHBox(
				avatar,
				widget.NewLabel("Hostname: "+agent.Hostname),
				widget.NewLabel("OS: "+agent.OS+" "+agent.Arch),
				widget.NewLabel("IP: "+agent.IP),
			),
		),

		// Bottom.
		nil,

		// Left.
		nil,

		// Right.
		nil,

		// Primary.
		container.NewBorder(
			// Top.
			nil,

			// Bottom.
			cmdInput,

			// Left.
			nil,

			// Right.
			container.NewVBox(layout.NewSpacer(), sendBtn),

			// Primary.
			messagesTxtScroll,
		),
	)
}
