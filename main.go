package main

import (
	"encoding/json"
	"fmt"
	"html/template"
	"io/ioutil"
	"log"
	"net/http"
)

type Page struct {
	Title  string
	Artists *[]Artist
}

type Artist struct {
	Id           int      `json:"id"`
	Image        string   `json:"image"`
	Name         string   `json:"name"`
	Members      []string `json:"members"`
	CreationDate int      `json:"creationDate"`
	FirstAlbum   string   `json:"firstAlbum"`
	Locations    string   `json:"locations"`
	ConcertDates string   `json:"concertDates"`
	Relations    string   `json:"relations"`
}

type Location struct {
	Id        int      `json:"id"`
	Locations []string `json:"locations"`
	Dates     string   `json:"dates"`
}

type Relation struct {
	Id            int                 `json:"id"`
	DatesLocation map[string][]string `json:"datesLocations"`
}

var imagesURLs []string
var artist []Artist
var p = Page{
	Title:  "Groupie Tracker",
	Artists: &artist,
}


func main() {
	artist = GetArtists()


	for i := range artist {
		datesLocations := GetRelation(artist[i].Relations).DatesLocation
		locations := GetLocations(artist[i].Locations).Locations
		locations = removeDuplicateStr(locations)
		fmt.Printf("ID = %v\n", artist[i].Id)
		fmt.Printf("Name = %s\n", artist[i].Name)
		fmt.Printf("Creation date = %d\n", artist[i].CreationDate)
		fmt.Printf("Members = %+q\n", artist[i].Members)
		fmt.Printf("First Album = %s\n", artist[i].FirstAlbum)
		fmt.Printf("Locations URL = %s\n", artist[i].Locations)
		fmt.Printf("Dates URL = %s\n", artist[i].ConcertDates)
		fmt.Printf("Relations URL = %s\n", artist[i].Relations)
		fmt.Printf("Locations = %+q\n", GetLocations(artist[i].Locations).Locations)
		for _, location := range locations {
			fmt.Printf("Location = %s | Dates = %+q\n", location, datesLocations[location])
		}
		fmt.Printf("Image URL = %s\n\n", artist[i].Image)
		imagesURLs = append(imagesURLs, artist[i].Image)
	}

	fs := http.FileServer(http.Dir("templates"))
	router := http.NewServeMux()

	//Create a server and listen on port 8080
	fmt.Println("Starting server on port 8080")

	//Serve template for the homepage
	router.HandleFunc("/", HandlerHomepage)

	//Handle requests for files in /templates (ex : style.css)
	router.Handle("/templates/", http.StripPrefix("/templates/", fs))

	//Start server on port 8080
	http.ListenAndServe(":8080", router)
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

func GetLocations(url string) Location {
	var artistLocation Location

	resp, err := http.Get(url)
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

func GetRelation(url string) Relation {
	var relationData Relation

	resp, err := http.Get(url)
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

	json.Unmarshal([]byte(sb), &relationData)

	return relationData
}

func removeDuplicateStr(strSlice []string) []string {
	var list []string
	allKeys := make(map[string]bool)
	for _, item := range strSlice {
		if _, value := allKeys[item]; !value {
			allKeys[item] = true
			list = append(list, item)
		}
	}
	return list
}

func HandlerHomepage(w http.ResponseWriter, r *http.Request) {
	t, _ := template.ParseGlob("templates/*.html")
	t.ExecuteTemplate(w, "index.html", p)
}
