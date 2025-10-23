package collector

import (
	"context"
	"log"
	"time"

	"github.com/linode/linodego"
	"github.com/prometheus/client_golang/prometheus"
)

type AccountCollector struct {
	client *linodego.Client

	up                 *prometheus.Desc
	info               *prometheus.Desc
	balance            *prometheus.Desc
	balance_uninvoiced *prometheus.Desc
}

func NewAccountCollector(client *linodego.Client) *AccountCollector {
	return &AccountCollector{
		client: client,
		up: prometheus.NewDesc(
			prometheus.BuildFQName(namespace, "account", "up"),
			"Was the last scrape of the Linode account endpoint successful.",
			nil,
			nil,
		),
		info: prometheus.NewDesc(
			prometheus.BuildFQName(namespace, "account", "info"),
			"Generic static information about the Linode account",
			[]string{"euuid", "active_since", "company", "country"},
			nil,
		),
		balance: prometheus.NewDesc(
			prometheus.BuildFQName(namespace, "account", "balance"),
			"Current account balance in USD (negative value indicates credit).",
			nil,
			nil,
		),
		balance_uninvoiced: prometheus.NewDesc(
			prometheus.BuildFQName(namespace, "account", "balance_uninvoiced"),
			"Estimated invoice in USD. Not final invoice balance. Transfer charges not included.",
			nil,
			nil,
		),
	}
}

func (c *AccountCollector) Describe(ch chan<- *prometheus.Desc) {
	ch <- c.up
	ch <- c.info
	ch <- c.balance
	ch <- c.balance_uninvoiced
}

func (c *AccountCollector) Collect(ch chan<- prometheus.Metric) {
	account, err := c.client.GetAccount(context.Background())
	if err != nil {
		log.Printf("Could not fetch account: %v", err)
		ch <- prometheus.MustNewConstMetric(c.up, prometheus.GaugeValue, 0)
		return
	}

	ch <- prometheus.MustNewConstMetric(c.up, prometheus.GaugeValue, 1)
	ch <- prometheus.MustNewConstMetric(
		c.info, prometheus.GaugeValue, 1,
		account.EUUID,
		account.ActiveSince.Format(time.RFC3339),
		account.Company,
		account.Country,
	)
	ch <- prometheus.MustNewConstMetric(c.balance, prometheus.GaugeValue, float64(account.Balance))
	ch <- prometheus.MustNewConstMetric(c.balance_uninvoiced, prometheus.GaugeValue, float64(account.BalanceUninvoiced))
}
