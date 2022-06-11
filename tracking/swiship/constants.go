package swiship

const (
	TrackingCodePrefix    = "TBA"
	TrackingCodeMinLength = 3
	deliveredStatus       = "DELIVERED"
	inTransitStatus       = "IN_TRANSIT"
)

var (
	// map "EVENTCode -> description
	// ref: https://github.com/amzn/selling-partner-api-docs/blob/main/references/fulfillment-outbound-api/fulfillmentOutbound_2020-07-01.md#"EVENTcode
	trackingEventDescription = map[string]string{
		"EVENT_101": "Carrier notified to pick up package",
		"EVENT_102": "Shipment picked up from seller's facility.",
		"EVENT_201": "Arrival scan.",
		"EVENT_202": "Departure scan.",
		"EVENT_203": "Arrived at destination country.",
		"EVENT_204": "Initiated customs clearance process.",
		"EVENT_205": "Completed customs clearance process.",
		"EVENT_206": "In transit to pickup location.",
		"EVENT_301": "Delivered.",
		"EVENT_302": "Out for delivery.",
		"EVENT_304": "Delivery attempted.",
		"EVENT_306": "Customer contacted to arrange delivery.",
		"EVENT_307": "Delivery appointment scheduled.",
		"EVENT_308": "Available for pickup.",
		"EVENT_309": "Returned to seller.",
		"EVENT_401": "Held by carrier - incorrect address.",
		"EVENT_402": "Customs clearance delay.",
		"EVENT_403": "Customer moved.",
		"EVENT_404": "Delay in delivery due to external factors.",
		"EVENT_405": "Shipment damaged.",
		"EVENT_406": "Held by carrier.",
		"EVENT_407": "Customer refused delivery.",
		"EVENT_408": "Returning to seller.",
		"EVENT_409": "Lost by carrier.",
		"EVENT_411": "Paperwork received - did not receive shipment.",
		"EVENT_412": "Shipment received - did not receive paperwork.",
		"EVENT_413": "Held by carrier - customer refused shipment due to customs charges.",
		"EVENT_414": "Missorted by carrier.",
		"EVENT_415": "Received from prior carrier.",
		"EVENT_416": "Undeliverable.",
		"EVENT_417": "Shipment missorted.",
		"EVENT_418": "Shipment delayed.",
		"EVENT_419": "Address corrected - delivery rescheduled.",
	}
)
