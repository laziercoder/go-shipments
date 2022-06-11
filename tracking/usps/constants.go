package usps

var (
	ApiPackageTracking       = "TrackV2"
	CharPrefix               = [3]string{"EU", "EC", "CP"}
	NumPrefix                = [2]int{8, 9}
	MinLengthTrackingHasChar = 13
	MinLengthTrackingHasNum  = 20
)
