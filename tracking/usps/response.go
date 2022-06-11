package usps

import (
	"autonomous-service/dao"
	"autonomous-service/pkg/third-party/shipment/tracking"
	"context"
)

func NewShipmentResponse(context context.Context, trackingCode string, gsDao dao.GlobalSettingDao) tracking.History {
	shipmentResponse := getShipmentResponse(context, trackingCode, gsDao)
	return shipmentResponse
}

type shipmentResponse struct {
	StatusCode      string        `json:"statusCode"`
	StatusText      string        `json:"statusText"`
	TrackDetails    []TrackDetail `json:"trackDetails"`
	IsBcdnMultiView bool          `json:"isBcdnMultiView"`
	IsLoggedInUser  bool          `json:"isLoggedInUser"`
	TrackedDateTime string        `json:"trackedDateTime"`
}

type TrackDetail struct {
	ErrorCode                  string                 `json:"errorCode"`
	ErrorText                  string                 `json:"errorText"`
	RequestedTrackingNumber    string                 `json:"requestedTrackingNumber"`
	IsMobileDevice             bool                   `json:"isMobileDevice"`
	ScheduledDeliveryDate      string                 `json:"scheduledDeliveryDate"`
	ShipmentProgressActivities []ShipmentActivity     `json:"shipmentProgressActivities"`
	DeliveredDate              string                 `json:"deliveredDate"`
	DeliveredTime              string                 `json:"deliveredTime"`
	PackageStatus              string                 `json:"packageStatus"`
	PackageStatusCode          string                 `json:"packageStatusCode"`
	TrackingNumber             string                 `json:"trackingNumber"`
	ProgressBarType            string                 `json:"progressBarType"`
	AdditionalInformation      map[string]interface{} `json:"additionalInformation"`
}

type AdditionalData struct {
	Weight int `json:"weight"`
}

type ShipmentActivity struct {
	Date         string `json:"date"`
	Time         string `json:"time"`
	Location     string `json:"location"`
	ActivityScan string `json:"activityScan"`
}

func (_this *shipmentResponse) MakeShipmentResponse(packageCode string) (result tracking.TrackShipmentResponse, err error) {
	return
}
