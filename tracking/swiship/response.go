package swiship

import (
	"context"
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
	DisplayableShipMethod string          `json:"displayableShipMethod"`
	EstimatedArrivalDate  time.Time       `json:"estimatedArrivalDate"`
	Histories             []trackingEvent `json:"trackingEvents"`
	TrackingCode          string          `json:"trackingNumber"`
	TransitState          string          `json:"transitState"`
}

type trackingEvent struct {
	EventAddress     string    `json:"eventAddress"`
	EventCode        string    `json:"eventCode"`
	EventDate        time.Time `json:"eventDate"`
	EventDescription string    `json:"eventDescription"`
}

// MakeShipmentResponse func
func (_this *shipmentResponse) MakeTrackShipmentResponse(packageCode string) (result tracking.TrackShipmentResponse, err error) {
	defaultErr := errors.New("Not found")
	if !_this.isSuccessful() {
		err = defaultErr
		return
	}

	if _this.TrackingCode == "" {
		err = defaultErr
		return
	}

	deliveredDate := _this.getEstimatedDeliveryDate()
	result.IsDelivered = _this.isDelivered()
	result.IsInTransit = _this.isInTransit()
	result.IsPickup = false
	result.IsPrePickup = false
	result.DeliveredDate = deliveredDate
	result.TrackingCode = _this.TrackingCode

	if result.PackageCode == "" {
		result.PackageCode = packageCode
	}

	if !deliveredDate.IsZero() {
		result.DeliveredDateString = deliveredDate.Format("Monday, January 2, 2006")
	}

	history := _this.getShipmentHistory()
	result.SetHistory(history)

	return
}

// isDelivered ...
func (_this *shipmentResponse) isDelivered() bool {
	return _this.TransitState == deliveredStatus
}

// isInTransit ...
func (_this *shipmentResponse) isInTransit() bool {
	return _this.TransitState == inTransitStatus
}

// getEstimatedDeliveryDate func
func (_this *shipmentResponse) getEstimatedDeliveryDate() (result time.Time) {
	return _this.EstimatedArrivalDate
}

// getShipmentHistory func
func (_this *shipmentResponse) getShipmentHistory() map[int64]tracking.TrackShipmentHistoryResponse {
	formatHistory := map[int64]tracking.TrackShipmentHistoryResponse{}
	for _, item := range _this.Histories {
		if !item.hasDate() {
			continue
		}
		keyDate := item.formatBeginDate().Unix()
		historyItem := tracking.TrackShipmentHistoryItemResponse{
			Time:   item.EventDate.Format(tracking.TimeFormat),
			Status: item.statusLabel(),
		}

		formatHistoryItem, ok := formatHistory[keyDate]
		if ok {
			formatHistoryItem.Details = append(formatHistoryItem.Details, historyItem)
			if formatHistoryItem.Location == "" && item.EventAddress != "" {
				formatHistoryItem.Location = item.EventAddress
			}

			formatHistory[keyDate] = formatHistoryItem
		} else {
			historyDetails := make([]tracking.TrackShipmentHistoryItemResponse, 0)
			historyDetails = append(historyDetails, historyItem)
			formatHistory[keyDate] = tracking.TrackShipmentHistoryResponse{
				Date:       item.EventDate,
				Location:   item.EventAddress,
				Details:    historyDetails,
				DateString: item.EventDate.Format("Monday, January 2, 2006"),
			}
		}
	}

	return formatHistory
}

// isSuccessful func: check the status of the response
func (_this *shipmentResponse) isSuccessful() bool {
	return _this.TransitState != ""
}

// formatBeginDate func
func (_this *trackingEvent) formatBeginDate() (date time.Time) {
	date = time.Date(_this.EventDate.Year(), _this.EventDate.Month(), _this.EventDate.Day(), 0, 0, 0, 0, time.Local)
	return
}

// hasDate func: check date of history item
func (_this *trackingEvent) hasDate() bool {
	return !_this.EventDate.IsZero()
}

func (_this *trackingEvent) statusLabel() string {
	return trackingEventDescription[_this.EventCode]
}
