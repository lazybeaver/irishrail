// Package irishrail implements a golang client for the IrishRail Realtime API
// Realtime API Documentation: http://api.irishrail.ie/realtime/
package irishrail

// Station information.
type Station struct {
	Code      string  `xml:"StationCode"`
	Name      string  `xml:"StationDesc"`
	Alias     string  `xml:"StationAlias"`
	Latitude  float64 `xml:"StationLatitude"`
	Longitude float64 `xml:"StationLongitude"`
}

// Train information.
type Train struct {
	Code      string  `xml:"TrainCode"`
	Date      string  `xml:"TrainDate"`
	Direction string  `xml:"Direction"`
	Status    string  `xml:"TrainStatus"`
	Message   string  `xml:"PublicMessage"`
	Latitude  float64 `xml:"TrainLatitude"`
	Longitude float64 `xml:"TrainLongitude"`
}

// StationDetail is a part of a sequence of StationDetails for a given train.
type StationDetail struct {
	StationCode            string `xml:"LocationCode"`
	StationName            string `xml:"LocationFullName"`
	StationOrder           int    `xml:"LocationOrder"`
	StationType            string `xml:"LocationType"`
	OriginName             string `xml:"TrainOrigin"`
	DestinationName        string `xml:"TrainDestination"`
	ScheduledArrivalTime   string `xml:"ScheduledArrival"`
	ScheduledDepartureTime string `xml:"ScheduledDeparture"`
	ExpectedArrivalTime    string `xml:"ExpectedArrival"`
	ExpectedDepartureTime  string `xml:"ExpectedDeparture"`
	ArrivalTime            string `xml:"Arrival"`
	DepartureTime          string `xml:"Departure"`
	StopType               string `xml:"StopType"`
}

// TrainDetail is a part of a sequence of TrainDetails for a given station.
type TrainDetail struct {
	TrainCode              string `xml:"Traincode"`
	TrainDate              string `xml:"Traindate"`
	TrainStatus            string `xml:"Status"`
	OriginName             string `xml:"Origin"`
	OriginTime             string `xml:"Origintime"`
	DestinationName        string `xml:"Destination"`
	DestinationTime        string `xml:"Destinationtime"`
	LastLocation           string `xml:"Lastlocation"`
	DueInMinutes           int    `xml:"Duein"`
	LateByMinutes          int    `xml:"Late"`
	ExpectedArrivalTime    string `xml:"Exparrival"`
	ExpectedDepartureTime  string `xml:"Expdepart"`
	ScheduledArrivalTime   string `xml:"Scharrival"`
	ScheduledDepartureTime string `xml:"Schdepart"`
	Direction              string `xml:"Direction"`
	TrainType              string `xml:"Traintype"`
	LocationType           string `xml:"Locationtype"`
}

// Client is the primary interface that offers access to the realtime API.
type Client interface {
	// Lookup station using name.
	LookupStation(string) (Station, error)

	// A list of stations.
	ListStations() ([]Station, error)

	// A list of trains (not all of these have running status).
	ListTrains() ([]Train, error)

	// A list of station details for a given train.
	GetStationDetails(Train) ([]StationDetail, error)

	// A list of train details for a given station.
	GetTrainDetails(Station) ([]TrainDetail, error)
}
