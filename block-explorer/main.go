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
	"io/ioutil"
	"log"
	"net/http"
	"os"
	"path/filepath"

	"github.com/btcsuite/btcd/rpcclient"
	"github.com/btcsuite/btcutil"
	"github.com/elastic/go-elasticsearch"
	"github.com/farukterzioglu/btc-scraper/block-explorer/api"
	_ "github.com/farukterzioglu/btc-scraper/block-explorer/swagger" // Required for Swagger to explore models
	"github.com/farukterzioglu/btc-scraper/services"
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

	rpcCLient := initRpcClient()
	elasticService := initElasticService()
	router := initRouter(rpcCLient, elasticService)

	// Host Swagger UI
	fs := http.FileServer(http.Dir("./swaggerui/"))
	router.PathPrefix("/swaggerui/").Handler(http.StripPrefix("/swaggerui/", fs))

	// Start to listen
	log.Fatal(http.ListenAndServe(fmt.Sprintf(":%s", *portNumber), router))
}

func initRpcClient() *rpcclient.Client {
	btcdHomeDir := btcutil.AppDataDir("btcd", false)
	certs, err := ioutil.ReadFile(filepath.Join(btcdHomeDir, "rpc.cert"))
	if err != nil {
		log.Fatal(err)
	}
	connCfg := &rpcclient.ConnConfig{
		Host:         "localhost:18556",
		Endpoint:     "ws",
		User:         "myuser",
		Pass:         "SomeDecentp4ssw0rd",
		Certificates: certs,
	}
	client, err := rpcclient.New(connCfg, nil)
	if err != nil {
		log.Fatal(err)
	}

	return client
}

func initElasticService() *services.ElasticService {
	// TODO : get url from params
	// TODO : better exception handling

	cfg := elasticsearch.Config{
		Addresses: []string{
			"http://localhost:9200",
		},
	}
	es, err := elasticsearch.NewClient(cfg)
	if err != nil {
		log.Fatalf("Error creating the client: %s", err)
	}

	return services.NewElasticService(es)
}

func initRouter(client *rpcclient.Client, elasticService *services.ElasticService) (router *mux.Router) {
	router = mux.NewRouter()

	v1 := router.PathPrefix("/v1").Subrouter()

	btcRoutes := api.NewBtcRoutes(client)
	btcRoutes.RegisterBtcRoutes(v1, "/btc-rpc")

	btcDbRoutes := api.NewBtcDbRoutes(elasticService)
	btcDbRoutes.RegisterRoutes(v1, "/btc")

	return
}

func commonMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Add("Content-Type", "application/json")
		next.ServeHTTP(w, r)
	})
}
