package metrics

import (
	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/promauto"
)

var DeliveryMetrics = struct {
	NewDeliveryTotal    prometheus.Counter
	FinishDeliveryTotal prometheus.Counter
}{
	NewDeliveryTotal: promauto.NewCounter(prometheus.CounterOpts{
		Namespace: "shop",
		Name:      "new_delivery_total",
		Help:      "new_delivery_total",
	}),
	FinishDeliveryTotal: promauto.NewCounter(prometheus.CounterOpts{
		Namespace: "shop",
		Name:      "finish_delivery_total",
		Help:      "finish_delivery_total",
	}),
}

var ItemsMetrics = struct {
	ItemsAddedTotal prometheus.Counter
	ItemsDeleted    prometheus.Counter
}{
	ItemsAddedTotal: promauto.NewCounter(prometheus.CounterOpts{
		Namespace: "shop",
		Name:      "items_add_total",
		Help:      "items_add_total",
	}),
	ItemsDeleted: promauto.NewCounter(prometheus.CounterOpts{
		Namespace: "shop",
		Name:      "items_deleted_total",
		Help:      "items_deleted_total",
	}),
}

func init() { // 2
	DeliveryMetrics.FinishDeliveryTotal.Inc()
	DeliveryMetrics.NewDeliveryTotal.Inc()

	ItemsMetrics.ItemsDeleted.Inc()
	ItemsMetrics.ItemsDeleted.Inc()
}
