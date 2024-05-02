package models

import "time"

type Exoplanet struct {
	Id              int       `json:"id"`
	PlanetName      string    `json:"planetName"`
	HostName        string    `json:"hostName"`
	SystemNumber    int       `json:"systemNumber"`
	DiscoveryMethod string    `json:"discoveryMethod"`
	YearDiscovered  time.Time `json:"yearDiscovered"`
}
