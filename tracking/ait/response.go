package ait

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

// shipmentResponse struct
type shipmentResponse struct {
	Histories             []shipmentHistory `json:"histories"`
	IsDelivered           bool              `json:"is_delivered"`
	Message               string            `json:"message"`
	Status                int               `json:"status"`
	TrackingCode          string            `json:"tracking_code"`
	EstimatedDeliveryDate string            `json:"estimated_delivery_date"` // format: 20060102
}

// shipmentHistory struct
type shipmentHistory struct {
	Code     string `json:"code"`
	Date     string `json:"date"` // format: Mon, 02 Jan 2006 15:04:05 MST
	Location string `json:"location"`
	Status   string `json:"status"`
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

	deliveredDate := _this.GetEstimatedDeliveryDate()
	result.IsDelivered = _this.IsDelivered
	result.IsInTransit = false
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

	history := _this.GetShipmentHistory()
	result.SetHistory(history)

	return
}

// GetEstimatedDeliveryDate func
func (_this *shipmentResponse) GetEstimatedDeliveryDate() (result time.Time) {
	if _this.EstimatedDeliveryDate == "" {
		return
	}

	result, _ = time.Parse(formatEstimatedDeliveryDate, _this.EstimatedDeliveryDate)
	return
}

// GetShipmentHistory func
func (_this *shipmentResponse) GetShipmentHistory() map[int64]tracking.TrackShipmentHistoryResponse {
	formatHistory := map[int64]tracking.TrackShipmentHistoryResponse{}
	for _, item := range _this.Histories {
		if !item.hasDate() {
			continue
		}
		formatDate, _ := time.Parse(aitFormatHistoriesDate, item.Date)
		keyDate := item.formatBeginDate().Unix()
		historyItem := tracking.TrackShipmentHistoryItemResponse{
			Time:   formatDate.Format(tracking.TimeFormat),
			Status: item.Status,
		}

		formatHistoryItem, ok := formatHistory[keyDate]
		if ok {
			formatHistoryItem.Details = append(formatHistoryItem.Details, historyItem)
			if formatHistoryItem.Location == "" && item.Location != "" {
				formatHistoryItem.Location = item.Location
			}

			formatHistory[keyDate] = formatHistoryItem
		} else {
			historyDetails := make([]tracking.TrackShipmentHistoryItemResponse, 0)
			historyDetails = append(historyDetails, historyItem)
			formatHistory[keyDate] = tracking.TrackShipmentHistoryResponse{
				Date:       formatDate,
				Location:   item.Location,
				Details:    historyDetails,
				DateString: formatDate.Format("Monday, January 2, 2006"),
			}
		}
	}

	return formatHistory
}

// isSuccessful func: check the status of the response
func (_this *shipmentResponse) isSuccessful() bool {
	return _this.Status == 1
}

// formatBeginDate func
func (_this *shipmentHistory) formatBeginDate() (date time.Time) {
	currentDate, _ := time.Parse(aitFormatHistoriesDate, _this.Date)
	date = time.Date(currentDate.Year(), currentDate.Month(), currentDate.Day(), 0, 0, 0, 0, time.Local)
	return
}

// hasDate func: check date of history item
func (_this *shipmentHistory) hasDate() bool {
	return !(_this.Date == "")
}
