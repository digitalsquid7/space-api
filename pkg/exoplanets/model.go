package exoplanets

type Exoplanet struct {
	Id         int    `json:"id"`
	PlanetName string `json:"planetName"`
	HostName   string `json:"hostName"`
}

type Response[T any] struct {
	Data []T `json:"data"`
}
