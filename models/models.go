package models

type IPResposne struct {
	IP  string
	Geo struct {
		Country    string
		City       string
		Latitude   float32
		Longitude  float32
		PostalCode string
		Region     string
	}
}

type IPPost struct {
	IPAddr []string
}
