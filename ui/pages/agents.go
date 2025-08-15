package pages

import (
	"bytes"
	"c2/core/listeners"
	"c2/models"
	"common/httpc"
	"fmt"
	"image/png"
	"math"
	"slices"
	"time"
	"ui/state"
	"ui/static"
	"ui/views"

	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/canvas"
	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/layout"
	"fyne.io/fyne/v2/theme"
	"fyne.io/fyne/v2/widget"
	dia "fyne.io/x/fyne/widget/diagramwidget"
)

var diagramWidget = dia.NewDiagramWidget("ENGAGEMENT")

var listenerAvatar *canvas.Image
var agentAvatar *canvas.Image
var osLogoWindows *canvas.Image
var osLogoLinux *canvas.Image
var osLogoDarwin *canvas.Image

const POS_OFFSET = 300

func generatePositions(initial fyne.Position, numPoints int) []fyne.Position {
	positions := make([]fyne.Position, numPoints)

	angle := 360.0 / float64(numPoints)

	for i := 0; i < numPoints; i++ {
		radians := angle * float64(i) * (math.Pi / 180.0)

		x := initial.X + float32(math.Cos(radians)*POS_OFFSET)
		y := initial.Y + float32(math.Sin(radians)*POS_OFFSET)

		positions[i] = fyne.Position{X: x, Y: y}
	}

	return positions
}

type TargetView struct {
	fyne.Container
	Node *dia.DiagramNode
}

func (mv *TargetView) Tapped(e *fyne.PointEvent) {
	agents, _ := httpc.GetAgents()
	menuItems := make([]*fyne.MenuItem, len(agents)+1)
	menuItems[0] = fyne.NewMenuItem("Inspect", func() {
		// TODO: Select the target in the inspector.
	})
	for i, agent := range agents {
		menuItems[i+1] = fyne.NewMenuItem("Atk with: "+agent.Hostname, func() {
		})
	}

	m := fyne.NewMenu(
		"Attack Menu",
		menuItems...,
	)
	widget.ShowPopUpMenuAtPosition(
		m,
		fyne.CurrentApp().Driver().CanvasForObject(diagramWidget),
		e.AbsolutePosition,
	)
}

func newAgentNode(agent models.Agent) *fyne.Container {
	return container.NewVBox(
		// 	// Primary.
		container.NewHBox(
			layout.NewSpacer(),
			widget.NewLabel(agent.Hostname),
			layout.NewSpacer(),
		),
		// ),

		agentAvatar,

		// BOTTOM.
		container.NewHBox(
			// osLogo,

			widget.NewLabel(agent.IP),

			widget.NewButtonWithIcon("", theme.MoreVerticalIcon(), func() {
				views.AgentWindow(agent)
			}),
		),
	)
}

func newListenerNode(listener listeners.Listener) *fyne.Container {
	return container.NewVBox(
		listenerAvatar,

		container.NewHBox(
			widget.NewLabel("Listener: "+listener.GetHost()),
		),
	)
}

