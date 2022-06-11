package fedex

import (
	"context"
	"errors"
	"github.com/laziercoder/go-payments/tracking"
	"gitlab.com/autonomous-ecm/bots/delivery-tracking-bot/pkg/carrier/fedex/soap_model"
	"sort"
)

var (
	_ tracking.History = (*shipmentResponse)(nil)
)

func NewShipmentResponse(context context.Context, trackingCode string) tracking.History {
	shipmentResponse := getShipmentResponse(context, trackingCode)
	return shipmentResponse
}

type shipmentResponse struct {
	detail *soap_model.TrackDetail
}

func (a *shipmentResponse) MakeTrackShipmentResponse(packageCode string) (result tracking.TrackShipmentResponse, err error) {
	if a.detail == nil {
		err = errors.New("Not found")
		return
	}

	deliveredDate := a.detail.GetDeliveryDate()
	result.IsDelivered = a.detail.IsDelivered()
	result.IsInTransit = a.detail.IsInTransit()
	result.IsPickup = a.detail.IsPickUp()
	result.IsPrePickup = a.detail.IsPrePickup()
	result.Status = a.detail.StatusDetail.Description
	result.StatusCode = a.detail.StatusDetail.Code
	result.DeliveredDate = deliveredDate
	result.Weight = int(a.detail.PackageWeight.Value)
	result.TrackingCode = a.detail.TrackingNumber

	if result.PackageCode == "" {
		result.PackageCode = packageCode
	}

	if !deliveredDate.IsZero() {
		result.DeliveredDateString = deliveredDate.Format("Monday, January 2, 2006")
	}
	// else if latestDeliveredPackage.DisplayEstDeliveryDateTime != "" {
	// 	result.DeliveredDateString = latestDeliveredPackage.DisplayEstDeliveryDateTime
	// }

	formatHistory := map[int64]tracking.TrackShipmentHistoryResponse{}
	for _, item := range a.detail.Events {
		if item.Date.IsZero() {
			continue
		}
		formatDate := item.BeginDate
		keyDate := formatDate.Unix()
		historyItem := tracking.TrackShipmentHistoryItemResponse{
			Time:   item.Date.Format(TimeFormat),
			Status: item.EventDescription,
		}

		scanLocation := item.Address.String()
		formatHistoryItem, ok := formatHistory[keyDate]
		if ok {
			formatHistoryItem.Details = append(formatHistoryItem.Details, historyItem)
			if formatHistoryItem.Location == "" && scanLocation != "" {
				formatHistoryItem.Location = scanLocation
			}

			formatHistory[keyDate] = formatHistoryItem
		} else {
			historyDetails := make([]tracking.TrackShipmentHistoryItemResponse, 0)
			historyDetails = append(historyDetails, historyItem)
			formatHistory[keyDate] = tracking.TrackShipmentHistoryResponse{
				Date:       formatDate,
				Location:   scanLocation,
				Details:    historyDetails,
				DateString: formatDate.Format("Monday, January 2, 2006"),
			}
		}
	}

	history := make([]tracking.TrackShipmentHistoryResponse, 0)
	for _, item := range formatHistory {
		history = append(history, item)
	}

	sort.Slice(history, func(i, j int) bool {
		return history[i].Date.After(history[j].Date)
	})
	result.History = history
	return
}
