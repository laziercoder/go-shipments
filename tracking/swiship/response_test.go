package swiship

import (
	"autonomous-service/dao"
	"autonomous-service/utils"
	"context"
	"testing"
)

func TestNewShipmentResponse(t *testing.T) {
	type args struct {
		context      context.Context
		trackingCode string
		gsDao        dao.GlobalSettingDao
	}
	tests := []struct {
		name string
		args args
	}{
		{
			name: "New",
			args: args{
				context:      utils.Background(nil),
				trackingCode: "TBA017695786904",
				gsDao:        dao.GlobalSettingDao{},
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := NewShipmentResponse(tt.args.context, tt.args.trackingCode, tt.args.gsDao)
			t.Logf("NewShipmentResponse() = %+v", got)

			resp, _ := got.MakeTrackShipmentResponse("")
			t.Logf("resp: %+v", resp)
		})
	}
}
