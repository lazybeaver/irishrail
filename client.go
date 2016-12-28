package irishrail

import (
	"bytes"
	"encoding/xml"
	"fmt"
	"net/http"
	"net/url"
	"strings"
	"time"
	"unicode"
)

// Options used to create an irishRailClient object.
type Options struct {
	URL     string
	Timeout time.Duration
}

var (
	// DefaultOptions points to the IrishRail Realtime API.
	DefaultOptions = Options{URL: "http://api.irishrail.ie/realtime/realtime.asmx", Timeout: 5 * time.Second}

	// LocalServerOptions serves from cached local files.
	LocalServerOptions = Options{URL: "http://127.0.0.1:8080", Timeout: 1 * time.Second}
)

// This type implements the methods required by the "Client" interface.
type irishRailClient struct {
	baseURL    string
	httpclient http.Client
}

// NewClient returns an irishrail.Client interface implementation.
func NewClient(options Options) Client {
	return &irishRailClient{baseURL: options.URL, httpclient: http.Client{Timeout: options.Timeout}}
}

func normalize(name string) string {
	buf := bytes.Buffer{}
	for _, r := range strings.ToLower(name) {
		if unicode.IsLower(r) {
			buf.WriteRune(r)
		}
	}
	return buf.String()
}

func (c *irishRailClient) LookupStation(name string) (Station, error) {
	stations, err := c.ListStations()
	if err != nil {
		return Station{}, err
	}
	normname := normalize(name)
	for _, station := range stations {
		if normname == normalize(station.Name) || normname == normalize(station.Alias) || normname == normalize(station.Code) {
			return station, nil
		}
	}
	return Station{}, fmt.Errorf("Failed to lookup station: %s", name)
}

func (c *irishRailClient) getXML(path string, result interface{}) error {
	response, err := c.httpclient.Get(c.baseURL + path)
	if err != nil {
		return err
	}
	defer response.Body.Close()
	return xml.NewDecoder(response.Body).Decode(result)
}

func (c *irishRailClient) ListStations() ([]Station, error) {
	result := struct {
		XMLName  xml.Name  `xml:"ArrayOfObjStation"`
		Stations []Station `xml:"objStation"`
	}{}
	if err := c.getXML("/getAllStationsXML", &result); err != nil {
		return nil, err
	}
	for i := 0; i < len(result.Stations); i++ {
		sanitizeStation(&result.Stations[i])
	}
	return result.Stations, nil
}

func sanitizeStation(station *Station) {
	station.Code = strings.TrimSpace(station.Code)
}

func (c *irishRailClient) ListTrains() ([]Train, error) {
	result := struct {
		XMLName xml.Name `xml:"ArrayOfObjTrainPositions"`
		Trains  []Train  `xml:"objTrainPositions"`
	}{}
	if err := c.getXML("/getCurrentTrainsXML", &result); err != nil {
		return nil, err
	}
	for i := 0; i < len(result.Trains); i++ {
		sanitizeTrain(&result.Trains[i])
	}
	return result.Trains, nil
}

func sanitizeTrain(train *Train) {
	train.Code = strings.TrimSpace(train.Code)
	train.Message = strings.TrimSpace(strings.TrimPrefix(strings.Replace(train.Message, "\\n", " | ", -1), train.Code))
}

func (c *irishRailClient) GetStationDetails(train Train) ([]StationDetail, error) {
	result := struct {
		XMLName        xml.Name        `xml:"ArrayOfObjTrainMovements"`
		StationDetails []StationDetail `xml:"objTrainMovements"`
	}{}
	path := fmt.Sprintf("/getTrainMovementsXML?TrainId=%s&TrainDate=%s", train.Code, url.QueryEscape(train.Date))
	if err := c.getXML(path, &result); err != nil {
		return nil, err
	}
	for i := 0; i < len(result.StationDetails); i++ {
		sanitizeStationDetail(&result.StationDetails[i])
	}
	return result.StationDetails, nil
}

func sanitizeStationDetail(sd *StationDetail) {
	if sd.StationType != "S" && sd.StationType != "T" && sd.StationType != "O" && sd.StationType != "D" {
		sd.StationType = ""
	}
	if sd.StopType != "C" && sd.StopType != "N" {
		sd.StopType = ""
	}
}

func (c *irishRailClient) GetTrainDetails(station Station) ([]TrainDetail, error) {
	result := struct {
		XMLName      xml.Name      `xml:"ArrayOfObjStationData"`
		TrainDetails []TrainDetail `xml:"objStationData"`
	}{}
	path := fmt.Sprintf("/getStationDataByCodeXML?StationCode=%s", station.Code)
	if err := c.getXML(path, &result); err != nil {
		return nil, err
	}
	for i := 0; i < len(result.TrainDetails); i++ {
		sanitizeTrainDetail(&result.TrainDetails[i])
	}
	return result.TrainDetails, nil
}

func sanitizeTrainDetail(td *TrainDetail) {
	td.TrainCode = strings.TrimSpace(td.TrainCode)
}
