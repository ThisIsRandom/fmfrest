package internal

import (
	"encoding/json"
	"errors"
	"math"
	"net/http"
	"net/url"
)

type Geo struct{}

type GetGeoConfig struct {
	Address string
	Postal  string
	Country string
}

type GR struct {
	Data []GL `json:"data"`
}

type GL struct {
	Latitude  float64 `json:"latitude"`
	Longitude float64 `json:"longitude"`
}

func (g *Geo) GetGeolocation(cfg *GetGeoConfig) (*GL, error) {
	baseURL, _ := url.Parse("http://api.positionstack.com")

	baseURL.Path += "v1/forward"

	params := url.Values{}
	params.Add("access_key", "d9d496301bba32aa799fb6d9775a1736")
	params.Add("query", "vestergade 42, 6051 almind")
	params.Add("output", "json")
	params.Add("limit", "1")

	baseURL.RawQuery = params.Encode()

	req, _ := http.NewRequest("GET", baseURL.String(), nil)

	res, _ := http.DefaultClient.Do(req)

	defer res.Body.Close()

	var x GR

	decoder := json.NewDecoder(res.Body)
	err := decoder.Decode(&x)

	if err != nil {
		return nil, err
	}

	if len(x.Data) < 1 {
		return nil, errors.New("dad")
	}

	v := x.Data[0].Latitude
	y := x.Data[0].Longitude

	return &GL{
		Longitude: y,
		Latitude:  v,
	}, nil
}

func distance(lat1 float64, lng1 float64, lat2 float64, lng2 float64, unit ...string) float64 {
	const PI float64 = 3.141592653589793

	radlat1 := float64(PI * lat1 / 180)
	radlat2 := float64(PI * lat2 / 180)

	theta := float64(lng1 - lng2)
	radtheta := float64(PI * theta / 180)

	dist := math.Sin(radlat1)*math.Sin(radlat2) + math.Cos(radlat1)*math.Cos(radlat2)*math.Cos(radtheta)

	if dist > 1 {
		dist = 1
	}

	dist = math.Acos(dist)
	dist = dist * 180 / PI
	dist = dist * 60 * 1.1515

	if len(unit) > 0 {
		if unit[0] == "K" {
			dist = dist * 1.609344
		} else if unit[0] == "N" {
			dist = dist * 0.8684
		}
	}

	return dist
}
