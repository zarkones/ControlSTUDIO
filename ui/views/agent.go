package views

import (
	"c2/models"
	"common/httpc"
	"encoding/hex"
	"fmt"
	"os"
	"path/filepath"
	"sync"
	"time"
	"ui/core"

	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/storage"
	"fyne.io/fyne/v2/widget"
)

const AGENT_WINDOW_WIDTH = 960
const AGENT_WINDOW_HEIGHT = 640

var (
	agentWindows = sync.Map{}
)

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
	if _, ok := agentWindows.Load(agent.ID); ok {
		return
	}
	agentWindows.Store(agent.ID, true)

	w := core.App.NewWindow("Agent: " + agent.ID)
	w.Resize(fyne.NewSize(AGENT_WINDOW_WIDTH, AGENT_WINDOW_HEIGHT))
	w.CenterOnScreen()

	w.SetContent(AgentDisplay(agent, &w))

	w.SetOnClosed(func() {
		agentWindows.Delete(agent.ID)
	})

	w.Show()
}

func AgentDisplay(agent models.Agent, w *fyne.Window) fyne.CanvasObject {
	cmdInput := widget.NewEntry()
	cmdInput.SetPlaceHolder("Enter Command Here")

	// Fixed: Use slice to maintain message order
	messages := []models.Message{}
	messageIDs := map[string]int{} // Helper map to find message index by ID

	messagesTxt := widget.NewRichText()
	messagesTxt.Wrapping = fyne.TextWrapWord

	messagesTxtScroll := container.NewVScroll(messagesTxt)

	fullyScrolled := true
	messagesTxtScroll.OnScrolled = func(p fyne.Position) {
		fmt.Println(p, messagesTxtScroll.Size())

		// Get the scrollable content's total height
		contentHeight := messagesTxtScroll.Content.Size().Height
		// Get the visible area height of the scroll container
		visibleHeight := messagesTxtScroll.Size().Height
		// Get the current vertical scroll position
		scrollY := p.Y

		// Check if the scroll position is at or near the bottom
		if scrollY >= contentHeight-visibleHeight {
			fmt.Println("Reached the bottom!")
			fullyScrolled = true
		} else {
			fullyScrolled = false
		}
	}

	var after *time.Time = nil
	var before *time.Time = nil
	page := 0

	updateMsg := func() {

		refresh := false
		scrollToBottom := false
		oldMsgLen := len(messages)
		defer func() {
			messagesTxt.Segments = make([]widget.RichTextSegment, len(messages)*2)
			i := 0
			for _, msg := range messages {
				messagesTxt.Segments[i] = &widget.TextSegment{
					Style: widget.RichTextStyleSubHeading,
					Text:  "_> " + msg.Request,
				}
				i++

				if msg.Request == "/screenshot" {
					screenshotsDir := "./screenshots"
					screenshotPath := filepath.Join(screenshotsDir, msg.ID)
					os.Mkdir(screenshotsDir, 0777)

					f, err := os.Create(screenshotPath)
					if err != nil {
						messagesTxt.Segments[i] = &widget.TextSegment{
							Style: widget.RichTextStyleSubHeading,
							Text:  "local error: " + err.Error(),
						}
						i++
					}
					hexDecoded, err := hex.DecodeString(msg.Response)
					if err != nil {
						messagesTxt.Segments[i] = &widget.TextSegment{
							Style: widget.RichTextStyleSubHeading,
							Text:  "local error: " + err.Error(),
						}
						i++
					}
					f.Write([]byte(hexDecoded))
					f.Close()
					// defer os.Remove(screenshotPath)
					uri := storage.NewFileURI(screenshotPath)

					messagesTxt.Segments[i] = &widget.ImageSegment{
						Source: uri,
					}
					i++
				} else {
					messagesTxt.Segments[i] = &widget.TextSegment{
						Style: widget.RichTextStyleCodeInline,
						Text:  msg.Response + "\n",
					}
					i++
				}
			}

			if oldMsgLen != len(messages) {
				refresh = true
			}

			fyne.DoAndWait(func() {
				if scrollToBottom || fullyScrolled {
					messagesTxtScroll.ScrollToBottom()
				}
				if refresh {
					messagesTxtScroll.Refresh()
				}
			})
		}()

		newMessagesCtx, err := httpc.GetMessages(agent.ID, before, after, &page)
		if err != nil {
			fmt.Println("error httpc.GetMessages:", err)
			return
		}

		fyne.DoAndWait(func() {
			if len(newMessagesCtx.Messages) != 0 {
				after = &newMessagesCtx.After
				for i := len(newMessagesCtx.Messages) - 1; i >= 0; i-- {
					msg := newMessagesCtx.Messages[i]
					// Check if message already exists
					if _, exists := messageIDs[msg.ID]; !exists {
						// Add new message
						messageIDs[msg.ID] = len(messages)
						messages = append(messages, msg)
					}
				}
				if fullyScrolled {
					scrollToBottom = true
					refresh = true
				}
			}
		})

		messagesWithoutResponsesIDs := []string{}
		for _, msg := range messages {
			if len(msg.Response) != 0 {
				continue
			}
			messagesWithoutResponsesIDs = append(messagesWithoutResponsesIDs, msg.ID)
		}

		if len(messagesWithoutResponsesIDs) == 0 {
			return
		}

		msgMapMaybeWithResponses, err := httpc.GetMessagesByIDs(&messagesWithoutResponsesIDs)
		if err != nil {
			fmt.Println("error httpc.GetMessagesByIDs:", err)
			scrollToBottom = true
			refresh = true
			return
		}

		for _, msgMaybeWithResponse := range msgMapMaybeWithResponses {
			if index, exists := messageIDs[msgMaybeWithResponse.ID]; exists {
				if len(messages[index].Response) == 0 && len(msgMaybeWithResponse.Response) != 0 {
					scrollToBottom = true
					refresh = true

					messages[index] = models.Message{
						ID:        msgMaybeWithResponse.ID,
						AgentID:   msgMaybeWithResponse.AgentID,
						Request:   msgMaybeWithResponse.Request,
						Response:  msgMaybeWithResponse.Response,
						CreatedAt: msgMaybeWithResponse.CreatedAt,
						UpdatedAt: msgMaybeWithResponse.UpdatedAt,
					}
				}
			}
		}
	}

	breakApiFetchLoop := false
	(*w).SetOnClosed(func() {
		breakApiFetchLoop = true
	})

	go func() {
		updateMsg()
		// time.Sleep(time.Second)
		fyne.DoAndWait(messagesTxtScroll.ScrollToBottom)
		for range time.Tick(time.Second * 4) {
			if breakApiFetchLoop {
				return
			}
			updateMsg()
		}
	}()

	cmdInput.OnSubmitted = func(s string) {
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
				updateMsg()
				if fullyScrolled {
					fyne.DoAndWait(func() {
						messagesTxtScroll.ScrollToBottom()
					})
				}
				cmdInput.SetText("")
			}
		}()
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

		// Primary.
		container.NewBorder(
			// Top.
			nil,

			// Bottom.
			cmdInput,

			// Left.
			nil,

			// Right.
			nil,

			// Primary.
			messagesTxtScroll,
		),
	)
}
