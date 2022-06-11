package tracking

import (
	"sort"
	"time"
)

type TrackShipmentResponse struct {
	Status              string                         `json:"status"`
	StatusCode          string                         `json:"status_code"`
	DeliveredDate       time.Time                      `json:"delivered_date"`
	DeliveredDateString string                         `json:"delivered_date_string"`
	IsDelivered         bool                           `json:"is_delivered"`
	IsCancel            bool                           `json:"is_cancel"`
	IsPrePickup         bool                           `json:"is_pre_pickup"`
	IsPickup            bool                           `json:"is_pickup"`
	IsInTransit         bool                           `json:"is_in_transit"`
	History             []TrackShipmentHistoryResponse `json:"history"`
	Weight              int                            `json:"weight"`
	TrackingCode        string                         `json:"tracking_code"`
	PackageCode         string                         `json:"package_code"`
	TrackingLink        string                         `json:"tracking_link"`
}

type TrackShipmentHistoryResponse struct {
	Date       time.Time                          `json:"date"`
	DateString string                             `json:"date_string"`
	Location   string                             `json:"location"`
	Details    []TrackShipmentHistoryItemResponse `json:"details"`
}

type TrackShipmentHistoryItemResponse struct {
	Time   string `json:"time"`
	Status string `json:"status"`
}

type EstimatedDeliveryDate struct {
	DeliveryDate        time.Time `json:"delivery_date"`
	TransitTime         int       `json:"transit_time"`
	IsDefault           bool      `json:"is_default"`
	WarehouseLocationId int       `json:"warehouse_location_id"`
	LastScan            time.Time `json:"last_scan"`
}

func (_this *TrackShipmentResponse) SetHistory(formatHistory map[int64]TrackShipmentHistoryResponse) {
	history := make([]TrackShipmentHistoryResponse, 0)
	for _, item := range formatHistory {
		history = append(history, item)
	}

	// sort history by latest date
	sort.Slice(history, func(i, j int) bool {
		return history[i].Date.After(history[j].Date)
	})

	_this.History = history
}
