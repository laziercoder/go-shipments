package swiship

import (
	"context"
	"encoding/json"
	"fmt"
	"gitlab.com/autonomous-ecm/backend/go-autonomous/pkg/logger"
	"io/ioutil"
	"net/http"
	"strings"
)

func getShipmentResponse(context context.Context, trackingCode string) *shipmentResponse {
	body := fmt.Sprintf(`{"trackingNumber":"%s","shipMethod":"AMZL_US_PREMIUM"}`, trackingCode)
	resp, err := http.Post(
		"https://www.swiship.com/api/getPackageTrackingDetails",
		"application/json",
		strings.NewReader(body),
	)
	result := &shipmentResponse{}
	if err != nil {
		logger.AtLog.Warn(err)
		return result
	}

	defer resp.Body.Close()
	data, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		logger.AtLog.Warn(err)
		return result
	}

	err = json.Unmarshal(data, result)
	if err != nil {
		logger.AtLog.Warn(err)
	}

	return result
}
