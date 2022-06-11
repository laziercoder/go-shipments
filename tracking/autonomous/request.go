package autonomous

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"gitlab.com/autonomous-ecm/backend/go-autonomous/pkg/logger"
	"io/ioutil"
	"net/http"
)

func getShipmentResponse(context context.Context, trackingCode string) *shipmentResponse {
	url := fmt.Sprintf("https://api.optimoroute.com/v1/get_completion_details?key=%s", "02bedd8727c6d4807838e26c8f7db5a4eiSLF5wXnug")
	dataPost := fmt.Sprintf(`{
	  "orders": [
		{
		  "orderNo": "%s"
		}
	  ]
	}`, trackingCode)

	result := &shipmentResponse{}
	resp, err := http.Post(url, "application/json", bytes.NewBufferString(dataPost))
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
