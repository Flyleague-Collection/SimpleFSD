package operation

type User struct {
	Cid      string
	Password string
	Rating   int
}

type FlightPlan struct {
	Cid              string
	Callsign         string
	FlightType       string
	AircraftType     string
	Tas              int
	DepartureAirport string
	DepartureTime    int
	AtcDepartureTime int
	CruiseAltitude   string
	ArrivalAirport   string
	RouteTimeHour    string
	RouteTimeMinute  string
	FuelTimeHour     string
	FuelTimeMinute   string
	AlternateAirport string
	Remarks          string
	Route            string
	Locked           bool
}