func Agents() fyne.CanvasObject {
	displayedAgentIDs := []string{}
	listenerNodes := map[string]*dia.DiagramNode{}

	// C2 sprite.
	c2Img, _ := png.Decode(bytes.NewReader(static.C2))
	c2Sprite := canvas.NewImageFromImage(c2Img)
	c2Sprite.FillMode = canvas.ImageFillOriginal
	// Listener Logo on agents.
	listenerAvatarImg, _ := png.Decode(bytes.NewReader(static.ListenerAvatar))
	listenerAvatar = canvas.NewImageFromImage(listenerAvatarImg)
	listenerAvatar.FillMode = canvas.ImageFillOriginal
	// Agent avatar.
	agentlogoImg, _ := png.Decode(bytes.NewReader(static.AgentAvatar))
	agentAvatar = canvas.NewImageFromImage(agentlogoImg)
	agentAvatar.FillMode = canvas.ImageFillOriginal
	// OS Avatars.
	osLogoWindowsImg, _ := png.Decode(bytes.NewReader(static.OsAvatarWindows))
	osLogoWindows = canvas.NewImageFromImage(osLogoWindowsImg)
	osLogoWindows.FillMode = canvas.ImageFillOriginal
	osLogoLinuxImg, _ := png.Decode(bytes.NewReader(static.OsAvatarLinux))
	osLogoLinux = canvas.NewImageFromImage(osLogoLinuxImg)
	osLogoLinux.FillMode = canvas.ImageFillOriginal
	osLogoDarwinImg, _ := png.Decode(bytes.NewReader(static.OsAvatarDarwin))
	osLogoDarwin = canvas.NewImageFromImage(osLogoDarwinImg)
	osLogoDarwin.FillMode = canvas.ImageFillOriginal

	// Main Diagram thingies.
	scrollContainer := container.NewScroll(diagramWidget)
	scrollContainer.Offset = fyne.NewPos(500, 500)

	// C2 Node.
	c2NodeID := "C2: " + httpc.BaseURL
	c2Node := dia.NewDiagramNode(
		diagramWidget,
		container.NewVBox(
			c2Sprite,
			container.NewHBox(
				widget.NewLabel(c2NodeID),

				widget.NewButtonWithIcon("", theme.MoreVerticalIcon(), func() {
					views.C2Window()
				}),
			),
		),
		c2NodeID,
	)

	c2NodeX := float32(300)
	c2NodeY := float32(300)
	c2Node.Move(fyne.NewPos(c2NodeX, c2NodeY))
	c2Node.SetProperties(dia.DiagramElementProperties{
		StrokeWidth: 0,
	})
	c2Node.Refresh()

	updateDiagram := func() {
		var err error
		state.Agents, err = httpc.GetAgents()
		if err != nil {
			fmt.Println("error httpc.GetAgents:", err)
			return
		}
		state.Listeners, err = httpc.GetListeners()
		if err != nil {
			fmt.Println("error httpc.GetListeners:", err)
			return
		}

		fyne.DoAndWait(func() {

			listenerPoints := generatePositions(c2Node.Position(), len(state.Listeners))

			for listenerIndex, listener := range state.Listeners {
				listenerNodeID := "LISTENER:" + listener.GetID()

				if _, ok := listenerNodes[listenerNodeID]; !ok {
					listenerNode := dia.NewDiagramNode(diagramWidget, nil, listenerNodeID)
					listenerNode.SetProperties(dia.DiagramElementProperties{
						StrokeWidth: 0,
					})
					listenerNode.Move(listenerPoints[listenerIndex])

					listenerNode.SetInnerObject(newListenerNode(listener))

					listenerLinkID := "NODE_LINK_LISTENER:" + c2NodeID + "->" + listenerNodeID
					listenerLink := dia.NewDiagramLink(diagramWidget, listenerLinkID)
					listenerLink.SetSourcePad(c2Node.GetEdgePad())
					listenerLink.SetTargetPad(listenerNode.GetEdgePad())
					listenerLink.AddSourceDecoration(dia.NewArrowhead())

					listenerNodes[listenerNodeID] = &listenerNode
				}

				for agentIndex, agent := range state.Agents {
					if slices.Contains(displayedAgentIDs, agent.ID) {
						continue
					}
					displayedAgentIDs = append(displayedAgentIDs, agent.ID)

					agentNodeID := "AGENT:" + agent.ID
					agentNode := dia.NewDiagramNode(diagramWidget, nil, agentNodeID)
					agentNode.SetProperties(dia.DiagramElementProperties{
						StrokeWidth: 0,
					})

					agentPoints := generatePositions((*listenerNodes[listenerNodeID]).Position(), len(state.Agents))
					agentNode.Move(agentPoints[agentIndex])

					agentNode.SetInnerObject(newAgentNode(agent))

					agentLinkID := "NODE_LINK_AGENT:" + c2NodeID + "->" + agentNodeID
					agentLink := dia.NewDiagramLink(diagramWidget, agentLinkID)
					agentLink.SetSourcePad((*listenerNodes[listenerNodeID]).GetEdgePad())
					agentLink.SetTargetPad(agentNode.GetEdgePad())
					agentLink.AddSourceDecoration(dia.NewArrowhead())
				}
			}
		})

	}

	go func() {
		updateDiagram()
		for range time.Tick(time.Second * 5) {
			// if state.AuthToken == "" {
			// continue
			// }
			updateDiagram()
		}
	}()

	return container.NewBorder(
		// views.AssetsToolbar(),
		nil,
		nil,
		nil,
		nil,
		scrollContainer,
	)
}
