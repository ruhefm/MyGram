//   Product Api:
//    version: 0.1
//    title: Product Api
//   Schemes: http, https
//   Host:
//   BasePath: /orders
//   BasePath: /items
//      Consumes:
//      - application/json
//   Produces:
//   - application/json
//   swagger:meta

package main

import (
	"mygram/database"
	"mygram/routers"
)

var PORT = ":8080"

func main() {
	database.StartDB()
	routers.StartServer().Run(PORT)
}
