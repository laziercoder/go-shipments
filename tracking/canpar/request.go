package canpar

import (
	"context"
)

func getShipmentResponse(context context.Context, trackingCode string) *shipmentResponse {
	// password := "autonomous123456789abc"
	// url := fmt.Sprintf("%s/api/canpar-tracking/%s/%s", setting.CurrentConfig().AutonomousAnalyticSite, trackingCode, password)
	//
	// result := &shipmentResponse{}
	// resp, err := http.Get(url)
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
	// logger.AtLog.Infof("DataResp: %+v, %+v", result, err)
	// return result
	return nil
}
