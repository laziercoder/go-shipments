package canpar

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
	EstimatedDeliveryDate string            `json:"estimated_delivery_date"`
	IsDelivered           bool              `json:"is_delivered"`
	TrackingCode          string            `json:"tracking_code"`
	Histories             []shipmentHistory `json:"histories"`
	Status                int               `json:"status"`
	Message               string            `json:"message"`
}

func (_this *shipmentResponse) MakeTrackShipmentResponse(packageCode string) (result tracking.TrackShipmentResponse, err error) {
	defaultErr := errors.New("Not found")
	if !_this.IsSuccessful() {
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
	// result.Weight = _this.GetWeight()
	result.TrackingCode = _this.TrackingCode

	if result.PackageCode == "" {
		result.PackageCode = packageCode
	}

	if !deliveredDate.IsZero() {
		result.DeliveredDateString = deliveredDate.Format("Monday, January 2, 2006")
	}

	// set history
	history := _this.GetShipmentHistory()
	result.SetHistory(history)

	return
}

func (_this *shipmentResponse) GetEstimatedDeliveryDate() (result time.Time) {
	if _this.EstimatedDeliveryDate == "" {
		return
	}

	result, _ = time.Parse(time.RFC3339, _this.EstimatedDeliveryDate)
	return
}

func (_this *shipmentResponse) GetShipmentHistory() map[int64]tracking.TrackShipmentHistoryResponse {
	formatHistory := map[int64]tracking.TrackShipmentHistoryResponse{}
	for _, item := range _this.Histories {
		if !item.HasDate() {
			continue
		}
		formatDate := item.Date
		keyDate := item.FormatBeginDate().Unix()
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

type shipmentHistory struct {
	Code     string    `json:"code"`
	Date     time.Time `json:"date"`
	Location string    `json:"location"`
	Status   string    `json:"status"`
}

func (_this *shipmentResponse) IsSuccessful() bool {
	return _this.Status == 1
}

func (_this *shipmentHistory) FormatBeginDate() (date time.Time) {
	// 2019-05-13T00:00:00-06:00
	currentDate := _this.Date
	date = time.Date(currentDate.Year(), currentDate.Month(), currentDate.Day(), 0, 0, 0, 0, time.Local)
	return
}

func (_this *shipmentHistory) HasDate() bool {
	return !_this.Date.IsZero()
}
