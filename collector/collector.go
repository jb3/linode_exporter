package collector

import (
	"context"
	"github.com/linode/linodego"
	"github.com/prometheus/client_golang/prometheus"
	"time"
)

const (
	namespace = "linode"
)

type LinodeCollector struct {
	client    *linodego.Client
	collector prometheus.Collector
	timeout   time.Duration
}

func NewLinodeCollector(client *linodego.Client, timeout time.Duration) *LinodeCollector {
	return &LinodeCollector{
		client:    client,
		timeout:   timeout,
		collector: NewAccountCollector(client),
	}
}

func (lc *LinodeCollector) Describe(ch chan<- *prometheus.Desc) {
	lc.collector.Describe(ch)
}

func (lc *LinodeCollector) Collect(ch chan<- prometheus.Metric) {
	ctx, cancel := context.WithTimeout(context.Background(), lc.timeout)
	defer cancel()

	done := make(chan struct{})
	go func() {
		lc.collector.Collect(ch)
		close(done)
	}()

	select {
	case <-done:
		return
	case <-ctx.Done():
		return
	}
}
