package metrics

import (
	"github.com/prometheus/client_golang/prometheus"
)

var (
	EmailsSentTotal = prometheus.NewCounter(
		prometheus.CounterOpts{
			Name: "notification_emails_sent_total",
			Help: "Total number of successfully sent emails",
		},
	)

	EmailsFailedTotal = prometheus.NewCounter(
		prometheus.CounterOpts{
			Name: "notification_emails_failed_total",
			Help: "Total number of failed email attempts",
		},
	)

	EmailSendDuration = prometheus.NewHistogram(
		prometheus.HistogramOpts{
			Name:    "notification_send_duration_seconds",
			Help:    "Time taken to send notification emails",
			Buckets: prometheus.DefBuckets,
		},
	)
)

func Register() {
	prometheus.MustRegister(EmailsSentTotal)
	prometheus.MustRegister(EmailsFailedTotal)
	prometheus.MustRegister(EmailSendDuration)
}
