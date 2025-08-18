package state

import (
	"c2/core/listeners"
	"c2/models"
	"crypto/rsa"
)

var (
	Username   string = ""
	PrivateKey *rsa.PrivateKey

	Agents    = []models.Agent{}
	Listeners = []listeners.Listener{}
)
