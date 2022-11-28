package main

import (
	//stuff
	"fmt"
	//log any errors
	"log"
	//encode my data into json when send to postman
	"encoding/json"
	//new id by this method
	"math/rand"
	//server
	"net/http"
	//id will convert to string
	"strconv"
	//package github, gorilla
	"github.com/gorilla/mux"
)

type Movie struct {
	ID       string    `json:"id"`
	Isbn     string    `json:"isbn"`
	Title    string    `json: "title"`
	Director *Director `json: "director"`
}

type Director struct {
	Firstname string `json: "firstname"`
	Lastname  string `json: "lastname"`
}

var movies []Movie

// get movies function
func getMovies(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	//respon yang akan dikirim diubah -> encode json
	json.NewEncoder(w).Encode(movies)
}

// get movie function
func getMovie(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	params := mux.Vars(r) //get access to id
	//loop movies one by one
	for _, item := range movies {
		//fint that one
		if item.ID == params["id"] {
			//encode to json
			json.NewEncoder(w).Encode(item)
			return
		}
	}
}

// create movie function
func createMovie(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	var movie Movie
	//we will send data from postman in the BODY
	_ = json.NewDecoder(r.Body).Decode(&movie)
	movie.ID = strconv.Itoa(rand.Intn(100000000)) //di dapat id movie
	movies = append(movies, movie)
	json.NewEncoder(w).Encode(movie)

}

// delete function: easy - passing the ID
func deleteMovie(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	//id yang dilalui dari postman, akan dinotice sbg params di function ini
	params := mux.Vars(r)
	for index, item := range movies {
		if item.ID == params["id"] {
			movies = append(movies[:index], movies[index+1:]...) //apapun id yang mau dihapus
			break
		}
	}
	json.NewEncoder(w).Encode(movies)
}

// update function: complicated --will passing the id
func updateMovie(w http.ResponseWriter, r *http.Request) {
	//set json content type
	w.Header().Set("Content-Type", "application/json")
	//params
	params := mux.Vars(r)
	//loop over the movies, range
	//delete the movie with the id that you've sent
	//add a new movie - the movie that we send in the body of postman

	//loop
	for index, item := range movies {
		//we find
		if item.ID == params["id"] {
			//we delete
			movies = append(movies[:index], movies[index+1:]...)
			var movie Movie
			_ = json.NewDecoder(r.Body).Decode(&movie)
			movie.ID = params["id"]
			movies = append(movies, movie)
			json.NewEncoder(w).Encode(movie)
			return
		}
	}
}

func main() {
	//this is from mux liblary
	r := mux.NewRouter()

	movies = append(movies, Movie{ID: "1", Isbn: "438227", Title: "Movie One", Director: &Director{Firstname: "Zia", Lastname: "Milda"}})
	movies = append(movies, Movie{ID: "2", Isbn: "438238", Title: "Movie Two", Director: &Director{Firstname: "Dean", Lastname: "Naisu"}})
	//this is the list function
	r.HandleFunc("/movies", getMovies).Methods("GET")
	r.HandleFunc("/movies/{id}", getMovie).Methods("GET")
	r.HandleFunc("/movies", createMovie).Methods("POST")
	r.HandleFunc("/movies/{id}", updateMovie).Methods("PUT")
	r.HandleFunc("/movies/{id}", deleteMovie).Methods("DELETE")

	fmt.Printf("Starting server at port 8080\n")
	log.Fatal(http.ListenAndServe(":8080", r))
}
