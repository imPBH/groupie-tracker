package main

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
)

type Artist struct {
	ID           int    `json:"id"`
	Image        string `json:"image"`
	Name         string `json:"name"`
	Members      string `json:"members"`
	CreationDate int    `json:"creationDate"`
	FirstAlbum   string `json:"firstAlbum"`
	Locations    string `json:"locations"`
	ConcertDates string `json:"concertDates"`
	Relations    string `json:"relations"`
}

type LocationIndex struct {
	Index []Location `json:"index"`
}

type Location struct {
	ID int `json:"id"`
	Locations []string `json:"locations"`
	Dates string `json:"dates"`
}

func main() {
	artist := GetArtists()
	locations := GetLocations()

	for i := range artist{
		fmt.Println(artist[i].Name)
	}

	for i := range locations.Index{
		fmt.Println(locations.Index[i].Locations)
	}

}

func GetArtists() []Artist {
	var artistData []Artist

	resp, err := http.Get("https://groupietrackers.herokuapp.com/api/artists")
	if err != nil {
		log.Fatalln(err)
	}

	//We Read the response body on the line below.
	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		log.Fatalln(err)
	}
	//Convert the body to type string
	sb := string(body)

	json.Unmarshal([]byte(sb), &artistData)

	return artistData
}

func GetLocations() LocationIndex{
	var artistLocation LocationIndex

	resp, err := http.Get("https://groupietrackers.herokuapp.com/api/locations")
	if err != nil {
		log.Fatalln(err)
	}

	//We Read the response body on the line below.
	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		log.Fatalln(err)
	}
	//Convert the body to type string
	sb := string(body)

	json.Unmarshal([]byte(sb), &artistLocation)

	return artistLocation
}