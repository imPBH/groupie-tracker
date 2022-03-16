package main

import (
	"encoding/json"
	"fmt"
	"html/template"
	"io/ioutil"
	"log"
	"net/http"
	"strconv"
	"strings"
)

type Page struct {
	Title   string
	Artists *[]Artist
}

type ProfilePage struct {
	ArtistId int
	Artist   Artist
	Albums   []AlbumsApiData
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

type SearchApi struct {
	Data  []SearchApiData `json:"data"`
	Total int             `json:"total"`
	Next  string          `json:"next"`
}

type SearchApiData struct {
	Id        int    `json:"id"`
	Name      string `json:"name"`
	Link      string `json:"link"`
	PictureXl string `json:"picture_xl"`
	NbAlbum   int    `json:"nb_album"`
	NbFan     int    `json:"nb_fan"`
}

type ArtistApiData struct {
	Id   int    `json:"id"`
	Name string `json:"name"`
	Link string `json:"link"`
}

type AlbumsApiData struct {
	Id          int    `json:"id"`
	Title       string `json:"title"`
	Cover       string `json:"cover_xl"`
	Fans        int    `json:"fans"`
	ReleaseDate string `json:"release_date"`
	RecordType  string `json:"record_type"`
	Tracklist   string `json:"tracklist"`
}

type AlbumsApi struct {
	Data  []AlbumsApiData `json:"data"`
	Total int             `json:"total"`
}

var imagesURLs []string
var artist []Artist
var p = Page{
	Title:   "Groupie Tracker",
	Artists: &artist,
}

func main() {
	artist = GetArtists()

	artist := GetArtistApi("PNL")
	albums := GetArtistAlbums(artist.Id)
	fmt.Println(albums)
	fs := http.FileServer(http.Dir("templates"))
	router := http.NewServeMux()

	//Create a server and listen on port 8080
	fmt.Println("Starting server on port 8080")

	//Serve template for pages
	router.HandleFunc("/", HandlerIndex)
	router.HandleFunc("/homepage", HandlerHomepage)
	router.HandleFunc("/profile", HandlerProfile)
	router.HandleFunc("/profiledates", HandlerProfiledates)

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

func GetArtistApi(name string) SearchApiData {
	var dataApi SearchApi

	url := "https://api.deezer.com/search/artist/?q=" + strings.Replace(name, " ", "%20", 10) + "&index=0&limit=1"
	resp, err := http.Get(url)
	if err != nil {
		log.Fatalln(err)
	}

	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		log.Fatalln(err)
	}
	//Convert the body to type string
	sb := string(body)
	json.Unmarshal([]byte(sb), &dataApi)
	return dataApi.Data[0]
}

func GetArtistAlbums(artistId int) []AlbumsApiData {
	var dataApi AlbumsApi

	url := "https://api.deezer.com/artist/" + strconv.Itoa(artistId) + "/albums"
	resp, err := http.Get(url)
	if err != nil {
		log.Fatalln(err)
	}

	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		log.Fatalln(err)
	}
	//Convert the body to type string
	sb := string(body)
	json.Unmarshal([]byte(sb), &dataApi)
	return dataApi.Data
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
	t.ExecuteTemplate(w, "homepage.html", p)
}

func HandlerIndex(w http.ResponseWriter, r *http.Request) {
	t, _ := template.ParseGlob("templates/*.html")
	t.ExecuteTemplate(w, "index.html", p)
}

func HandlerProfile(w http.ResponseWriter, r *http.Request) {
	switch r.Method {
	case "GET":
		if err := r.ParseForm(); err != nil {
			fmt.Fprintf(w, "ParseForm() err: %v", err)
			return
		}
		artistIdString := r.FormValue("id")
		artistId, _ := strconv.Atoi(artistIdString)

		pProfile := ProfilePage{
			ArtistId: artistId,
			Artist:   artist[artistId-1],
			Albums:   GetArtistAlbums(GetArtistApi(artist[artistId-1].Name).Id),
		}
		t, _ := template.ParseGlob("templates/*.html")
		t.ExecuteTemplate(w, "profile.html", pProfile)
	}
}

func HandlerProfiledates(w http.ResponseWriter, r *http.Request) {
	switch r.Method {
	case "GET":
		if err := r.ParseForm(); err != nil {
			fmt.Fprintf(w, "ParseForm() err: %v", err)
			return
		}
		artistIdString := r.FormValue("id")
		artistId, _ := strconv.Atoi(artistIdString)
		pProfile := ProfilePage{
			Artist: artist[artistId-1],
		}
		t, _ := template.ParseGlob("templates/*.html")
		t.ExecuteTemplate(w, "profiledates.html", pProfile)
	}
}
