package operation

type User struct {
	Cid      string
	Password string
	Rating   int
}

type FlightPlan struct {
	Cid              string `json:"cid"`
	Callsign         string `json:"callsign"`
	FlightType       string `json:"flight_type"`
	AircraftType     string `json:"aircraft_type"`
	Tas              int    `json:"tas"`
	DepartureAirport string `json:"departure_airport"`
	DepartureTime    int    `json:"departure_time"`
	AtcDepartureTime int    `json:"atc_departure_time"`
	CruiseAltitude   string `json:"cruise_altitude"`
	ArrivalAirport   string `json:"arrival_airport"`
	RouteTimeHour    string `json:"route_time_hour"`
	RouteTimeMinute  string `json:"route_time_minute"`
	FuelTimeHour     string `json:"fuel_time_hour"`
	FuelTimeMinute   string `json:"fuel_time_minute"`
	AlternateAirport string `json:"alternate_airport"`
	Remarks          string `json:"remarks"`
	Route            string `json:"route"`
	Locked           bool   `json:"-"`
}
