// Package shipment
// For handle business for shipment like track shipment, ...
package shipment

import (
	"context"
	"github.com/laziercoder/go-payments/carriers"
	"github.com/laziercoder/go-payments/tracking"
	"github.com/laziercoder/go-payments/tracking/ait"
	"github.com/laziercoder/go-payments/tracking/autonomous"
	"github.com/laziercoder/go-payments/tracking/canpar"
	"github.com/laziercoder/go-payments/tracking/fedex"
	"github.com/laziercoder/go-payments/tracking/swiship"
	"github.com/laziercoder/go-payments/tracking/ups"
)

type (
	HistoryConstructor = func(context context.Context, trackingCode string) tracking.History
)

var (
	constructorByCarrier = map[string]HistoryConstructor{
		carriers.Fedex:   fedex.NewShipmentResponse,
		carriers.UPS:     ups.NewShipmentResponse,
		carriers.AIT:     ait.NewShipmentResponse,
		carriers.Canpar:  canpar.NewShipmentResponse,
		carriers.Auto:    autonomous.NewShipmentResponse,
		carriers.Swiship: swiship.NewShipmentResponse,
	}
)
