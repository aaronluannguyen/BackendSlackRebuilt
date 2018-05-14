package summary

import (
	"os"
	"net/http"
	"github.com/labstack/gommon/log"
)

func main() {
	addr := os.Getenv("ADDR")
	if len(addr) == 0 {
		addr = ":80"
	}

	mux := http.NewServeMux()
	mux.HandleFunc("/v1/summary", SummaryHandler)

	log.Printf("listening at http://%s...", addr)
	log.Fatal(http.ListenAndServe(addr, mux))
}