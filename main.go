package main

import (
	"fmt"
	"log"
	"net/http"
	"time"

	"github.com/cedricleblond35/HTTPserver/cmd/collector"
	"github.com/cedricleblond35/HTTPserver/cmd/publisher"
	"github.com/prometheus/client_golang/prometheus/promhttp"
	"github.com/spf13/cobra"
)

var (
	version  = "dev"
	revision = "none"
	date     = "unknown"

	// CLI flags
	metrics string
)

func main() {
	c := &cobra.Command{
		Use:     "serversTest",
		Version: fmt.Sprintf("%s (revision %.7s @ %s)", version, revision, date),
		PersistentPreRun: func(_ *cobra.Command, _ []string) {
			if metrics == "" {
				return
			}

			server := http.NewServeMux()
			server.HandleFunc("/ping", func(w http.ResponseWriter, _ *http.Request) {
				w.WriteHeader(http.StatusOK)
			})
			server.Handle("/metrics", promhttp.Handler())

			go func() {
				log.Println("Healthcheck & metrics will start when service is ready")

				time.Sleep(20 * time.Second) // Delay the start of healthcheck for blackbox_exporter
				log.Printf("Healthcheck & metrics listening on %s", metrics)
				http.ListenAndServe(metrics, server) //nolint: errcheck
			}()
		},
	}

	// Our subcommands
	c.PersistentFlags().StringVarP(&metrics, "metrics", "", "", "Enable metrics & ping endpoint (e.g. ':5555', 'localhost:5555')")
	c.AddCommand(publisher.Command)
	c.AddCommand(collector.Command)

}
