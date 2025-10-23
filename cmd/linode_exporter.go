package main

import (
	"context"
	"flag"
	"github.com/jb3/linode_exporter/collector"
	"github.com/linode/linodego"
	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/collectors"
	"github.com/prometheus/client_golang/prometheus/promhttp"
	"golang.org/x/oauth2"
	"log"
	"net/http"
	"os"
	"time"
)

var (
	listenAddress = flag.String("web.listen-address", ":9388", "Address to listen on for web interface and telemetry.")
	metricsPath   = flag.String("web.telemetry-path", "/metrics", "Path under which to expose metrics.")
	timeout       = flag.Duration("timeout", 10*time.Second, "Timeout for collecting metrics.")
	apiToken      = flag.String("linode.token", "", "Linode API token (can also use LINODE_TOKEN env var).")
)

func main() {
	flag.Parse()

	token := *apiToken
	if token == "" {
		token = os.Getenv("LINODE_TOKEN")
	}
	if token == "" {
		log.Fatal("Linode API token must be provided via -linode.token flag or LINODE_TOKEN environment variable")
	}

	tokenSource := oauth2.StaticTokenSource(&oauth2.Token{AccessToken: token})
	oauth2Client := oauth2.NewClient(context.Background(), tokenSource)
	linodeClient := linodego.NewClient(oauth2Client)

	registry := prometheus.NewRegistry()

	registry.MustRegister(collectors.NewGoCollector())
	registry.MustRegister(collectors.NewProcessCollector(collectors.ProcessCollectorOpts{}))

	linodeCollector := collector.NewLinodeCollector(&linodeClient, *timeout)
	registry.MustRegister(linodeCollector)

	http.Handle(*metricsPath, promhttp.HandlerFor(registry, promhttp.HandlerOpts{
		ErrorLog:      log.New(os.Stderr, "", log.LstdFlags),
		ErrorHandling: promhttp.ContinueOnError,
		Registry:      registry,
	}))

	http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "text/html")
		_, err := w.Write([]byte(`<html>
<head><title>Linode Exporter</title></head>
<body>
<h1>Linode Exporter</h1>
<p><a href='` + *metricsPath + `'>Metrics</a></p>
</body>
</html>`))

		if err != nil {
			log.Fatal("Could not return response to client for index route.")
		}
	})

	http.HandleFunc("/health", func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
		_, err := w.Write([]byte("OK"))

		if err != nil {
			log.Fatal("Could not return response to client for healthcheck.")
		}
	})

	log.Printf("Starting Linode exporter on %s", *listenAddress)
	log.Printf("Metrics available at %s%s", *listenAddress, *metricsPath)

	if err := http.ListenAndServe(*listenAddress, nil); err != nil {
		log.Fatal(err)
	}
}
