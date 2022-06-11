package autonomous

import (
	"context"
	"encoding/json"
	"errors"
	"github.com/laziercoder/go-payments/tracking"
	"time"
)

var (
	_ tracking.History = (*shipmentResponse)(nil)
)

func NewShipmentResponse(context context.Context, trackingCode string) tracking.History {
	shipmentResponse := getShipmentResponse(context, trackingCode)
	return shipmentResponse
}

type shipmentResponse struct {
	Success bool           `json:"success"`
	Orders  []shipmentAuto `json:"orders"`
}

type shipmentAuto struct {
	Success bool             `json:"success"`
	OrderNo string           `json:"orderNo"`
	Data    shipmentAutoData `json:"data"`
}

type shipmentAutoData struct {
	Status    string               `json:"status"`
	StartTime shipmentAutoDataTime `json:"startTime"`
	EndTime   shipmentAutoDataTime `json:"endTime"`
	Form      json.RawMessage      `json:"form"`
}

type shipmentAutoDataTime struct {
	UnixTimestamp int64  `json:"unixTimestamp"`
	UtcTime       string `json:"utcTime"`
	LocalTime     string `json:"localTime"`
}

func (_this *shipmentResponse) MakeTrackShipmentResponse(packageCode string) (result tracking.TrackShipmentResponse, err error) {
	if !_this.Success {
		err = errors.New("Not found")
		return
	}
	if result.PackageCode == "" {
		result.PackageCode = packageCode
	}

	orders := _this.Orders
	if len(orders) == 0 {
		err = errors.New("Not found")
		return
	}

	item := orders[0]
	if !item.Success {
		err = errors.New("Not found")
		return
	}

	result.Status = item.Data.Status
	result.TrackingCode = item.OrderNo
	// result.TrackingCode = "AUTO751460"
	history := make([]tracking.TrackShipmentHistoryResponse, 0)

	temp := tracking.TrackShipmentHistoryItemResponse{}
	switch item.Data.Status {
	case "scheduled", "servicing", "failed":
		{
			temp.Status = "Shipped"
		}
	case "success":
		{
			temp.Status = "Delivered"
		}
	default:
		{
			temp.Status = "Order confirm"
		}
	}
	switch item.Data.Status {
	case "servicing":
		{
			temp.Time = time.Unix(item.Data.StartTime.UnixTimestamp, 0).Format("15:04")
		}
	case "failed", "success":
		{
			if item.Data.EndTime.UnixTimestamp > 0 {
				temp.Time = time.Unix(item.Data.EndTime.UnixTimestamp, 0).Format("15:04")
			} else {
				temp.Time = time.Unix(item.Data.StartTime.UnixTimestamp, 0).Format("15:04")
			}
		}
	}
	historyDetails := make([]tracking.TrackShipmentHistoryItemResponse, 0)
	historyDetails = append(historyDetails, temp)
	historyItem := tracking.TrackShipmentHistoryResponse{
		DateString: time.Unix(item.Data.StartTime.UnixTimestamp, 0).Format("Monday, January 2, 2006"),
		Location:   "AutoExpress",
		Details:    historyDetails,
		Date:       time.Unix(item.Data.StartTime.UnixTimestamp, 0),
	}

	history = append(history, historyItem)
	result.History = history
	return
}
