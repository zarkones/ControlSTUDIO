package main

import (
	"c2/core"
	"c2/core/listeners"
	"c2/db"
	"c2/repos"
	"c2/repos/operatorsRepo"
	"c2/repos/permissionsRepo"
	"common/utils"
	"crypto/rsa"
	"encoding/hex"
	"flag"
	"fmt"
	"log"
	"net"
	"net/http"
	"os"

	access "github.com/zarkones/ControlACCESS"
	cry "github.com/zarkones/xena-crypto"
)

var (
	host   = flag.String("host", "0.0.0.0", "host of the c2 server")
	port   = flag.String("port", "8000", "port of the c2 server")
	dbName = flag.String("db-name", "database.sqlite", "path to the database")

	userInsertAdmin = flag.Bool("create-admin", false, "creates an administrative operator account with all permissions, requires --username and --password arguments")
	userName        = flag.String("username", "", "specify username, used by --create-admin")
)

func main() {
	flag.Parse()

	utils.MaybeFatal(db.Init(*dbName))

	// Creation of administrative operator account.
	if *userInsertAdmin {
		operator, hexEncodedPrivateKey, err := core.CreateAdminOperator(*userName)
		if err != nil {
			log.Println("failed to create admin operator account:", err)
			os.Exit(1)
		}

		if err := operatorsRepo.Insert(&operator); err != nil {
			log.Println("failed to inser admin operator account into database:", err)
			os.Exit(1)
		}

		log.Println("Username:", operator.UserID)
		log.Println("Private Key:", hexEncodedPrivateKey)

		os.Exit(0)
	}

	profilesWithMeta, err := repos.GetProfiles()
	utils.MaybeFatal(err)
	for _, profileWithMeta := range profilesWithMeta {
		profile, err := profileWithMeta.GetProfile()
		utils.MaybeFatal(err)

		hosts := profile.GetHosts()

		for _, listenerService := range hosts {
			if err := listeners.Insert(
				listeners.Listener{
					ProfileID: profileWithMeta.ID,
					Address:   listenerService.Address,
					Port:      listenerService.Port,
				},
				&profile,
			); err != nil {
				fmt.Println(
					"failed to start a listener service",
					profileWithMeta.ID,
					listenerService.Address,
					listenerService.Port,
					"error:",
					err,
				)
				continue
			}
		}
	}

	r := http.NewServeMux()

	initRouting(r)

	access.HandlerGetPermissionsByUserID = func(userID string) ([]access.Permission, error) {
		dbPermissions, err := permissionsRepo.GetMultipleByUserID(userID)
		if err != nil {
			return nil, err
		}
		// permissions := make([]access.Permission, len(dbPermissions))
		// for i, p := range dbPermissions {
		// 	permissions[i] = access.Permission{
		// 		Key:       p.Key,
		// 		UserID:    p.UserID,
		// 		Metadata:  p.Metadata,
		// 		CreatedAt: p.CreatedAt,
		// 	}
		// }
		return dbPermissions, nil
	}

	access.HandlerGetUserPublicKey = func(userID string) (*rsa.PublicKey, error) {
		operator, err := operatorsRepo.Get(userID)
		if err != nil {
			return nil, err
		}
		pemEncodedPubKey, err := hex.DecodeString(operator.PublicKeyHex)
		if err != nil {
			return nil, err
		}
		return cry.ImportPubKeyPEM(pemEncodedPubKey)
	}

	utils.MaybeFatal(http.ListenAndServe(net.JoinHostPort(*host, *port), r))
}
