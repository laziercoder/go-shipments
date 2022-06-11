package tracking

type History interface {
	MakeTrackShipmentResponse(packageCode string) (result TrackShipmentResponse, err error)
}
