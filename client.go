package shipment

type Client struct {
	TrackingCode                  string
	shipmentStrategy map[string]
}

func(c *Client) TrackShipment(trackingCode string, carrier string) {}