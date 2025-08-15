package main

import (
	"c2/core/listeners"
	"c2/db"
	"c2/repos"
	"common/utils"
	"flag"
	"fmt"
	"net"
	"net/http"
)

var (
	host   = flag.String("host", "0.0.0.0", "host of the c2 server")
	port   = flag.String("port", "8000", "port of the c2 server")
	dbName = flag.String("db-name", "database.sqlite", "path to the database")
)

func main() {
	flag.Parse()

	utils.MaybeFatal(db.Init(*dbName))

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

	utils.MaybeFatal(http.ListenAndServe(net.JoinHostPort(*host, *port), r))
}
