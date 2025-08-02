package main

import (
	"c2/db"
	"common/utils"
	"flag"
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

	r := http.NewServeMux()

	initRouting(r)

	utils.MaybeFatal(http.ListenAndServe(net.JoinHostPort(*host, *port), r))
}
