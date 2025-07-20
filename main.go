package main

import (
	"fmt"
	"log"
	"net/http"
	"os"

	"github.com/PeterBocan/p-demo/api"
	"github.com/PeterBocan/p-demo/ledger"
)

func main() {
	port := os.Getenv("DEMO_PORT")
	if port == "" {
		port = "6734"
	}

	ledger := ledger.New()
	server := api.NewServer(ledger)

	log.Printf("listening on port %s", port)
	http.ListenAndServe(fmt.Sprintf(":%s", port), server)
}
