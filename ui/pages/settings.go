package pages

import (
	"common/httpc"
	"encoding/hex"
	"ui/core"
	"ui/state"

	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/layout"
	"fyne.io/fyne/v2/widget"
	"github.com/golang-jwt/jwt/v5"
	cry "github.com/zarkones/xena-crypto"
)

func Settings() fyne.CanvasObject {
	c2Host := widget.NewEntry()
	c2Host.Text = httpc.BaseURL
	c2Host.OnChanged = func(s string) {
		httpc.BaseURL = s
		c2NodeLabel.SetText("C2: " + s)
	}

	username := widget.NewEntry()

	privKey := widget.NewMultiLineEntry()
	privKey.Disable()
	privKey.SetPlaceHolder("")
	privKey.Validator = func(s string) error {
		privateKeyPem, err := hex.DecodeString(s)
		if err != nil {
			return err
		}
		if _, err := cry.ImportPrivKeyPEM(privateKeyPem); err != nil {
			return err
		}
		return nil
	}
	privKey.OnChanged = func(s string) {
		privateKeyPem, err := hex.DecodeString(s)
		if err != nil {
			return
		}
		key, err := cry.ImportPrivKeyPEM(privateKeyPem)
		if err != nil {
			return
		}
		state.PrivateKey = key

		token := jwt.NewWithClaims(jwt.SigningMethodRS512, jwt.MapClaims{
			"id": username.Text,
		})

		stringifiedToken, err := token.SignedString(state.PrivateKey)
		if err != nil {
			privKey.SetValidationError(err)
			return
		}

		httpc.AuthHeaderValue = stringifiedToken
	}

	username.OnChanged = func(s string) {
		state.Username = s
		if len(s) != 0 {
			privKey.Enable()
		} else {
			privKey.Disable()
		}
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
			widget.NewLabel("Username:"),
			username,
			widget.NewLabel("(Hex -> PEM encoded) Private Key:"),
			privKey,
			layout.NewSpacer(),
			widget.NewLabel("ControlSTUDIO "+core.VERSION),
		),
	)
}
