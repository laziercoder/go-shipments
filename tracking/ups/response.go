package ups

import (
	"context"
	"errors"
	"fmt"
	"github.com/laziercoder/go-payments/tracking"
	"sort"
	"strconv"
	"strings"
	"time"

	"gitlab.com/autonomous-ecm/backend/go-autonomous/pkg/logger"
)

func NewShipmentResponse(context context.Context, trackingCode string) tracking.History {
	shipmentResponse := getShipmentResponse(context, trackingCode)
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

func (_this *shipmentResponse) MakeTrackShipmentResponse(packageCode string) (result tracking.TrackShipmentResponse, err error) {
	defaultErr := errors.New("Not Found")
	if !_this.IsSuccessful() {
		err = defaultErr
		return
	}

	latestDeliveredPackage := _this.GetLatestDeliveredPackage()
	if latestDeliveredPackage.TrackingNumber == "" {
		err = defaultErr
		return
	}

	deliveredDate := latestDeliveredPackage.GetDeliveryDate()
	result.IsDelivered = latestDeliveredPackage.isDelivered()
	result.IsInTransit = latestDeliveredPackage.isInTransit()
	result.IsPickup = latestDeliveredPackage.isPickup()
	result.IsPrePickup = latestDeliveredPackage.IsPrePickup()
	result.Status = latestDeliveredPackage.PackageStatus
	result.StatusCode = latestDeliveredPackage.PackageStatusCode
	result.DeliveredDate = deliveredDate
	result.Weight = _this.GetWeight()
	result.TrackingCode = latestDeliveredPackage.TrackingNumber

	if result.PackageCode == "" {
		result.PackageCode = packageCode
	}

	if !deliveredDate.IsZero() {
		result.DeliveredDateString = deliveredDate.Format("Monday, January 2, 2006")
	}

	formatHistory := map[int64]tracking.TrackShipmentHistoryResponse{}
	for _, item := range latestDeliveredPackage.ShipmentProgressActivities {
		if !item.HasData() {
			continue
		}
		formatDate := item.FormatDate()
		keyDate := item.FormatBeginDate().Unix()
		historyItem := tracking.TrackShipmentHistoryItemResponse{
			Time:   formatDate.Format(tracking.TimeFormat),
			Status: item.ActivityScan,
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

func (_this *shipmentResponse) GetLatestDeliveredPackage() (result TrackDetail) {
	if len(_this.TrackDetails) > 0 {
		sort.Slice(_this.TrackDetails, func(i, j int) bool {
			return _this.TrackDetails[i].GetDeliveryDate().After(_this.TrackDetails[j].GetDeliveryDate())
		})
		trackPackage := _this.TrackDetails[0]
		if trackPackage.TrackingNumber != "" {
			result = trackPackage
			return
		}
	}

	return
}

func (_this *shipmentResponse) IsSuccessful() bool {
	return _this.StatusCode == "200" && _this.StatusText == "Successful"
}

func (_this *shipmentResponse) GetWeight() int {
	if len(_this.TrackDetails) > 0 {
		additionalData := _this.TrackDetails[0].AdditionalInformation
		weightValue := additionalData["weight"]

		if weightValue == nil {
			return 0
		}

		weightStr := additionalData["weight"].(string)
		if weightStr == "" {
			return 0
		}

		weightFloat, err := strconv.ParseFloat(weightStr, 32)
		if err != nil {
			logger.AtLog.Warn(err)
		}
		return int(weightFloat)
	}
	return 0
}

func (_this *TrackDetail) GetDeliveryDate() (deliveryDate time.Time) {
	if _this.isDelivered() {
		deliveryDate = FormatDate(_this.DeliveredDate, _this.DeliveredTime)
		return
	}

	if _this.ScheduledDeliveryDate != "" {
		deliveryDate = FormatDate(_this.ScheduledDeliveryDate, TimeStrDefault)
	}

	return
}

func (_this *TrackDetail) isDelivered() bool {
	return _this.PackageStatus == DeliveredStatusStr && _this.DeliveredDate != "" && _this.DeliveredTime != ""
}

func (_this *TrackDetail) isInTransit() bool {
	if _this.isDelivered() {
		return false
	}
	return _this.ProgressBarType == InTransitText && _this.PackageStatusCode == InTransitCode
}

func (_this *TrackDetail) isPickup() bool {
	if _this.isDelivered() || _this.isInTransit() {
		return false
	}

	return _this.ShipmentProgressActivities != nil && len(_this.ShipmentProgressActivities) > 0
}

func (_this *TrackDetail) IsPrePickup() bool {
	if _this.isDelivered() || _this.isInTransit() || _this.isPickup() {
		return false
	}

	return _this.PackageStatusCode == LabelCreatedCode && _this.ShipmentProgressActivities == nil
}

func (_this *ShipmentActivity) HasData() bool {
	return _this.Date != "" && _this.Time != ""
}

func (_this *ShipmentActivity) FormatDate() (date time.Time) {
	if !_this.HasData() {
		return
	}

	return FormatDate(_this.Date, _this.Time)
}

func (_this *ShipmentActivity) FormatBeginDate() (date time.Time) {
	return FormatDate(_this.Date, TimeStrDefault)
}

func FormatDate(dateStr, timeStr string) time.Time {
	/*
	* dateStr := "04/07/2020"
	* timeStr := "3:58 A.M."
	*
	 */
	timeFormatted := strings.ReplaceAll(timeStr, ".", "")
	dateTimeStr := fmt.Sprintf("%s %s", dateStr, timeFormatted)

	date, err := time.Parse("01/02/2006 3:04 PM", dateTimeStr)
	if err != nil {
		logger.AtLog.Warn(err)
	}

	return date
}
