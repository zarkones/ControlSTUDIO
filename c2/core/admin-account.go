package core

import (
	"c2/core/permissions"
	"c2/models"
	"c2/repos/permissionsRepo"
	"encoding/hex"
	"errors"
	"time"

	access "github.com/zarkones/ControlACCESS"
	cry "github.com/zarkones/xena-crypto"
)

var (
	ErrUsername = errors.New("username is empty or too long, must be equal or less to 32 characters")
)

// CreateAdminOperator creates a really powerful administrative account with all permissions assigned.
// Recommendation is to have only one such account, and make second account for yourself which is less
// powerful, so that you can use it as a daily driver.
func CreateAdminOperator(username string) (operator models.Operator, hexEncodedPrivateKey string, err error) {
	if len(username) == 0 || len(username) > 32 {
		return operator, "", ErrUsername
	}

	privateKey, err := cry.GenPrivKey()
	if err != nil {
		return operator, "", err
	}

	serializedPrivateKey, err := cry.PrivKeyToPEM(privateKey)
	if err != nil {
		return operator, "", err
	}

	serializedPublicKey, err := cry.PubKeyToPEM(&privateKey.PublicKey)
	if err != nil {
		return operator, "", err
	}

	for _, key := range permissions.AllPermissions {
		if err := permissionsRepo.Insert(&access.Permission{
			UserID:    username,
			Key:       key,
			CreatedAt: time.Now(),
		}); err != nil {
			return operator, "", err
		}
	}

	return models.Operator{
			UserID:       username,
			PublicKeyHex: hex.EncodeToString([]byte(serializedPublicKey)),
		},
		hex.EncodeToString([]byte(serializedPrivateKey)),
		nil
}
