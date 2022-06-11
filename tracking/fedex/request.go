package fedex

import (
	"context"
	"gitlab.com/autonomous-ecm/bots/delivery-tracking-bot/pkg/carrier/fedex"
)

func getShipmentResponse(context context.Context, trackingCode string) *shipmentResponse {
	result := &shipmentResponse{}

	// GetFedExServiceInstance
	fedexAPI := fedex.GetFedExServiceInstance()
	result.detail = fedexAPI.TrackPackageForSaveShipmentHistories(trackingCode)

	return result
}
