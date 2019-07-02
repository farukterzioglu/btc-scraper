// BlockExplorer API
//
// Block explorer for multiple crypto currency
//
//     Schemes: http
//     BasePath: /v1
//     Version: 1.0.0
//     License: MIT http://opensource.org/licenses/MIT
//     Contact: Faruk Terzioğlu <faruk.terzioglu@hotmail.com>
//     Host:
//
//     Consumes:
//     - application/json
//
//     Produces:
//     - application/json
//
//     Security:
//     - bearer
//
//     SecurityDefinitions:
//     bearer:
//          type: apiKey
//          name: Authorization
//          in: header
//
// swagger:meta
package main

import (
	"flag"
	"fmt"
	"log"
	"net/http"
	"os"

	"github.com/farukterzioglu/btc-scraper/block-explorer/api"
	_ "github.com/farukterzioglu/btc-scraper/block-explorer/swagger" // Required for Swagger to explore models
	kitlog "github.com/go-kit/kit/log"
	"github.com/gorilla/mux"
)

var (
	portNumber = flag.String("port", "8000", "HTTP listen port")
)

func main() {
	// Logging domain.
	var logger kitlog.Logger
	{
		logger = kitlog.NewLogfmtLogger(os.Stderr)
		logger = kitlog.With(logger, "ts", kitlog.DefaultTimestampUTC)
		logger = kitlog.With(logger, "caller", "block-explorer") //kitlog.DefaultCaller)
	}

	flag.Parse()
	logger.Log("Application port", *portNumber)

	router := initRouter()

	// Host Swagger UI
	fs := http.FileServer(http.Dir("./swaggerui/"))
	router.PathPrefix("/swaggerui/").Handler(http.StripPrefix("/swaggerui/", fs))

	// Start to listen
	log.Fatal(http.ListenAndServe(fmt.Sprintf(":%s", *portNumber), router))
}

func initRouter() (router *mux.Router) {
	router = mux.NewRouter()
	// TODO : Causes Swagger UI to be parsed as JSON instead of html
	// router.Use(commonMiddleware)

	v1 := router.PathPrefix("/v1").Subrouter()

	btcRoutes := api.NewBtcRoutes()
	btcRoutes.RegisterBtcRoutes(v1, "/btc")

	return
}

func commonMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Add("Content-Type", "application/json")
		next.ServeHTTP(w, r)
	})
}
