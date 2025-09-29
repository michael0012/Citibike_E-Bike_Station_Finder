package citibike

import (
	"net/http"
	"io"
	"encoding/json"
	"sort"
	"math"
	"fmt"
)


type Coordinate struct {
	Lat float64
	Lng float64
}


var startingCoordinate Coordinate

func distance(lat1 float64, lng1 float64, lat2 float64, lng2 float64, unit ...string) float64 {	
	radlat1 := float64(math.Pi * lat1 / 180)
	radlat2 := float64(math.Pi * lat2 / 180)
	
	theta := float64(lng1 - lng2)
	radtheta := float64(math.Pi * theta / 180)
	
	dist := math.Sin(radlat1) * math.Sin(radlat2) + math.Cos(radlat1) * math.Cos(radlat2) * math.Cos(radtheta);
	if dist > 1 {
		dist = 1
	}
	
	dist = math.Acos(dist)
	dist = dist * 180 / math.Pi
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

type CitiBikeStatusFeed struct{
	Name string `json:"name"`
	Url string	`json:"url"`
}

type CitiBikeStatus struct{
	LastUpdated int `json:"last_updated"`
	Ttl int `json:"ttl"`
	Data struct{
		En struct{
			Feeds []CitiBikeStatusFeed `json:"feeds"`
		} `json:"en"`
/*
		Es struct{
			Feeds []Feed `json:"feeds"`
			} `json:"es"`
		Fr struct{
			Feeds []Feed `json:"feeds"`
		} `json:"fr"`
*/
	}`json:"data"`
}

type StationInformation struct{
	LastUpdated int `json:"last_updated"`
	Ttl int `json:"ttl"`
	Data struct{
		Stations []struct{
			/*
			RentalUris struct{
				Ios string `json:"ios"`
				Android string `json:"android"`
			} `json:"rental_uris"`
			*/
			//ExternalId string `json:"external_id"`
			//LegacyId string `json:"legacy_id"`
			//EightdStationServices []string `json:"eightd_station_services"`
			//EightdHasKeyDispenser bool `json:"eightd_has_key_dispenser"`
			StationId string `json:"station_id"`
			Lon float64 `json:"lon"`
			//ShortName string `json:"short_name"`
			StationType string `json:"station_type"`
			Name string `json:"name"`
			Lat float64 `json:"lat"`
			//HasKiosk bool `json:"has_kiosk"`
			//ElectricBikeSurchargeWaiver bool `json:"electric_bike_surcharge_waiver"`
			//Capacity int `json:"capacity"`
			RentalMethods []string `json:"rental_methods"`
		} `json:"stations"`
	} `json:"data"`

}

type StationStatus struct{
	LastUpdated int `json:"last_updated"`
	Ttl int `json:"ttl"`
	Data struct{
		Stations []struct{
			//IsReturning int `json:"is_returning"`
			NumBikesAvailable int `json:"num_bikes_available"`
			NumEBikesAvailable int `json:"num_ebikes_available"`
			//LastReported int `json:"last_reported"`
			//NumDocksDisabled int `json:"num_docks_disabled"`
			//IsRenting int `json:"is_renting"`
			StationStatus string `json:"station_status"`
			//NumDocksAvailable int `json:"num_docks_available"`
			//NumBikesDisabled int `json:"num_bikes_disabled"`
			StationId string `json:"station_id"`
			//IsInstalled int `json:"is_installed"`
			//EightdHasAvailableKeys bool `json:"eightd_has_available_keys"`
			//LegacyId string `json:"legacy_id"`
		}`json:"stations"`
	}`json:"data"`
}

type StationDataMerged struct {
	StationId string
	Name string
	Lon float64
	Lat float64
	NumBikesAvailable int
	NumEBikesAvailable int
	StationStatus string
	RentalMethods []string
}

type StationDataArray []*StationDataMerged

func (stations StationDataArray) Less(i, j int) bool {
	startingLocation := startingCoordinate
	distanceI := distance(stations[i].Lat, stations[i].Lon, startingLocation.Lat, startingLocation.Lng, "N")
	distanceJ := distance(stations[j].Lat, stations[j].Lon, startingLocation.Lat, startingLocation.Lng, "N")
	return distanceI < distanceJ
}

func (stations StationDataArray) Len() int{
	return len(stations)
}

func (stations StationDataArray) Swap(i, j int) {
	stations[i], stations[j] = stations[j], stations[i]
}

func getJson(url string) []byte{
	resp, err := http.Get(url)
	if err != nil {
		panic(err)
	}
	body, err := io.ReadAll(resp.Body)
	if err != nil {
		panic(err)
	}
	resp.Body.Close()
	return body
}

var stationStatusURL string

// TODO: Cache Citibike data in db to pull later
func getCitiBikeStationsData() *StationInformation {
	baseURL := "https://gbfs.citibikenyc.com/gbfs/gbfs.json"
	var citiBikeStatus CitiBikeStatus
	var stationInformationBody StationInformation
	var stationInformationURL string
	body := getJson(baseURL)
	if err := json.Unmarshal(body, &citiBikeStatus); err != nil{
		panic(err)
	}
	for _, value := range citiBikeStatus.Data.En.Feeds{
		if value.Name == "station_information"{
			stationInformationURL = value.Url
		}else if value.Name == "station_status"{
			stationStatusURL = value.Url
		}
	}
	
	body = getJson(stationInformationURL)
	if err := json.Unmarshal(body, &stationInformationBody); err != nil{
		panic(err)
	}
	
	return &stationInformationBody
}


func getCitiBikeStationsWithAllElectric(numClassicBikes_ int, lat, long float64) []*StationDataMerged{
	var stationStatus StationStatus
	stationDataMerged := map[string]*StationDataMerged{}
	var freeBikeStations StationDataArray
	var body []byte
	var numClassicBikes int
	numClassicBikes = numClassicBikes_
	startingCoordinate = Coordinate{lat, long}
	stationInformation := getCitiBikeStationsData()

	for _, station := range stationInformation.Data.Stations{
		_, ok := stationDataMerged[station.StationId]
		if !ok{
			stationDataMerged[station.StationId] = &StationDataMerged{
				StationId: station.StationId,
				Name: station.Name,
				Lon: station.Lon,
				Lat: station.Lat,
				NumBikesAvailable: 0,
				NumEBikesAvailable: 0,
				StationStatus: "",
				RentalMethods: station.RentalMethods,
			}
		}
	}
	
	body = getJson(stationStatusURL)
	
	if err := json.Unmarshal(body, &stationStatus); err != nil{
		panic(err)
	}
	
	freeBikeStations = make([]*StationDataMerged, 0, len(stationStatus.Data.Stations))

	for _, station := range stationStatus.Data.Stations{
		_, ok := stationDataMerged[station.StationId]
		if ok {
			if station.StationStatus != "out_of_service"{
				stationDataMerged[station.StationId].NumBikesAvailable = station.NumBikesAvailable
				stationDataMerged[station.StationId].NumEBikesAvailable = station.NumEBikesAvailable
				stationDataMerged[station.StationId].StationStatus = station.StationStatus
				if classicBikesRemaining := stationDataMerged[station.StationId].NumBikesAvailable - stationDataMerged[station.StationId].NumEBikesAvailable; (stationDataMerged[station.StationId].NumEBikesAvailable != 0) && (classicBikesRemaining <= numClassicBikes){
					freeBikeStations = append(freeBikeStations, stationDataMerged[station.StationId])
				}
			} else{
				delete(stationDataMerged, station.StationId)
			}
		}
		
	}

	sort.Sort(StationDataArray(freeBikeStations))
	return freeBikeStations
}

func PrintCitiBikeStationsWithElectric(numClassicBikes int, lat, long float64, top int) {
	var freeBikeStations []*StationDataMerged
	freeBikeStations = getCitiBikeStationsWithAllElectric(numClassicBikes, lat, long)
	counter := 0
	for _, freeBikeStation := range freeBikeStations{
		if counter >= top {
			break
		}else{
			counter++
		}
		fmt.Printf("%s | classic bikes remaining: %d, electric bikes: %d\n", freeBikeStation.Name, freeBikeStation.NumBikesAvailable - freeBikeStation.NumEBikesAvailable, freeBikeStation.NumEBikesAvailable)
	}

}