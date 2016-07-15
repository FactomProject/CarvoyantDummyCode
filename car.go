package main

import (
	"crypto/sha256"
	"encoding/json"
	"io/ioutil"
	"log"
	"net/http"

	"github.com/FactomProject/factom"
)

type CarDataSet struct {
	DataSet []struct {
		ID             int    `json:"id"`
		VehicleID      int    `json:"vehicleId"`
		TripID         int    `json:"tripId"`
		Timestamp      string `json:"timestamp"`
		IgnitionStatus string `json:"ignitionStatus"`
		Datum          []struct {
			ID              int    `json:"id"`
			Timestamp       string `json:"timestamp"`
			Key             string `json:"key"`
			Value           string `json:"value"`
			TranslatedValue string `json:"translatedValue"`
		} `json:"datum"`
	} `json:"dataSet"`
	TotalRecords int `json:"totalRecords"`
}

func main() {
	// Construct API request to gather most recent data
	req, err := http.NewRequest("GET", "https://sandbox-api.carvoyant.com/sandbox/api/vehicle/252773/dataSet/?mostRecentOnly=1", nil)
	if err != nil {
		log.Fatal(err)
	}
	// Replace with correct authorization token
	req.Header.Set("Authorization", "Bearer XXXXXXXXXXXXXXXXXXXXXXXX")

	// Send http request
	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		log.Fatal(err)
	}
	defer resp.Body.Close()

	rawVehicleData, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		log.Fatal(err)
	}

	vehicleData := new(CarDataSet)
	json.Unmarshal(rawVehicleData, vehicleData)

	ec := factom.NewECAddress()

	// Hash the data, record which data points are hashed by their timestamps.
	entry := new(factom.Entry)
	entry.ExtIDs = make([][]byte, 2)
	entry.ExtIDs[0] = []byte(vehicleData.DataSet[0].Timestamp)
	entry.ExtIDs[1] = []byte(vehicleData.DataSet[len(vehicleData.DataSet)-1].Timestamp)

	contentHash := sha256.Sum256(rawVehicleData)
	entry.Content = contentHash[:]

	factom.CommitEntry(entry, ec)
}
