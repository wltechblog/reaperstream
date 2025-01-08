package main

type Car struct {
	ColorName        string `json:"ColorName"`
	EngineTimeDelay  string `json:"EngineTimeDelay"`
	MakerName        string `json:"MakerName"`
	ModelName        string `json:"ModelName"`
	NumSatellitesGPS string `json:"NumSatellitesGPS"`
	UseCacheGPS      string `json:"UseCacheGPS"`
	LicensePlate     string
	UUID             string
	Server           string
	Filename         string
}
