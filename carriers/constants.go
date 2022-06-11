package carriers

const (
	Fedex   = "FEDEX"
	UPS     = "UPS"
	Canpar  = "CANPAR"
	AIT     = "AIT"
	Auto    = "AUTO"
	Swiship = "swiship"
	DHL     = "dhl"
	USPS    = "usps"
)

var (
	trackingLink = map[string]string{
		Fedex:   "https://www.fedex.com/apps/fedextrack/?action=track&trackingnumber=%s&cntry_code=us&locale=en_US",
		UPS:     "https://www.ups.com/track?loc=en_US&tracknum=%s",
		Canpar:  "http://www.canpar.com/en/track/TrackingAction.do?locale=en&type=0&reference=%s",
		AIT:     "https://fastrak.aitworldwide.com/?TrackingNumber=%s",
		Swiship: "https://www.swiship.com/track?loc=en-US&id=%s",
	}

	CanScanByBot = []string{Fedex, Swiship}
)
