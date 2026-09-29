package types

import "github.com/prometheus/client_golang/prometheus"

type Metrics struct {
	ChatsInFlight    prometheus.Gauge
	TimeToFirstToken prometheus.Histogram
	TokensStreamed   prometheus.Counter
}
