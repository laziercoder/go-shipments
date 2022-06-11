package ait

import (
	"context"
)

func getShipmentResponse(context context.Context, trackingCode string) *shipmentResponse {
	// url := fmt.Sprintf("%s/api/ait-carrier/%s/%s", setting.CurrentConfig().AutonomousAnalyticSite, trackingCode, setting.CurrentConfig().AutonomousAdminSecureCode)
	// resp, err := http.Get(url)
	// result := &shipmentResponse{}
	// if err != nil {
	// 	logger.AtLog.Warn(err)
	// 	return result
	// }
	//
	// defer resp.Body.Close()
	// data, err := ioutil.ReadAll(resp.Body)
	// if err != nil {
	// 	logger.AtLog.Warn(err)
	// 	return result
	// }
	//
	// err = json.Unmarshal(data, result)
	// if err != nil {
	// 	logger.AtLog.Warn(err)
	// }
	//
	// return result
	return nil
}
